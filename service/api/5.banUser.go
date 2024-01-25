package api

/*
go run ./cmd/webapi/
curl -v \
	-X PUT \
	-H 'Authorization: 2' \
	localhost:3000/banned/{1}
*/

import (
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "banUser: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	banner, present, err := rt.db.SearchUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "banUser: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	var pathUid uint64
	pathUid, err = strconv.ParseUint(ps.ByName("uid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "banUser: invalid path parameter uid"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	banned, present, err := rt.db.SearchUserByID(pathUid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "banUser: path parameter uid not matching any existing user"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	if banner.ID == banned.ID {
		stringErr := "banUser: requesting user trying to ban himself"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	isBanned, err := rt.db.CheckBan(banner.ID, banned.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isBanned {
		stringErr := "banUser: requesting user already banned user"
		http.Error(w, stringErr, http.StatusForbidden)
		return
	}

	// database section
	err = rt.db.BanUser(banner.ID, banned.ID)

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// cascade ban
	err = rt.db.CascadeBanBothDirections(banner.ID, banned.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
