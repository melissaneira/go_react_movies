package db

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// OpenDB initializes a new database connection using the provided
// connection string. It returns a *sql.DB instance and any error encountered.
//
// It's important to close the database connection when it's no longer needed.
// This is typically done using defer db.Close() immediately after opening the connection.
func OpenDB(dataSourceName string) (*sql.DB, error) {
	// sql.Open only validates its arguments without creating a connection.
	// The connection is lazily established by the first query.
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Verify the database connection is alive.
	// This step is important to check if the database is accessible.
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
