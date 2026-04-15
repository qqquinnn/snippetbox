package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/qqquinnn/snippetbox/internal/models"
)

// Acts as the holding structure for any dynamic data that needs
// to be passed to the HTML templates.
type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
}

// Returns formatted string representation of a time.Time value.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04 UTC")
}

// Initialize a template.FuncMap value and store it in a global variable.
// Acts as a lookup table mapping names to functions.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize new map to act as the cache.
	cache := map[string]*template.Template{}

	// Get a slice of all filepaths matching the pattern "./ui/html/pages/*.html".
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Loop through filepaths one by one.
	for _, page := range pages {
		// Extract filename from full filepath and assign it to a variable.
		name := filepath.Base(page)

		// Parse base template file into a template set.
		templateSet, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() on this template set to add partials.
		templateSet, err = templateSet.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() on this template set to add the page tempate.
		templateSet, err = templateSet.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add template set to map.
		cache[name] = templateSet
	}

	return cache, nil
}
