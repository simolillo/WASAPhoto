package api

/*
go run ./cmd/webapi/
curl -v \
	-X GET \
	-H 'Authorization: 1' \
	localhost:3000/users/{1}
*/

import (
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"net/http"
	"strconv"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var token uint64
	token, err := strconv.ParseUint(r.Header.Get("Authorization"), 10, 64)

	// Unauthorized check
	if err != nil {
		stringErr := "getUserProfile: invalid authorization token"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}
	requestingUser, present, err := rt.db.SearchUserByID(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "getUserProfile: authorization token not matching any existing user"
		http.Error(w, stringErr, http.StatusUnauthorized)
		return
	}

	var pathUid uint64
	pathUid, err = strconv.ParseUint(ps.ByName("uid"), 10, 64)

	// BadRequest check
	if err != nil {
		stringErr := "getUserProfile: invalid path parameter uid"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}
	requestedUser, present, err := rt.db.SearchUserByID(pathUid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !present {
		stringErr := "getUserProfile: path parameter uid not matching any existing user"
		http.Error(w, stringErr, http.StatusBadRequest)
		return
	}

	// Forbidden check
	someoneIsBanned, err := rt.db.CheckBanBothDirections(requestingUser.ID, requestedUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if someoneIsBanned {
		stringErr := "getUserProfile: someone has banned the other"
		http.Error(w, stringErr, http.StatusForbidden)
		return
	}

	var requestedProfile Profile
	isFollowing, err := rt.db.CheckFollow(requestingUser.ID, requestedUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// database section
	dbProfile, err := rt.db.GetUserProfile(requestedUser.ID)
	requestedProfile.FromDatabase(dbProfile)
	requestedProfile.IsFollowedByViewer = isFollowing

	// InternalServerError check
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(requestedProfile)
}