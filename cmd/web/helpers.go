package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
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

// Returns a templateData struct initialized with the current year and any flash messages.
func (app *application) newTemplateData(request *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
		Flash:       app.sessionManager.PopString(request.Context(), "flash"),
	}
}

// Retrieves the appropriate template set from cache and writes to the ResponseWriter.
func (app *application) render(writer http.ResponseWriter, request *http.Request, status int, page string, data templateData) {
	// Retrieve template set based on page name.
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

// Handles decoding of HTML form data to a target destination.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call method on the decoder instance, with the target destination as
	// the first parameter.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If target destination is invalid, Decode() returns form.InvalidDecoderError.
		// Panic if this occurs.
		if _, ok := errors.AsType[*form.InvalidDecoderError](err); ok {
			panic(err)
		}

		// Return all other errors normally.
		return err
	}

	return nil
}
