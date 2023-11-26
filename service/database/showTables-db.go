package database

// This function returns all the records present in the "users" table of the database.
func (db *appdbimpl) ShowUsersTable() ([]User, error) {
	rows, err := db.c.Query(`SELECT * FROM users`)

	if err != nil {
		return nil, err
	}

	// wait for the function to finish before closing rows
	defer func() { _ = rows.Close() }()

	// read all the users in the resulset
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name)
		
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if rows.Err() != nil { // rows.Err() returns the last error that occurred during the iteration of rows
		return nil, err
	}

	return users, nil
}

