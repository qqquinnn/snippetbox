package models

import (
	"database/sql"
	"errors"
	"time"
)

// Holds data for an individual snippet. Corresponds to MySQL table.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Describes methods of SnippetModel struct.
type SnippetModelInterface interface {
	Insert(title, content string, expires int) (int, error)
	Get(id int) (Snippet, error)
	Latest() ([]Snippet, error)
}

// Insert a new snippet into the database.
func (model *SnippetModel) Insert(title, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES($1, $2, now() AT TIME ZONE 'utc', now() AT TIME ZONE 'utc' + (INTERVAL '1 day' * $3))
	RETURNING id`

	var id int
	err := model.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Return a specific snippet based on its id.
func (model *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > (now() AT TIME ZONE 'utc') AND id = $1`

	// Use method on connection pool to execute SQL statement.
	row := model.DB.QueryRow(stmt, id)

	var snippet Snippet

	// Copy the values from each field in row to the corresponding
	// field in the Snippet struct.
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		// Return custom error if the query returns no rows.
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return snippet, nil
}

// Return the 10 most recently created snippets.
func (model *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > (now() AT TIME ZONE 'utc') ORDER BY id DESC LIMIT 10`

	// Use method on connection pool to execute SQL statement.
	rows, err := model.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Defer rows.Close() to ensure resultset is closed before the method returns.
	defer rows.Close()

	// Initialize empty slice to hold Snippet structs.
	var snippets []Snippet

	// Iterate through the rows in the resultset. Automatically closes
	// and frees up connection upon completion.
	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// Retrieve any encountered error.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
