package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Register routes

	// Session
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	
	// Settings
	rt.router.PUT("/settings", rt.wrap(rt.setMyUserName))

	// Following
	rt.router.PUT("/following/:uid", rt.wrap(rt.followUser))
	rt.router.DELETE("/following/:uid", rt.wrap(rt.unfollowUser))

	// Banned
	rt.router.PUT("/banned/:uid", rt.wrap(rt.banUser))
	rt.router.DELETE("/banned/:uid", rt.wrap(rt.unbanUser))

	// Photos
	rt.router.POST("/photos/", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/photos/:pid/", rt.wrap(rt.deletePhoto))

	// Likes
	rt.router.PUT("/likes/:pid", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/likes/:pid", rt.wrap(rt.unlikePhoto))

	// Comments
	rt.router.POST("/photos/:pid/comments/", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/photos/:pid/comments/:cid", rt.wrap(rt.uncommentPhoto))

	// Getters
	rt.router.GET("/users/:uid/", rt.wrap(rt.getUserProfile))
	rt.router.GET("/users/:uid/photos/", rt.wrap(rt.getPhotosList))
	rt.router.GET("/users/:uid/followers/", rt.wrap(rt.getFollowersList))
	rt.router.GET("/users/:uid/followings/", rt.wrap(rt.getFollowingsList))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
