package api

/*
go run ./cmd/webapi/
curl -v \
	-X POST \
	-H 'Authorization: 1' \
	-H 'Content-Type: image/*' \
	--data-binary "@./photo-samples/dog/dog1.jpg" \
	localhost:3000/photos/
*/

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
	"github.com/simolillo/WASAPhoto/service/fileSystem"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "uploadPhoto: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	author, present, err := rt.db.SearchUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "uploadPhoto: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	binaryData, err := io.ReadAll(r.Body)

	// BadRequest check
	if err != nil {
		stringErr := "uploadPhoto: invalid binary data"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	format := http.DetectContentType(binaryData)
	switch format {
	case "image/png":
		format = "png"
	case "image/jpeg":
		format = "jpg"
	default:
		stringErr := "uploadPhoto: binary data not png/jpg"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	binaryImage := binaryData
	photo := Photo{
		AuthorID: author.ID,
		Format:   format,
		Date:     time.Now().Format("2006-01-02 15:04:05"),
	}

	// database section
	dbPhoto, err := rt.db.CreatePhoto(photo.ToDatabase())
	photo.FromDatabase(dbPhoto)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fileSystem section
	err = fs.CreatePhotoFile(photo.ToFileSystem(), binaryImage)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(photo)
}
