package fsw

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/log"
)

func prepareTestDir(testdir string) {
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 2; j++ {
			dir := filepath.Join(testdir, strconv.Itoa(i), strconv.Itoa(j))
			os.MkdirAll(dir, os.FileMode(0770))
			for k := 1; k <= 2; k++ {
				fn, _ := filepath.Abs(filepath.Join(dir, fmt.Sprintf("t%d.txt", k)))
				log.Info("Prepare ", fn)
				ioutil.WriteFile(fn, []byte("init"), os.FileMode(0660))
			}
		}
	}
}

func changeTestFiles(testdir string) {
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 2; j++ {
			dir := filepath.Join(testdir, strconv.Itoa(i), strconv.Itoa(j))
			os.MkdirAll(dir, os.FileMode(0770))
			for k := 1; k <= 2; k++ {
				time.Sleep(time.Millisecond * 500)
				fn, _ := filepath.Abs(filepath.Join(dir, fmt.Sprintf("t%d.txt", k)))
				log.Info("Change ", fn)
				ioutil.WriteFile(fn, []byte("test"), os.FileMode(0660))
			}
		}
	}
	time.Sleep(time.Second * 3)
}

func assertTestFiles(t *testing.T, testdir string, files map[string]Op) {
	if 3*2*2 != len(files) {
		t.Errorf("len(files) = %v, want %v", len(files), 3*2*2)
		return
	}

	for i := 1; i <= 3; i++ {
		for j := 1; j <= 2; j++ {
			dir := filepath.Join(testdir, strconv.Itoa(i), strconv.Itoa(j))
			for k := 1; k <= 2; k++ {
				fn := filepath.Join(dir, fmt.Sprintf("t%d.txt", k))
				_, ok := files[fn]
				if !ok {
					t.Errorf("file %q not exists", fn)
				}
			}
		}
	}
}

// func TestWatchDirOnly(t *testing.T) {
// 	//TODO: fix test on go1.17
// 	fmt.Println(runtime.GOOS, runtime.Version())
// 	if strings.HasPrefix(runtime.Version(), "go1.17") {
// 		return
// 	}
// 	if runtime.GOOS != "windows" && runtime.GOOS != "linux" {
// 		return
// 	}

// 	testdir := "TestWatchDirOnly-" + strconv.Itoa(rand.Int())

// 	os.RemoveAll(testdir)
// 	defer os.RemoveAll(testdir)

// 	os.MkdirAll(testdir, os.FileMode(0770))

// 	fw := NewFileWatcher()

// 	log.SetWriter(&log.StreamWriter{Color: true})
// 	fw.Logger = log.GetLogger("fswatch")

// 	files := make(map[string]Op)
// 	fw.Add(testdir, OpALL, func(path string, op Op) {
// 		files[path] = op
// 		fw.Logger.Infof("%q [%v]", path, op)
// 	})

// 	fw.Logger.Info("------------ Prepare ----------------")
// 	prepareTestDir(testdir)

// 	fw.Logger.Info("------------ Start ----------------")
// 	err := fw.Start()
// 	if err != nil {
// 		t.Errorf("fw.Start() = %v", err)
// 		return
// 	}

// 	changeTestFiles(testdir)

// 	fw.Logger.Info("------------ Stop ----------------")
// 	err = fw.Stop()
// 	if err != nil {
// 		t.Errorf("fw.Stop() = %v", err)
// 		return
// 	}

// 	if 3 != len(files) {
// 		t.Errorf("len(files) = %v, want %v", len(files), 3)
// 		return
// 	}
// 	for i := 1; i <= 3; i++ {
// 		dir := filepath.Join(testdir, strconv.Itoa(i))
// 		_, ok := files[dir]
// 		if !ok {
// 			t.Errorf("dir %v not exists", dir)
// 		}
// 	}
// }

func TestWatchRecursive(t *testing.T) {
	testdir := "TestWatchRecursive-" + strconv.Itoa(rand.Int())

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	os.MkdirAll(testdir, os.FileMode(0770))

	fw := NewFileWatcher()

	log.SetWriter(&log.StreamWriter{Color: true})
	fw.Logger = log.GetLogger("fswatch")

	files := make(map[string]Op)
	fw.AddRecursive(testdir, OpALL, "", func(path string, op Op) {
		if err := iox.FileExists(path); err == nil {
			files[path] = op
		}
		fw.Logger.Infof("%q [%v]", path, op)
	})

	fw.Logger.Info("------------ Prepare ----------------")
	prepareTestDir(testdir)

	fw.Logger.Info("------------ Start ----------------")
	err := fw.Start()
	if err != nil {
		t.Errorf(`fw.Start() = %v`, err)
		return
	}

	changeTestFiles(testdir)

	fw.Logger.Info("------------ Stop ----------------")
	err = fw.Stop()
	if err != nil {
		t.Errorf(`fw.Stop() = %v`, err)
		return
	}

	assertTestFiles(t, testdir, files)
}

func TestWatchAgain(t *testing.T) {
	testdir := "TestWatchAgain-" + strconv.Itoa(rand.Int())

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	os.MkdirAll(testdir, os.FileMode(0770))

	fw := NewFileWatcher()

	log.SetWriter(&log.StreamWriter{Color: true})
	fw.Logger = log.GetLogger("fswatch")

	files := make(map[string]Op)
	fw.AddRecursive(testdir, OpALL, "", func(path string, op Op) {
		if err := iox.FileExists(path); err == nil {
			files[path] = op
		}
		fw.Logger.Infof("%q [%v]", path, op)
	})

	fw.Logger.Info("------------ Prepare ----------------")
	prepareTestDir(testdir)

	fw.Logger.Info("------------ Start ----------------")
	err := fw.Start()
	if err != nil {
		t.Errorf(`fw.Start() = %v`, err)
		return
	}

	fw.Logger.Info("------------ Stop ----------------")
	err = fw.Stop()
	if err != nil {
		t.Errorf(`fw.Stop() = %v`, err)
		return
	}

	fw.Logger.Info("------------ Start Again ----------------")
	err = fw.Start()
	if err != nil {
		t.Errorf(`fw.Start() = %v`, err)
		return
	}

	changeTestFiles(testdir)

	fw.Logger.Info("------------ Stop Again ----------------")
	err = fw.Stop()
	if err != nil {
		t.Errorf(`fw.Stop() = %v`, err)
		return
	}

	assertTestFiles(t, testdir, files)
}
