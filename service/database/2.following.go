package database

func (db *appdbimpl) FollowUser(followerID uint64, followedID uint64) (err error) {

	query := "INSERT INTO following (followerID, followedID) VALUES (?,?);"

	_, err = db.c.Exec(query, followerID, followedID)
	return err
}

func (db *appdbimpl) UnfollowUser(followerID uint64, followedID uint64) (err error) {

	query := "DELETE FROM following WHERE (followerID = ? AND followedID = ?);"

	_, err = db.c.Exec(query, followerID, followedID)
	return err
}

func (db *appdbimpl) CheckFollow(followerID uint64, followedID uint64) (isFollowing bool, err error) {

	query := "SELECT EXISTS (SELECT '_' FROM following WHERE followerID = ? AND followedID = ?);"

	err = db.c.QueryRow(query, followerID, followedID).Scan(&isFollowing)
	return
}

func (db *appdbimpl) RemoveFollowBothDirections(user1ID uint64, user2ID uint64) (err error) {
	err = db.UnfollowUser(user1ID, user2ID)
	if err != nil {
		return err
	}
	err = db.UnfollowUser(user2ID, user1ID)
	return err
}
