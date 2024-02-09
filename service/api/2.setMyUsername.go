package api

/*
go run ./cmd/webapi/
curl -v \
	-X PUT \
	-H 'Authorization: 1' \
	-H 'Content-Type: application/json' \
	-d '{"username": "lillo"}' \
	localhost:3000/settings
*/

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
	"net/http"
	"strconv"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "setMyUserName: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	dbUser, present, err := rt.db.SearchUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "setMyUserName: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	var updatedUser User
	updatedUser.FromDatabase(dbUser)
	err = json.NewDecoder(r.Body).Decode(&updatedUser)

	// BadRequest check
	if err != nil {
		stringErr := "setMyUserName: invalid JSON object"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	if !updatedUser.HasValidUsername() {
		stringErr := "setMyUserName: invalid username"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	_, present, err = rt.db.SearchUserByUsername(updatedUser.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if present {
		stringErr := "setMyUserName: username already exists"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	// database section
	err = rt.db.UpdateUsername(updatedUser.ToDatabase())

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(updatedUser)
}
