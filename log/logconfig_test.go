package log

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/num"
)

func TestLogConfigJSON(t *testing.T) {
	log := NewLog()
	if err := log.Config("testdata/log.json"); err != nil {
		t.Fatalf(`log.Config("testdata/log.json") = %v`, err)
	}
	assertLogConfig(t, log)
}

func TestLogConfigINI(t *testing.T) {
	log := NewLog()
	if err := log.Config("testdata/log.ini"); err != nil {
		t.Fatalf(`log.Config("testdata/log.ini") = %v`, err)
	}
	assertLogConfig(t, log)
}

func assertLogEqual(t *testing.T, msg string, want any, val any) {
	if want != val {
		t.Errorf("%s: actual = %v, want %v", msg, val, want)
	}
}

func assertLogConfig(t *testing.T, log *Log) {
	assertLogEqual(t, `log.GetLevel()`, LevelInfo, log.GetLevel())
	assertLogEqual(t, `len(log.levels)`, 2, len(log.levels))
	assertLogEqual(t, `log.levels["sql"]`, LevelDebug, log.levels["sql"])
	assertLogEqual(t, `log.levels["http"]`, LevelTrace, log.levels["http"])

	lgsql := log.GetLogger("sql")
	assertLogEqual(t, `lgsql.GetLevel()`, LevelDebug, lgsql.GetLevel())

	lghttp := log.GetLogger("http")
	assertLogEqual(t, `lghttp.GetLevel()`, LevelTrace, lghttp.GetLevel())

	if _, ok := log.GetFormatter().(*TextFormatter); !ok {
		t.Fatalf("Not TextFormatter")
	}

	aw, ok := log.writer.(*AsyncWriter)
	if !ok {
		t.Fatalf("Not AsyncWriter")
	}

	mw, ok := aw.writer.(*MultiWriter)
	if !ok {
		t.Fatalf("Not MultiWriter")
	}

	assertLogEqual(t, `len(mw.Writers)`, 8, len(mw.Writers))

	i := 0
	{
		w, ok := mw.Writers[i].(*StreamWriter)
		if !ok {
			t.Fatalf("Not StreamWriter")
		}
		assertLogEqual(t, `w.Color`, false, w.Color)

		f, ok := w.Logfil.(*MultiFilter)
		if !ok {
			t.Fatalf("Not MultiFilter")
		}
		assertLogEqual(t, `len(f.Filters)`, 2, len(f.Filters))

		nf, ok := f.Filters[0].(*NameFilter)
		if !ok {
			t.Fatalf("Not NameFilter")
		}
		assertLogEqual(t, `nf.Name`, "out", nf.Name)

		lf, ok := f.Filters[1].(*LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `lf.Level`, LevelTrace, lf.Level)
	}

	i++
	{
		_, ok := mw.Writers[i].(*StreamWriter)
		if !ok {
			t.Fatalf("Not StreamWriter")
		}
	}

	i++
	{
		aw, ok := mw.Writers[i].(*AsyncWriter)
		if !ok {
			t.Fatalf("Not AsyncWriter")
		}

		w, ok := aw.writer.(*ConnWriter)
		if !ok {
			t.Fatalf("Not ConnWriter")
		}
		assertLogEqual(t, `w.Net`, "tcp", w.Net)
		assertLogEqual(t, `w.Addr`, "localhost:9999", w.Addr)
		assertLogEqual(t, `w.Timeout`, time.Second*5, w.Timeout)

		f, ok := w.Logfil.(*LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `f.Level`, LevelDebug, f.Level)
	}

	i++
	{
		aw, ok := mw.Writers[i].(*AsyncWriter)
		if !ok {
			t.Fatalf("Not AsyncWriter")
		}

		w, ok := aw.writer.(*FileWriter)
		if !ok {
			t.Fatalf("Not FileWriter")
		}
		assertLogEqual(t, `w.DirPerm`, uint32(0777), w.DirPerm)
		assertLogEqual(t, `w.MaxDays`, 7, w.MaxDays)
		assertLogEqual(t, `w.MaxSize`, int64(num.MB*4), w.MaxSize)
		assertLogEqual(t, `w.SyncLevel`, LevelError, w.SyncLevel)

		f, ok := w.Logfil.(*LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `f.Level`, LevelInfo, f.Level)
	}

	i++
	{
		aw, ok := mw.Writers[i].(*AsyncWriter)
		if !ok {
			t.Fatalf("Not AsyncWriter")
		}

		w, ok := aw.writer.(*SlackWriter)
		if !ok {
			t.Fatalf("Not SlackWriter")
		}
		assertLogEqual(t, `w.Webhook`, "https://hooks.slack.com/services/...", w.Webhook)
		assertLogEqual(t, `w.Timeout`, time.Second*5, w.Timeout)

		f, ok := w.Logfil.(*LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `f.Level`, LevelWarn, f.Level)
	}

	i++
	{
		aw, ok := mw.Writers[i].(*AsyncWriter)
		if !ok {
			t.Fatalf("Not AsyncWriter")
		}

		fw, ok := aw.writer.(*FailoverWriter)
		if !ok {
			t.Fatalf("Not FailoverWriter")
		}

		w, ok := fw.writer.(*SMTPWriter)
		if !ok {
			t.Fatalf("Not SMTPWriter")
		}
		assertLogEqual(t, `w.Host`, "localhost", w.Host)
		assertLogEqual(t, `w.Port`, 25, w.Port)
		assertLogEqual(t, `w.Username`, "-----", w.Username)
		assertLogEqual(t, `w.Password`, "xxxxxxx", w.Password)
		assertLogEqual(t, `w.From`, "pango@google.com", w.From)
		assertLogEqual(t, `w.Tos`, "to1@test.com to2@test.com", strings.Join(w.Tos, " "))
		assertLogEqual(t, `w.Ccs`, "cc1@test.com cc2@test.com", strings.Join(w.Ccs, " "))
		assertLogEqual(t, `w.Timeout`, time.Second*5, w.Timeout)

		f, ok := w.Logfil.(*LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `f.Level`, LevelError, f.Level)
	}

	i++
	{
		aw, ok := mw.Writers[i].(*AsyncWriter)
		if !ok {
			t.Fatalf("Not AsyncWriter")
		}

		w, ok := aw.writer.(*TeamsWriter)
		if !ok {
			t.Fatalf("Not TeamsWriter")
		}
		assertLogEqual(t, `w.Webhook`, "https://xxx.webhook.office.com/webhookb2/...", w.Webhook)
		assertLogEqual(t, `w.Timeout`, time.Second*3, w.Timeout)

		f, ok := w.Logfil.(*LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `f.Level`, LevelFatal, f.Level)
	}

	i++
	{
		aw, ok := mw.Writers[i].(*AsyncWriter)
		if !ok {
			t.Fatalf("Not AsyncWriter")
		}

		w, ok := aw.writer.(*WebhookWriter)
		if !ok {
			t.Fatalf("Not WebhookWriter")
		}
		assertLogEqual(t, `w.Webhook`, "http://localhost:9200/pango/logs", w.Webhook)
		assertLogEqual(t, `w.ContentType`, "application/json", w.ContentType)
		assertLogEqual(t, `w.Timeout`, time.Second*5, w.Timeout)

		jf, ok := w.Logfmt.(*JSONFormatter)
		if jf == nil || !ok {
			t.Fatalf("Not JSONFormatter")
		}

		f, ok := w.Logfil.(*LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `f.Level`, LevelFatal, f.Level)
	}
}

