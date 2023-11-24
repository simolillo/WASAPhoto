package api

import "github.com/simolillo/WASAPhoto/service/database"

type User struct {
	ID int64    `json:"userID"`
	Name string `json:"username"`
}

// _________________ DB-Conversion Methods _________________

func (user User) ToDatabase() database.User {
	return database.User{ID: user.ID, Name: user.Name}
}
