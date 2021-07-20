package logwatch

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/iox/fswatch"
	"github.com/pandafw/pango/log"
)

func TestLogConfigFile1toFile2(t *testing.T) {
	os.RemoveAll("conftest")
	defer os.RemoveAll("conftest")

	path := "conftest/log.json"

	iox.CopyFile("../testdata/log-file1.json", path)
	lg := log.NewLog()
	err := lg.Config(path)
	if err != nil {
		t.Errorf("lg.Config(%q) = %v", path, err)
	}

	fw := fswatch.NewFileWatcher()
	fw.Start()
	defer fw.Stop()

	reconfiged := false
	err = fw.Add(path, fswatch.OpWrite, func(path string, _ fswatch.Op) {
		err := lg.Config(path)
		if err != nil {
			t.Fatalf("Failed to config log by %q: %v\n", path, err)
		}
		reconfiged = true
	})
	if err != nil {
		t.Errorf("fw.Add(%q) = %v", path, err)
	}

	lg.Info("This is info.")
	lg.Warn("This is warn.")
	lg.Error("This is error.")
	lg.Flush()

	// Sleep for async flush
	time.Sleep(time.Second)

	bs, _ := ioutil.ReadFile("conftest/logs/file1.log")
	a := string(bs)
	w := "ERROR - This is error." + iox.EOL
	if a != w {
		t.Errorf(`%q = %v, want %v`, "conftest/logs/file1.log", a, w)
	}

	err = iox.FileExists("conftest/logs/file2.log")
	if err == nil {
		t.Errorf("%q not exists", "conftest/logs/file2.log")
	}

	// Sleep 1s for log watch
	time.Sleep(time.Second * 1)
	fmt.Println("Change config file")
	err = iox.CopyFile("../testdata/log-file2.json", path)
	if err != nil {
		fmt.Printf("Failed to change config %v\n", err)
		t.Fatalf("Failed to change config %v", err)
		return
	}

	// wait for file change event and log config reload
	for i := 0; i < 50; i++ {
		if reconfiged {
			break
		}
		if i%10 == 0 {
			fmt.Println(strconv.Itoa(i) + " - Sleep 1s for log config reload")
		}
		time.Sleep(time.Millisecond * 100)
	}

	lg.Info("This is info.")
	lg.Warn("This is warn.")
	lg.Error("This is error.")

	tl := lg.GetLogger("test")
	tl.Warn("This is WARN.")
	tl.Error("This is ERROR.")

	// Close log
	lg.Close()

	bs, _ = ioutil.ReadFile("conftest/logs/file1.log")
	a = string(bs)
	w = "ERROR - This is error." + iox.EOL + "ERROR - This is error." + iox.EOL + "ERROR - This is ERROR." + iox.EOL
	if a != w {
		t.Errorf(`%q = %v, want %v`, "conftest/logs/file1.log", a, w)
	}

	bs, _ = ioutil.ReadFile("conftest/logs/file2.log")
	a = string(bs)
	w = "WARN - This is WARN." + iox.EOL + "ERROR - This is ERROR." + iox.EOL
	if a != w {
		t.Errorf(`%q = %v, want %v`, "conftest/logs/file2.log", a, w)
	}
}
