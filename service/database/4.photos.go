package database

import "database/sql"

func (db *appdbimpl) CreatePhoto(photo Photo) (dbPhoto Photo, err error) {

	query := "INSERT INTO photos (authorID, format, date) VALUES (?,?,?);"

	sqlResult, err := db.c.Exec(query, photo.AuthorID, photo.Format, photo.Date)
	if err != nil {
		return
	}
	dbPhoto = photo
	photoID, err := sqlResult.LastInsertId()
	dbPhoto.ID = uint64(photoID)
	return
}

func (db *appdbimpl) DeletePhoto(ID uint64) (err error) {

	query := "DELETE FROM photos WHERE photoID = ?;"

	_, err = db.c.Exec(query, ID)
	return err
}

func (db *appdbimpl) SearchPhotoByID(ID uint64) (dbPhoto Photo, present bool, err error) {

	query := "SELECT * FROM photos WHERE photoID = ?;"

	row := db.c.QueryRow(query, ID)
	err = row.Scan(&dbPhoto.ID, &dbPhoto.AuthorID, &dbPhoto.Format, &dbPhoto.Date)
	if err != nil && err != sql.ErrNoRows {
		return
	} else if err == sql.ErrNoRows {
		err = nil
		return
	} else {
		err = nil
		present = true
		return
	}
}
