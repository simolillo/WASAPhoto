package fs

import (
	"fmt"
	"path/filepath"
)

type Photo struct {
	ID       uint64 `json:"photoID"`
	AuthorID uint64 `json:"authorID"`
	Format   string `json:"format"`
	Date     string `json:"date"`
}

func (p *Photo) Path() (photoPath string) {
	return filepath.Join(Root, fmt.Sprintf("%d.%s", p.ID, p.Format))
}
