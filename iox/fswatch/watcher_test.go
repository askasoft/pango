package fswatch

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/log"
	"github.com/stretchr/testify/assert"
)

func prepareTestDir() {
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 2; j++ {
			dir := filepath.Join("testdir", strconv.Itoa(i), strconv.Itoa(j))
			os.MkdirAll(dir, os.FileMode(0770))
			for k := 1; k <= 2; k++ {
				fn := filepath.Join(dir, fmt.Sprintf("t%d.txt", k))
				ioutil.WriteFile(fn, []byte("init"), os.FileMode(0660))
			}
		}
	}
}

func changeTestFiles() {
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 2; j++ {
			dir := filepath.Join("testdir", strconv.Itoa(i), strconv.Itoa(j))
			os.MkdirAll(dir, os.FileMode(0770))
			for k := 1; k <= 2; k++ {
				time.Sleep(time.Second * 1)
				fn := filepath.Join(dir, fmt.Sprintf("t%d.txt", k))
				ioutil.WriteFile(fn, []byte("test"), os.FileMode(0660))
			}
		}
	}
	time.Sleep(time.Second * 3)
}

func assertTestFiles(t *testing.T, files map[string]Op) {
	assert.Equal(t, 3*2*2, len(files))

	for i := 1; i <= 3; i++ {
		for j := 1; j <= 2; j++ {
			dir := filepath.Join("testdir", strconv.Itoa(i), strconv.Itoa(j))
			for k := 1; k <= 2; k++ {
				fn := filepath.Join(dir, fmt.Sprintf("t%d.txt", k))
				_, ok := files[fn]
				assert.True(t, ok, fn)
			}
		}
	}
}

func TestWatchDirOnly(t *testing.T) {
	os.RemoveAll("testdir")
	defer os.RemoveAll("testdir")

	os.MkdirAll("testdir", os.FileMode(0770))

	fw := NewFileWatcher()

	log.SetWriter(&log.StreamWriter{Color: true})
	fw.Logger = log.GetLogger("fswatch")

	files := make(map[string]Op)
	fw.Add("testdir", OpALL, func(path string, op Op) {
		files[path] = op
		fw.Logger.Infof("%q [%v]", path, op)
	})

	fw.Logger.Info("------------ Prepare ----------------")
	prepareTestDir()

	fw.Logger.Info("------------ Start ----------------")
	err := fw.Start()
	assert.Nil(t, err)

	changeTestFiles()

	fw.Logger.Info("------------ Stop ----------------")
	err = fw.Stop()
	assert.Nil(t, err)

	assert.Equal(t, 3, len(files))
	for i := 1; i <= 3; i++ {
		dir := filepath.Join("testdir", strconv.Itoa(i))
		_, ok := files[dir]
		assert.True(t, ok, dir)
	}
}

func TestWatchRecursive(t *testing.T) {
	os.RemoveAll("testdir")
	defer os.RemoveAll("testdir")

	os.MkdirAll("testdir", os.FileMode(0770))

	fw := NewFileWatcher()

	log.SetWriter(&log.StreamWriter{Color: true})
	fw.Logger = log.GetLogger("fswatch")

	files := make(map[string]Op)
	fw.AddRecursive("testdir", OpALL, "", func(path string, op Op) {
		if err := iox.FileExists(path); err == nil {
			files[path] = op
		}
		fw.Logger.Infof("%q [%v]", path, op)
	})

	fw.Logger.Info("------------ Prepare ----------------")
	prepareTestDir()

	fw.Logger.Info("------------ Start ----------------")
	err := fw.Start()
	assert.Nil(t, err)

	changeTestFiles()

	fw.Logger.Info("------------ Stop ----------------")
	err = fw.Stop()
	assert.Nil(t, err)

	assertTestFiles(t, files)
}

func TestWatchAgain(t *testing.T) {
	os.RemoveAll("testdir")
	defer os.RemoveAll("testdir")

	os.MkdirAll("testdir", os.FileMode(0770))

	fw := NewFileWatcher()

	log.SetWriter(&log.StreamWriter{Color: true})
	fw.Logger = log.GetLogger("fswatch")

	files := make(map[string]Op)
	fw.AddRecursive("testdir", OpALL, "", func(path string, op Op) {
		if err := iox.FileExists(path); err == nil {
			files[path] = op
		}
		fw.Logger.Infof("%q [%v]", path, op)
	})

	fw.Logger.Info("------------ Prepare ----------------")
	prepareTestDir()

	fw.Logger.Info("------------ Start ----------------")
	err := fw.Start()
	assert.Nil(t, err)

	fw.Logger.Info("------------ Stop ----------------")
	err = fw.Stop()
	assert.Nil(t, err)

	fw.Logger.Info("------------ Start Again ----------------")
	err = fw.Start()
	assert.Nil(t, err)

	changeTestFiles()

	fw.Logger.Info("------------ Stop Again ----------------")
	err = fw.Stop()
	assert.Nil(t, err)

	assertTestFiles(t, files)
}
