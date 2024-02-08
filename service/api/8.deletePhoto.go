package api

/*
go run ./cmd/webapi/
curl -v \
	-X DELETE \
	-H 'Authorization: 1' \
	localhost:3000/photos/{1}/
*/

import (
	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
	"github.com/simolillo/WASAPhoto/service/fileSystem"
	"net/http"
	"strconv"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "deletePhoto: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	requestingUser, present, err := rt.db.SearchUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "deletePhoto: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	var pathPid uint64
	pathPid, err = strconv.ParseUint(ps.ByName("pid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "deletePhoto: invalid path parameter pid"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	photo, present, err := rt.db.SearchPhotoByID(pathPid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "deletePhoto: path parameter pid not matching any existing photo"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	// Forbidden check
	if requestingUser.ID != photo.AuthorID {
		stringErr := "deletePhoto: requesting user not author of the photo"
		http.Error(w, stringErr, http.StatusForbidden)
		return
	}

	// database section
	err = rt.db.DeletePhoto(photo.ID)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fileSystem section
	err = fs.DeletePhotoFile(fs.Photo(photo))

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
