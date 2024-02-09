package database

import "database/sql"

func (db *appdbimpl) CommentPhoto(comment Comment) (dbComment Comment, err error) {

	query := "INSERT INTO comments (photoID, authorID, authorUsername, commentText, date) VALUES (?,?,?,?,?);"

	sqlResult, err := db.c.Exec(query, comment.PhotoID, comment.AuthorID, comment.AuthorUsername, comment.Text, comment.Date)
	if err != nil {
		return
	}
	dbComment = comment
	commentID, err := sqlResult.LastInsertId()
	dbComment.ID = uint64(commentID)
	return
}

func (db *appdbimpl) UncommentPhoto(ID uint64) (err error) {

	query := "DELETE FROM comments WHERE commentID = ?;"

	_, err = db.c.Exec(query, ID)
	return
}

func (db *appdbimpl) SearchCommentByID(ID uint64) (dbComment Comment, present bool, err error) {

	query := "SELECT * FROM comments WHERE commentID = ?;"

	row := db.c.QueryRow(query, ID)
	err = row.Scan(&dbComment.ID, &dbComment.PhotoID, &dbComment.AuthorID, &dbComment.AuthorUsername, &dbComment.Text, &dbComment.Date)
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

func (db *appdbimpl) RemoveCommentsBothDirections(user1ID uint64, user2ID uint64) (err error) {

	query := `
		DELETE FROM comments WHERE 
		photoID IN (SELECT photoID FROM photos WHERE authorID = ?) 
		AND  authorID = ?;`

	_, err = db.c.Exec(query, user1ID, user2ID)
	if err != nil {
		return err
	}
	_, err = db.c.Exec(query, user2ID, user1ID)
	return
}
