package fswatch

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/log"
	"github.com/pandafw/pango/str"
	"github.com/pandafw/pango/str/wildcard"
)

// Op describes a set of file operations.
type Op = fsnotify.Op

const (
	// OpCreate create operation
	OpCreate = Op(fsnotify.Create)

	// OpWrite write operation
	OpWrite = Op(fsnotify.Write)

	// OpRemove remove operation
	OpRemove = Op(fsnotify.Remove)

	// OpRename rename operation
	OpRename = Op(fsnotify.Rename)

	// OpChmod chmod operation
	OpChmod = Op(fsnotify.Chmod)

	// OpModifies modifies operations (OpCreate | OpWrite | OpRemove | OpRename)
	OpModifies = OpCreate | OpWrite | OpRemove | OpRename

	// OpALL all operations
	OpALL = Op(0xFFFFFFFF)
)

type fileitem struct {
	OpMask   Op
	FnMask   string
	Callback func(string, Op)
}

type fileevent struct {
	Name      string
	Op        Op
	Time      time.Time
	Callbacks []func(string, Op)
}

// FileWatcher struct for file watching
type FileWatcher struct {
	Delay     time.Duration // Event delay
	Logger    log.Logger    // Error logger
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

// Add add a file to watch on specified operation op occured
func (fw *FileWatcher) Add(path string, op Op, callback func(string, Op)) error {
	path = filepath.Clean(path)
	err := fw.fswatcher.Add(path)
	if err != nil {
		return err
	}
	fw.items[path] = &fileitem{OpMask: op, Callback: callback}
	return nil
}

// Remove stop watching the file
func (fw *FileWatcher) Remove(path string) error {
	path = filepath.Clean(path)
	delete(fw.items, path)
	return fw.fswatcher.Remove(path)
}

// AddRecursive add files and all sub-directories under the path to watch
// op: operation mask
// fn: file path wildcard mask, "" or "*" means no mask
func (fw *FileWatcher) AddRecursive(path string, op Op, fn string, cb func(string, Op)) error {
	path = filepath.Clean(path)
	if err := fw.doRecursive(path, true); err != nil {
		return err
	}
	fw.items[path] = &fileitem{OpMask: op, FnMask: fn, Callback: cb}
	return nil
}

// RemoveRecursive stops watching the named directory and all sub-directories.
func (fw *FileWatcher) RemoveRecursive(name string) error {
	if err := fw.doRecursive(name, false); err != nil {
		return err
	}
	return nil
}

// doRecursive adds all directories under the given one to the watch list.
// this is probably a very racey process. What if a file is added to a folder before we get the watch added?
func (fw *FileWatcher) doRecursive(root string, watch bool) error {
	err := filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			if watch {
				if err = fw.fswatcher.Add(path); err != nil {
					return err
				}
			} else {
				if err = fw.fswatcher.Remove(path); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

func (fw *FileWatcher) delayCallbacks(fe *fileevent) {
	if fw.Logger != nil {
		fw.Logger.Infof("fswatch delay callback %q [%v]", fe.Name, fe.Op)
	}
	time.Sleep(fw.Delay)

	fw.mu.Lock()
	defer fw.mu.Unlock()

	for _, cb := range fe.Callbacks {
		cb(fe.Name, Op(fe.Op))
	}

	delete(fw.events, fe.Name)
}

func (fw *FileWatcher) findCallbacks(fe *fileevent) (cbs []func(string, Op)) {
	for k, i := range fw.items {
		if fe.Op&i.OpMask != 0 && str.StartsWith(fe.Name, k) {
			if i.FnMask == "" || wildcard.Match(i.FnMask, fe.Name) {
				cbs = append(cbs, i.Callback)
			}
		}
	}
	return
}

// procEvent add or remove watch file
func (fw *FileWatcher) procEvent(evt *fsnotify.Event) {
	if evt.Op&OpCreate != 0 {
		if err := iox.DirExists(evt.Name); err == nil {
			if err = fw.doRecursive(evt.Name, true); err != nil && fw.Logger != nil {
				fw.Logger.Errorf("fswatch add %q error: %v", evt.Name, err)
			}
		}
	}

	if evt.Op&OpRemove != 0 {
		if err := fw.fswatcher.Remove(evt.Name); err != nil && fw.Logger != nil {
			fw.Logger.Errorf("fswatch remove %q error: %v", evt.Name, err)
		}
	}
}

func (fw *FileWatcher) sendEvent(event *fsnotify.Event) {
	fe := fw.events[event.Name]
	if fe == nil {
		fe := &fileevent{Name: event.Name, Op: Op(event.Op), Time: time.Now()}
		fe.Callbacks = fw.findCallbacks(fe)
		if len(fe.Callbacks) > 0 {
			fw.events[event.Name] = fe
			go fw.delayCallbacks(fe)
		}
		return
	}

	// some editor use create->rename to save file,
	// this cloud raise 2 OpWrite event continously,
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
		fe.Op |= Op(event.Op)
		fe.Callbacks = fw.findCallbacks(fe)
	}
}

func (fw *FileWatcher) doEvent(evt *fsnotify.Event) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	fw.procEvent(evt)
	fw.sendEvent(evt)
}

func (fw *FileWatcher) watch() {
	for {
		select {
		case evt, ok := <-fw.fswatcher.Events:
			if !ok {
				return
			}
			fw.doEvent(&evt)
		case err, ok := <-fw.fswatcher.Errors:
			if !ok {
				return
			}
			if fw.Logger != nil {
				fw.Logger.Error("fswatch ", err)
			}
		}
	}
}
