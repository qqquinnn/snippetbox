package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/qqquinnn/snippetbox/internal/models"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold application-wide dependencies.
type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	// Define command-line flags.
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:userpassword@/snippetbox?parseTime=true", "MySQL data source name")

	// Parse the command-line flags.
	flag.Parse()

	// Initialize a structured logger which writes to stdout with default settings.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Pass openDB() the DSN from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Defer a call to db.Close() so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize decoder instance.
	formDecoder := form.NewDecoder()

	// Initialize new instance of the application struct.
	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	// Print log message to indicate server is starting.
	logger.Info("starting server", "addr", *addr)

	// Use the http.ListenAndServe() function to start a new web server.
	// We pass two parameters: the TCP network address to listen on (default: ":4000")
	// and the servemux from routes.go.
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// Wraps sql.Open() and returns a sql.DB connection pool.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
