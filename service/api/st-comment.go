package api

import (
	"github.com/simolillo/WASAPhoto/service/database"
	"time"
)

type Comment struct {
	ID uint64       `json:"commentID"`
	PhotoID uint64  `json:"photoID"`
	AuthorID uint64 `json:"authorID"`
	Text string     `json:"commentText"`
	Date time.Time  `json:"date"`
}

func (c *Comment) ToDatabase() database.Comment {
	return database.Comment{
		ID:       c.ID,
		PhotoID:  c.PhotoID,
		AuthorID: c.AuthorID,
		Text:     c.Text,
		Date:     c.Date,
	}
}

func (c *Comment) FromDatabase(comment database.Comment) {
	c.ID = comment.ID
	c.PhotoID = comment.PhotoID
	c.AuthorID = comment.AuthorID
	c.Text = comment.Text
	c.Date = comment.Date
}
