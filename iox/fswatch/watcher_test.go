package fswatch

import (
	"fmt"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
)

func xTestWatchFolder(t *testing.T) {
	w, err := fsnotify.NewWatcher()
	assert.Nil(t, err)

	w.Add("testdata")

	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}
			fmt.Printf("%v\n", event)
		}
	}
}
