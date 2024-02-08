package api

import (
	"github.com/simolillo/WASAPhoto/service/database"
	"github.com/simolillo/WASAPhoto/service/fileSystem"
)

type Photo struct {
	ID       uint64 `json:"photoID"`
	AuthorID uint64 `json:"authorID"`
	Format   string `json:"format"`
	Date     string `json:"date"`
}

func (p *Photo) ToDatabase() database.Photo {
	return database.Photo{
		ID:       p.ID,
		AuthorID: p.AuthorID,
		Format:   p.Format,
		Date:     p.Date,
	}
}

func (p *Photo) FromDatabase(photo database.Photo) {
	p.ID = photo.ID
	p.AuthorID = photo.AuthorID
	p.Format = photo.Format
	p.Date = photo.Date
}

func (p *Photo) ToFileSystem() fs.Photo {
	return fs.Photo{
		ID:       p.ID,
		AuthorID: p.AuthorID,
		Format:   p.Format,
		Date:     p.Date,
	}
}

func (p *Photo) FromFileSystem(photo fs.Photo) {
	p.ID = photo.ID
	p.AuthorID = photo.AuthorID
	p.Format = photo.Format
	p.Date = photo.Date
}
