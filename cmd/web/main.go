package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/qqquinnn/snippetbox/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
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
	// Load .env file if it exists.
	_ = godotenv.Load()

	// Define command-line flags using env variables as defaults.
	defaultDSN := os.Getenv("SNIPPETBOX_DSN")
	dsn := flag.String("dsn", defaultDSN, "MySQL data source name")
	addr := flag.String("addr", ":4000", "HTTP network address")
	debug := flag.Bool("debug", false, "Enable debug mode")

	// Parse the command-line flags.
	flag.Parse()

	// Panic if DSN config is missing.
	if *dsn == "" {
		panic("database DSN must be provided via -dsn flag or SNIPPETBOX_DSN env var")
	}

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

	// Initialize and configure new session manager.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// Initialize new instance of the application struct.
	app := &application{
		debug:          *debug,
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize a struct to hold non-default TLS settings.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
	}

	// Initialize a new http.Server struct.
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// Writes log entries at the Error level.
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,
		// Idle, Read & Write timeouts.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Fetch TLS cert & key; panic if missing.
	tlsCert := os.Getenv("TLS_CERT")
	tlsKey := os.Getenv("TLS_KEY")
	if tlsCert == "" || tlsKey == "" {
		panic("Paths to TLS certificate and private key must be provided via TLS_CERT and TLS_KEY env vars")
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
	err = srv.ListenAndServeTLS(tlsCert, tlsKey)
	if err == http.ErrServerClosed {
		logger.Info("closed gracefully")
	} else {
		logger.Error(err.Error())
		os.Exit(1)
	}
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
