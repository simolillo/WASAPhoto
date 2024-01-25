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

func (db *appdbimpl) GetPhotosList(authorID uint64) (photosList []Photo, err error) {

	query := "SELECT * FROM photos WHERE authorID = ?;"

	rows, err := db.c.Query(query, authorID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.ID, &photo.AuthorID, &photo.Format, &photo.Date)
		if err != nil {
			return
		}
		photosList = append(photosList, photo)
	}

	err = rows.Err()
	return
}

func (db *appdbimpl) GetFollowersList(followedID uint64) (followersList []User, err error) {

	query := "SELECT followerID FROM following WHERE followedID = ?;"

	rows, err := db.c.Query(query, followedID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID)
		if err != nil {
			return
		}
		user, _, err = db.SearchUserByID(user.ID)
		if err != nil {
			return
		}
		followersList = append(followersList, user)
	}

	err = rows.Err()
	return
}

func (db *appdbimpl) GetFollowingsList(followerID uint64) (followingsList []User, err error) {

	query := "SELECT followedID FROM following WHERE followerID = ?;"

	rows, err := db.c.Query(query, followerID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID)
		if err != nil {
			return
		}
		user, _, err = db.SearchUserByID(user.ID)
		if err != nil {
			return
		}
		followingsList = append(followingsList, user)
	}

	err = rows.Err()
	return
}
