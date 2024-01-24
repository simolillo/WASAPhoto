package api

import (
	"time"

	"github.com/simolillo/WASAPhoto/service/database"
	"github.com/simolillo/WASAPhoto/service/fileSystem"
)

type User struct {
	ID int64    `json:"userID"`
	Name string `json:"username"`
}

type Photo struct {
	ID int64                 `json:"photoID"`
	AuthorID int64           `json:"authorID"`
	Path string              `json:"path"`
	Format string            `json:"format"`
	UploadDateTime time.Time    `json:"uploadDateTime"`
}

type Comment struct {
	ID int64                 `json:"commentID"`
	Text string           `json:"commentText"`
	PhotoID int64              `json:"photoID"`
	AuthorID int64            `json:"authorID"`
	PublishDate time.Time    `json:"publishDate"`
}

type Profile struct {
	Username string                 `json:"username"`
	Posts int64           `json:"posts"`
	Followers int64              `json:"followers"`
	Following int64            `json:"following"`
	IsFollowedByViewer bool   `json:"isFollowedByViewer"`
}

// _________________ DB-Conversion Methods _________________

func (user User) ToDatabase() database.User {
	return database.User{ID: user.ID, Name: user.Name}
}

func (photo Photo) ToDatabase() database.Photo {
	return database.Photo{ID: photo.ID, AuthorID: photo.AuthorID, Path: photo.Path, Format: photo.Format, UploadDateTime: photo.UploadDateTime}
}


// _________________ FS-Conversion Methods _________________

func (user User) ToFileSystem() fs.User {
	return fs.User{ID: user.ID, Name: user.Name}
}

func (photo Photo) ToFileSystem() fs.Photo {
	return fs.Photo{ID: photo.ID, AuthorID: photo.AuthorID, Path: photo.Path, Format: photo.Format, UploadDateTime: photo.UploadDateTime}
}