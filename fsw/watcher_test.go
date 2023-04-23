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

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/log"
)

func testSleep() {
	time.Sleep(time.Second * 3)
}

func prepareTestDir(log log.Logger, testdir string) {
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
	testSleep()
}

func changeTestFiles(log log.Logger, testdir string) {
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 2; j++ {
			dir := filepath.Join(testdir, strconv.Itoa(i), strconv.Itoa(j))
			os.MkdirAll(dir, os.FileMode(0770))
			for k := 1; k <= 2; k++ {
				fn, _ := filepath.Abs(filepath.Join(dir, fmt.Sprintf("t%d.txt", k)))

				time.Sleep(time.Millisecond * 200)
				log.Info("Change 1 - ", fn)
				ioutil.WriteFile(fn, []byte("test1"), os.FileMode(0660))

				time.Sleep(time.Millisecond * 200)
				log.Info("Change 2 - ", fn)
				ioutil.WriteFile(fn, []byte("test2"), os.FileMode(0660))
			}
		}
	}
	testSleep()
}

func assertTestFiles(t *testing.T, testdir string, files map[string]int) {
	if 3*2*2 != len(files) {
		t.Errorf("len(files) = %v, want %v", len(files), 3*2*2)
		return
	}

	for i := 1; i <= 3; i++ {
		for j := 1; j <= 2; j++ {
			dir := filepath.Join(testdir, strconv.Itoa(i), strconv.Itoa(j))
			for k := 1; k <= 2; k++ {
				fn := filepath.Join(dir, fmt.Sprintf("t%d.txt", k))
				c, ok := files[fn]
				if !ok {
					t.Errorf("file %q not exists", fn)
				}
				if c != 1 {
					t.Errorf("file %q event count = %d, want 1", fn, c)
				}
			}
		}
	}
}

func testCreateWatcher() (*FileWatcher, *log.Log) {
	lg := log.NewLog()
	sw := &log.StreamWriter{Color: true}
	aw := log.NewAsyncWriter(sw, 100)
	log.SetWriter(aw)

	fw := NewFileWatcher()
	fw.Logger = lg.GetLogger("FSW")
	return fw, lg
}

func TestWatchRecursive(t *testing.T) {
	testdir := "TestWatchRecursive-" + strconv.Itoa(rand.Int())

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	os.MkdirAll(testdir, os.FileMode(0770))

	fw, lg := testCreateWatcher()

	files := make(map[string]int)
	fw.AddRecursive(testdir, OpModifies, func(path string, op Op) {
		if err := fsu.FileExists(path); err == nil {
			files[path]++
		}
		fw.Logger.Infof("%q [%v]", path, op)
	})

	fw.Logger.Info("------------ Prepare ----------------")
	prepareTestDir(fw.Logger, testdir)

	fw.Logger.Info("------------ Start ----------------")
	err := fw.Start()
	if err != nil {
		t.Errorf(`fw.Start() = %v`, err)
		return
	}

	changeTestFiles(fw.Logger, testdir)

	fw.Logger.Info("------------ Stop ----------------")
	err = fw.Stop()
	if err != nil {
		t.Errorf(`fw.Stop() = %v`, err)
		return
	}

	assertTestFiles(t, testdir, files)

	lg.Close()
}

func TestWatchAgain(t *testing.T) {
	testdir := "TestWatchAgain-" + strconv.Itoa(rand.Int())

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	os.MkdirAll(testdir, os.FileMode(0770))

	fw, lg := testCreateWatcher()

	files := make(map[string]int)
	fw.AddRecursive(testdir, OpModifies, func(path string, op Op) {
		if err := fsu.FileExists(path); err == nil {
			files[path]++
		}
		fw.Logger.Infof("%q [%v]", path, op)
	})

	fw.Logger.Info("------------ Prepare ----------------")
	prepareTestDir(fw.Logger, testdir)

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

	// wait for stop()
	testSleep()

	fw.Logger.Info("------------ Start Again ----------------")
	err = fw.Start()
	if err != nil {
		t.Errorf(`fw.Start() = %v`, err)
		return
	}

	changeTestFiles(fw.Logger, testdir)

	fw.Logger.Info("------------ Stop Again ----------------")
	err = fw.Stop()
	if err != nil {
		t.Errorf(`fw.Stop() = %v`, err)
		return
	}

	assertTestFiles(t, testdir, files)

	lg.Close()
}
