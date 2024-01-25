package database

func (db *appdbimpl) FollowUser(followerID uint64, followedID uint64) (err error) {

	query := "INSERT INTO following (followerID, followedID) VALUES (?,?)"

	_, err = db.c.Exec(query, followerID, followedID)
	return err
}

func (db *appdbimpl) RemoveFollow(followerID uint64, followedID uint64) (err error) {

	query := "DELETE FROM following WHERE (followerID = ? AND followedID = ?)"

	_, err = db.c.Exec(query, followerID, followedID)
	return err
}
