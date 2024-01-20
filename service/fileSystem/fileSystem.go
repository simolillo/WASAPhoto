package fs

// "github.com/simolillo/WASAPhoto/service/fileSystem"

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

var Root = filepath.Join("/tmp", "photos")

// Funzione per creare una nuova cartella per un specifico utente
func CreateUserFolder(userID int64, username string) error {

	userFolderName := UserFolderName(userID, username)

	// Creao il path media/useridentifier/ dentro il project dir
	path := filepath.Join(Root, userFolderName)

	err := os.MkdirAll(path, os.ModePerm)

	return err
}

// Funzione per creare una nuova cartella per un specifico utente
func UserFolderName(userID int64, username string) string {

	userFolderName := fmt.Sprint(userID) + "." + username

	return userFolderName
}

func UpdateUserFolderName(userID int64, oldUsername string, newUsername string) error {
	

	oldPath := UserFolderName(userID, oldUsername)
	newPath := UserFolderName(userID, newUsername)
	err :=moveFiles(oldPath, newPath)

	return err
}

func moveFiles(oldPath, newPath string) error {
	return filepath.Walk(oldPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(oldPath, path)
		if err != nil {
			return err
		}

		newFilePath := filepath.Join(newPath, relPath)

		if info.IsDir() {
			return os.MkdirAll(newFilePath, info.Mode())
		}

		return os.Rename(path, newFilePath)
	})
}

// Funzione per creare una nuova photo file
func CreatePhotoFile(photo Photo, imageFile image.Image) error {

	createdFile, err := os.Create(photo.Path)
	if err != nil {
		return err
	}
	defer createdFile.Close()

	switch photo.Format {
	case "jpg":
		err = jpeg.Encode(createdFile, imageFile, nil)
	case "png":
		err = png.Encode(createdFile, imageFile)
	}

	return err
}