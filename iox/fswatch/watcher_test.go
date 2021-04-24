package fswatch

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pandafw/pango/log"
	"github.com/stretchr/testify/assert"
)

func xTestNofityFolder(t *testing.T) {
	w, err := fsnotify.NewWatcher()
	assert.Nil(t, err)

	p, _ := filepath.Abs("testdir")
	w.Add(p)

	go func() {
		defer fmt.Printf("%v %v\n", time.Now(), "End")
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}
				fmt.Printf("%v %v\n", time.Now(), event)
			}
		}
	}()

	time.Sleep(time.Second * 5)
	w.Close()
	time.Sleep(time.Second * 5)
}

func TestWatchFolder(t *testing.T) {
	fw, err := NewFileWatcher()
	assert.Nil(t, err)

	log.SetWriter(&log.StreamWriter{Color: true})
	fw.Logger = log.GetLogger("fswatch")

	fw.AddRecursive("testdir", OpALL, "", func(path string, op Op) {
		fw.Logger.Infof("%q [%v]", path, op)
	})

	fw.StartWatch()

	time.Sleep(time.Minute * 1)

	fw.StopWatch()
}
