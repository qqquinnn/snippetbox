package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qqquinnn/snippetbox/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	// Initialize response recorder & dummy http.Request.
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP handler; writes 200 status code & "OK" response body.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Pass mock handler to commonHeaders middleware.
	commonHeaders(next).ServeHTTP(rr, req)

	// Get results of the test.
	res := rr.Result()
	defer res.Body.Close()

	// CHECK FOR CORRECT HEADER SETTINGS...

	// Content-Security-Policy
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, res.Header.Get("Content-Security-Policy"), expectedValue)

	// Referrer-Policy
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, res.Header.Get("Referrer-Policy"), expectedValue)

	// X-Content-Type-Options
	expectedValue = "nosniff"
	assert.Equal(t, res.Header.Get("X-Content-Type-Options"), expectedValue)

	// X-Frame-Options
	expectedValue = "deny"
	assert.Equal(t, res.Header.Get("X-Frame-Options"), expectedValue)

	// X-XSS-Protection
	expectedValue = "0"
	assert.Equal(t, res.Header.Get("X-XSS-Protection"), expectedValue)

	// Check that the next handler in line is called and the response status code
	// & body are as expected.
	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
