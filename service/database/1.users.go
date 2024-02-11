package database

import "database/sql"

func (db *appdbimpl) CreateUser(username string) (dbUser User, err error) {

	query := "INSERT INTO users (username) VALUES (?);"

	sqlResult, err := db.c.Exec(query, username)
	if err != nil {
		return
	}
	dbUser.Name = username
	userID, err := sqlResult.LastInsertId()
	dbUser.ID = uint64(userID)
	return
}

func (db *appdbimpl) UpdateUsername(dbUser User) (err error) {

	query := "UPDATE users SET username = ? WHERE userID = ?;"
	_, err = db.c.Exec(query, dbUser.Name, dbUser.ID)
	if err != nil {
		return
	}

	query = "UPDATE comments SET authorUsername = ? WHERE authorID = ?;"
	_, err = db.c.Exec(query, dbUser.Name, dbUser.ID)

	return
}

func (db *appdbimpl) SearchUserByUsername(username string) (dbUser User, present bool, err error) {

	query := "SELECT * FROM users WHERE username = ?;"

	err = db.c.QueryRow(query, username).Scan(&dbUser.ID, &dbUser.Name)
	if err != nil && err != sql.ErrNoRows {
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

func (db *appdbimpl) SearchUserByID(ID uint64) (dbUser User, present bool, err error) {

	query := "SELECT * FROM users WHERE userID = ?;"

	err = db.c.QueryRow(query, ID).Scan(&dbUser.ID, &dbUser.Name)
	if err != nil && err != sql.ErrNoRows {
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
