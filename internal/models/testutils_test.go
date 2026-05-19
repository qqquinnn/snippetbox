package models

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Creates new connection pool for test database, executes setup.sql
// script, and registers 'cleanup' function that executes teardown.sql script.
func newTestDB(t *testing.T) *sql.DB {
	// Fetch test DSN.
	dsn := os.Getenv("SNIPPETBOX_TEST_DSN")
	if dsn == "" {
		t.Fatal("DSN missing from environment.")
	}

	// Establish connection pool for test database.
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}

	// Read & execute setup SQL script, and handle any errors.
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	// Register function to read & execute teardown script. Called automatically.
	t.Cleanup(func() {
		defer db.Close()

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
	})

	// Return database connection pool.
	return db
}
