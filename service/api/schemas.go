package api

import (
	

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
	UploadDateTime string    `json:"uploadDateTime"`
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