package api

import "strings"

/*
This function checks if the username is valid or not, hence:

	- if the client has entered white spaces only or;
	- if the username is not matching the required string pattern or;
	- if the username hasn't got the required length.
*/
func isValid(username string) bool {
	username = strings.TrimSpace(username) // leading and trailing white spaces removed

	// checking if username is present (the client may have entered white spaces only)
	if username == "" {
		return false
	}
	
	// checking if username matches string pattern: '^.*?$' (it must not contain \n)
	if strings.ContainsAny(username, "\n") {
		return false
	}

	// checking if username has required length
	if len(username)<3 || len(username)>16 {
		return false
	}

	return true
}

