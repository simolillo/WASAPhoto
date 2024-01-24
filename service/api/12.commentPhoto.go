package api

/*
go run ./cmd/webapi/
curl -v \
	-X POST \
	-H 'Content-Type: text/plain' \
	-H 'Authorization: 2' \
	-d "Wow! This photo looks amazing" \
	localhost:3000/photos/{1}/comments/
*/

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
	"github.com/simolillo/WASAPhoto/service/database"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1.
	// checking if the request is valid

	// extracting {pid} parameter from the path
	var pathPid int64
	pathPid, err := strconv.ParseInt(ps.ByName("pid"), 10, 64)

	if err != nil {
		// the path parameter {uid} was not a parseable int64 or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("commentPhoto: the path parameter {pid} was not a parseable int64 or is missing")
		fmt.Fprint(w, "\ncommentPhoto: the path parameter {pid} was not a parseable int64 or is missing\n\n")
		return
	}
	photo, present := rt.db.SearchPByID(pathPid)
	if !present {
		// the {pid} path parameter is not matching any existing photo, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("commentPhoto: the path parameter {pid} is not matching any existing photo")
		fmt.Fprint(w, "\ncommentPhoto: the path parameter {pid} is not matching any existing photo\n\n")
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
		ctx.Logger.WithError(err).Error("commentPhoto: the user is not authenticated")
		fmt.Fprint(w, "\ncommentPhoto: the user is not authenticated\n\n")
		return
	}
	requestingUser, present := rt.db.SearchUByID(authorizationUid)
	if !present {
		// the Authorization ID is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("commentPhoto: the Authorization ID is not matching any existing user")
		fmt.Fprint(w, "\ncommentPhoto: the Authorization ID is not matching any existing user\n\n")
		return
	}


	// a user cannot comment to photo of someone who banned him
	banned := rt.db.CheckBan(requestingUser.ID, photo.AuthorID)
	if banned {
		// rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("commentPhoto: the user is trying to comment photo of someone who banned him")
		fmt.Fprint(w, "\ncommentPhoto: the user is trying to comment photo of someone who banned him\n\n")
		return
	}
	// both directions
	banned = rt.db.CheckBan(photo.AuthorID, requestingUser.ID)
	if banned {
		// rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("commentPhoto: the user is trying to comment photo of someone in his banned list")
		fmt.Fprint(w, "\ncommentPhoto: the user is trying to comment photo of someone in his banned list\n\n")
		return
	}


	// extracting text of comment from the request body
	body, err := io.ReadAll(r.Body)
	commentText := string(body)

	// 4.
	// checking if decoding operation of commentText ended successfully
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("commentPhoto: error extracting the commentText from the request body")
		fmt.Fprint(w, "\ncommentPhoto: error extracting the commentText from the request body\n\n")
		return
	}

	currentTime := time.Now()
	comment := database.Comment{Text: commentText, PhotoID: pathPid, AuthorID: authorizationUid, PublishDate: currentTime}

	// add like record to database
	_, err = rt.db.CommentPhoto(comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("commentPhoto: unalbe to comment photo")
		fmt.Fprint(w, "\ncommentPhoto: unalbe to comment photo\n\n")
		return
	}

	fmt.Fprintf(w, "\nCongrats %s! You commented photo number %d.\n\n", requestingUser.Name, photo.ID)
}