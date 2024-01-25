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
	rt.router.PUT("/session", rt.wrap(rt.setMyUserName))

	// Following
	rt.router.PUT("/following/:uid", rt.wrap(rt.followUser))
	rt.router.DELETE("/following/:uid", rt.wrap(rt.unfollowUser))

	// Banned
	rt.router.PUT("/banned/:uid", rt.wrap(rt.banUser))
	rt.router.DELETE("/banned/:uid", rt.wrap(rt.unbanUser))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
