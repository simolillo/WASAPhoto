package database

type User struct {
	ID uint64   `json:"userID"`
	Name string `json:"username"`
}
