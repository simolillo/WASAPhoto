package fs

import "time"


type User struct {
	ID int64    `json:"userID"`
	Name string `json:"username"`
}

type Photo struct {
	ID int64                 `json:"photoID"`
	AuthorID int64           `json:"authorID"`
	Path string              `json:"path"`
	Format string            `json:"format"`
	UploadDateTime time.Time `json:"uploadDateTime"`
}