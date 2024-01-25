package api

/*
go run ./cmd/webapi/
curl -v \
	-X PUT \
	-H 'Authorization: 2' \
	localhost:3000/likes/{1}
*/

import (
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "likePhoto: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	liker, present, err := rt.db.SearchUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "likePhoto: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	var pathPid uint64
	pathPid, err = strconv.ParseUint(ps.ByName("pid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "likePhoto: invalid path parameter pid"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	photo, present, err := rt.db.SearchPhotoByID(pathPid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "likePhoto: path parameter pid not matching any existing photo"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	// Forbidden check
	isBanned, err := rt.db.CheckBan(liker.ID, photo.AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isBanned {
		stringErr := "likePhoto: liker banned photo author"
		http.Error(w, stringErr, http.StatusForbidden)
		return
	}
	isBanned, err = rt.db.CheckBan(photo.AuthorID, liker.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isBanned {
		stringErr := "likePhoto: photo author banned liker"
		http.Error(w, stringErr, http.StatusForbidden)
		return
	}

	// database section
	err = rt.db.LikePhoto(liker.ID, photo.ID)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
