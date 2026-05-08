package main

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/qqquinnn/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	// Create application struct (w/ mocked dependencies) & new test server.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Execute HTTP request against test server.
	res := ts.get(t, "/ping")

	// Check value of response status code & body.
	assert.Equal(t, res.status, http.StatusOK)
	assert.Equal(t, res.body, "OK")
}

func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		ts.resetClientCookieJar(t)

		res := ts.get(t, "/snippet/create")
		assert.Equal(t, res.status, http.StatusSeeOther)
		assert.Equal(t, res.headers.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		ts.resetClientCookieJar(t)

		// Make GET /user/login request.
		res := ts.get(t, "/user/login")

		// Make POST /user/login request using credentials from mock user
		// model & CSRF token from previous response.
		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "pa$$word")
		form.Add("csrf_token", extractCSRFToken(t, res.body))

		ts.postForm(t, "/user/login", form)

		// Check that authenticated user is shown create snippet form.
		res = ts.get(t, "/snippet/create")
		assert.Equal(t, res.status, http.StatusOK)
		assert.True(t, strings.Contains(res.body, `<form action="/snippet/create" method="POST">`))
	})
}

func TestSnippetView(t *testing.T) {
	// Create application struct (w/ mocked dependencies) & new test server.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Table-driven tests to check application responses for different URLs.
	tests := []struct {
		name       string
		urlPath    string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Valid ID",
			urlPath:    "/snippet/view/1",
			wantStatus: http.StatusOK,
			wantBody:   "An old silent pond...",
		},
		{
			name:       "Non-existent ID",
			urlPath:    "/snippet/view/2",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Negative ID",
			urlPath:    "/snippet/view/-1",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Decimal ID",
			urlPath:    "/snippet/view/1.23",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "String ID",
			urlPath:    "/snippet/view/foo",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Empty ID",
			urlPath:    "/snippet/view/",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset test server cookie jar at start of each sub-test.
			ts.resetClientCookieJar(t)

			// Make GET request to specified URL path.
			res := ts.get(t, tt.urlPath)

			// Check response status & response body.
			assert.Equal(t, res.status, tt.wantStatus)
			assert.True(t, strings.Contains(res.body, tt.wantBody))
		})
	}
}

func TestUserSignup(t *testing.T) {
	// Create application struct (w/ mocked dependencies) & new test server.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Field values for table-driven tests.
	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = `<form action="/user/signup" method="POST" novalidate>`
	)

	// Table-driven tests to check behavior of POST /user/signup route.
	tests := []struct {
		name              string
		userName          string
		userEmail         string
		userPassword      string
		useValidCSRFToken bool
		wantStatus        int
		wantFormTag       string
	}{
		{
			name:              "Valid submission",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusSeeOther,
		},
		{
			name:              "Invalid CSRF Token",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      validPassword,
			useValidCSRFToken: false,
			wantStatus:        http.StatusBadRequest,
		},
		{
			name:              "Empty name",
			userName:          "",
			userEmail:         validEmail,
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Empty email",
			userName:          validName,
			userEmail:         "",
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Empty password",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      "",
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Invalid email",
			userName:          validName,
			userEmail:         "bob@example.",
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Short password",
			userName:          validName,
			userEmail:         validEmail,
			userPassword:      "pa$$",
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
		{
			name:              "Duplicate email",
			userName:          validName,
			userEmail:         "dupe@example.com",
			userPassword:      validPassword,
			useValidCSRFToken: true,
			wantStatus:        http.StatusUnprocessableEntity,
			wantFormTag:       formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset test server cookie jar at start of each sub-test.
			ts.resetClientCookieJar(t)

			// Make GET /user/signup request; adds CSRF cookie to cookie jar.
			res := ts.get(t, "/user/signup")

			// Build up form values for sub-test.
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			if tt.useValidCSRFToken {
				form.Add("csrf_token", extractCSRFToken(t, res.body))
			}

			// Make POST /user/signup request using form values.
			res = ts.postForm(t, "/user/signup", form)

			// Test response data.
			assert.Equal(t, res.status, tt.wantStatus)
			assert.True(t, strings.Contains(res.body, tt.wantFormTag))
		})
	}
}
