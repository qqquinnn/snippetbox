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

	// Unprotected application routes using "dynamic" middleware chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave, preventCSRF)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected (auth-only) application routes using "protected" middleware chain.
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// Create and return middleware chain containing 'standard' middleware
	// which will be used for every request the application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
