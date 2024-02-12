package gormfile

import (
	"errors"
	"io/fs"
	"net/http"

	"gorm.io/gorm"
)

// hfs implements http.FileSystem interface
type hfs struct {
	db *gorm.DB
	tn string // table name
}

func FS(db *gorm.DB, table string) http.FileSystem {
	return &hfs{db, table}
}

func (hfs *hfs) Open(name string) (http.File, error) {
	db := hfs.db

	f := &File{ID: name}
	r := db.Table(hfs.tn).Omit("data").Take(f)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, fs.ErrNotExist
	}
	if r.Error != nil {
		return nil, r.Error
	}

	hf := &file{hfs: hfs, f: f}
	return hf, nil
}
