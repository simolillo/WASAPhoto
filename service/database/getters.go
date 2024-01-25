package database

func (db *appdbimpl) GetUserProfile(ID uint64) (dbProfile Profile, err error) {

	queries := [4]string{
		"SELECT username FROM users WHERE userID = ?",
		"SELECT COUNT(*) FROM photos WHERE authorID = ?",
		"SELECT COUNT(*) FROM following WHERE followedID = ?",
		"SELECT COUNT(*) FROM following WHERE followerID = ?",
	}

	err = db.c.QueryRow(queries[0], ID).Scan(&dbProfile.Username)
	if err != nil {
		return
	}
	err = db.c.QueryRow(queries[1], ID).Scan(&dbProfile.PhotosCount)
	if err != nil {
		return
	}
	err = db.c.QueryRow(queries[2], ID).Scan(&dbProfile.FollowersCount)
	if err != nil {
		return
	}
	err = db.c.QueryRow(queries[3], ID).Scan(&dbProfile.FollowingCount)
	return
}
