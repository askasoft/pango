package fswatch

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pandafw/pango/str"
)

const (
	// CREATE create event
	CREATE = uint32(fsnotify.Create)

	// WRITE write event
	WRITE = uint32(fsnotify.Write)

	// REMOVE remove event
	REMOVE = uint32(fsnotify.Remove)

	// RENAME rename event
	RENAME = uint32(fsnotify.Rename)

	// CHMOD chmod event
	CHMOD = uint32(fsnotify.Chmod)

	// ALL all events
	ALL = uint32(0xFFFFFFFF)
)

type fileitem struct {
	OpMask   uint32
	FnMask   string
	Callback func(string, uint32)
}

type fileevent struct {
	Name      string
	Op        uint32
	Time      time.Time
	Callbacks []func(string, uint32)
}

// FileWatcher struct for file watching
type FileWatcher struct {
	Delay     time.Duration // Event delay
	fswatcher *fsnotify.Watcher
	items     map[string]*fileitem
	events    map[string]*fileevent
	mu        sync.Mutex
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
		items:     make(map[string]*fileitem),
		events:    make(map[string]*fileevent),
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

// AddFile add a file to watch on specified event op occurs
func (fw *FileWatcher) AddFile(path string, op uint32, callback func(string, uint32)) error {
	path = filepath.Clean(path)
	err := fw.fswatcher.Add(path)
	if err != nil {
		return err
	}
	fw.items[path] = &fileitem{OpMask: op, Callback: callback}
	return nil
}

// AddFiles add files (match filename fn) under the path to watch on specified event op occurs
func (fw *FileWatcher) AddFiles(path string, op uint32, fn string, cb func(string, uint32)) error {
	path = filepath.Clean(path)
	err := fw.fswatcher.Add(path)
	if err != nil {
		return err
	}
	fw.items[path] = &fileitem{OpMask: op, FnMask: fn, Callback: cb}
	return nil
}

// RemoveFile stop watching the file
func (fw *FileWatcher) RemoveFile(path string) error {
	path = filepath.Clean(path)
	delete(fw.items, path)
	return fw.fswatcher.Remove(path)
}

func (fw *FileWatcher) delayCallbacks(fe *fileevent) {
	time.Sleep(fw.Delay)

	fw.mu.Lock()
	defer fw.mu.Unlock()

	for _, cb := range fe.Callbacks {
		cb(fe.Name, uint32(fe.Op))
	}

	delete(fw.events, fe.Name)
}

func (fw *FileWatcher) findCallbacks(fe *fileevent) (cbs []func(string, uint32)) {
	for k, i := range fw.items {
		if fe.Op&i.OpMask != 0 && str.StartsWith(fe.Name, k) {
			if i.FnMask == "" {
				cbs = append(cbs, i.Callback)
				continue
			}

			if ok, _ := filepath.Match(i.FnMask, fe.Name); ok {
				cbs = append(cbs, i.Callback)
				continue
			}
		}
	}
	return
}

func (fw *FileWatcher) doEvent(event *fsnotify.Event) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	fe := fw.events[event.Name]
	if fe == nil {
		fe := &fileevent{Name: event.Name, Op: uint32(event.Op), Time: time.Now()}
		fe.Callbacks = fw.findCallbacks(fe)
		if len(fe.Callbacks) > 0 {
			go fw.delayCallbacks(fe)
		}
		return
	}

	// some editor use create->rename to save file,
	// this cloud raise 2 WRITE event continously,
	// delay 1s for prevent duplicated event
	due := fe.Time.Add(fw.Delay)
	now := time.Now()
	if due.Before(now) {
		fe.Time = now
		fe.Callbacks = fw.findCallbacks(fe)
		if len(fe.Callbacks) > 0 {
			go fw.delayCallbacks(fe)
		}
	} else {
		fe.Op |= uint32(event.Op)
		fe.Callbacks = fw.findCallbacks(fe)
	}
}

func (fw *FileWatcher) watch() {
	for {
		select {
		case event, ok := <-fw.fswatcher.Events:
			if !ok {
				return
			}
			fw.doEvent(&event)
		}
	}
}
