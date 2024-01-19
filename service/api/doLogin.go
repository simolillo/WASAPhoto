package api

/*
go run ./cmd/webapi/
curl -v \
	-X POST \
	-H 'Content-Type: text/plain' \
	-d "Giachi" \
	localhost:3000/session
*/

/*
Possible outcomes:

1. checking if decoding operation of username ended successfully
   curl -v -X POST -H 'Content-Type: text/plain' -d "Lillo localhost:3000/session
   (the text/plain data is missing a closing double-quote resulting in an invalid text/plain)

2. checking if the username is valid
   a. curl -v -X POST -H 'Content-Type: text/plain' -d "     " localhost:3000/session
      (the client has enterd white spaces only, hence the username is not valid)

   b. (username doesn't match string pattern: '^.*?$': it contains a new line)

   c. curl -v -X POST -H 'Content-Type: text/plain' -d "ab" localhost:3000/session
      (username hasn't got required length: is <3 or >16)

3. if the user altready exists, return the ID
   (post an alredy existing username)

	4. if encoding operation is unsuccessful though the user is present
	   (server error)

5. if the user doesn't exist yet, create it and return the ID
   (post a new username)

	6. if user creation or ID retrieval is unsuccessful
	   (server error)

	7. if encoding operation is unsuccessful though the user has been created
	   (server error)
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/simolillo/WASAPhoto/service/fileSystem"

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
	body, err := io.ReadAll(r.Body)
	username := string(body)

	// 1.
	// checking if decoding operation of username ended successfully
	if err != nil {
		// the request body (the username) was not a valid text/plain or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("doLogin: the request body (the username) was not a valid text/plain or is missing")
		fmt.Fprint(w, "\ndoLogin: the request body (the username) was not a valid text/plain or is missing\n\n")
		return
	}

	// 2.
	// checking if the username is valid
	if !isValid(username) {
		// the username is not valid, rejecting request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("doLogin: the username is not valid")
		fmt.Fprint(w, "\ndoLogin: the username is not valid\n\n")
		return
	}

	// moving on to the database section
	// first of all we search the user to see if it has alredy been created
	selectedUser, present := rt.db.SearchByUsername(username)

	// 3.
	// if the user already exists, return the ID
	if present {
		fmt.Fprintln(w)
		err = json.NewEncoder(w).Encode(selectedUser)

		// 4.
		// if encoding operation is unsuccessful though the user is present
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) //500
			ctx.Logger.WithError(err).Error("doLogin: unable to encode JSON response though the user is present")
			fmt.Fprint(w, "\ndoLogin: unable to encode JSON response though the user is present\n\n")
			return
		}

		w.WriteHeader(http.StatusOK) //200
		fmt.Fprint(w, "\nUser log-in action successful.\nThe user ID is returned in the content.\n\n")
		return
	}

	// 5.
	// if the user doesn't exist yet, create it and return the ID
	createdUser, err := rt.db.CreateUser(username)

	// 6.
	// if user creation or ID retrieval is unsuccessful
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("doLogin: user creation or ID retrieval was unsuccessful")
		fmt.Fprint(w, "\ndoLogin: user creation or ID retrieval was unsuccessful\n\n")
		return
	}

	err = fs.CreateUserFolder(createdUser.ID, createdUser.Name)
	if err != nil {
		ctx.Logger.WithError(err).Error("session/createUserFolder:: error creating directories for user")
		return
	}

	fmt.Fprintln(w)
	err = json.NewEncoder(w).Encode(createdUser)

	// 7.
	// if encoding operation is unsuccessful though the user has been created
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("doLogin: unable to encode JSON response though the user has been created")
		fmt.Fprint(w, "\ndoLogin: unable to encode JSON response though the user has been created\n\n")
		return
	}

	w.WriteHeader(http.StatusCreated) //201
	fmt.Fprint(w, "\nUser sign-up action successful.\nThe user ID has been created and is returned in the content.\n\n")

}

func createUserFolder(i int64, s string) {
	panic("unimplemented")
}

