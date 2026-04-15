package main

import (
	"errors"
	"fmt"

	//"html/template"
	"net/http"
	"strconv"

	"github.com/qqquinnn/snippetbox/internal/models"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Server", "Go")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(writer, request, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(writer, "%+v\n", snippet)
	}

	/*
		// Initialize a slice containing the template file paths.
		// The file containing the base template must be the first file.
		templateFiles := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			"./ui/html/pages/home.html",
		}

		// Read HTML template files into a template set. If there's an error, we log
		// the error message, send an Internal Server Error response, and return
		// from the handler.
		templateSet, err := template.ParseFiles(templateFiles...)
		if err != nil {
			app.serverError(writer, request, err)
			return
		}

		// We then write the template content as the response body.
		err = templateSet.ExecuteTemplate(writer, "base", nil)
		if err != nil {
			app.serverError(writer, request, err)
		}
	*/
}

func (app *application) snippetView(writer http.ResponseWriter, request *http.Request) {
	// Extract value of "id" wildcard from request and try to convert to integer.
	// If conversion unsuccessful or value < 1, return 404 response.
	id, err := strconv.Atoi(request.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(writer, request)
		return
	}

	// Retrieve data from SnippetModel.Get() method for a specific record.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(writer, request)
		} else {
			app.serverError(writer, request, err)
		}
		return
	}

	// Write snippet data as plain-text HTTP response body.
	fmt.Fprintf(writer, "%+v", snippet)
}

func (app *application) snippetCreate(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Display a form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(writer http.ResponseWriter, request *http.Request) {
	// Dummy data for development.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	// Pass data to the SnippetModel.Insert() method.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(writer, request, err)
		return
	}

	// Redirect user to relevant page for the snippet.
	http.Redirect(writer, request, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
