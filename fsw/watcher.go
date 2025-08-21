package fsw

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/log"
	"github.com/fsnotify/fsnotify"
)

// Op describes a set of file operations.
type Op = fsnotify.Op

const (
	// OpNone none operation
	OpNone = Op(0)

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
	OpMask    Op
	Recursive bool
	Callback  func(string, Op)
}

type fileevent struct {
	Name      string
	Op        Op
	Time      time.Time
	Callbacks []func(string, Op)
}

// FileWatcher struct for file watching
type FileWatcher struct {
	Delay  time.Duration
	Logger log.Logger

	fsnotify *fsnotify.Watcher
	items    map[string]*fileitem
	events   map[string]*fileevent
	mutex    sync.RWMutex
	timer    *time.Timer
}

// NewFileWatcher create a FileWatcher
func NewFileWatcher() *FileWatcher {
	fw := &FileWatcher{
		Delay:  time.Second,
		items:  make(map[string]*fileitem),
		events: make(map[string]*fileevent),
	}
	return fw
}

// Start start file watching go-routine
func (fw *FileWatcher) Start() error {
	if fw.fsnotify != nil {
		return nil
	}

	log := fw.Logger
	if log != nil {
		log.Info("fswatch: start")
	}

	fsn, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	fw.fsnotify = fsn
	for path, fi := range fw.items {
		if fi.Recursive {
			if err = fw.doRecursive(path, true); err != nil {
				return err
			}
		} else {
			if log != nil {
				log.Debugf("fswatch: watch file '%s'", path)
			}
			err = fsn.Add(path)
			if err != nil {
				return err
			}
		}
	}

	// some editor use create->rename to save file,
	// this could raise 2 OpWrite event continuously,
	// delay some time to prevent duplicated event.
	if fw.Delay.Milliseconds() > 0 {
		fw.timer = time.AfterFunc(time.Hour, func() {
			fw.delayCallbacks()
		})
	}

	go fw.watch(fsn)
	return nil
}

// Stop stop file watching go-routine
func (fw *FileWatcher) Stop() error {
	fsn := fw.fsnotify
	if fsn == nil {
		return nil
	}

	timer, log := fw.timer, fw.Logger

	fw.timer = nil
	fw.fsnotify = nil

	if log != nil {
		log.Info("fswatch: stop")
	}
	if timer != nil {
		timer.Stop()
	}

	return fsn.Close()
}

// Add add a file to watch on specified operation op occurred
func (fw *FileWatcher) Add(path string, op Op, callback func(string, Op)) error {
	path = filepath.Clean(path)
	fsn, log := fw.fsnotify, fw.Logger
	if fsn != nil {
		if log != nil {
			log.Debugf("fswatch: watch file '%s' ", path)
		}
		if err := fsn.Add(path); err != nil {
			return err
		}
	}
	fw.items[path] = &fileitem{OpMask: op, Callback: callback}
	return nil
}

// Remove stop watching the file
func (fw *FileWatcher) Remove(path string) error {
	path = filepath.Clean(path)
	fsn, log := fw.fsnotify, fw.Logger
	if fsn != nil {
		if log != nil {
			log.Debugf("fswatch: unwatch file '%s'", path)
		}
		if err := fsn.Remove(path); err != nil {
			return err
		}
	}
	delete(fw.items, path)
	return nil
}

// AddRecursive add files and all sub-directories under the path to watch
func (fw *FileWatcher) AddRecursive(path string, op Op, cb func(string, Op)) error {
	path = filepath.Clean(path)
	if err := fw.doRecursive(path, true); err != nil {
		return err
	}
	fw.items[path] = &fileitem{OpMask: op, Recursive: true, Callback: cb}
	return nil
}

// RemoveRecursive stops watching the directory and all sub-directories.
func (fw *FileWatcher) RemoveRecursive(path string) error {
	path = filepath.Clean(path)
	if err := fw.doRecursive(path, false); err != nil {
		return err
	}
	delete(fw.items, path)
	return nil
}

