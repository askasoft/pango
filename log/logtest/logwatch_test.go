package logtest

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/fsw"
	"github.com/askasoft/pango/log"
)

const (
	LOGCONF1 = `
{
	"async": ASYNC,
	"format": "text:%l %S %F() - %m%n%T",
	"level": {
		"*": "info",
		"sql": "debug",
		"http": "trace"
	},
	"writer": [{
		"_": "file",
		"path": "LOGFILE1",
		"maxDays": 7,
		"format": "%l - %m%n",
		"filter": "level:error"
	}]
}
`

	LOGCONF2 = `
{
	"async": ASYNC,
	"format": "text:%l %S %F() - %m%n%T",
	"level": {
		"*": "info",
		"sql": "debug",
		"http": "trace"
	},
	"writer": [{
		"_": "file",
		"path": "LOGFILE1",
		"maxDays": 7,
		"format": "%l - %m%n",
		"filter": "level:error"
	}, {
		"_": "file",
		"path": "LOGFILE2",
		"maxDays": 7,
		"format": "%l - %m%n",
		"filter": "level:warn name:test"
	}]
}
`
)

func TestLogConfigAsyncFile1toAsyncFile2(t *testing.T) {
	testLogConfigFile1toFile2(t, "1000", "1000")
}

func TestLogConfigAsyncFile1toSyncFile2(t *testing.T) {
	testLogConfigFile1toFile2(t, "1000", "0")
}

func TestLogConfigSyncFile1toAsyncFile2(t *testing.T) {
	testLogConfigFile1toFile2(t, "0", "1000")
}

func TestLogConfigSyncFile1toSyncFile2(t *testing.T) {
	testLogConfigFile1toFile2(t, "0", "0")
}

func testLogConfigFile1toFile2(t *testing.T, async1, async2 string) {
	testdir := "conftest-" + strconv.Itoa(rand.Int())

	os.RemoveAll(testdir)
	os.MkdirAll(testdir, os.FileMode(0777))
	defer os.RemoveAll(testdir)

	path := testdir + "/log.json"
	logfile1 := testdir + "/logs/file1.log"
	logfile2 := testdir + "/logs/file2.log"

	logconf1 := strings.ReplaceAll(LOGCONF1, "LOGFILE1", logfile1)
	logconf1 = strings.ReplaceAll(logconf1, "ASYNC", async1)

	logconf2 := strings.ReplaceAll(LOGCONF2, "LOGFILE1", logfile1)
	logconf2 = strings.ReplaceAll(logconf2, "LOGFILE2", logfile2)
	logconf2 = strings.ReplaceAll(logconf2, "ASYNC", async2)

	os.WriteFile(path, ([]byte)(logconf1), os.FileMode(0666))

	lg := log.NewLog()
	err := lg.Config(path)
	if err != nil {
		t.Fatalf("lg.Config(%q) = %v", path, err)
		return
	}

	fw := fsw.NewFileWatcher()
	fw.Start()
	defer fw.Stop()

	reconfiged := false
	err = fw.Add(path, fsw.OpWrite, func(path string, _ fsw.Op) {
		fmt.Println("Reload config file: ", path)
		err := lg.Config(path)
		if err != nil {
			t.Fatalf("Failed to config log by %q: %v\n", path, err)
		}
		reconfiged = true
	})
	if err != nil {
		t.Fatalf("fw.Add(%q) = %v", path, err)
		return
	}

	lg.Info("This is info.")
	lg.Warn("This is warn.")
	lg.Error("This is error.")
	lg.Flush()

	// Sleep for async flush
	time.Sleep(time.Second)

	bs, _ := os.ReadFile(logfile1)
	a := string(bs)
	w := "ERROR - This is error." + eol
	if a != w {
		t.Errorf(`%q = %v, want %v`, logfile1, a, w)
	}

	err = fsu.FileExists(logfile2)
	if err == nil {
		t.Errorf("%q not exists", logfile2)
	}

	// Sleep 1s for log watch
	time.Sleep(time.Second)

	fmt.Println("Change config file")

	err = os.WriteFile(path, ([]byte)(logconf2), os.FileMode(0666))
	if err != nil {
		fmt.Printf("Failed to change config %v\n", err)
		t.Fatalf("Failed to change config %v", err)
		return
	}

	// wait for file change event and log config reload
	for i := 0; i < 10000; i++ {
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

	bs, _ = os.ReadFile(logfile1)
	a = string(bs)
	w = "ERROR - This is error." + eol + "ERROR - This is error." + eol + "ERROR - This is ERROR." + eol
	if a != w {
		t.Errorf(`%q = %v, want %v`, logfile1, a, w)
	}

	bs, _ = os.ReadFile(logfile2)
	a = string(bs)
	w = "WARN - This is WARN." + eol + "ERROR - This is ERROR." + eol
	if a != w {
		t.Errorf(`%q = %v, want %v`, logfile2, a, w)
	}
}
