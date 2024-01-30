package xfu

import (
	"mime"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/xin"
	"github.com/google/uuid"
)

type FileItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}

type FileResult struct {
	File *FileItem `json:"file"`
}

type FilesResult struct {
	Files []*FileItem `json:"files"`
}

func SaveUploadedFile(c *xin.Context, dir string, file *multipart.FileHeader) (*FileItem, error) {
	ext := str.IfEmpty(path.Ext(file.Filename), ".x")

	name := str.RemoveByte(uuid.New().String(), '-') + ext

	if err := c.SaveUploadedFile(file, path.Join(dir, name)); err != nil {
		return nil, err
	}

	fi := &FileItem{
		ID:   name,
		Name: file.Filename,
		Size: file.Size,
		Type: mime.TypeByExtension(ext),
	}
	return fi, nil
}

func getLogger(loggers ...log.Logger) log.Logger {
	if len(loggers) > 0 {
		return loggers[0]
	}
	return log.GetLogger("XFU")
}

func CleanOutdatedFiles(dir string, due time.Time, loggers ...log.Logger) {
	logger := getLogger(loggers...)

	logger.Debugf("CleanOutdatedFiles('%s', '%v')", dir, due)

	f, err := os.Open(dir)
	if err != nil {
		logger.Errorf("Open('%s') failed: %v", dir, err)
		return
	}
	defer f.Close()

	des, err := f.ReadDir(-1)
	if err != nil {
		logger.Error("ReadDir('%s') failed: %v", dir, err)
		return
	}

	for _, de := range des {
		path := filepath.Join(dir, de.Name())

		if de.IsDir() {
			CleanOutdatedFiles(path, due, logger)
			if err := fsu.DirIsEmpty(path); err != nil {
				log.Errorf("DirIsEmpty('%s') failed: %v", path, err)
			} else {
				if err := os.Remove(path); err != nil {
					log.Errorf("Remove('%s') failed: %v", path, err)
				} else {
					log.Debugf("Remove('%s') OK", path)
				}
			}
			continue
		}

		if fi, err := de.Info(); err != nil {
			log.Errorf("DirEntry('%s').Info() failed: %v", path, err)
		} else {
			if fi.ModTime().Before(due) {
				if err := os.Remove(path); err != nil {
					log.Errorf("Remove('%s') failed: %v", path, err)
				} else {
					log.Debugf("Remove('%s') OK", path)
				}
			}
		}
	}
}
