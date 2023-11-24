package api

// curl -X POST -H 'Content-Type: application/json' -d 'Lillo' https://localhost:3000/session


import (
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
)

/*
Summary: logs in the user

Description:
The login endpoint accepts a username like “Maria” without any password.
If the username does not exist, it will be created, and an identifier is returned.
If the username exists, the user identifier is returned.
*/
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// setting response header
	w.Header().Set("Content-Type", "application/json")

	// extracting username from the request
	var user User
	err := json.NewDecoder(r.Body).Decode(&user.Name)

	// checking if decoding operation ended successfully
	if err != nil {
		// the request body was not a parseable JSON, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("qui")
		return
	}

	// checking if the username is valid
	if !isValid(user.Name) {
		// the username is not valid, rejecting request
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	
	// moving on to the database section
	// first of all we search the user to see if it has alredy been created
	selectedUser, present := rt.db.SearchByUsername(user.ToDatabase())

	// if the user altready exists, return the ID
	if present {
		w.WriteHeader(http.StatusOK) //200
		err = json.NewEncoder(w).Encode(selectedUser)

		// if encoding operation is unsuccessful
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) //500
			ctx.Logger.WithError(err).Error("doLogin: unable to encode JSON response.")
			return
		}
		return
	}

	// if the user doesn't exist yet, create it and return the ID
	createdUser, err := rt.db.CreateUser(user.ToDatabase())

	// if user creation or ID retrieval is unsuccessful
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		return
	}

	w.WriteHeader(http.StatusCreated) //201
	err = json.NewEncoder(w).Encode(createdUser)

	// if encoding operation is unsuccessful
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("doLogin: unable to encode JSON response.")
	}

}
