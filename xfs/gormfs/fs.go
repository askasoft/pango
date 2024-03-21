package gormfs

import (
	"errors"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/askasoft/pango/sqx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xfs"
	"gorm.io/gorm"
)

// gfs implements fs.FS interface
type gfs struct {
	db *gorm.DB
	tn string // table name
}

func FS(db *gorm.DB, table string) xfs.XFS {
	return &gfs{db, table}
}

func (gfs *gfs) Open(name string) (fs.File, error) {
	f := &xfs.File{}
	r := gfs.db.Table(gfs.tn).Omit("data").Where("id = ?", name).First(f)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, fs.ErrNotExist
	}
	if r.Error != nil {
		return nil, r.Error
	}

	hf := &xfs.FSFile{XFS: gfs, File: f}
	return hf, nil
}

func (gfs *gfs) SaveFile(id string, filename string, data []byte, modTime ...time.Time) (*xfs.File, error) {
	name := filepath.Base(filename)
	fext := str.ToLower(filepath.Ext(filename))

	fi := &xfs.File{
		ID:   id,
		Name: name,
		Ext:  fext,
		Size: int64(len(data)),
		Data: data,
	}

	if len(modTime) > 0 {
		fi.Time = modTime[0]
	}
	if fi.Time.IsZero() {
		fi.Time = time.Now()
	}

	r := gfs.db.Table(gfs.tn).Save(fi)
	return fi, r.Error
}

func (gfs *gfs) ReadFile(id string) ([]byte, error) {
	f := &xfs.File{}
	r := gfs.db.Table(gfs.tn).Where("id = ?", id).First(f)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, fs.ErrNotExist
	}
	if r.Error != nil {
		return nil, r.Error
	}
	return f.Data, nil
}

func (gfs *gfs) DeleteFile(id string) error {
	r := gfs.db.Table(gfs.tn).Where("id = ?", id).Delete(&xfs.File{})
	return r.Error
}

func (gfs *gfs) DeleteFiles(ids ...string) (int64, error) {
	r := gfs.db.Table(gfs.tn).Where("id IN ?", ids).Delete(&xfs.File{})
	return r.RowsAffected, r.Error
}

func (gfs *gfs) DeletePrefix(prefix string) (int64, error) {
	r := gfs.db.Table(gfs.tn).Where("id LIKE ?", sqx.StartsLike(prefix)).Delete(&xfs.File{})
	return r.RowsAffected, r.Error
}

func (gfs *gfs) DeleteBefore(before time.Time) (int64, error) {
	r := gfs.db.Table(gfs.tn).Where("time < ?", before).Delete(&xfs.File{})
	return r.RowsAffected, r.Error
}

func (gfs *gfs) DeletePrefixBefore(prefix string, before time.Time) (int64, error) {
	r := gfs.db.Table(gfs.tn).Where("id LIKE ? AND time < ?", sqx.StartsLike(prefix), before).Delete(&xfs.File{})
	return r.RowsAffected, r.Error
}

func (gfs *gfs) DeleteWhere(where string, args ...any) (int64, error) {
	r := gfs.db.Table(gfs.tn).Where(where, args...).Delete(&xfs.File{})
	return r.RowsAffected, r.Error
}

// DeleteAll use "DELETE FROM files" to delete all files
func (gfs *gfs) DeleteAll() (int64, error) {
	r := gfs.db.Exec("DELETE FROM " + gfs.tn)
	return r.RowsAffected, r.Error
}

// Truncate use "TRUNCATE TABLE files" to truncate files
func (gfs *gfs) Truncate() error {
	r := gfs.db.Exec("TRUNCATE TABLE " + gfs.tn)
	return r.Error
}
