package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/qqquinnn/snippetbox/internal/models"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(writer, request, err)
		return
	}

	// Call helper to get a templateData struct containing the default data,
	// and add snippets slice.
	data := app.newTemplateData(request)
	data.Snippets = snippets

	app.render(writer, request, http.StatusOK, "home.html", data)
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

	// Call helper to get a templateData struct containing the default data,
	// and add snippets slice.
	data := app.newTemplateData(request)
	data.Snippet = snippet

	app.render(writer, request, http.StatusOK, "view.html", data)
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
