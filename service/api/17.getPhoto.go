package api

/*
go run ./cmd/webapi/
curl -v \
	-X GET \
	localhost:3000/photos/{1}/
*/

import (
	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
	"github.com/simolillo/WASAPhoto/service/fileSystem"
	"net/http"
	"strconv"
)

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var pathPid uint64
	pathPid, err := strconv.ParseUint(ps.ByName("pid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "getPhoto: invalid path parameter pid"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	photo, present, err := rt.db.SearchPhotoByID(pathPid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "getPhoto: path parameter pid not matching any existing photo"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	fsPhoto := fs.Photo(photo)
	photoPath := fsPhoto.Path()

	// serving photo
	w.Header().Set("Content-Type", "image/*")
	// w.WriteHeader(http.StatusOK) // superfluous
	http.ServeFile(w, r, photoPath)
}
