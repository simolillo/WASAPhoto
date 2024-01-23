package api

/*
go run ./cmd/webapi/
curl -v \
	-X DELETE \
	-H 'Authorization: 2' \
	localhost:3000/banned/{1}
*/

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1.
	// checking if the request is valid

	// extracting {uid} parameter from the path
	var pathUid int64
	pathUid, err := strconv.ParseInt(ps.ByName("uid"), 10, 64)

	if err != nil {
		// the path parameter {uid} was not a parseable int64 or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unbanUser: the path parameter {uid} was not a parseable int64 or is missing")
		fmt.Fprint(w, "\nunbanUser: the path parameter {uid} was not a parseable int64 or is missing\n\n")
		return
	}
	banned, present := rt.db.SearchUByID(pathUid)
	if !present {
		// the {uid} path parameter is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unbanUser: the path parameter {uid} is not matching any existing user")
		fmt.Fprint(w, "\nunbanUser: the path parameter {uid} is not matching any existing user\n\n")
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
		ctx.Logger.WithError(err).Error("unbanUser: the user is not authenticated")
		fmt.Fprint(w, "\nunbanUser: the user is not authenticated\n\n")
		return
	}
	banner, present := rt.db.SearchUByID(authorizationUid)
	if !present {
		// the Authorization ID is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unbanUser: the Authorization ID is not matching any existing user")
		fmt.Fprint(w, "\nunbanUser: the Authorization ID is not matching any existing user\n\n")
		return
	}

	// a user cannot unban someone who has not banned already
	banning := rt.db.CheckBan(banner.ID, banned.ID)
	if !banning {
		// rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unbanUser: the user is trying to unfollow someone who is not following")
		fmt.Fprint(w, "\nunbanUser: the user is trying to unfollow someone who is not following\n\n")
		return
	}

	// remove ban record from database
	err = rt.db.UnbanUser(banner.ID, banned.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("unbanUser: unalbe to delete following relation")
		fmt.Fprint(w, "\nunbanUser: unalbe to delete following relation\n\n")
		return
	}

	fmt.Fprintf(w, "\nHello %s! You have unbanned %s.\n\n", banner.Name, banned.Name)
}
