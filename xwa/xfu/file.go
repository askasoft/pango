package xfu

import (
	"os"
	"path/filepath"
	"time"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/log"
)

type File struct {
	ID   string    `gorm:"size:256;not null;primaryKey" json:"id"`
	Name string    `gorm:"not null;" json:"name"`
	Ext  string    `gorm:"not null;" json:"ext"`
	Size int64     `gorm:"not null;" json:"size"`
	Data []byte    `gorm:"not null" json:"-"`
	Time time.Time `gorm:"not null" json:"time,omitempty"`
}

type FileResult struct {
	File *File `json:"file"`
}

type FilesResult struct {
	Files []*File `json:"files"`
}

func CleanOutdatedFiles(dir string, before time.Time, loggers ...log.Logger) {
	logger := getLogger(loggers...)

	logger.Debugf("CleanOutdatedFiles('%s', '%v')", dir, before)

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
			CleanOutdatedFiles(path, before, logger)
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
			if fi.ModTime().Before(before) {
				if err := os.Remove(path); err != nil {
					log.Errorf("Remove('%s') failed: %v", path, err)
				} else {
					log.Debugf("Remove('%s') OK", path)
				}
			}
		}
	}
}

func getLogger(loggers ...log.Logger) log.Logger {
	if len(loggers) > 0 {
		return loggers[0]
	}
	return log.GetLogger("XFU")
}
