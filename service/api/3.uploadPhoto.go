package api

/*
go run ./cmd/webapi/
curl -v \
	-X POST \
	-H 'Content-Type: image/*' \
	-H 'Authorization: 2' \
	--data-binary "@/Users/simonerussolillo/Pictures/Random/rainforest.png" \
	localhost:3000/photos/
*/

/*
Possible outcomes:

1. checking if the request is valid                                                                     400 BadRequest
   a. curl -v -X POST -H 'Content-Type: image/*' -H 'Authorization: 1' \
      --data-binary "@/Users/simonerussolillo/Pictures/Random/rainforest.png" \
      localhost:3000/users/{aaa}/photos/
      (the path parameter {uid} is not a parseable int64)

   b. curl -v -X POST -H 'Content-Type: image/*' -H 'Authorization: 1' \
      --data-binary "@/Users/simonerussolillo/Pictures/Random/rainforest.png" \
      localhost:3000/users/{1000}/photos/
      (the {uid} path parameter is not matching any existing user)

2. authentication phase
   a. curl -v -X POST -H 'Content-Type: image/*' \                                                      401 Unauthorized
      --data-binary "@/Users/simonerussolillo/Pictures/Random/rainforest.png" \
	  localhost:3000/users/{1}/photos/
      (the Authorization header is not present or no value is specified)

   b. curl -v -X POST -H 'Content-Type: image/*' -H 'Authorization: 1000' \                             400 BadRequest
      --data-binary "@/Users/simonerussolillo/Pictures/Random/rainforest.png" \
	  localhost:3000/users/{1}/photos/
      (the Authorization ID is not matching any existing user)

3. authorization phase                                                                                  403 Forbidden
   curl -v -X POST -H 'Content-Type: image/*' -H 'Authorization: 1' \
   --data-binary "@/Users/simonerussolillo/Pictures/Random/rainforest.png" \
   localhost:3000/users/{2}/photos/
   (the {uid} path parameter is different from the Authorization ID)

4. extracting the binary data from the request body                                                     500 InternalServerError
   (server error)

5. checking the format of the binary data                                                               400 BadRequest
   curl -v -X POST -H 'Content-Type: image/*' -H 'Authorization: 1' \
   --data-binary "@/Users/simonerussolillo/Pictures/Random/pretty please kitty softpaws.gif" \
   localhost:3000/users/{1}/photos/
   (the path leads to binary data which is not a png not jpg)

6. creating the photo                                                                                   201 Created
   curl -v -X POST -H 'Content-Type: image/*' -H 'Authorization: 1' \
   --data-binary "@/Users/simonerussolillo/Pictures/Random/rainforest.png" \
   localhost:3000/users/{1}/photos/

	7. if photo creation or ID retrieval is unsuccessful                                                500 InternalServerError
	   (server error)

	8. if encoding operation is unsuccessful though the photo has been created                          500 InternalServerError
	   (server error)

9. if error occurred during the creation of the photo file                                              500 InternalServerError
   (server error)
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
	"github.com/simolillo/WASAPhoto/service/database"
	fs "github.com/simolillo/WASAPhoto/service/fileSystem"
)

/*
Summary: upload a new photo on personal account

Description:
Upload a new photo on personal account.
The server will create a new unique ID, the client can find it in the response.
*/
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// setting response header
	w.Header().Set("Content-Type", "application/json")

	// 1.
	// checking if the request is valid

	// 2.
	// authentication phase

	// extracting authorizationUid from the Authorization header
	var authorizationUid int64
	authorizationUid, err := strconv.ParseInt(r.Header.Get("Authorization"), 10, 64)

	if err != nil {
		// the Authorization header is not present or no value is specified, rejecting the request
		w.WriteHeader(http.StatusUnauthorized) //401
		ctx.Logger.WithError(err).Error("uploadPhoto: the user is not authenticated")
		fmt.Fprint(w, "\nuploadPhoto: the user is not authenticated\n\n")
		return
	}
	selectedUser, present := rt.db.SearchUByID(authorizationUid)
	if !present {
		// the Authorization ID is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("uploadPhoto: the Authorization ID is not matching any existing user")
		fmt.Fprint(w, "\nuploadPhoto: the Authorization ID is not matching any existing user\n\n")
		return
	}

	// the user is both authenticated and authorized
	requestingUser := selectedUser

	// 4.
	// extracting the binary data from the request body
	binaryData, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("uploadPhoto: error extracting binary data from the request body")
		fmt.Fprint(w, "\nuploadPhoto:\nerror extracting binary data from the request body\n\n")
		return
	}

	// 5.
	// checking the format of the binary data
	format := http.DetectContentType(binaryData)
	switch format {
	case "image/png":
		format = "png"
	case "image/jpeg":
		format = "jpg"
	default:
		// the request body contains binary data which is not image/png nor image/jpeg, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("uploadPhoto: the request body contains binary data which is not image/png nor image/jpeg")
		fmt.Fprint(w, "\nuploadPhoto: the request body contains binary data which is not image/png nor image/jpeg\n\n")
		return
	}

	// the binary data is an image of the supperted type (png/jpg)
	binaryImage := binaryData

	// building the photo instance
	currentTime := time.Now()
	photo := Photo{AuthorID: requestingUser.ID, Format: format, UploadDateTime: currentTime}

	// 6.
	// creating the photo

	// in the database
	createdPhoto, err := rt.db.CreatePhoto(database.Photo(photo))

	// 7.
	// if photo creation or ID retrieval is unsuccessful
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("uploadPhoto: photo creation or ID retrieval was unsuccessful")
		fmt.Fprint(w, "\nuploadPhoto: photo creation or ID retrieval was unsuccessful\n\n")
		return
	}

	w.WriteHeader(http.StatusCreated) //201
	fmt.Fprintln(w)
	err = json.NewEncoder(w).Encode(createdPhoto)

	// 8.
	// if encoding operation is unsuccessful though the photo has been created
	if err != nil {
		http.Error(w, "StatusInternalServerError", http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("uploadPhoto: unable to encode JSON response though the photo has been created")
		fmt.Fprint(w, "\nuploadPhoto: unable to encode JSON response though the photo has been created\n\n")
		return
	}

	fmt.Fprint(w, "\nuploadPhoto:\nThe photo has been successfully created, it is returned in the content.\n\n")

	// in the file system
	err = fs.CreatePhotoFile(fs.Photo(createdPhoto), binaryImage)

	// 9.
	// if error occurred during the creation of the photo file
	if err != nil {
		http.Error(w, "StatusInternalServerError", http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("uploadPhoto: an error occurred during the creation of the photo file though the photo has been created in the database")
		fmt.Fprint(w, "\nuploadPhoto: an error occurred during the creation of the photo file though the photo has been created in the database\n\n")
		return
	}

}
