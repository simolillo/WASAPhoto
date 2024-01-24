package api

import (
	"github.com/simolillo/WASAPhoto/service/database"
	"strings"
)

type User struct {
	ID uint64   `json:"userID"`
	Name string `json:"username"`
}

func (u *User) IsValid() bool {
	username := strings.TrimSpace(u.Name)
	if username == "" {
		return false
	}
	if len(username)<3 || len(username)>16 {
		return false
	}
	return true
}

func (u *User) ToDatabase() database.User {
	return database.User{
		ID:   u.ID,
		Name: u.Name,
	}
}

func (u *User) FromDatabase(user database.User) {
	u.ID = user.ID
	u.Name = user.Name
}