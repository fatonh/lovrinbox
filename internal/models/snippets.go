package models

import (
	"database/sql"
	"errors"
	"time"
)

// deffine a Snippet struct to hold data for an individual snippet
// Notice how the fields of the struct corepond to the feilds in our
// Mysql snippets table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// define a SnippetModel struct which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// define a Insert() method on SnippetModel which inserts a new snippet
// into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// define the SQL statement for inserting a new snippet record
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// use the Exec() method on the embedded DB field to execute the
	// SQL statement. we pass in the title, content and expires values
	// as parameters
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// use the LastInsertId() method to get the ID of the newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned by LastInsertId() is of type int64,
	// so we convert to int type before returning
	return int(id), nil
}

// This will return a specific snippet based on its ID
func (m *SnippetModel) Get(id int) (Snippet, error) {
	// define the SQL statement for getting the snippet
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// use the QueryRow() method on the embedded DB field to execute
	// the SQL statement, passing in the id variable as a parameter.
	// this returns a pointer to a sql.Row object
	row := m.DB.QueryRow(stmt, id)

	//initialize a new zeroed Snippet struct.
	var s Snippet

	// use row.Scan() to copy the values from each field in sql.Row
	// to the corresponding field in the Snippet struct
	// Notice that the arguements
	// to row.Scan are *pointers* to the fields in the Snippet struct
	// and the number of the arguments must be exactly the same as the number of
	// selected columns in the SQL statement
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		// if the query returns no rows, then row.Scan will return
		// a sql.ErrNoRows error. we use errors.Is() to check for this
		// and return our own ErrNoRecord error
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

// This will return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	// Write the SQL statment we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP()
	ORDER BY id DESC LIMIT 10`

	// Use the Query() method on connection pool to execute our
	// SQL statment. This reutrns a sql.Rows resultset containing the result
	// of our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// We defer a call to rows.Close() to ensure that the
	// resultset is closed before Latest() exits.
	// This defer shoud come *after* you check for an error from the Query()
	//method. Otherwise if Query() returns an error, we will get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initalize an emptu slice to hold the Snippet structs.
	var snippets []Snippet

	// Use rows.Next() to iterate through the rows in the resultset.
	// This prepares the first (and then each subsequent) row to be
	// acted on by the rows.Scan() method. if iteration over all
	// the rows complete then the resultset automatically closes itself
	// and frees up the underlying database connection.
	for rows.Next() {
		// Create a new zero value Snippet struct.
		var s Snippet

		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet struct that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// Append the Snippet struct to the slice.
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished, call rows.Err()
	// to retrieve any error that encountered during iteration.
	// It's important to call this -don't assume the iteration completed
	// successfully over the entrie resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the slice of snippets.
	return snippets, nil

}
