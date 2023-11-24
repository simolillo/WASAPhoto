package database

import "errors"

// This function searches for a specific user in the database given its username.
// It retruns the user if present and a boolean.
func (db *appdbimpl) SearchByUsername(targetUser User) (selectedUser User, present bool) {
	err := db.c.QueryRow("SELECT * FROM users WHERE username = ?;", targetUser.Name).Scan(&selectedUser)

	// if the query selects no rows
	var ErrNoRows = errors.New("sql: no rows in result set")
	if err == ErrNoRows {
		present = false
		return
	}

	present = true
	return
}


// This function inserts a new user (record) into the "users" table of the database.
// It returns the newly created user and possibly an error.
func (db *appdbimpl) CreateUser(user User) (User, error) {
	// since userID value is not specified, SQLite automatically assigns the next sequential integer
	sqlResult, err := db.c.Exec("INSERT INTO users (username) VALUES (?)", user.Name)

	if err != nil {
		return user, err
	}

	user.ID, err = sqlResult.LastInsertId()
	return user, err
}
