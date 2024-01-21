package fs

// "github.com/simolillo/WASAPhoto/service/fileSystem"

import (

	"os"
	"path/filepath"
)

var Root = filepath.Join("/tmp", "photos")




// Funzione per creare una nuova photo file
func CreatePhotoFile(photo Photo, binaryImage []byte) error {

	err := os.MkdirAll(Root, 0755)
	if err != nil {
		return err
	}

	createdFile, err := os.Create(photo.Path)
	if err != nil {
		return err
	}
	defer createdFile.Close()

	_, err = createdFile.Write(binaryImage)

	return err
}