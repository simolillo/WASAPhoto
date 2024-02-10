package database

type User struct {
	ID   uint64 `json:"userID"`
	Name string `json:"username"`
}

type Photo struct {
	ID             uint64    `json:"photoID"`
	AuthorID       uint64    `json:"authorID"`
	AuthorUsername string    `json:"authorUsername"`
	Format         string    `json:"format"`
	Date           string    `json:"date"`
	LikesList      []User    `json:"likesList"`
	CommentsList   []Comment `json:"commentsList"`
}

type Comment struct {
	ID             uint64 `json:"commentID"`
	PhotoID        uint64 `json:"photoID"`
	AuthorID       uint64 `json:"authorID"`
	AuthorUsername string `json:"authorUsername"`
	Text           string `json:"commentText"`
	Date           string `json:"date"`
}

type Profile struct {
	Username         string `json:"username"`
	PhotosCount      uint64 `json:"photosCount"`
	FollowersCount   uint64 `json:"followersCount"`
	FollowingCount   uint64 `json:"followingCount"`
	IsItMe           bool   `json:"isItMe"`
	DoIFollowUser    bool   `json:"doIFollowUser"`
	IsInMyBannedList bool   `json:"isInMyBannedList"`
	AmIBanned        bool   `json:"amIBanned"`
}
