package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// Returns a servemux containing the application routes.
func (app *application) routes() http.Handler {
	// Initialize a new servemux.
	mux := http.NewServeMux()

	// Create file server to serve files out of "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Register file server as the handler for all paths starting with "/static".
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// Register other functions and URL patterns.
	mux.HandleFunc("GET /{$}", app.home)                          // Display home page
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)     // Display specific snippet
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)      // Display form for creating new snippet
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost) // Save new snippet

	// Create and return middleware chain containing 'standard' middleware
	// which will be used for every request the application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
