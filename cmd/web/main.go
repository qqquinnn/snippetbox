package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Define an application struct to hold application-wide dependencies.
type application struct {
	logger *slog.Logger
}

func main() {
	// Define a command-line flag with the name 'addr', a default value of ":4000",
	// and some help text.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Parse the command-line flag.
	flag.Parse()

	// Initialize a structured logger which writes to stdout with default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize new instance of the application struct.
	app := &application{
		logger: logger,
	}

	// Print log message to indicate server is starting.
	logger.Info("starting server", "addr", *addr)

	// Use the http.ListenAndServe() function to start a new web server.
	// We pass two parameters: the TCP network address to listen on (default: ":4000")
	// and the servemux from routes.go.
	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
