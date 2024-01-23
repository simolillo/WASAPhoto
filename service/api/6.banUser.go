package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1.
	// checking if the request is valid

	// extracting {uid} parameter from the path
	var pathUid int64
	pathUid, err := strconv.ParseInt(ps.ByName("uid"), 10, 64)

	if err != nil {
		// the path parameter {uid} was not a parseable int64 or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("banUser: the path parameter {uid} was not a parseable int64 or is missing")
		fmt.Fprint(w, "\nbanUser: the path parameter {uid} was not a parseable int64 or is missing\n\n")
		return
	}
	banner, present := rt.db.SearchByID(pathUid)
	if !present {
		// the {uid} path parameter is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("banUser: the path parameter {uid} is not matching any existing user")
		fmt.Fprint(w, "\nbanUser: the path parameter {uid} is not matching any existing user\n\n")
		return
	}
	
	// extracting {uid2} parameter from the path
	var pathUid2 int64
	pathUid2, err = strconv.ParseInt(ps.ByName("uid"), 10, 64)

	if err != nil {
		// the path parameter {uid2} was not a parseable int64 or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("banUser: the path parameter {uid2} was not a parseable int64 or is missing")
		fmt.Fprint(w, "\nbanUser: the path parameter {uid2} was not a parseable int64 or is missing\n\n")
		return
	}
	banned, present := rt.db.SearchByID(pathUid2)
	if !present {
		// the {uid2} path parameter is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("banUser: the path parameter {uid2} is not matching any existing user")
		fmt.Fprint(w, "\nbanUser: the path parameter {uid2} is not matching any existing user\n\n")
		return
	}

	// a user cannot ban himself
	if banner.ID == banned.ID {
		// rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("banUser: the user is trying to follow himself")
		fmt.Fprint(w, "\nbanUser: the user is trying to follow himself\n\n")
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
		ctx.Logger.WithError(err).Error("banUser: the user is not authenticated")
		fmt.Fprint(w, "\nbanUser: the user is not authenticated\n\n")
		return
	}
	selectedUser, present := rt.db.SearchByID(authorizationUid)
	if !present {
		// the Authorization ID is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("banUser: the Authorization ID is not matching any existing user")
		fmt.Fprint(w, "\nbanUser: the Authorization ID is not matching any existing user\n\n")
		return
	}

	// 3.
	// authorization phase
	if banner.ID != selectedUser.ID {
		// the ID of the user attempting the request is different from the one he wants to edit the account of, rejecting the request
		w.WriteHeader(http.StatusForbidden) //403
		ctx.Logger.WithError(err).Error("banUser: the ID of the user attempting the request is different from the one he wants to edit the account of")
		fmt.Fprint(w, "\nbanUser: the ID of the user attempting the request is different from the one he wants to to edit the account of\n\n")
		return		
	}

	// add ban record to database
	err = rt.db.BanUser(banner.ID, banned.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("banUser: unalbe to create new ban")
		fmt.Fprint(w, "\nbanUser: unalbe to create new ban\n\n")
	}

	// once banned remove follow if present
	err = rt.db.UnfollowUser(banner.ID, banned.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("banUser: unalbe to unfollow")
		fmt.Fprint(w, "\nbanUser: unalbe to unfollow\n\n")
		return
	}
	// in both directions
	err = rt.db.UnfollowUser(banned.ID, banner.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("banUser: unalbe to unfollow")
		fmt.Fprint(w, "\nbanUser: unalbe to unfollow\n\n")
		return
	}

	fmt.Fprintf(w, "\nHello %s! You banned %s.\n\n", banner.Name, banned.Name)
}