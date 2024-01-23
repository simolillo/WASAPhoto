package api

/*
go run ./cmd/webapi/
curl -v \
	-X GET \
	-H 'Authorization: 1' \
	localhost:3000/users/{1}
*/

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1.
	// checking if the request is valid

	// extracting {uid} parameter from the path
	var pathUid int64
	pathUid, err := strconv.ParseInt(ps.ByName("uid"), 10, 64)

	if err != nil {
		// the path parameter {uid} was not a parseable int64 or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("getUserProfile: the path parameter {uid} was not a parseable int64 or is missing")
		fmt.Fprint(w, "\ngetUserProfile: the path parameter {uid} was not a parseable int64 or is missing\n\n")
		return
	}
	requestedUser, present := rt.db.SearchUByID(pathUid)
	if !present {
		// the {uid} path parameter is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("getUserProfile: the path parameter {uid} is not matching any existing user")
		fmt.Fprint(w, "\ngetUserProfile: the path parameter {uid} is not matching any existing user\n\n")
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
		ctx.Logger.WithError(err).Error("getUserProfile: the user is not authenticated")
		fmt.Fprint(w, "\ngetUserProfile: the user is not authenticated\n\n")
		return
	}
	requestingUser, present := rt.db.SearchUByID(authorizationUid)
	if !present {
		// the Authorization ID is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("getUserProfile: the Authorization ID is not matching any existing user")
		fmt.Fprint(w, "\ngetUserProfile: the Authorization ID is not matching any existing user\n\n")
		return
	}

	// a user cannot get the profile of someone who banned him
	banning := rt.db.CheckBan(requestedUser.ID, requestingUser.ID)
	if banning {
		// rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("getUserProfile: the user is trying to get the profile of someone who banned him")
		fmt.Fprint(w, "\ngetUserProfile: the user is trying to get the profile of someone who banned him\n\n")
		return
	}
	// both directions
	banning = rt.db.CheckBan(requestingUser.ID, requestedUser.ID)
	if banning {
		// rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("getUserProfile: the user is trying to get the profile of someone who banned")
		fmt.Fprint(w, "\ngetUserProfile: the user is trying to get the profile of someone who banned\n\n")
		return
	}

	// remove ban record from database
	profile, err := rt.db.GetUserProfile(requestedUser.ID, requestingUser.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("getUserProfile: unalbe to get user profile")
		fmt.Fprint(w, "\ngetUserProfile: unalbe to get user profile\n\n")
		return
	}

	w.WriteHeader(http.StatusOK) //200
	fmt.Fprintln(w)
	err = json.NewEncoder(w).Encode(profile)

	// 7.
	// if encoding operation is unsuccessful though the user has been created
	if err != nil {
		http.Error(w, "StatusInternalServerError", http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("getUserProfile: unable to encode JSON profile response")
		fmt.Fprint(w, "\ngetUserProfile: unable to encode JSON profile response\n\n")
		return
	}

	fmt.Fprint(w, "\ngetUserProfile:\nAction successful.\nUser profile info returned in the content.\n\n")
}
