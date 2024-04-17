package sqlxfs

import (
	"database/sql"
	"errors"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/askasoft/pango/sqx"
	"github.com/askasoft/pango/sqx/sqlx"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xfs"
)

// sfs implements xfs.XFS interface
type sfs struct {
	db sqlx.Sqlx
	tn string // table name
}

func FS(db sqlx.Sqlx, table string) xfs.XFS {
	return &sfs{db, table}
}

func (sfs *sfs) Open(name string) (fs.File, error) {
	f, err := sfs.FindFile(name)
	if err != nil {
		return nil, err
	}

	hf := &xfs.FSFile{XFS: sfs, File: f}
	return hf, nil
}

// FindFile find a file
func (sfs *sfs) FindFile(id string) (*xfs.File, error) {
	s := sfs.db.Rebind("SELECT id, name, ext, size, time FROM " + sfs.tn + " WHERE id = ?")

	f := &xfs.File{}
	err := sfs.db.Get(f, s, id)
	if err != nil {
		if errors.Is(err, sqlx.ErrNoRows) {
			return nil, fs.ErrNotExist
		}
		return nil, err
	}

	return f, nil
}

func (sfs *sfs) SaveFile(id string, filename string, modTime time.Time, data []byte) (*xfs.File, error) {
	name := filepath.Base(filename)
	fext := str.ToLower(filepath.Ext(filename))

	fi := &xfs.File{
		ID:   id,
		Name: name,
		Ext:  fext,
		Size: int64(len(data)),
		Time: modTime,
		Data: data,
	}

	var r sql.Result
	var err error

	if _, err = sfs.FindFile(id); err == nil {
		su := sfs.db.Rebind("UPDATE " + sfs.tn + " SET name = ?, ext = ?, size = ?, time = ?, data = ? WHERE id = ?")
		ps := []any{fi.Name, fi.Ext, fi.Size, fi.Time, fi.Data, fi.ID}
		r, err = sfs.db.Exec(su, ps...)
	} else {
		si := sfs.db.Rebind("INSERT INTO " + sfs.tn + " (id, name, ext, size, time, data) VALUES (?, ?, ?, ?, ?, ?)")
		ps := []any{fi.ID, fi.Name, fi.Ext, fi.Size, fi.Time, fi.Data}
		r, err = sfs.db.Exec(si, ps...)
	}

	if err != nil {
		return fi, err
	}

	n, err := r.RowsAffected()
	if err != nil {
		return fi, err
	}

	if n != 1 {
		return fi, fs.ErrNotExist
	}
	return fi, nil
}

func (sfs *sfs) ReadFile(id string) ([]byte, error) {
	s := sfs.db.Rebind("SELECT * FROM " + sfs.tn + " WHERE id = ?")

	f := &xfs.File{}
	err := sfs.db.Get(f, s, id)
	if err != nil {
		if errors.Is(err, sqlx.ErrNoRows) {
			return nil, fs.ErrNotExist
		}
		return nil, err
	}

	return f.Data, nil
}

func (sfs *sfs) DeleteFile(id string) error {
	s := sfs.db.Rebind("DELETE FROM " + sfs.tn + " WHERE id = ?")
	_, err := sfs.db.Exec(s, id)
	return err
}

func (sfs *sfs) DeleteFiles(ids ...string) (int64, error) {
	sql, args := sqx.In("id", ids)
	return sfs.DeleteWhere(sql, args...)
}

func (sfs *sfs) DeletePrefix(prefix string) (int64, error) {
	return sfs.DeleteWhere("id LIKE ?", sqx.StartsLike(prefix))
}

func (sfs *sfs) DeleteBefore(before time.Time) (int64, error) {
	return sfs.DeleteWhere("time < ?", before)
}

func (sfs *sfs) DeletePrefixBefore(prefix string, before time.Time) (int64, error) {
	return sfs.DeleteWhere("id LIKE ? AND time < ?", sqx.StartsLike(prefix), before)
}

func (sfs *sfs) DeleteWhere(where string, args ...any) (int64, error) {
	s := sfs.db.Rebind("DELETE FROM " + sfs.tn + " WHERE " + where)
	r, err := sfs.db.Exec(s, args...)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

// DeleteAll use "DELETE FROM files" to delete all files
func (sfs *sfs) DeleteAll() (int64, error) {
	r, err := sfs.db.Exec("DELETE FROM " + sfs.tn)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

// Truncate use "TRUNCATE TABLE files" to truncate files
func (sfs *sfs) Truncate() error {
	_, err := sfs.db.Exec("TRUNCATE TABLE " + sfs.tn)
	return err
}
