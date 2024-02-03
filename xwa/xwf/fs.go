package xwf

import (
	"errors"
	"io/fs"
	"net/http"

	"gorm.io/gorm"
)

// hfs implements http.FileSystem interface
type hfs struct {
	db *gorm.DB
}

func FS(db *gorm.DB) http.FileSystem {
	return &hfs{db}
}

func (hfs *hfs) Open(name string) (http.File, error) {
	f := &File{ID: name}
	r := hfs.db.Omit("data").Take(f)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, fs.ErrNotExist
	}
	if r.Error != nil {
		return nil, r.Error
	}

	hf := &file{f: f, db: hfs.db}
	return hf, nil
}
