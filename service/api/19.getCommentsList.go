package api

/*
go run ./cmd/webapi/
curl -v \
	-X GET \
	-H 'Authorization: 1' \
	localhost:3000/photos/{1}/comments/
*/

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
	"net/http"
	"strconv"
)

func (rt *_router) getCommentsList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "getCommentsList: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	requestingUser, present, err := rt.db.SearchUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "getCommentsList: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	var pathPid uint64
	pathPid, err = strconv.ParseUint(ps.ByName("pid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "getCommentsList: invalid path parameter pid"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	photo, present, err := rt.db.SearchPhotoByID(pathPid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "getCommentsList: path parameter pid not matching any existing photo"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	// Forbidden check
	someoneIsBanned, err := rt.db.CheckBanBothDirections(requestingUser.ID, photo.AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if someoneIsBanned {
		stringErr := "getCommentsList: someone has banned the other"
		http.Error(w, stringErr, http.StatusForbidden)
		return
	}

	// database section
	commentsList, err := rt.db.GetCommentsList(photo.ID)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(commentsList)
}
