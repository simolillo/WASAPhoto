package api

import (
	"github.com/simolillo/WASAPhoto/service/database"
)

type Comment struct {
	ID             uint64 `json:"commentID"`
	PhotoID        uint64 `json:"photoID"`
	AuthorID       uint64 `json:"authorID"`
	AuthorUsername string `json:"authorUsername"`
	Text           string `json:"commentText"`
	Date           string `json:"date"`
}

func (c *Comment) ToDatabase() database.Comment {
	return database.Comment{
		ID:             c.ID,
		PhotoID:        c.PhotoID,
		AuthorID:       c.AuthorID,
		AuthorUsername: c.AuthorUsername,
		Text:           c.Text,
		Date:           c.Date,
	}
}

func (c *Comment) FromDatabase(comment database.Comment) {
	c.ID = comment.ID
	c.PhotoID = comment.PhotoID
	c.AuthorID = comment.AuthorID
	c.AuthorUsername = comment.AuthorUsername
	c.Text = comment.Text
	c.Date = comment.Date
}