func TestLogConfigFile1(t *testing.T) {
	os.RemoveAll("conftest")
	defer os.RemoveAll("conftest")

	log := NewLog()
	if err := log.Config("testdata/log-file1.json"); err != nil {
		t.Fatalf(`log.Config("testdata/log-file1.json") = %v`, err)
	}
	log.Info("This is info.")
	log.Warn("This is warn.")
	log.Error("This is error.")
	log.Close()

	bs, _ := os.ReadFile("conftest/logs/file1.log")
	a := string(bs)
	w := "ERROR - This is error." + EOL
	if a != w {
		t.Errorf("\n actual = %v\n expect = %v", a, w)
	}
}

func TestLogConfigFile2(t *testing.T) {
	os.RemoveAll("conftest")
	defer os.RemoveAll("conftest")

	log := NewLog()
	if err := log.Config("testdata/log-file2.json"); err != nil {
		t.Fatalf(`log.Config("testdata/log-file2.json") = %v`, err)
	}
	log.Info("This is info.")
	log.Warn("This is warn.")
	log.Error("This is error.")

	tl := log.GetLogger("test")
	tl.Warn("This is WARN.")
	tl.Error("This is ERROR.")
	log.Close()

	bs, _ := os.ReadFile("conftest/logs/file1.log")
	a := string(bs)
	w := "ERROR - This is error." + EOL + "ERROR - This is ERROR." + EOL
	if a != w {
		t.Errorf("\n actual = %v\n expect = %v", a, w)
	}

	bs, _ = os.ReadFile("conftest/logs/file2.log")
	a = string(bs)
	w = "WARN - This is WARN." + EOL + "ERROR - This is ERROR." + EOL
	if a != w {
		t.Errorf("\n actual = %v\n expect = %v", a, w)
	}
}

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

	lg := NewLog()
	err := lg.Config(path)
	if err != nil {
		t.Fatalf("lg.Config(%q) = %v", path, err)
		return
	}

	lg.Info("This is info.")
	lg.Warn("This is warn.")
	lg.Error("This is error.")
	lg.Flush()

	// Sleep for async flush
	time.Sleep(time.Second)

	os.WriteFile(path, ([]byte)(logconf1), os.FileMode(0666))

	// Sleep for file sync
	time.Sleep(time.Second)

	bs, _ := os.ReadFile(logfile1)
	a := string(bs)
	w := "ERROR - This is error." + EOL
	if a != w {
		t.Errorf(`%q = %v, want %v`, logfile1, a, w)
	}

	err = fsu.FileExists(logfile2)
	if err == nil {
		t.Errorf("%q should not exists", logfile2)
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

	fmt.Println("Reload config file: ", path)
	err = lg.Config(path)
	if err != nil {
		t.Fatalf("Failed to config log by %q: %v\n", path, err)
		return
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
	w = "ERROR - This is error." + EOL + "ERROR - This is error." + EOL + "ERROR - This is ERROR." + EOL
	if a != w {
		t.Errorf(`%q = %v, want %v`, logfile1, a, w)
	}

	bs, _ = os.ReadFile(logfile2)
	a = string(bs)
	w = "WARN - This is WARN." + EOL + "ERROR - This is ERROR." + EOL
	if a != w {
		t.Errorf(`%q = %v, want %v`, logfile2, a, w)
	}
}
