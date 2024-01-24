package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes

	// Login endpoint
	// All operations related to users logging in.
	rt.router.POST("/session", rt.wrap(rt.doLogin))


	// Settings
	rt.router.PUT("/settings", rt.wrap(rt.setMyUserName))

	// Following endpoint
	// All operations realated to users following other profiles.
	rt.router.PUT("/following/:uid", rt.wrap(rt.followUser))
	rt.router.DELETE("/following/:uid", rt.wrap(rt.unfollowUser))

	// Banned endpoint
	// All operations realated to users banned list.
	rt.router.PUT("/banned/:uid", rt.wrap(rt.banUser))
	rt.router.DELETE("/banned/:uid", rt.wrap(rt.unbanUser))
	
	// Photo endpoint
	// All operations realated to users managing their own profile.
	rt.router.POST("/photos/", rt.wrap(rt.uploadPhoto))

	// Browsing endpoint
	rt.router.GET("/users/:uid", rt.wrap(rt.getUserProfile))

	// Likes endpoint
	rt.router.PUT("/photos/:pid/likes/:uid", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/photos/:pid/likes/:uid", rt.wrap(rt.unlikePhoto))
	
	// Comments endpoint
	rt.router.POST("/photos/:pid/comments/", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/photos/:pid/comments/:cid", rt.wrap(rt.uncommentPhoto))
	
	// Getters
	rt.router.GET("/photos/:pid", rt.wrap(rt.getPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
