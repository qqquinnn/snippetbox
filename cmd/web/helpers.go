package main

import (
	"net/http"
)

// Writes a log entry at Error level, then sends generic 500 Internal Server
// Error response to user.
func (app *application) serverError(writer http.ResponseWriter, request *http.Request, err error) {
	var (
		method = request.Method
		uri    = request.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Sends a specific status code & corresponding description to user.
func (app *application) clientError(writer http.ResponseWriter, status int) {
	http.Error(writer, http.StatusText(status), status)
}
