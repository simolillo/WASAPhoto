package database

import "database/sql"

type User struct {
	ID uint64   `json:"userID"`
	Name string `json:"username"`
}

func (db *appdbimpl) CreateUser(username string) (user User, err error) {

	query := "INSERT INTO users (username) VALUES (?)"

	sqlResult, err := db.c.Exec(query, username)
	if err != nil {
		return
	}
	user.Name = username
	userID, err := sqlResult.LastInsertId()
	user.ID = uint64(userID)
	return user, err
}

func (db *appdbimpl) SearchUserByUsername(username string) (user User, present bool, err error) {

	query := "SELECT * FROM users WHERE username = ?;"

	err = db.c.QueryRow(query, username).Scan(&user)
	if err != nil && err != sql.ErrNoRows{
		return
	} else if err == sql.ErrNoRows {
		err = nil
		return
	} else {
		err = nil
		present = true
		return
	}
}