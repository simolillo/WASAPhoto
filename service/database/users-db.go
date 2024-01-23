package database

import (
	"database/sql"
)

// This function creates a new user with the provided username in the database.
// It returns the newly created user and possibly an error.
func (db *appdbimpl) CreateUser(username string) (User, error) {
	var user User
	// since userID value is not specified, SQLite automatically assigns the next sequential integer
	sqlResult, err := db.c.Exec("INSERT INTO users (username) VALUES (?)", username)

	if err != nil {
		return user, err
	}

	user.Name = username
	user.ID, err = sqlResult.LastInsertId()
	return user, err
}

func (db *appdbimpl) FollowUser(followerID int64, followedID int64) error {
	_, err := db.c.Exec("INSERT INTO follow (followerID, followedID) VALUES (?,?)", followerID, followedID)
	return err
}

func (db *appdbimpl) BanUser(bannerID int64, bannedID int64) error {
	_, err := db.c.Exec("INSERT INTO ban (bannerID, bannedID) VALUES (?,?)", bannerID, bannedID)
	return err
}

// This function checks follow presence
func (db *appdbimpl) CheckFollow(userID1 int64, userID2 int64) (following bool) {
	_ = db.c.QueryRow("SELECT EXISTS (SELECT 'ok' FROM follow WHERE followerID = ? AND followedID = ?);", userID1, userID2).Scan(&following)
	return
}

// This function checks ban presence
func (db *appdbimpl) CheckBan(userID1 int64, userID2 int64) (banned bool) {
	_ = db.c.QueryRow("SELECT EXISTS (SELECT 'ok' FROM ban WHERE bannerID = ? AND bannedID = ?);", userID1, userID2).Scan(&banned)
	return
}

func (db *appdbimpl) SearchByUsername(targetUsername string) (selectedUser User, present bool) {
	err := db.c.QueryRow("SELECT * FROM users WHERE username = ?;", targetUsername).Scan(&selectedUser.ID, &selectedUser.Name)

	// if the query selects no rows
	if err == sql.ErrNoRows {
		present = false
		return
	}

	present = true
	return
}

// This function searches for a specific user in the database given its user ID.
// It retruns the user if present and a boolean.
func (db *appdbimpl) SearchUByID(targetUserID int64) (selectedUser User, present bool) {
	err := db.c.QueryRow("SELECT * FROM users WHERE UserID = ?;", targetUserID).Scan(&selectedUser.ID, &selectedUser.Name)

	// if the query selects no rows
	if err == sql.ErrNoRows {
		present = false
		return
	}

	present = true
	return
}

// This function updates the username of the specified user.
// It returns the updated user and possibly an error.
func (db *appdbimpl) UpdateUsername(userID int64, newUsername string) (User, error) {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE userID = ?", newUsername, userID)
	return User{userID, newUsername}, err
}

func (db *appdbimpl) UnfollowUser(followerID int64, followedID int64) error {
	_, err := db.c.Exec("DELETE FROM follow WHERE (followerID = ? AND followedID = ?)", followerID, followedID)
	return err
}

func (db *appdbimpl) UnbanUser(bannerID int64, bannedID int64) error {
	_, err := db.c.Exec("DELETE FROM ban WHERE (bannerID = ? AND bannedID = ?)", bannerID, bannedID)
	return err
}

func (db *appdbimpl) GetUserProfile(userID int64, viewerID int64) (Profile, error) {
	var profile Profile

	// Get the number of posts
	row := db.c.QueryRow("SELECT COUNT(*) FROM photos WHERE authorID = ?", userID)
	if err := row.Scan(&profile.Posts); err != nil {
		return profile, err
	}

	// Get the number of followers
	row = db.c.QueryRow("SELECT COUNT(*) FROM follow WHERE followedID = ?", userID)
	if err := row.Scan(&profile.Followers); err != nil {
		return profile, err
	}

	// Get the number of following
	row = db.c.QueryRow("SELECT COUNT(*) FROM follow WHERE followerID = ?", userID)
	if err := row.Scan(&profile.Following); err != nil {
		return profile, err
	}

	// Get the username
	row = db.c.QueryRow("SELECT username FROM users WHERE userID = ?", userID)
	if err := row.Scan(&profile.Username); err != nil {
		return profile, err
	}

	// Check if the viewer follows the user
	row = db.c.QueryRow("SELECT EXISTS (SELECT 1 FROM follow WHERE followerID = ? AND followedID = ?)", viewerID, userID)
	if err := row.Scan(&profile.IsFollowedByViewer); err != nil {
		return profile, err
	}

	return profile, nil
}
