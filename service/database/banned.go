package database

func (db *appdbimpl) BanUser(bannerID uint64, bannedID uint64) (err error) {

	query := "INSERT INTO banned (bannerID, bannedID) VALUES (?,?)"

	_, err = db.c.Exec(query, bannerID, bannedID)
	return err
}

func (db *appdbimpl) CheckBan(bannerID uint64, bannedID uint64) (isBanned bool, err error) {

	query := "SELECT EXISTS (SELECT _ FROM banned WHERE bannerID = ? AND bannedID = ?);"
	
	err = db.c.QueryRow(query, bannerID, bannedID).Scan(&isBanned)
	return
}
