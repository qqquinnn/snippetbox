package main

import (
	"context"
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/qqquinnn/snippetbox/internal/models"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv5"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

// Define an application struct to hold application-wide dependencies.
type application struct {
	debug          bool
	logger         *slog.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	ctx := context.Background()
	// Load .env file & define variables for local development.
	_ = godotenv.Load()
	dsn := os.Getenv("SNIPPETBOX_DSN")
	addr := ":" + os.Getenv("PORT")
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	prodStr := os.Getenv("IS_PROD")
	prod, err := strconv.ParseBool(prodStr)

	// Panic if DSN config is missing.
	if dsn == "" {
		panic("database DSN must be provided via SNIPPETBOX_DSN env var")
	}

	// Initialize a structured logger.
	var logger *slog.Logger

	if !prod {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	} else {
		// Configure minimum log level.
		logLevel := slog.LevelInfo
		if debug {
			logLevel = slog.LevelDebug
		}

		// Initialize a structured logger which writes JSON to stdout.
		// Replace the default "level" key with "severity" for Google Cloud Logging compatibility.
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					a.Key = "severity"
				}
				return a
			},
		}))
	}

	// Pass openDB() the DSN from the command-line flag.
	db, cleanup, err := openDB(dsn, prod)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Defer a call to close the connection pool & terminate Cloud SQL connector's
	// background goroutines before the main function exits.
	defer db.Close()
	defer cleanup()

	// Initialize template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize decoder instance.
	formDecoder := form.NewDecoder()

	// Initialize and configure new session manager.
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// Initialize new instance of the application struct.
	app := &application{
		debug:          debug,
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize a new http.Server struct.
	srv := &http.Server{
		Addr:     addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
		// Idle, Read & Write timeouts.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Print log message to indicate server is starting.
	logger.Info("starting server", "addr", srv.Addr)

	var c chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logger.Info("beginning graceful shutdown")
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	// Use the ListenAndServeTLS() function on the http.Server struct
	// to start the server.
	err = srv.ListenAndServe()
	if err == http.ErrServerClosed {
		logger.Info("closed gracefully")
	} else {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// Wraps sql.Open() registers Cloud SQL driver (if prod), and returns a sql.DB connection pool & cleanup func.
func openDB(dsn string, prod bool) (*sql.DB, func() error, error) {
	if prod {
		// Register custom driver with IAM Auth enabled.
		cleanup, err := pgxv5.RegisterDriver("cloudsql-postgres", cloudsqlconn.WithIAMAuthN())
		if err != nil {
			return nil, nil, err
		}

		// Open SQL connection.
		db, err := sql.Open("cloudsql-postgres", dsn)
		if err != nil {
			cleanup()
			return nil, nil, err
		}

		// Verify database connection.
		err = db.Ping()
		if err != nil {
			db.Close()
			cleanup()
			return nil, nil, err
		}

		return db, cleanup, nil
	}

	// Local dev setup
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, nil, err
	}

	return db, func() error { return nil }, nil
}
