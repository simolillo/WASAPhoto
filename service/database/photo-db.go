package database

import (
	"fmt"
	"path/filepath"
)

// This function creates a new photo with the provided authorID and uploadDateTime in the database.
// It returns the newly created photo and possibly an error.
func (db *appdbimpl) CreatePhoto(photo Photo) (Photo, error) {

	// since photoID value is not specified, SQLite automatically assigns the next sequential integer
	sqlResult, err := db.c.Exec("INSERT INTO photos (authorID, format, uploadDateTime) VALUES (?,?,?)", photo.AuthorID, photo.Format, photo.UploadDateTime)

	if err != nil {
		return photo, err
	}

	
	photo.ID, err = sqlResult.LastInsertId()
	if err != nil {
		return photo, err
	}
	photo.Path = filepath.Join(photo.Path + fmt.Sprint(photo.ID) + photo.Format)
	
	_, err = db.c.Exec("UPDATE photos SET path = ? WHERE photoID = ?", photo.Path, photo.ID)

	return photo, err
}