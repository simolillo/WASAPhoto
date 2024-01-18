/*
The database can be found here:
/private/tmp/decaf.db
open /tmp
*/

/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	
	// Creates a new user in the database.
	CreateUser(username string) (User, error)

	// Searches for a specific user in the database given its username.
	SearchByUsername(targetUsername string) (selectedUser User, present bool)

	// Searches for a specific user in the database given its user ID.
	SearchByID(targetUserID int64) (selectedUser User, present bool)

	// Updates the username of a specific user in the database.
	UpdateUsername(userID int64, newUsername string) (User, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = createDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// This function creates the entire structure of the database (tables and relations) through SQL statements.
func createDatabase(db *sql.DB) error {
	tables := [1]string{
		`CREATE TABLE IF NOT EXISTS users (
			userID INTEGER NOT NULL PRIMARY KEY,
			username VARCHAR(16) NOT NULL UNIQUE
			);`,
	}

	// execute each SQL statement
	for t := 0; t < len(tables); t++ {
		sqlStmt := tables[t]
		_, err := db.Exec(sqlStmt)
		
		if err != nil {
			return err
		}
	}

	return nil
} 