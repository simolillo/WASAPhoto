package database

import (
	"database/sql"
)

// This function creates a new user with the provided username in the database.
// It returns the newly created user and possibly an error.
func (db *appdbimpl) CreateUser(username string) (User, error) {
	var user User
	// since userID value is not specified, SQLite automatically assigns the next sequential integer
	sqlResult, err := db.c.Exec("INSERT INTO users (username) VALUES (?)", username)
	
	if err != nil {
		return user, err
	}
	
	user.Name = username
	user.ID, err = sqlResult.LastInsertId()
	return user, err
}


// This function searches for a specific user in the database given its username.
// It retruns the user if present and a boolean.
func (db *appdbimpl) SearchByUsername(targetUsername string) (selectedUser User, present bool) {
	err := db.c.QueryRow("SELECT * FROM users WHERE username = ?;", targetUsername).Scan(&selectedUser.ID, &selectedUser.Name)

	// if the query selects no rows
	if err == sql.ErrNoRows {
		present = false
		return
	}

	present = true
	return
}


// This function searches for a specific user in the database given its user ID.
// It retruns the user if present and a boolean.
func (db *appdbimpl) SearchByID(targetUserID int64) (selectedUser User, present bool) {
	err := db.c.QueryRow("SELECT * FROM users WHERE UserID = ?;", targetUserID).Scan(&selectedUser.ID, &selectedUser.Name)

	// if the query selects no rows
	if err == sql.ErrNoRows {
		present = false
		return
	}

	present = true
	return
}


// This function updates the username of the specified user.
// It returns the updated user and possibly an error.
func (db *appdbimpl) UpdateUsername(userID int64, newUsername string) (User, error) {
	user := User{userID, newUsername}
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE userID = ?", newUsername, userID)
	return user, err
}