// doRecursive adds all directories under the given one to the watch list.
// this is probably a very racey process. What if a file is added to a folder before we get the watch added?
func (fw *FileWatcher) doRecursive(root string, watch bool) error {
	fsn, log := fw.fsnotify, fw.Logger
	if fsn == nil {
		return nil
	}

	return filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			if watch {
				if log != nil {
					log.Debugf("fswatch: watch dir '%s'", path)
				}
				if err = fsn.Add(path); err != nil {
					return err
				}
			} else {
				if log != nil {
					log.Debugf("fswatch: unwatch dir '%s'", path)
				}
				if err = fsn.Remove(path); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (fw *FileWatcher) findCallbacks(fe *fileevent) (cbs []func(string, Op)) {
	for k, i := range fw.items {
		if fe.Op&i.OpMask != 0 && strings.HasPrefix(fe.Name, k) {
			cbs = append(cbs, i.Callback)
		}
	}
	return
}

func (fw *FileWatcher) isRecursive(name string) bool {
	for k, i := range fw.items {
		if strings.HasPrefix(name, k) {
			if i.Recursive {
				return true
			}
		}
	}
	return false
}

func (fw *FileWatcher) watch(fsn *fsnotify.Watcher) {
	for {
		select {
		case evt, ok := <-fsn.Events:
			if !ok {
				return
			}
			if fsn == fw.fsnotify {
				fw.doEvent(&evt)
			}
		case err, ok := <-fsn.Errors:
			if !ok {
				return
			}
			log := fw.Logger
			if log != nil {
				log.Errorf("fswatch: watch error: %v", err)
			}
		}
	}
}

func (fw *FileWatcher) doEvent(evt *fsnotify.Event) {
	fw.procEvent(evt)
	fw.sendEvent(evt)
}

// procEvent add or remove watch file
func (fw *FileWatcher) procEvent(evt *fsnotify.Event) {
	if evt.Op&OpCreate != 0 {
		if err := fsu.DirExists(evt.Name); err == nil {
			if fw.isRecursive(evt.Name) {
				if err = fw.doRecursive(evt.Name, true); err != nil {
					if log := fw.Logger; log != nil {
						log.Errorf("fswatch: watch dir '%s' error: %v", evt.Name, err)
					}
				}
			}
		}
	}
}

func (fw *FileWatcher) sendEvent(event *fsnotify.Event) {
	timer := fw.timer
	if timer == nil {
		fe := &fileevent{Name: event.Name, Op: Op(event.Op), Time: time.Now()}
		fe.Callbacks = fw.findCallbacks(fe)
		if len(fe.Callbacks) > 0 {
			fw.execCallbacks(fe)
		}
		return
	}

	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	log := fw.Logger

	fe := fw.events[event.Name]
	if fe != nil {
		fe.Time = time.Now()
		fe.Op |= Op(event.Op)
		if log != nil {
			log.Debugf("fswatch: delay '%s' [%v] (%s) ", fe.Name, fe.Op, fe.Time)
		}
		return
	}

	fe = &fileevent{Name: event.Name, Op: Op(event.Op), Time: time.Now()}
	fe.Callbacks = fw.findCallbacks(fe)
	if len(fe.Callbacks) > 0 {
		if log != nil {
			log.Debugf("fswatch: queue '%s' [%v] (%s) ", fe.Name, fe.Op, fe.Time)
		}
		fw.events[event.Name] = fe
		if len(fw.events) == 1 {
			if log != nil {
				log.Debugf("fswatch: reset timer (%s) ", fw.Delay)
			}
			timer.Reset(fw.Delay)
		}
	}
}

func (fw *FileWatcher) delayCallbacks() {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	due := time.Now().Add(-fw.Delay)
	for _, fe := range fw.events {
		if fe.Time.Before(due) {
			fw.execCallbacks(fe)
			delete(fw.events, fe.Name)
		}
	}

	if len(fw.events) > 0 {
		if log := fw.Logger; log != nil {
			log.Debugf("fswatch: reset timer (%s) ", fw.Delay)
		}
		fw.timer.Reset(fw.Delay)
	}
}

func (fw *FileWatcher) execCallbacks(fe *fileevent) {
	if log := fw.Logger; log != nil {
		log.Debugf("fswatch: execute callback '%s' [%v]", fe.Name, fe.Op)
	}

	for _, cb := range fe.Callbacks {
		cb(fe.Name, Op(fe.Op))
	}
}
