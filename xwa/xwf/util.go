package xwf

import (
	"mime"
	"mime/multipart"
	"path"
	"path/filepath"
	"time"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/squ"
	"gorm.io/gorm"
)

func SaveLocalFile(db *gorm.DB, table string, id string, filename string) (*File, error) {
	base := filepath.Base(filename)
	ext := path.Ext(filename)

	data, err := fsu.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fi := &File{
		ID:        id,
		Name:      base,
		Size:      int64(len(data)),
		Type:      mime.TypeByExtension(ext),
		Data:      data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	r := db.Table(table).Save(fi)
	return fi, r.Error
}

func SaveUploadedFile(db *gorm.DB, table string, id string, file *multipart.FileHeader) (*File, error) {
	ext := path.Ext(file.Filename)

	fr, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fr.Close()

	data := make([]byte, file.Size)
	_, err = fr.Read(data)
	if err != nil {
		return nil, err
	}

	fi := &File{
		ID:        id,
		Name:      file.Filename,
		Size:      file.Size,
		Type:      mime.TypeByExtension(ext),
		Data:      data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	r := db.Table(table).Save(fi)
	return fi, r.Error
}

func DeleteFile(db *gorm.DB, table string, id string) error {
	r := db.Table(table).Where("id = ?", id).Delete(&File{})
	return r.Error
}

func DeleteFilesByPrefix(db *gorm.DB, table string, prefix string) error {
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

	r := db.Table(table).Where("updated_at < ?", before).Delete(&File{})
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
