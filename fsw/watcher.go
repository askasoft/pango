package fsw

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/gog"
	"github.com/askasoft/pango/log"
	"github.com/fsnotify/fsnotify"
)

// Op describes a set of file operations.
type Op = fsnotify.Op

// Callback callback function
type Callback func(string, Op)

const (
	// OpNone none operation
	OpNone = Op(0)

	// OpCreate create operation
	OpCreate = fsnotify.Create

	// OpWrite write operation
	OpWrite = fsnotify.Write

	// OpRemove remove operation
	OpRemove = fsnotify.Remove

	// OpRename rename operation
	OpRename = fsnotify.Rename

	// OpChmod chmod operation
	OpChmod = fsnotify.Chmod

	// OpModifies modifies operations (OpCreate | OpWrite | OpRemove | OpRename)
	OpModifies = OpCreate | OpWrite | OpRemove | OpRename

	// OpALL all operations
	OpALL = Op(0xFFFFFFFF)
)

type filewatch struct {
	OpMask    Op
	Recursive bool
	Callback  Callback
}

type fileevent struct {
	Name      string
	Op        Op
	Time      time.Time
	Callbacks []Callback
}

// FileWatcher struct for file watching
type FileWatcher struct {
	Logger   log.Logger
	fsnotify *fsnotify.Watcher
	mutex    sync.RWMutex // watchs lock
	watchs   map[string]*filewatch
	events   map[string]*fileevent // buffered delay events
	delay    time.Duration
}

// NewFileWatcher create a FileWatcher with default 1sec delay.
// Some editor use create->rename to save file,
// This could raise 2 OpWrite events continuously,
// delay some time to prevent duplicated event.
func NewFileWatcher(delays ...time.Duration) *FileWatcher {
	fw := &FileWatcher{
		watchs: make(map[string]*filewatch),
		events: make(map[string]*fileevent),
		delay:  time.Second,
	}

	if len(delays) > 0 {
		fw.delay = delays[0]
	}

	return fw
}

// Add add a file to watch on specified operation op occurred
func (fw *FileWatcher) Add(path string, op Op, callback Callback) error {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	path = filepath.Clean(path)
	fsn, log := fw.fsnotify, fw.Logger
	if fsn != nil {
		if log != nil {
			log.Debugf("fsw: %p watch file '%s' ", fsn, path)
		}
		if err := fsn.Add(path); err != nil {
			return fmt.Errorf("fsw: Add('%s') - %w", path, err)
		}
	}

	fw.watchs[path] = &filewatch{OpMask: op, Callback: callback}
	return nil
}

// Remove stop watching the file
func (fw *FileWatcher) Remove(path string) error {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	path = filepath.Clean(path)
	fsn, log := fw.fsnotify, fw.Logger
	if fsn != nil {
		if log != nil {
			log.Debugf("fsw: %p unwatch file '%s'", fsn, path)
		}
		if err := fsn.Remove(path); err != nil {
			return fmt.Errorf("fsw: Remove('%s') - %w", path, err)
		}
	}

	delete(fw.watchs, path)
	return nil
}

// AddRecursive add files and all sub-directories under the path to watch
func (fw *FileWatcher) AddRecursive(path string, op Op, cb Callback) error {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	path = filepath.Clean(path)
	if err := fw.doRecursive(path, true); err != nil {
		return err
	}

	fw.watchs[path] = &filewatch{OpMask: op, Recursive: true, Callback: cb}
	return nil
}

// RemoveRecursive stops watching the directory and all sub-directories.
func (fw *FileWatcher) RemoveRecursive(path string) error {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	path = filepath.Clean(path)
	if err := fw.doRecursive(path, false); err != nil {
		return err
	}

	delete(fw.watchs, path)
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
					log.Debugf("fsw: %p watch folder '%s'", fsn, path)
				}
				if err = fsn.Add(path); err != nil {
					return fmt.Errorf("fsw: Add('%s') - %w", path, err)
				}
			} else {
				if log != nil {
					log.Debugf("fsw: %p unwatch folder '%s'", fsn, path)
				}
				if err = fsn.Remove(path); err != nil {
					return fmt.Errorf("fsw: Remove('%s') - %w", path, err)
				}
			}
		}
		return nil
	})
}

// Start start file watching go-routine
func (fw *FileWatcher) Start() error {
	fw.mutex.RLock()
	defer fw.mutex.RUnlock()

	if fw.fsnotify != nil {
		return nil
	}

	fsn, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("fsw: create fsnotify - %w", err)
	}

	fw.fsnotify = fsn
	for p, w := range fw.watchs {
		if w.Recursive {
			if err = fw.doRecursive(p, true); err != nil {
				return err
			}
		} else {
			if log := fw.Logger; log != nil {
				log.Debugf("fsw: %p watch file '%s'", fsn, p)
			}
			if err = fsn.Add(p); err != nil {
				return fmt.Errorf("fsw: Add('%s) - %w", p, err)
			}
		}
	}

	go fw.watch(fsn)

	return nil
}

