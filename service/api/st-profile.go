package api

import "github.com/simolillo/WASAPhoto/service/database"

type Profile struct {
	Username           string `json:"username"`
	PhotosCount        uint64 `json:"PhotosCount"`
	FollowersCount     uint64 `json:"FollowersCount"`
	FollowingCount     uint64 `json:"FollowingCount"`
	IsFollowedByViewer bool   `json:"IsFollowedByViewer"`
}

func (p *Profile) ToDatabase() database.Profile {
	return database.Profile{
		Username:           p.Username,
		PhotosCount:        p.PhotosCount,
		FollowersCount:     p.FollowersCount,
		FollowingCount:     p.FollowingCount,
		IsFollowedByViewer: p.IsFollowedByViewer,
	}
}

func (p *Profile) FromDatabase(profile database.Profile) {
	p.Username = profile.Username
	p.PhotosCount = profile.PhotosCount
	p.FollowersCount = profile.FollowersCount
	p.FollowingCount = profile.FollowingCount
	p.IsFollowedByViewer = profile.IsFollowedByViewer
}
