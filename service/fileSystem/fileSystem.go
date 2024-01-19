package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

var Root = filepath.Join("/tmp", "photos")

// Funzione per creare una nuova cartella per un specifico utente
func CreateUserFolder(userID int64, username string) error {

	userFolder := fmt.Sprint(userID) + "." + username

	// Creao il path media/useridentifier/ dentro il project dir
	path := filepath.Join(Root, userFolder)

	// nel path creato aggiungo un subdir "photos"
	err := os.MkdirAll(path, os.ModePerm)

	return err
}