// Stop stop file watching go-routine
func (fw *FileWatcher) Stop() error {
	fw.mutex.RLock()
	defer fw.mutex.RUnlock()

	return fw.stop()
}

// Close stop file watching go-routine and removes all watches
func (fw *FileWatcher) Close() error {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	clear(fw.watchs)

	return fw.stop()
}

func (fw *FileWatcher) stop() error {
	fsn := fw.fsnotify
	if fsn == nil {
		return nil
	}

	fw.fsnotify = nil
	if err := fsn.Close(); err != nil {
		return fmt.Errorf("fsw: close fsnotify - %w", err)
	}
	return nil
}

func (fw *FileWatcher) watch(fsn *fsnotify.Watcher) {
	timer := time.NewTimer(gog.If(fw.delay <= 0, time.Minute, fw.delay))

	defer func() {
		timer.Stop()
		clear(fw.events)

		if log := fw.Logger; log != nil {
			log.Infof("fsw: %p watching stopped", fsn)
		}
	}()

	if log := fw.Logger; log != nil {
		log.Infof("fsw: %p start watching", fsn)
	}

	for {
		select {
		case <-timer.C:
			if fw.delay > 0 {
				fw.delayCallbacks()
				timer.Reset(fw.delay)
			}
		case evt, ok := <-fsn.Events:
			if !ok {
				return
			}

			// check fsn to discard events after Stop()
			if fsn == fw.fsnotify {
				fw.doEvent(evt)
			}
		case err, ok := <-fsn.Errors:
			if !ok {
				return
			}
			if log := fw.Logger; log != nil {
				log.Errorf("fsw: %p watch error: %v", fsn, err)
			}
		}
	}
}

func (fw *FileWatcher) doEvent(evt fsnotify.Event) {
	fw.procEvent(evt)
	fw.sendEvent(evt)
}

// procEvent add or remove watch file
func (fw *FileWatcher) procEvent(evt fsnotify.Event) {
	if evt.Op&OpCreate != 0 {
		if err := fsu.DirExists(evt.Name); err == nil {
			if fw.isRecursive(evt.Name) {
				if err = fw.doRecursive(evt.Name, true); err != nil {
					if log := fw.Logger; log != nil {
						log.Errorf("fsw: watch dir '%s' error: %v", evt.Name, err)
					}
				}
			}
		}
	}
}

func (fw *FileWatcher) isRecursive(name string) bool {
	fw.mutex.RLock()
	defer fw.mutex.RUnlock()

	for p, w := range fw.watchs {
		if strings.HasPrefix(name, p) {
			if w.Recursive {
				return true
			}
		}
	}
	return false
}

func (fw *FileWatcher) sendEvent(evt fsnotify.Event) {
	if fw.delay <= 0 {
		cbs := fw.findCallbacks(evt)
		if len(cbs) > 0 {
			fw.execCallbacks(evt.Name, evt.Op, cbs)
		}
		return
	}

	if fe, ok := fw.events[evt.Name]; ok {
		fe.Time = time.Now()
		fe.Op |= evt.Op

		if log := fw.Logger; log != nil {
			log.Debugf("fsw: delay '%s' [%v] (%s) ", fe.Name, fe.Op, fe.Time.Format(time.RFC3339))
		}
		return
	}

	cbs := fw.findCallbacks(evt)
	if len(cbs) > 0 {
		fe := &fileevent{Name: evt.Name, Op: evt.Op, Time: time.Now(), Callbacks: cbs}
		fw.events[evt.Name] = fe

		if log := fw.Logger; log != nil {
			log.Debugf("fsw: queue '%s' [%v] (%s) ", fe.Name, fe.Op, fe.Time.Format(time.RFC3339))
		}
	}
}

func (fw *FileWatcher) findCallbacks(evt fsnotify.Event) (cbs []Callback) {
	fw.mutex.RLock()
	defer fw.mutex.RUnlock()

	for p, w := range fw.watchs {
		if evt.Op&w.OpMask != 0 && strings.HasPrefix(evt.Name, p) {
			cbs = append(cbs, w.Callback)
		}
	}
	return
}

func (fw *FileWatcher) delayCallbacks() {
	due := time.Now().Add(-fw.delay)
	for _, fe := range fw.events {
		if fe.Time.Before(due) {
			fw.execCallbacks(fe.Name, fe.Op, fe.Callbacks)
			delete(fw.events, fe.Name)
		}
	}
}

func (fw *FileWatcher) execCallbacks(name string, op Op, cbs []Callback) {
	if log := fw.Logger; log != nil {
		log.Debugf("fsw: execute callback '%s' [%v]", name, op)
	}

	for _, cb := range cbs {
		cb(name, op)
	}
}
