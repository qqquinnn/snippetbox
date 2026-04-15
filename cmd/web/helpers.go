package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
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

// Returns a templateData struct initialized with the current year.
func (app *application) newTemplateData(request *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *application) render(writer http.ResponseWriter, request *http.Request, status int, page string, data templateData) {
	// Retrieve appropriate template set from cache based on page name.
	templateSet, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(writer, request, err)
		return
	}

	// Initialize new buffer.
	buf := new(bytes.Buffer)

	// Execute template set and write to buffer.
	err := templateSet.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(writer, request, err)
		return
	}

	// If template is written to buffer without errors,
	// write out provided HTTP status code to ResponseWriter.
	writer.WriteHeader(status)

	// Write contents of buffer to ResponseWriter
	buf.WriteTo(writer)
}
