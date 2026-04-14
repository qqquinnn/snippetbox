package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize a new servemux.
	mux := http.NewServeMux()

	// Create file server to serve files out of "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Register file server as the handler for all paths starting with "/static".
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// Register other functions and URL patterns.
	mux.HandleFunc("GET /{$}", home)                          // Display home page
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)     // Display specific snippet
	mux.HandleFunc("GET /snippet/create", snippetCreate)      // Display form for creating new snippet
	mux.HandleFunc("POST /snippet/create", snippetCreatePost) // Save new snippet

	// Print log message to indicate server is starting.
	log.Print("starting server on :4000")

	// Use the http.ListenAndServe() function to start a new web server.
	// We pass two parameters: the TCP network address to listen on (":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and terminate the
	// program. Note that any error returned by http.ListenAndServe() is always
	// non-nil.
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
