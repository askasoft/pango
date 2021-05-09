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
	"github.com/stretchr/testify/assert"
)

func TestLogConfigFile1toFile2(t *testing.T) {
	os.RemoveAll("conftest")
	defer os.RemoveAll("conftest")

	path := "conftest/log.json"

	iox.CopyFile("../testdata/log-file1.json", path)
	lg := log.NewLog()
	assert.Nil(t, lg.Config(path))

	fw := fswatch.NewFileWatcher()
	fw.Start()
	defer fw.Stop()

	assert.Nil(t, fw.Add(path, fswatch.OpWrite, func(path string, _ fswatch.Op) {
		err := lg.Config(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to config log by %q: %v\n", path, err)
		}
	}))

	lg.Info("This is info.")
	lg.Warn("This is warn.")
	lg.Error("This is error.")
	lg.Flush()

	bs, _ := ioutil.ReadFile("conftest/logs/file1.log")
	assert.Equal(t, "ERROR - This is error."+iox.EOL, string(bs))

	assert.NotNil(t, iox.FileExists("conftest/logs/file2.log"))

	// Sleep 1s for log watch
	time.Sleep(time.Second * 1)
	fmt.Println("Change config file")
	err := iox.CopyFile("../testdata/log-file2.json", path)
	if err != nil {
		fmt.Printf("Failed to change config %v\n", err)
		assert.Fail(t, "Failed to change config %v", err)
		return
	}

	// wait for file change event and log config reload
	for i := 0; i < 10; i++ {
		_, ok := lg.GetWriter().(*log.MultiWriter)
		if ok {
			break
		}
		fmt.Println(strconv.Itoa(i) + " - Sleep 1s for log config reload")
		time.Sleep(time.Second * 1)
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
	if !assert.Equal(t, "ERROR - This is error."+iox.EOL+"ERROR - This is error."+iox.EOL+"ERROR - This is ERROR."+iox.EOL, string(bs)) {
		return
	}

	bs, _ = ioutil.ReadFile("conftest/logs/file2.log")
	if !assert.Equal(t, "WARN - This is WARN."+iox.EOL+"ERROR - This is ERROR."+iox.EOL, string(bs)) {
		return
	}
}
