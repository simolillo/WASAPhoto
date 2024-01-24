package api

/*
go run ./cmd/webapi/
curl -v \
	-X DELETE \
	-H 'Authorization: 2' \
	localhost:3000/photos/{1}/comments/{1}
*/

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1.
	// checking if the request is valid
	
		// extracting {pid} parameter from the path
		var pathPid int64
		pathPid, err := strconv.ParseInt(ps.ByName("pid"), 10, 64)
	
		if err != nil {
			// the path parameter {uid} was not a parseable int64 or is missing, rejecting the request
			w.WriteHeader(http.StatusBadRequest) //400
			ctx.Logger.WithError(err).Error("uncommentPhoto: the path parameter {pid} was not a parseable int64 or is missing")
			fmt.Fprint(w, "\nuncommentPhoto: the path parameter {pid} was not a parseable int64 or is missing\n\n")
			return
		}
		photo, present := rt.db.SearchPByID(pathPid)
		if !present {
			// the {pid} path parameter is not matching any existing photo, rejecting the request
			w.WriteHeader(http.StatusBadRequest) //400
			ctx.Logger.WithError(err).Error("uncommentPhoto: the path parameter {pid} is not matching any existing photo")
			fmt.Fprint(w, "\nuncommentPhoto: the path parameter {pid} is not matching any existing photo\n\n")
			return
		}

	// extracting {cid} parameter from the path
	var pathCid int64
	pathCid, err = strconv.ParseInt(ps.ByName("cid"), 10, 64)

	if err != nil {
		// the path parameter {cid} was not a parseable int64 or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("uncommentPhoto: the path parameter {cid} was not a parseable int64 or is missing")
		fmt.Fprint(w, "\nuncommentPhoto: the path parameter {cid} was not a parseable int64 or is missing\n\n")
		return
	}
	selectedComment, present := rt.db.SearchCByID(pathCid)
	if !present || selectedComment.PhotoID != photo.ID{
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("uncommentPhoto: comment not present or not matching the photo")
		fmt.Fprint(w, "\nuncommentPhoto: comment not present or not matching the photo\n\n")
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
		ctx.Logger.WithError(err).Error("uncommentPhoto: the user is not authenticated")
		fmt.Fprint(w, "\nuncommentPhoto: the user is not authenticated\n\n")
		return
	}
	requestingUser, present := rt.db.SearchUByID(authorizationUid)
	if !present || requestingUser.ID != selectedComment.AuthorID {
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("uncommentPhoto: the Authorization ID is not matching any existing user or not author")
		fmt.Fprint(w, "\nuncommentPhoto: the Authorization ID is not matching any existing user or not author\n\n")
		return
	}

	// remove like record to database
	err = rt.db.UncommentPhoto(selectedComment.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("uncommentPhoto: unalbe to remove comment")
		fmt.Fprint(w, "\nuncommentPhoto: unalbe to remove comment\n\n")
		return
	}

	fmt.Fprintf(w, "\nCongrats %s! You removed comment number %d.\n\n", requestingUser.Name, selectedComment.ID)
}