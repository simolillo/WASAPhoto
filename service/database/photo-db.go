package database

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/simolillo/WASAPhoto/service/fileSystem"
)

// This function creates a new photo with the provided specifications in the database.
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

	photo.Path = filepath.Join(fs.Root, fmt.Sprint(photo.ID) + "." + photo.Format)
	_, err = db.c.Exec("UPDATE photos SET path = ? WHERE photoID = ?", photo.Path, photo.ID)
	return photo, err
}


func (db *appdbimpl) GetFromDatabase(photoID int64) (Photo, error) {
	var photo Photo
	err := db.c.QueryRow("SELECT * FROM photos WHERE photoID = ?;", photoID).Scan(&photo.ID, &photo.AuthorID, &photo.Path, &photo.Format, &photo.UploadDateTime)
	return photo, err
}

// This function searches for a specific photo in the database given its photo ID.
// It retruns the user if present and a boolean.
func (db *appdbimpl) SearchPByID(targetPhotoID int64) (selectedPhoto Photo, present bool) {
	err := db.c.QueryRow("SELECT * FROM photos WHERE photoID = ?;", targetPhotoID).Scan(&selectedPhoto.ID, &selectedPhoto.AuthorID, &selectedPhoto.Path, &selectedPhoto.Format, &selectedPhoto.UploadDateTime)

	// if the query selects no rows
	if err == sql.ErrNoRows {
		present = false
		return
	}

	present = true
	return
}

func (db *appdbimpl) LikePhoto(photoID int64, userID int64) error {
	_, err := db.c.Exec("INSERT INTO likes (photoID, userID) VALUES (?,?)", photoID, userID)
	return err
}