package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Server", "Go")

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
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// We then write the template content as the response body.
	err = templateSet.ExecuteTemplate(writer, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func snippetView(writer http.ResponseWriter, request *http.Request) {
	// Extract value of "id" wildcard from request and try to convert to integer.
	// If conversion unsuccessful or value < 1, return 404 response.
	id, err := strconv.Atoi(request.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(writer, request)
		return
	}

	// Interpolate id value with a message, then write as HTTP response.
	fmt.Fprintf(writer, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Display a form for creating a new snippet..."))
}

func snippetCreatePost(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte("Save a new snippet..."))
}
