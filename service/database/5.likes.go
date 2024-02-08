package database

func (db *appdbimpl) LikePhoto(likerID uint64, photoID uint64) (err error) {

	query := "INSERT INTO likes (likerID, photoID) VALUES (?,?);"

	_, err = db.c.Exec(query, likerID, photoID)
	return
}

func (db *appdbimpl) UnlikePhoto(likerID uint64, photoID uint64) (err error) {

	query := "DELETE FROM likes WHERE (likerID = ? AND photoID = ?);"

	_, err = db.c.Exec(query, likerID, photoID)
	return
}

func (db *appdbimpl) RemoveLikesBothDirections(user1ID uint64, user2ID uint64) (err error) {

	query := `
		DELETE FROM likes WHERE 
		likerID = ? 
		AND photoID IN (SELECT photoID FROM photos WHERE authorID = ?);`

	_, err = db.c.Exec(query, user1ID, user2ID)
	if err != nil {
		return err
	}
	_, err = db.c.Exec(query, user2ID, user1ID)
	return
}
