package api

/*
go run ./cmd/webapi/
curl -X DELETE "http://localhost:3000/tables?tableName=&recordID="

ex.
curl -X DELETE "http://localhost:3000/tables?tableName=users&recordID=4"
*/

import (
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
)

// Deletes a specific record of a specific table, both (tableName and recordID) are specified in the quey of the URL.
//
// ex.: curl -X DELETE "http://localhost:3000/tables?tableName=users&recordID=3"
func (rt *_router) deleteRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	tableName := r.URL.Query().Get("tableName")
	recordID, err := strconv.ParseInt(r.URL.Query().Get("recordID"), 10, 64)
	
	// checking if int conversion ended successfully
	if err != nil {
		// invalid parameter value for recordID
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("deleteRecord: invalid parameter value for recordID")
		fmt.Fprint(w, "\ndeleteRecord: invalid parameter value for recordID\n\n")
		return
	}
	
	switch tableName {
	case "":
		w.WriteHeader(http.StatusBadRequest) //400
		fmt.Fprint(w, "\ndeleteRecord: no table specified\n\n")
	case "users" :
		userID := recordID
		var toDelete = User{ID: userID}

		selectedUser, present := rt.db.SearchByID(toDelete.ToDatabase())
		if !present {
			w.WriteHeader(http.StatusBadRequest) //400
			fmt.Fprint(w, "\ndeleteRecord: the user you are trying to delete doesn't exist\n\n")
			return
		}
		toDelete.Name = selectedUser.Name

		err := rt.db.DeleteUsersRecord(toDelete.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) //500
			ctx.Logger.WithError(err).Error("deleteRecord: unable to delete specified user though the user exists")
			fmt.Fprint(w, "\ndeleteRecord: unable to delete specified user though the user exists\n\n")
			return
		}

		fmt.Fprintln(w)
		err = json.NewEncoder(w).Encode(toDelete)

		// if encoding operation is unsuccessful though the user has been deleted
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) //500
			ctx.Logger.WithError(err).Error("deleteRecord: unable to encode JSON response though the user has been deleted")
			fmt.Fprint(w, "deleteRecord: unable to encode JSON response though the user has been deleted\n\n")
			return
		}
	
		w.WriteHeader(http.StatusOK) //200
		fmt.Fprint(w, "\nThe above user has been successfully deleted.\n\n")
		
	default:
		fmt.Fprintf(w, "\ndeleteRecord: the \"%s\" table does not exist\n\n", tableName)
	}

}