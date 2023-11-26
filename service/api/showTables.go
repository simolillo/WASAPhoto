package api

/*
go run ./cmd/webapi/
curl "http://localhost:3000/tables?tableName="

ex.
curl "http://localhost:3000/tables?tableName=users"
*/

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
)

// Shows the table specified in the quey of the URL.
//
// ex.: curl "http://localhost:3000/tables?tableName=users"
func (rt *_router) showTables(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	tableName := r.URL.Query().Get("tableName")
	
	switch tableName {
	case "":
		w.WriteHeader(http.StatusBadRequest) //400
		fmt.Fprintln(w, "showTable: no table specified")
	case "users" :
		users, err := rt.db.ShowUsersTable()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) //500
			ctx.Logger.WithError(err).Error("showTables: unable to show \"users\" table")
			fmt.Fprintln(w, "showTables: unable to show \"users\" table")
			return
		}
		fmt.Fprint(w, "\n\"users\" table:\n\n")
		for _, user := range users {
			fmt.Fprintf(w, "UserID: %d, Username: %s\n", user.ID, user.Name)
		}
		fmt.Fprintln(w)
	default:
		fmt.Fprintf(w, "showTable: the \"%s\" table does not exist\n", tableName)
	}

}