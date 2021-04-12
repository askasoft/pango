package iox

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pandafw/pango/str"
)

// OnModified file change callback function
type OnModified func(path string)

// FileWatcher struct for file watching
type FileWatcher struct {
	Delay     time.Duration
	fswatcher *fsnotify.Watcher
	callbacks map[string]OnModified
}

// NewFileWatcher create a FileWatcher
func NewFileWatcher() (*FileWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fw := &FileWatcher{
		Delay:     time.Second,
		fswatcher: w,
		callbacks: make(map[string]OnModified),
	}

	return fw, nil
}

// StartWatch start file watching go-routine
func (fw *FileWatcher) StartWatch() {
	go fw.watch()
}

// StopWatch stop file watching go-routine
func (fw *FileWatcher) StopWatch() error {
	return fw.fswatcher.Close()
}

// AddFile add a file to watch
func (fw *FileWatcher) AddFile(path string, callback OnModified) error {
	path = filepath.Clean(path)
	err := fw.fswatcher.Add(path)
	if err != nil {
		return err
	}
	fw.callbacks[path] = callback
	return nil
}

// RemoveFile stop watching the file
func (fw *FileWatcher) RemoveFile(path string) error {
	path = filepath.Clean(path)
	delete(fw.callbacks, path)
	return fw.fswatcher.Remove(path)
}

// AddFolder add files in the folder and sub folders to watch
func (fw *FileWatcher) AddFolder(path string, callback OnModified) error {
	path = filepath.Clean(path)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return fw.fswatcher.Add(path)
		}
		return nil
	})

	if err != nil {
		fw.callbacks[path] = callback
	}
	return err
}

// RemoveFolder stop watching the files in the folder and sub folders
func (fw *FileWatcher) RemoveFolder(path string) error {
	path = filepath.Clean(path)
	delete(fw.callbacks, path)

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return fw.fswatcher.Remove(path)
		}
		return nil
	})

	return err
}

func (fw *FileWatcher) watch() {
	lastWrites := make(map[string]time.Time)
	for {
		select {
		case event, ok := <-fw.fswatcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				last := lastWrites[event.Name]

				// some editor use create->rename to save file,
				// this cloud raise 2 WRITE event continously,
				// delay 1s for prevent duplicated event
				due := last.Add(fw.Delay)
				now := time.Now()
				if due.Before(now) {
					lastWrites[event.Name] = now
					var callback OnModified
					for k, v := range fw.callbacks {
						if str.StartsWith(event.Name, k) {
							callback = v
							break
						}
					}
					if callback != nil {
						go func(callback OnModified, path string) {
							time.Sleep(fw.Delay)
							callback(path)
						}(callback, event.Name)
					}
				}
			}
		}
	}
}
