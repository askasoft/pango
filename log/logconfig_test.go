package log

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogConfigJSON(t *testing.T) {
	log := Default()
	assert.Nil(t, log.Config("testdata/log.json"))
	assertLogConfig(t, log)
}

func TestLogConfigINI(t *testing.T) {
	log := Default()
	assert.Nil(t, log.Config("testdata/log.ini"))
	assertLogConfig(t, log)
}

func assertLogConfig(t *testing.T, log *Log) {
	assert.Equal(t, LevelInfo, log.level)
	assert.Equal(t, 2, len(log.levels))
	assert.Equal(t, LevelDebug, log.levels["sql"])
	assert.Equal(t, LevelTrace, log.levels["http"])

	_, ok := log.GetFormatter().(*TextFormatter)
	assert.True(t, ok)

	assert.NotNil(t, log.writer)
	mw, ok := log.writer.(*MultiWriter)
	assert.True(t, ok)

	assert.Equal(t, 7, len(mw.Writers))

	i := 0
	{
		w, ok := mw.Writers[i].(*StreamWriter)
		assert.NotNil(t, w)
		assert.True(t, ok)
		assert.False(t, w.Color)

		f, ok := w.Logfil.(*MultiFilter)
		assert.NotNil(t, f)
		assert.True(t, ok)
		assert.Equal(t, 2, len(f.Filters))

		nf, ok := f.Filters[0].(*NameFilter)
		assert.NotNil(t, nf)
		assert.True(t, ok)
		assert.Equal(t, "out", nf.Name)

		lf, ok := f.Filters[1].(*LevelFilter)
		assert.NotNil(t, lf)
		assert.True(t, ok)
		assert.Equal(t, LevelDebug, lf.Level)
	}

	i++
	{
		w, ok := mw.Writers[i].(*StreamWriter)
		assert.NotNil(t, w)
		assert.True(t, ok)
	}

	i++
	{
		w, ok := mw.Writers[i].(*ConnWriter)
		assert.NotNil(t, w)
		assert.True(t, ok)
		assert.Equal(t, "tcp", w.Net)
		assert.Equal(t, "localhost:9999", w.Addr)
		assert.Equal(t, time.Second*5, w.Timeout)

		f, ok := w.Logfil.(*LevelFilter)
		assert.NotNil(t, f)
		assert.True(t, ok)
		assert.Equal(t, LevelError, f.Level)
	}

	i++
	{
		w, ok := mw.Writers[i].(*FileWriter)
		assert.NotNil(t, w)
		assert.True(t, ok)
		assert.Equal(t, uint32(0777), w.DirPerm)
		assert.True(t, w.Daily)
		assert.Equal(t, 7, w.MaxDays)
		assert.Equal(t, LevelError, w.FlushLevel)

		f, ok := w.Logfil.(*LevelFilter)
		assert.NotNil(t, f)
		assert.True(t, ok)
		assert.Equal(t, LevelError, f.Level)
	}

	i++
	{
		w, ok := mw.Writers[i].(*SlackWriter)
		assert.NotNil(t, w)
		assert.True(t, ok)
		assert.Equal(t, "develop", w.Channel)
		assert.Equal(t, "gotest", w.Username)
		assert.Equal(t, "https://hooks.slack.com/services/...", w.Webhook)
		assert.Equal(t, time.Second*5, w.Timeout)

		f, ok := w.Logfil.(*LevelFilter)
		assert.NotNil(t, f)
		assert.True(t, ok)
		assert.Equal(t, LevelError, f.Level)
	}

	i++
	{
		w, ok := mw.Writers[i].(*SMTPWriter)
		assert.NotNil(t, w)
		assert.True(t, ok)
		assert.Equal(t, "localhost", w.Host)
		assert.Equal(t, 25, w.Port)
		assert.Equal(t, "-----", w.Username)
		assert.Equal(t, "xxxxxxx", w.Password)
		assert.Equal(t, "pango@google.com", w.From)
		assert.Equal(t, "to1@test.com to2@test.com", strings.Join(w.Tos, " "))
		assert.Equal(t, "cc1@test.com cc2@test.com", strings.Join(w.Ccs, " "))
		assert.Equal(t, time.Second*5, w.Timeout)

		f, ok := w.Logfil.(*LevelFilter)
		assert.NotNil(t, f)
		assert.True(t, ok)
		assert.Equal(t, LevelError, f.Level)
	}

	i++
	{
		w, ok := mw.Writers[i].(*WebhookWriter)
		assert.True(t, ok)
		assert.Equal(t, "http://localhost:9200/pango/logs", w.Webhook)
		assert.Equal(t, "application/json", w.ContentType)
		assert.Equal(t, time.Second*5, w.Timeout)

		o, ok := w.Logfmt.(*JSONFormatter)
		assert.NotNil(t, o)
		assert.True(t, ok)

		f, ok := w.Logfil.(*LevelFilter)
		assert.True(t, ok)
		assert.Equal(t, LevelError, f.Level)
	}
}

func TestLogConfigFile1(t *testing.T) {
	os.RemoveAll("conftest")
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
	os.RemoveAll("conftest")
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
