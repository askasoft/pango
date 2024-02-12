package xwf

import (
	"errors"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/net/httpx"
	"github.com/askasoft/pango/squ"
	"github.com/askasoft/pango/str"
	"gorm.io/gorm"
)

func SaveFile(db *gorm.DB, table string, id string, filename string, data []byte, modTime ...time.Time) (*File, error) {
	name := filepath.Base(filename)
	fext := str.ToLower(filepath.Ext(filename))

	fi := &File{
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

	r := db.Table(table).Save(fi)
	return fi, r.Error
}

func SaveLocalFile(db *gorm.DB, table string, id string, filename string) (*File, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	data, err := fsu.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return SaveFile(db, table, id, filename, data, fi.ModTime())
}

func SaveUploadedFile(db *gorm.DB, table string, id string, file *multipart.FileHeader) (*File, error) {
	data, err := httpx.ReadMultipartFile(file)
	if err != nil {
		return nil, err
	}

	return SaveFile(db, table, id, file.Filename, data)
}

func ReadFile(db *gorm.DB, table string, id string) ([]byte, error) {
	f := &File{}
	r := db.Table(table).Where("id = ?", id).Take(f)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, fs.ErrNotExist
	}
	if r.Error != nil {
		return nil, r.Error
	}
	return f.Data, nil
}

func DeleteFile(db *gorm.DB, table string, id string) error {
	r := db.Table(table).Where("id = ?", id).Delete(&File{})
	return r.Error
}

func DeleteFiles(db *gorm.DB, table string, ids ...string) error {
	r := db.Table(table).Where("id IN ?", ids).Delete(&File{})
	return r.Error
}

func DeleteWhere(db *gorm.DB, table string, where string, args ...any) error {
	r := db.Table(table).Where(where, args...).Delete(&File{})
	return r.Error
}

func DeletePrefix(db *gorm.DB, table string, prefix string) error {
	r := db.Table(table).Where("id LIKE ?", squ.StartsLike(prefix)).Delete(&File{})
	return r.Error
}

// CleanFiles use "DELETE FROM files" to clean files
func CleanFiles(db *gorm.DB, table string) error {
	r := db.Exec("DELETE FROM " + table)
	return r.Error
}

// TruncateFiles use "TRUNCATE TABLE files" to truncate files
func TruncateFiles(db *gorm.DB, table string) error {
	r := db.Exec("TRUNCATE TABLE " + table)
	return r.Error
}

func CleanOutdatedFiles(db *gorm.DB, table string, before time.Time, loggers ...log.Logger) {
	logger := getLogger(loggers...)

	logger.Debugf("CleanOutdatedFiles('%s', '%v')", table, before)

	r := db.Table(table).Where("time < ?", before).Delete(&File{})
	if r.Error != nil {
		logger.Errorf("CleanOutdatedFiles('%s', '%v') failed: %v", table, before, r.Error)
		return
	}

	logger.Infof("CleanOutdatedFiles('%s', '%v'): %d", table, before, r.RowsAffected)
}

func getLogger(loggers ...log.Logger) log.Logger {
	if len(loggers) > 0 {
		return loggers[0]
	}
	return log.GetLogger("XWF")
}
