package api

/*
go run ./cmd/webapi/
curl -v \
	-X DELETE \
	-H 'Authorization: 2' \
	localhost:3000/following/{1}
*/

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1.
	// checking if the request is valid

	// extracting {uid} parameter from the path
	var pathUid int64
	pathUid, err := strconv.ParseInt(ps.ByName("uid"), 10, 64)

	if err != nil {
		// the path parameter {uid} was not a parseable int64 or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unfollowUser: the path parameter {uid} was not a parseable int64 or is missing")
		fmt.Fprint(w, "\nunfollowUser: the path parameter {uid} was not a parseable int64 or is missing\n\n")
		return
	}
	followed, present := rt.db.SearchUByID(pathUid)
	if !present {
		// the {uid} path parameter is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unfollowUser: the path parameter {uid} is not matching any existing user")
		fmt.Fprint(w, "\nunfollowUser: the path parameter {uid} is not matching any existing user\n\n")
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
		ctx.Logger.WithError(err).Error("unfollowUser: the user is not authenticated")
		fmt.Fprint(w, "\nunfollowUser: the user is not authenticated\n\n")
		return
	}
	follower, present := rt.db.SearchUByID(authorizationUid)
	if !present {
		// the Authorization ID is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unfollowUser: the Authorization ID is not matching any existing user")
		fmt.Fprint(w, "\nunfollowUser: the Authorization ID is not matching any existing user\n\n")
		return
	}

	// a user cannot unfollow someone who is not following already
	following := rt.db.CheckFollow(follower.ID, followed.ID)
	if !following {
		// rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unfollowUser: the user is trying to unfollow someone who is not following")
		fmt.Fprint(w, "\nunfollowUser: the user is trying to unfollow someone who is not following\n\n")
		return
	}

	// remove follow record from database
	err = rt.db.UnfollowUser(follower.ID, followed.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("unfollowUser: unalbe to delete following relation")
		fmt.Fprint(w, "\nunfollowUser: unalbe to delete following relation\n\n")
		return
	}

	fmt.Fprintf(w, "\nHello %s! You do not follow %s any longer.\n\n", follower.Name, followed.Name)
}
