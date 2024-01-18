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

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
