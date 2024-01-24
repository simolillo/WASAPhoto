package api

/*
go run ./cmd/webapi/
curl -v \
	-X DELETE \
	-H 'Authorization: 1' \
	localhost:3000/photos/{1}
*/

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/simolillo/WASAPhoto/service/fileSystem"
	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1.
	// checking if the request is valid
	
		// extracting {pid} parameter from the path
		var pathPid int64
		pathPid, err := strconv.ParseInt(ps.ByName("pid"), 10, 64)
	
		if err != nil {
			// the path parameter {uid} was not a parseable int64 or is missing, rejecting the request
			w.WriteHeader(http.StatusBadRequest) //400
			ctx.Logger.WithError(err).Error("deletePhoto: the path parameter {pid} was not a parseable int64 or is missing")
			fmt.Fprint(w, "\ndeletePhoto: the path parameter {pid} was not a parseable int64 or is missing\n\n")
			return
		}
		photo, present := rt.db.SearchPByID(pathPid)
		if !present {
			// the {pid} path parameter is not matching any existing photo, rejecting the request
			w.WriteHeader(http.StatusBadRequest) //400
			ctx.Logger.WithError(err).Error("deletePhoto: the path parameter {pid} is not matching any existing photo")
			fmt.Fprint(w, "\ndeletePhoto: the path parameter {pid} is not matching any existing photo\n\n")
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
		ctx.Logger.WithError(err).Error("deletePhoto: the user is not authenticated")
		fmt.Fprint(w, "\ndeletePhoto: the user is not authenticated\n\n")
		return
	}
	requestingUser, present := rt.db.SearchUByID(authorizationUid)
	if !present || requestingUser.ID != photo.AuthorID {
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("deletePhoto: the Authorization ID is not matching any existing user or not author")
		fmt.Fprint(w, "\ndeletePhoto: the Authorization ID is not matching any existing user or not author\n\n")
		return
	}

	// remove photo record from database
	err = rt.db.DeletePhoto(photo.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("deletePhoto: unalbe to remove photo")
		fmt.Fprint(w, "\ndeletePhoto: unalbe to remove photo\n\n")
		return
	}

	// remove photo from file system
	err = fs.DeletePhotoFile(fs.Photo(photo))

	// 9.
	// if error occurred during the creation of the photo file
	if err != nil {
		http.Error(w, "StatusInternalServerError", http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("deletePhoto: an error occurred during the deletion of the photo file though the photo has been removed from the database")
		fmt.Fprint(w, "\ndeletePhoto: an error occurred during the deletion of the photo file though the photo has been removed from the database\n\n")
		return
	}

	fmt.Fprintf(w, "\nCongrats %s! You deleted photo number %d.\n\n", requestingUser.Name, photo.ID)
}