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

	// User endpoint
	// All operations realated to users managing their own profile.
	rt.router.PUT("/users/:uid/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/users/:uid/following/:uid2", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:uid/following/:uid2", rt.wrap(rt.unfollowUser))
	rt.router.PUT("/users/:uid/banned/:uid2", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:uid/banned/:uid2", rt.wrap(rt.unbanUser))
	
	// Photo endpoint
	// All operations realated to users managing their own profile.
	rt.router.POST("/users/:uid/photos/", rt.wrap(rt.uploadPhoto))

	// Getters
	rt.router.GET("/photos/:pid", rt.wrap(rt.getPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
