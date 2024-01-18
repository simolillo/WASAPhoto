package database

// This function deletes the specified user from the "users" table of the database.
func (db *appdbimpl) DeleteUsersRecord(userID int64) error {
	_, err := db.c.Exec("DELETE FROM users WHERE UserID = ?", userID)
	return err
}

