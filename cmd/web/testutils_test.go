package main

import (
	"bytes"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/qqquinnn/snippetbox/internal/models/mocks"
)

// Returns an instance of application struct w/ mocked dependencies.
func newTestApplication(t *testing.T) *application {
	// Create instances of template cache, form decoder, & session manager.

	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// Return struct containing mocked dependencies.
	return &application{
		logger:         slog.New(slog.DiscardHandler),
		snippets:       &mocks.SnippetModel{},
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

// Embeds an httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// Initializes & returns new instance of testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	// Initialize test server.
	ts := httptest.NewTLSServer(h)

	// Initialize cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add cookie jar to test server client. Response cookies will be
	// stored in jar and sent w/ subsequent requests.
	ts.Client().Jar = jar

	// Prevent test server client from following redirects.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// Resets test server client to use a new & empty cookie jar.
func (ts *testServer) resetClientCookieJar(t *testing.T) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar
}

// Holds data about responses from test server.
type testResponse struct {
	status  int
	headers http.Header
	cookies []*http.Cookie
	body    string
}

// Makes a GET request to a given URL path using test server client.
// Returns testResponse struct containing response data.
func (ts *testServer) get(t *testing.T, urlPath string) testResponse {
	req, err := http.NewRequest(http.MethodGet, ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return testResponse{
		status:  res.StatusCode,
		headers: res.Header,
		cookies: res.Cookies(),
		body:    string(bytes.TrimSpace(body)),
	}
}

// Makes a POST request to a given URL path using test server client.
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) testResponse {
	req, err := http.NewRequest(http.MethodPost, ts.URL+urlPath, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	// Set appropriate Content-Type header for form data, & Sec-Fetch-Site header.
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	return testResponse{
		status:  res.StatusCode,
		headers: res.Header,
		cookies: res.Cookies(),
		body:    string(bytes.TrimSpace(body)),
	}
}

// Extracts the CSRF token from an HTML response body.
func extractCSRFToken(t *testing.T, body string) string {
	// Regex to capture CSRF token value from HTML for user signup page.
	csrfTokenRX := regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)">`)

	// Extract token from HTML body.
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(matches[1])
}
