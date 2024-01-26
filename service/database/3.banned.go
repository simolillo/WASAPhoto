package database

func (db *appdbimpl) BanUser(bannerID uint64, bannedID uint64) (err error) {

	query := "INSERT INTO banned (bannerID, bannedID) VALUES (?,?);"

	_, err = db.c.Exec(query, bannerID, bannedID)
	return err
}

func (db *appdbimpl) UnbanUser(bannerID uint64, bannedID uint64) (err error) {

	query := "DELETE FROM banned WHERE (bannerID = ? AND bannedID = ?);"

	_, err = db.c.Exec(query, bannerID, bannedID)
	return err
}

func (db *appdbimpl) CheckBan(bannerID uint64, bannedID uint64) (isBanned bool, err error) {

	query := "SELECT EXISTS (SELECT '_' FROM banned WHERE bannerID = ? AND bannedID = ?);"

	err = db.c.QueryRow(query, bannerID, bannedID).Scan(&isBanned)
	return
}

func (db *appdbimpl) CheckBanBothDirections(user1ID uint64, user2ID uint64) (someoneIsBanned bool, err error) {
	someoneIsBanned, err = db.CheckBan(user1ID, user2ID)
	if err != nil || someoneIsBanned {
		return
	}
	someoneIsBanned, err = db.CheckBan(user2ID, user1ID)
	return
}

func (db *appdbimpl) CascadeBanBothDirections(user1ID uint64, user2ID uint64) (err error) {
	err = db.RemoveFollowBothDirections(user1ID, user2ID)
	if err != nil {
		return
	}
	err = db.RemoveLikesBothDirections(user1ID, user2ID)
	if err != nil {
		return
	}
	err = db.RemoveCommentsBothDirections(user1ID, user2ID)
	return
}
