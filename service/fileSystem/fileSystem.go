package fs

// "github.com/simolillo/WASAPhoto/service/fileSystem"

import (
	"fmt"
	"os"
	"path/filepath"
)

var Root = filepath.Join("/private", "tmp", "photos")


// Funzione per creare una nuova photo file
func CreatePhotoFile(photo Photo, binaryImage []byte) error {

	err := os.MkdirAll(Root, 0755)
	if err != nil {
		return err
	}

	photoPath := filepath.Join(Root, fmt.Sprint(photo.ID) + "." + photo.Format)

	createdFile, err := os.Create(photoPath)
	if err != nil {
		return err
	}
	defer createdFile.Close()

	_, err = createdFile.Write(binaryImage)

	return err
}

// Function to delete a photo file
func DeletePhotoFile(photo Photo) error {
    photoPath := filepath.Join(Root, fmt.Sprintf("%d.%s", photo.ID, photo.Format))

    err := os.Remove(photoPath)
    return err
}

