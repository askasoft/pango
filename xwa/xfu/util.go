package xfu

import (
	"os"
	"path/filepath"
	"time"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/log"
)

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