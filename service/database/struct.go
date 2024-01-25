package database

import "time"

type User struct {
	ID uint64   `json:"userID"`
	Name string `json:"username"`
}

type Photo struct {
	ID uint64       `json:"photoID"`
	AuthorID uint64 `json:"authorID"`
	Format string   `json:"format"`
	Date time.Time  `json:"date"`
}

type Comment struct {
	ID uint64       `json:"commentID"`
	PhotoID uint64  `json:"photoID"`
	AuthorID uint64 `json:"authorID"`
	Text string     `json:"commentText"`
	Date time.Time  `json:"date"`
}
