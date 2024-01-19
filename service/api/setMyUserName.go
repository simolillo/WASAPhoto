package api

/*
go run ./cmd/webapi/
curl -v \
	-X PUT \
	-H 'Content-Type: text/plain' \
	-H 'Authorization: 1' \
	-d "Lillino" \
	localhost:3000/users/{1}/username
*/

/*
suppose user: 3 Anna
curl -v	-X POST	-H 'Content-Type: text/plain' -d "Anna" localhost:3000/session

Possible outcomes:

1. checking if the request is valid
   a. curl -v -X PUT -H 'Content-Type: text/plain' -H 'Authorization: 3' -d "Fabia" localhost:3000/users/{aaa}/username
      (the path parameter {uid} is not a parseable int64)

   b. curl -v -X PUT -H 'Content-Type: text/plain' -H 'Authorization: 3' -d "Fabia" localhost:3000/users/{1000}/username
      (the {uid} path parameter is not matching any existing user)

2. authentication phase
   a. curl -v -X PUT -H 'Content-Type: text/plain' -d "Fabia" localhost:3000/users/{3}/username
      (the Authorization header is not present or no value is specified)

   b. curl -v -X PUT -H 'Content-Type: text/plain' -H 'Authorization: 1000' -d "Fabia" localhost:3000/users/{3}/username
      (the Authorization ID is not matching any existing user)

3. authorization phase
   curl -v -X PUT -H 'Content-Type: text/plain' -H 'Authorization: 1' -d "Fabia" localhost:3000/users/{3}/username
   (the ID of the user attempting the request is different from the one he wants to update the username of)

4. checking if decoding operation of username ended successfully
   curl -v -X PUT -H 'Content-Type: text/plain' -H 'Authorization: 3' -d "Fabia localhost:3000/users/{3}/username
   (the text/plain data is missing a closing double-quote resulting in an invalid text/plain)

5. checking if the new username is valid
   a. curl -v -X PUT -H 'Content-Type: text/plain' -H 'Authorization: 3' -d "     " localhost:3000/users/{3}/username
      (the client has enterd white spaces only, hence the username is not valid)

   b. (username doesn't match string pattern: '^.*?$': it contains a new line)

   c. curl -v -X PUT -H 'Content-Type: text/plain' -H 'Authorization: 3' -d "Fa" localhost:3000/users/{3}/username
      (username hasn't got required length: is <3 or >16)

6. updating the username
   curl -v -X PUT -H 'Content-Type: text/plain' -H 'Authorization: 3' -d "Fabia" localhost:3000/users/{3}/username
   curl -v -X PUT -H 'Content-Type: text/plain' -H 'Authorization: 3' -d "Anna" localhost:3000/users/{3}/username

	7. if oldUsername = newUsername
	   (update same username twice)

	8. if user update is unsuccessful
	   (server error)

	9. if encoding operation is unsuccessful though the user has been updated
	   (server error)
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/simolillo/WASAPhoto/service/api/reqcontext"
)

/*
Summary: update personal username

Description:
Update personal username with the new string provided in the request body.
*/
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// extracting username from the request
	body, err := io.ReadAll(r.Body)
	newUsername := string(body)

	// 4.
	// checking if decoding operation of username ended successfully
	if err != nil {
		// the request body (the username) was not a valid text/plain or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("setMyUserName: the request body (the username) was not a valid text/plain or is missing")
		fmt.Fprint(w, "\nsetMyUserName: the request body (the username) was not a valid text/plain or is missing\n\n")
		return
	}

	// 5.
	// checking if the new username is valid
	if !isValid(newUsername) {
		// the new username is not valid, rejecting request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("setMyUserName: the new username is not valid")
		fmt.Fprint(w, "\nsetMyUserName: the new username is not valid\n\n")
		return
	}

	// moving on to the database section

	// 6.
	// updating the username
	userID := selectedUser1.ID
	oldUsername := selectedUser1.Name

	// 7.
	// if oldUsername = newUsername
	if oldUsername == newUsername {
		w.WriteHeader(http.StatusOK) //200
		fmt.Fprint(w, "\nUsername already updated: ", oldUsername, ".\nNo need to change it.\n\n")
		return
	}

	updatedUser, err := rt.db.UpdateUsername(userID, newUsername)

	// 8.
	// if user update is unsuccessful
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("setMyUserName: user update was unsuccessful")
		fmt.Fprint(w, "\nsetMyUserName: user update was unsuccessful\n\n")
		return
	}

	fmt.Fprintln(w)
	err = json.NewEncoder(w).Encode(updatedUser)

	// 9.
	// if encoding operation is unsuccessful though the user has been updated
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("setMyUserName: unable to encode JSON response though the user has been updated")
		fmt.Fprint(w, "\nsetMyUserName: unable to encode JSON response though the user has been updated\n\n")
		return
	}

	w.WriteHeader(http.StatusOK) //200
	fmt.Fprint(w, "\nUsername successfully updated.\nThe updated user is returned in the content.\n\n")

}
