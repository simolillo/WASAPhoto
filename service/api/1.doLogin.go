package api

/*
go run ./cmd/webapi/
curl -v \
	-X POST \
	-H 'Content-Type: application/json' \
	-d '{"username": "Simo"}' \
	localhost:3000/session
*/

import (
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"net/http"
	"fmt"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	// BadRequest check
	if err != nil {
		stringErr := "doLogIn: invalid JSON object"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	if !user.HasValidUsername() {
		stringErr := "doLogIn: invalid username"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	// database section
	dbUser, present, err := rt.db.SearchUserByUsername(user.Name)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if present {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		user.FromDatabase(dbUser)
		_ = json.NewEncoder(w).Encode(user)
		fmt.Fprint(w, "\ndoLogIn: log-in successful\n\n")
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		dbUser, err = rt.db.CreateUser(user.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.FromDatabase(dbUser)
		_ = json.NewEncoder(w).Encode(user)
		fmt.Fprint(w, "\ndoLogIn: sign-up successful\n\n")
		return
	}
}
