package api

import "github.com/simolillo/WASAPhoto/service/database"

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

func (p *Profile) ToDatabase() database.Profile {
	return database.Profile{
		Username:         p.Username,
		PhotosCount:      p.PhotosCount,
		FollowersCount:   p.FollowersCount,
		FollowingCount:   p.FollowingCount,
		IsItMe:           p.IsItMe,
		DoIFollowUser:    p.DoIFollowUser,
		IsInMyBannedList: p.IsInMyBannedList,
		AmIBanned:        p.AmIBanned,
	}
}

func (p *Profile) FromDatabase(profile database.Profile) {
	p.Username = profile.Username
	p.PhotosCount = profile.PhotosCount
	p.FollowersCount = profile.FollowersCount
	p.FollowingCount = profile.FollowingCount
	// p.IsItMe = profile.IsItMe
	// p.DoIFollowUser = profile.DoIFollowUser
	// p.IsInMyBannedList = profile.IsInMyBannedList
	// p.AmIBanned = profile.AmIBanned
}
