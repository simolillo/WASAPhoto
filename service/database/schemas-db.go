package database

type User struct {
	ID int64    `json:"userID"`
	Name string `json:"username"`
}
