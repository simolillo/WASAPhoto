package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"

	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
	"github.com/simolillo/WASAPhoto/service/fileSystem"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// setting response header
	w.Header().Set("Content-Type", "application/json")

	// extracting {uid} parameter from the path
	var pathUid int64
	pathUid, err := strconv.ParseInt(ps.ByName("uid"), 10, 64)

	// 1.
	// checking if the request is valid
	if err != nil {
		// the path parameter {uid} was not a parseable int64 or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("setMyUserName: the path parameter {uid} was not a parseable int64 or is missing")
		fmt.Fprint(w, "\nsetMyUserName: the path parameter {uid} was not a parseable int64 or is missing\n\n")
		return
	}
	selectedUser1, present := rt.db.SearchByID(pathUid)
	if !present {
		// the {uid} path parameter is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("setMyUserName: the path parameter {uid} is not matching any existing user")
		fmt.Fprint(w, "\nsetMyUserName: the path parameter {uid} is not matching any existing user\n\n")
		return
	}

	// extracting authorizationUid from the Authorization header
	var authorizationUid int64
	authorizationUid, err = strconv.ParseInt(r.Header.Get("Authorization"), 10, 64)

	// 2.
	// authentication phase
	if err != nil {
		// the Authorization header is not present or no value is specified, rejecting the request
		w.WriteHeader(http.StatusUnauthorized) //401
		ctx.Logger.WithError(err).Error("setMyUserName: the user is not authenticated")
		fmt.Fprint(w, "\nsetMyUserName: the user is not authenticated\n\n")
		return
	}
	selectedUser2, present := rt.db.SearchByID(authorizationUid)
	if !present {
		// the Authorization ID is not matching any existing user, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("setMyUserName: the Authorization ID is not matching any existing user")
		fmt.Fprint(w, "\nsetMyUserName: the Authorization ID is not matching any existing user\n\n")
		return
	}

	// 3.
	// authorization phase
	if selectedUser1.ID != selectedUser2.ID {
		// the ID of the user attempting the request is different from the one he wants to update the username of, rejecting the request
		w.WriteHeader(http.StatusForbidden) //403
		ctx.Logger.WithError(err).Error("setMyUserName: the ID of the user attempting the request is different from the one he wants to update the username of")
		fmt.Fprint(w, "\nsetMyUserName: the ID of the user attempting the request is different from the one he wants to update the username of\n\n")
		return		
	}

	// Legge il body della richiesta e verifica se ci sono errori durante la lettura.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error reading body content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Reimposta il body della richiesta in modo da poterlo leggere di nuovo in seguito
	// Dopo aver letto il body bisogna riassegnare un io.ReadCloser per poterlo rileggere
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// verifico se il contenuto del body è una immagine png o jpeg(in caso di errore:400 badrequest)
	imageFile, err := jpeg.Decode(r.Body)
	format := "jpg"
	if err != nil {
		r.Body = io.NopCloser(bytes.NewBuffer(data))
		imageFile, err = png.Decode(r.Body)
		format = "png"
		if err != nil {
			return
		}
	}



	path := fs.UserFolderName(selectedUser1.ID, selectedUser1.Name)

	// extracting uploadDateTime from the Date header
	uploadDateTime, err := time.Parse(http.TimeFormat, r.Header.Get("Date"))
	if err != nil {
		return
	}
	photo := Photo{AuthorID: selectedUser1.ID, Path: path, Format: format, UploadDateTime: uploadDateTime}


	createdPhoto, err := rt.db.CreatePhoto(photo.ToDatabase())

	// 6.
	// if photo creation or ID retrieval is unsuccessful
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("doLogin: user creation or ID retrieval was unsuccessful")
		fmt.Fprint(w, "\ndoLogin: user creation or ID retrieval was unsuccessful\n\n")
		return
	}

	err = fs.CreatePhotoFile(fs.Photo(createdPhoto), imageFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("photo-upload: error copying body content into file photo")
		return
	}

	// Invia una risposta con stato "Created" e un oggetto JSON che rappresenta la foto appena caricata.
	fmt.Fprintln(w)
	err = json.NewEncoder(w).Encode(createdPhoto)

	// 7.
	// if encoding operation is unsuccessful though the user has been created
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("doLogin: unable to encode JSON response though the user has been created")
		fmt.Fprint(w, "\ndoLogin: unable to encode JSON response though the user has been created\n\n")
		return
	}

	w.WriteHeader(http.StatusCreated)
}
