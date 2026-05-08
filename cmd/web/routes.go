package main

import (
	"net/http"

	"github.com/qqquinnn/snippetbox/ui"

	"github.com/justinas/alice"
)

// Returns a servemux containing the application routes.
func (app *application) routes() http.Handler {
	// Initialize a new servemux.
	mux := http.NewServeMux()

	// Create file server to serve the embedded files in ui.Files, and
	// register file server as the handler for all paths starting with "/static".
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// Route for health check.
	mux.HandleFunc("GET /ping", ping)

	// Unprotected application routes using "dynamic" middleware chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave, preventCSRF, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected (auth-only) application routes using "protected" middleware chain.
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("GET /account/view", protected.ThenFunc(app.accountView))
	mux.Handle("GET /account/password/update", protected.ThenFunc(app.accountPasswordUpdate))
	mux.Handle("POST /account/password/update", protected.ThenFunc(app.accountPasswordUpdatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// Create and return middleware chain containing 'standard' middleware
	// which will be used for every request the application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
