package api

import (
	"github.com/simolillo/WASAPhoto/service/database"
	"regexp"
	"strings"
)

type User struct {
	ID   uint64 `json:"userID"`
	Name string `json:"username"`
}

func (u *User) HasValidUsername() bool {
	username := strings.TrimSpace(u.Name)
	if username == "" {
		return false
	}
	if len(username) < 3 || len(username) > 16 {
		return false
	}
	match, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_]{2,15}$", username)
	return match
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
