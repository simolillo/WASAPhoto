package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
	"github.com/simolillo/WASAPhoto/service/fileSystem"
)

// Funzione che restituisce la foto richiesta
// viene servito il file con il metodo http
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	
	// extracting {pid} parameter from the path
	var pathPid int64
	pathPid, err := strconv.ParseInt(ps.ByName("pid"), 10, 64)

	// 1.
	// checking if the request is valid
	if err != nil {
		// the path parameter {pid} was not a parseable int64 or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("getPhoto: the path parameter {pid} was not a parseable int64 or is missing")
		fmt.Fprint(w, "\ngetPhoto: the path parameter {pid} was not a parseable int64 or is missing\n\n")
		return
	}
	
	photo, err := rt.db.GetFromDatabase(pathPid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("getPhoto: unable to get photo from database")
		fmt.Fprint(w, "\ngetPhoto: unable to get photo from database\n\n")
		return
	}

	photoPath := filepath.Join(fs.Root, fmt.Sprint(photo.ID) + "." + photo.Format)
	http.ServeFile(w, r, photoPath)
}
