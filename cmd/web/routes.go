package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// Returns a servemux containing the application routes.
func (app *application) routes() http.Handler {
	// Initialize a new servemux.
	mux := http.NewServeMux()

	// Create file server to serve files out of "./ui/static" directory, and
	// register file server as the handler for all paths starting with "/static".
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// Create new middleware chain for dynamic application routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Register other functions and URL patterns.
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))                          // Display home page
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))     // Display specific snippet
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))      // Display form for creating new snippet
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost)) // Save new snippet

	// Create and return middleware chain containing 'standard' middleware
	// which will be used for every request the application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
