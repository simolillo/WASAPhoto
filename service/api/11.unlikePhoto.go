package api

/*
go run ./cmd/webapi/
curl -v \
	-X DELETE \
	-H 'Authorization: 1' \
	localhost:3000/photos/{1}/likes/{1}
*/

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1.
	// checking if the request is valid

	// extracting {uid} parameter from the path
	var pathUid int64
	pathUid, err := strconv.ParseInt(ps.ByName("uid"), 10, 64)

	if err != nil {
		// the path parameter {uid} was not a parseable int64 or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unlikePhoto: the path parameter {uid} was not a parseable int64 or is missing")
		fmt.Fprint(w, "\nunlikePhoto: the path parameter {uid} was not a parseable int64 or is missing\n\n")
		return
	}
	liker, present := rt.db.SearchUByID(pathUid)
	if !present {
		// the {uid} path parameter is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unlikePhoto: the path parameter {uid} is not matching any existing user")
		fmt.Fprint(w, "\nunlikePhoto: the path parameter {uid} is not matching any existing user\n\n")
		return
	}

	// extracting {pid} parameter from the path
	var pathPid int64
	pathPid, err = strconv.ParseInt(ps.ByName("pid"), 10, 64)

	if err != nil {
		// the path parameter {uid} was not a parseable int64 or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unlikePhoto: the path parameter {pid} was not a parseable int64 or is missing")
		fmt.Fprint(w, "\nunlikePhoto: the path parameter {pid} was not a parseable int64 or is missing\n\n")
		return
	}
	photo, present := rt.db.SearchPByID(pathPid)
	if !present {
		// the {pid} path parameter is not matching any existing photo, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unlikePhoto: the path parameter {pid} is not matching any existing photo")
		fmt.Fprint(w, "\nunlikePhoto: the path parameter {pid} is not matching any existing photo\n\n")
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
		ctx.Logger.WithError(err).Error("unlikePhoto: the user is not authenticated")
		fmt.Fprint(w, "\nunlikePhoto: the user is not authenticated\n\n")
		return
	}
	requestingUser, present := rt.db.SearchUByID(authorizationUid)
	if !present {
		// the Authorization ID is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unlikePhoto: the Authorization ID is not matching any existing user")
		fmt.Fprint(w, "\nunlikePhoto: the Authorization ID is not matching any existing user\n\n")
		return
	}

	// requesting user can only put like for himself
	if requestingUser.ID != liker.ID {
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("unlikePhoto: the Authorization ID is not matching path uid")
		fmt.Fprint(w, "\nunlikePhoto: the Authorization ID is not matching path uid\n\n")
		return
	}

	// remove like record to database
	err = rt.db.UnlikePhoto(photo.ID, requestingUser.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("unlikePhoto: unalbe to remove like")
		fmt.Fprint(w, "\nunlikePhoto: unalbe to remove like\n\n")
		return
	}

	fmt.Fprintf(w, "\nCongrats %s! You removed like to photo number %d.\n\n", requestingUser.Name, photo.ID)
}