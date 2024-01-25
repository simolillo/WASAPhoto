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
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "commentPhoto: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	author, present, err := rt.db.SearchUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "commentPhoto: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	var pathPid uint64
	pathPid, err = strconv.ParseUint(ps.ByName("pid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "commentPhoto: invalid path parameter pid"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	photo, present, err := rt.db.SearchPhotoByID(pathPid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "commentPhoto: path parameter pid not matching any existing photo"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	// Forbidden check
	isBanned, err := rt.db.CheckBan(author.ID, photo.AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isBanned {
		stringErr := "commentPhoto: comment author banned photo author"
		http.Error(w, stringErr, http.StatusForbidden)
		return
	}
	isBanned, err = rt.db.CheckBan(photo.AuthorID, author.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isBanned {
		stringErr := "commentPhoto: photo author banned comment author"
		http.Error(w, stringErr, http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		stringErr := "commentPhoto: invalid request body"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	commentText := string(body)
	comment := Comment{
		PhotoID: photo.ID,
		AuthorID: author.ID,
		Text: commentText,
		Date: time.Now(),
	}

	// database section
	dbComment, err := rt.db.CommentPhoto(comment.ToDatabase())
	comment.FromDatabase(dbComment)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(comment)
}
