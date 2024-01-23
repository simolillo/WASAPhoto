package api

/*
go run ./cmd/webapi/
curl -v \
	-X PUT \
	-H 'Authorization: 3' \
	localhost:3000/following/{1}
*/

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1.
	// checking if the request is valid

	// extracting {uid} parameter from the path
	var pathUid int64
	pathUid, err := strconv.ParseInt(ps.ByName("uid"), 10, 64)

	if err != nil {
		// the path parameter {uid} was not a parseable int64 or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("followUser: the path parameter {uid} was not a parseable int64 or is missing")
		fmt.Fprint(w, "\nfollowUser: the path parameter {uid} was not a parseable int64 or is missing\n\n")
		return
	}
	followed, present := rt.db.SearchUByID(pathUid)
	if !present {
		// the {uid} path parameter is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("followUser: the path parameter {uid} is not matching any existing user")
		fmt.Fprint(w, "\nfollowUser: the path parameter {uid} is not matching any existing user\n\n")
		return
	}

	// 2.
	// authentication phase

	// extracting authorizationUid from the Authorization header
	var authorizationUid int64
	authorizationUid, err = strconv.ParseInt(r.Header.Get("Authorization"), 10, 64)

	if err != nil {
		// the Authorization header is not present or no value is specified, rejecting the request
		w.WriteHeader(http.StatusUnauthorized) //401
		ctx.Logger.WithError(err).Error("followUser: the user is not authenticated")
		fmt.Fprint(w, "\nfollowUser: the user is not authenticated\n\n")
		return
	}
	follower, present := rt.db.SearchUByID(authorizationUid)
	if !present {
		// the Authorization ID is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("followUser: the Authorization ID is not matching any existing user")
		fmt.Fprint(w, "\nfollowUser: the Authorization ID is not matching any existing user\n\n")
		return
	}

	// a user cannot follow himself
	if follower.ID == followed.ID {
		// rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("followUser: the user is trying to follow himself")
		fmt.Fprint(w, "\nfollowUser: the user is trying to follow himself\n\n")
		return
	}

	// a user cannot follow someone who banned him
	banned := rt.db.CheckBan(followed.ID, follower.ID)
	if banned {
		// rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("followUser: the user is trying to follow someone who banned him")
		fmt.Fprint(w, "\nfollowUser: the user is trying to follow someone who banned him\n\n")
		return
	}
	// both directions
	banned = rt.db.CheckBan(follower.ID, followed.ID)
	if banned {
		// rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("followUser: the user is trying to follow someone in his banned list")
		fmt.Fprint(w, "\nfollowUser: the user is trying to follow someone in his banned list\n\n")
		return
	}

	// add follow record to database
	err = rt.db.FollowUser(follower.ID, followed.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("followUser: unalbe to create new following relation")
		fmt.Fprint(w, "\nfollowUser: unalbe to create new following relation\n\n")
		return
	}

	fmt.Fprintf(w, "\nCongrats %s! You started following %s.\n\n", follower.Name, followed.Name)
}
