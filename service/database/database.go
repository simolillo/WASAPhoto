/*
The database can be found here:
/private/tmp/decaf.db
open /tmp
*/

/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// Users table
	CreateUser(username string) (dbUser User, err error)
	UpdateUsername(dbUser User) (err error)
	SearchUserByUsername(username string) (dbUser User, present bool, err error)
	SearchUserByID(ID uint64) (dbUser User, present bool, err error)

	// Following table
	FollowUser(followerID uint64, followedID uint64) (err error)
	UnfollowUser(followerID uint64, followedID uint64) (err error)
	CheckFollow(followerID uint64, followedID uint64) (isFollowing bool, err error)
	RemoveFollowBothDirections(user1ID uint64, user2ID uint64) (err error)

	// Banned table
	BanUser(bannerID uint64, bannedID uint64) (err error)
	UnbanUser(bannerID uint64, bannedID uint64) (err error)
	CheckBan(bannerID uint64, bannedID uint64) (isBanned bool, err error)
	CheckBanBothDirections(user1ID uint64, user2ID uint64) (someoneIsBanned bool, err error)
	CascadeBanBothDirections(user1ID uint64, user2ID uint64) (err error)

	// Photos table
	CreatePhoto(photo Photo) (dbPhoto Photo, err error)
	DeletePhoto(ID uint64) (err error)
	SearchPhotoByID(ID uint64) (dbPhoto Photo, present bool, err error)

	// Likes table
	LikePhoto(likerID uint64, photoID uint64) (err error)
	UnlikePhoto(likerID uint64, photoID uint64) (err error)
	RemoveLikesBothDirections(user1ID uint64, user2ID uint64) (err error)

	// Comments table
	CommentPhoto(comment Comment) (dbComment Comment, err error)
	UncommentPhoto(ID uint64) (err error)
	SearchCommentByID(ID uint64) (dbComment Comment, present bool, err error)
	RemoveCommentsBothDirections(user1ID uint64, user2ID uint64) (err error)

	// Getters
	GetUserProfile(ID uint64) (dbProfile Profile, err error)
	GetPhotosList(authorID uint64) (photosList []Photo, err error)
	GetFollowersList(followedID uint64) (followersList []User, err error)
	GetFollowingsList(followerID uint64) (followingsList []User, err error)
	GetLikesList(photoID uint64) (likesList []User, err error)
	GetCommentsList(photoID uint64) (commentsList []Comment, err error)
	GetMyStream(requestingUserID uint64) (stream []Photo, err error)

	// Search
	SearchUser(usernameToSearch string) (usersList []User, err error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = createDatabase(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// This function creates the entire structure of the database (tables and relations) through SQL statements.
func createDatabase(db *sql.DB) error {
	tables := [6]string{
		`CREATE TABLE IF NOT EXISTS users (
			userID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(16) NOT NULL UNIQUE
			);`,
		`CREATE TABLE IF NOT EXISTS following (
			followerID INTEGER NOT NULL REFERENCES users (userID),
			followedID INTEGER NOT NULL REFERENCES users (userID),
			PRIMARY KEY (followerID, followedID)
			);`,
		`CREATE TABLE IF NOT EXISTS banned (
			bannerID INTEGER NOT NULL REFERENCES users (userID),
			bannedID INTEGER NOT NULL REFERENCES users (userID),
			PRIMARY KEY (bannerID, bannedID)
			);`,
		`CREATE TABLE IF NOT EXISTS photos (
			photoID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			authorID INTEGER NOT NULL REFERENCES users (userID),
			format VARCHAR(3) NOT NULL,
			date TEXT NOT NULL
			);`,
		`CREATE TABLE IF NOT EXISTS likes (
			likerID INTEGER NOT NULL REFERENCES users (userID),
			photoID INTEGER NOT NULL REFERENCES photos (photoID),
			PRIMARY KEY (likerID, photoID)
			);`,
		`CREATE TABLE IF NOT EXISTS comments (
			commentID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			photoID INTEGER NOT NULL REFERENCES photos (photoID),
			authorID INTEGER NOT NULL REFERENCES users (userID),
			authorUsername VARCHAR(16) NOT NULL REFERENCES users (username),
			commentText TEXT NOT NULL,
			date TEXT NOT NULL
			);`,
	}

	// execute each SQL statement
	for t := 0; t < len(tables); t++ {
		sqlStmt := tables[t]
		_, err := db.Exec(sqlStmt)

		if err != nil {
			return err
		}
	}

	return nil
}
