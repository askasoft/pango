package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/pandafw/pango/iox"
	"github.com/stretchr/testify/assert"
)

func TestLogConfig(t *testing.T) {
	log := Default()
	assert.Nil(t, log.Config("testdata/log.json"))
}

func TestLogConfigFile1(t *testing.T) {
	defer os.RemoveAll("conftest")

	log := Default()
	assert.Nil(t, log.Config("testdata/log-file1.json"))
	log.Info("This is info.")
	log.Warn("This is warn.")
	log.Error("This is error.")
	log.Close()

	bs, _ := ioutil.ReadFile("conftest/logs/file1.log")
	assert.Equal(t, "ERROR - This is error."+eol, string(bs))
}

func TestLogConfigFile2(t *testing.T) {
	defer os.RemoveAll("conftest")

	log := Default()
	assert.Nil(t, log.Config("testdata/log-file2.json"))
	log.Info("This is info.")
	log.Warn("This is warn.")
	log.Error("This is error.")

	tl := log.GetLogger("test")
	tl.Warn("This is WARN.")
	tl.Error("This is ERROR.")
	log.Close()

	bs, _ := ioutil.ReadFile("conftest/logs/file1.log")
	assert.Equal(t, "ERROR - This is error."+eol+"ERROR - This is ERROR."+eol, string(bs))

	bs, _ = ioutil.ReadFile("conftest/logs/file2.log")
	assert.Equal(t, "WARN - This is WARN."+eol+"ERROR - This is ERROR."+eol, string(bs))
}

func TestLogConfigFile1toFile2(t *testing.T) {
	path := "conftest/log.json"

	iox.CopyFile("testdata/log-file1.json", path)
	log := Default()
	assert.Nil(t, log.Config(path))

	assert.Nil(t, log.Watch(path))

	log.Info("This is info.")
	log.Warn("This is warn.")
	log.Error("This is error.")
	log.Flush()

	bs, _ := ioutil.ReadFile("conftest/logs/file1.log")
	assert.Equal(t, "ERROR - This is error."+eol, string(bs))

	assert.NotNil(t, iox.FileExists("conftest/logs/file2.log"))

	// Sleep 1s for log watch
	time.Sleep(time.Second * 1)
	fmt.Println("Change config file")
	err := iox.CopyFile("testdata/log-file2.json", path)
	if err != nil {
		fmt.Printf("Failed to change config %v\n", err)
		assert.Fail(t, "Failed to change config %v", err)
		return
	}

	// wait for file change event and log config reload
	for i := 0; i < 30; i++ {
		_, ok := log.writer.(*MultiWriter)
		if ok {
			break
		}
		fmt.Println(strconv.Itoa(i) + " - Sleep 1s for log config reload")
		time.Sleep(time.Second * 1)
	}

	log.Info("This is info.")
	log.Warn("This is warn.")
	log.Error("This is error.")

	tl := log.GetLogger("test")
	tl.Warn("This is WARN.")
	tl.Error("This is ERROR.")
	log.Close()

	bs, _ = ioutil.ReadFile("conftest/logs/file1.log")
	if !assert.Equal(t, "ERROR - This is error."+eol+"ERROR - This is error."+eol+"ERROR - This is ERROR."+eol, string(bs)) {
		return
	}

	bs, _ = ioutil.ReadFile("conftest/logs/file2.log")
	if !assert.Equal(t, "WARN - This is WARN."+eol+"ERROR - This is ERROR."+eol, string(bs)) {
		return
	}

	os.RemoveAll("conftest")
}
