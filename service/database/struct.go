package database

type User struct {
	ID   uint64 `json:"userID"`
	Name string `json:"username"`
}

type Photo struct {
	ID       uint64 `json:"photoID"`
	AuthorID uint64 `json:"authorID"`
	Format   string `json:"format"`
	Date     string `json:"date"`
}

type Comment struct {
	ID       uint64 `json:"commentID"`
	PhotoID  uint64 `json:"photoID"`
	AuthorID uint64 `json:"authorID"`
	Text     string `json:"commentText"`
	Date     string `json:"date"`
}

type Profile struct {
	Username           string `json:"username"`
	PhotosCount        uint64 `json:"PhotosCount"`
	FollowersCount     uint64 `json:"FollowersCount"`
	FollowingCount     uint64 `json:"FollowingCount"`
	IsFollowedByViewer bool   `json:"IsFollowedByViewer"`
}
