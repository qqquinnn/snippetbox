package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/qqquinnn/snippetbox/internal/models"
	"github.com/qqquinnn/snippetbox/ui"
)

// Acts as the holding structure for any dynamic data that needs
// to be passed to the HTML templates.
type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
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

	// Get a slice of all filepaths in the embedded filesystem
	// matching the pattern "./ui/html/pages/*.html".
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Loop through filepaths one by one.
	for _, page := range pages {
		// Extract filename from full filepath and assign it to a variable.
		name := filepath.Base(page)

		// Create a slice containing filepath patterns for the desired templates.
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		// Parse template files into a template set.
		templateSet, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add template set to map.
		cache[name] = templateSet
	}

	return cache, nil
}
