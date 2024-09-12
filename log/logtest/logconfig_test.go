package logtest

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/log/httplog"
	"github.com/askasoft/pango/log/slacklog"
	"github.com/askasoft/pango/log/smtplog"
	"github.com/askasoft/pango/log/teamslog"
	"github.com/askasoft/pango/num"
)

var eol = log.EOL

func TestLogConfigJSON(t *testing.T) {
	log := log.NewLog()
	if err := log.Config("testdata/log.json"); err != nil {
		t.Fatalf(`log.Config("testdata/log.json") = %v`, err)
	}
	assertLogConfig(t, log)
}

func TestLogConfigINI(t *testing.T) {
	log := log.NewLog()
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

func testGetLogLevels(i any) map[string]log.Level {
	v := reflect.ValueOf(i).Elem()
	f := v.FieldByName("levels")
	p := (*map[string]log.Level)(unsafe.Pointer(f.UnsafeAddr()))
	return *p
}

func testGetLogWriter(i any) log.Writer {
	v := reflect.ValueOf(i).Elem()
	f := v.FieldByName("writer")
	p := (*log.Writer)(unsafe.Pointer(f.UnsafeAddr()))
	return *p
}

func assertLogConfig(t *testing.T, lg *log.Log) {
	assertLogEqual(t, `log.GetLevel()`, log.LevelInfo, lg.GetLevel())

	levels := testGetLogLevels(lg)
	assertLogEqual(t, `len(log.levels)`, 2, len(levels))
	assertLogEqual(t, `log.levels["sql"]`, log.LevelDebug, levels["sql"])
	assertLogEqual(t, `log.levels["http"]`, log.LevelTrace, levels["http"])

	lgsql := lg.GetLogger("sql")
	assertLogEqual(t, `lgsql.GetLevel()`, log.LevelDebug, lgsql.GetLevel())

	lghttp := lg.GetLogger("http")
	assertLogEqual(t, `lghttp.GetLevel()`, log.LevelTrace, lghttp.GetLevel())

	if _, ok := lg.GetFormatter().(*log.TextFormatter); !ok {
		t.Fatalf("Not TextFormatter")
	}

	writer := testGetLogWriter(lg)
	aw, ok := writer.(*log.AsyncWriter)
	if !ok {
		t.Fatalf("Not AsyncWriter")
	}

	writer = testGetLogWriter(aw)
	mw, ok := writer.(*log.MultiWriter)
	if !ok {
		t.Fatalf("Not MultiWriter")
	}

	assertLogEqual(t, `len(mw.Writers)`, 9, len(mw.Writers))

	i := 0
	{
		w, ok := mw.Writers[i].(*log.StreamWriter)
		if !ok {
			t.Fatalf("Not StreamWriter")
		}
		assertLogEqual(t, `w.Color`, false, w.Color)

		f, ok := w.Filter.(*log.MultiFilter)
		if !ok {
			t.Fatalf("Not MultiFilter")
		}
		assertLogEqual(t, `len(f.Filters)`, 2, len(f.Filters))

		nf, ok := f.Filters[0].(*log.NameFilter)
		if !ok {
			t.Fatalf("Not NameFilter")
		}
		assertLogEqual(t, `nf.Name`, "out", nf.Name)

		lf, ok := f.Filters[1].(*log.LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `lf.Level`, log.LevelTrace, lf.Level)
	}

	i++
	{
		_, ok := mw.Writers[i].(*log.StreamWriter)
		if !ok {
			t.Fatalf("Not StreamWriter")
		}
	}

	i++
	{
		aw, ok := mw.Writers[i].(*log.AsyncWriter)
		if !ok {
			t.Fatalf("Not AsyncWriter")
		}

		writer := testGetLogWriter(aw)
		w, ok := writer.(*log.ConnWriter)
		if !ok {
			t.Fatalf("Not ConnWriter")
		}
		assertLogEqual(t, `w.Net`, "tcp", w.Net)
		assertLogEqual(t, `w.Addr`, "localhost:9999", w.Addr)
		assertLogEqual(t, `w.Timeout`, time.Second*5, w.Timeout)

		f, ok := w.Filter.(*log.LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `f.Level`, log.LevelDebug, f.Level)
	}

	i++
	{
		aw, ok := mw.Writers[i].(*log.AsyncWriter)
		if !ok {
			t.Fatalf("Not AsyncWriter")
		}

		writer := testGetLogWriter(aw)
		w, ok := writer.(*log.FileWriter)
		if !ok {
			t.Fatalf("Not FileWriter")
		}
		assertLogEqual(t, `w.DirPerm`, uint32(0777), w.DirPerm)
		assertLogEqual(t, `w.MaxDays`, 7, w.MaxDays)
		assertLogEqual(t, `w.MaxSize`, int64(num.MB*4), w.MaxSize)
		assertLogEqual(t, `w.SyncLevel`, log.LevelError, w.SyncLevel)

		f, ok := w.Filter.(*log.LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `f.Level`, log.LevelInfo, f.Level)
	}

	i++
	{
		aw, ok := mw.Writers[i].(*log.AsyncWriter)
		if !ok {
			t.Fatalf("Not AsyncWriter")
		}

		writer := testGetLogWriter(aw)
		w, ok := writer.(*slacklog.SlackWriter)
		if !ok {
			t.Fatalf("Not SlackWriter")
		}
		assertLogEqual(t, `w.Webhook`, "https://hooks.slack.com/services/...", w.Webhook)
		assertLogEqual(t, `w.Timeout`, time.Second*5, w.Timeout)

		f, ok := w.Filter.(*log.LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `f.Level`, log.LevelWarn, f.Level)
	}

	i++
	{
		aw, ok := mw.Writers[i].(*log.AsyncWriter)
		if !ok {
			t.Fatalf("Not AsyncWriter")
		}

		writer := testGetLogWriter(aw)
		fw, ok := writer.(*log.FailoverWriter)
		if !ok {
			t.Fatalf("Not FailoverWriter")
		}

		writer = testGetLogWriter(fw)
		w, ok := writer.(*smtplog.SMTPWriter)
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

		f, ok := w.Filter.(*log.LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `f.Level`, log.LevelError, f.Level)
	}

	i++
	{
		aw, ok := mw.Writers[i].(*log.AsyncWriter)
		if !ok {
			t.Fatalf("Not AsyncWriter")
		}

		writer := testGetLogWriter(aw)
		w, ok := writer.(*teamslog.TeamsWriter)
		if !ok {
			t.Fatalf("Not TeamsWriter")
		}
		assertLogEqual(t, `w.Webhook`, "https://xxx.webhook.office.com/webhookb2/...", w.Webhook)
		assertLogEqual(t, `w.Timeout`, time.Second*3, w.Timeout)

		f, ok := w.Filter.(*log.LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `f.Level`, log.LevelFatal, f.Level)
	}

	i++
	{
		aw, ok := mw.Writers[i].(*log.FailoverWriter)
		if !ok {
			t.Fatalf("Not FailoverWriter")
		}

		writer := testGetLogWriter(aw)
		w, ok := writer.(*httplog.HTTPWriter)
		if !ok {
			t.Fatalf("Not HTTPWriter")
		}
		assertLogEqual(t, `w.Webhook`, "http://localhost:9200/pango_logs/_doc", w.URL)
		assertLogEqual(t, `w.ContentType`, "application/json", w.ContentType)
		assertLogEqual(t, `w.Timeout`, time.Second*5, w.Timeout)

		jf, ok := w.Formatter.(*log.JSONFormatter)
		if jf == nil || !ok {
			t.Fatalf("Not JSONFormatter")
		}

		f, ok := w.Filter.(*log.LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `f.Level`, log.LevelFatal, f.Level)
	}

	i++
	{
		aw, ok := mw.Writers[i].(*log.AsyncWriter)
		if !ok {
			t.Fatalf("Not AsyncWriter")
		}

		writer := testGetLogWriter(aw)
		w, ok := writer.(*httplog.HTTPWriter)
		if !ok {
			t.Fatalf("Not HTTPWriter")
		}
		assertLogEqual(t, `w.URL`, "http://localhost:9200/pango_logs/_bulk", w.URL)
		assertLogEqual(t, `w.ContentType`, "application/json", w.ContentType)
		assertLogEqual(t, `w.Timeout`, time.Second*5, w.Timeout)

		jf, ok := w.Formatter.(*log.JSONFormatter)
		if jf == nil || !ok {
			t.Fatalf("Not JSONFormatter")
		}

		f, ok := w.Filter.(*log.LevelFilter)
		if !ok {
			t.Fatalf("Not LevelFilter")
		}
		assertLogEqual(t, `f.Level`, log.LevelDebug, f.Level)

		assertLogEqual(t, `w.BatchCount`, 5, w.BatchCount)
		assertLogEqual(t, `w.CacheCount`, 10, w.CacheCount)
		assertLogEqual(t, `w.FlushLevel`, log.LevelError, w.FlushLevel)
		assertLogEqual(t, `w.FlushDelta`, time.Second*60, w.FlushDelta)
	}
}

func TestLogConfigFile1(t *testing.T) {
	os.RemoveAll("conftest")
	defer os.RemoveAll("conftest")

	lg := log.NewLog()
	if err := lg.Config("testdata/log-file1.json"); err != nil {
		t.Fatalf(`log.Config("testdata/log-file1.json") = %v`, err)
	}
	lg.Info("This is info.")
	lg.Warn("This is warn.")
	lg.Error("This is error.")
	lg.Close()

	bs, _ := os.ReadFile("conftest/logs/file1.log")
	a := string(bs)
	w := "ERROR - This is error." + eol
	if a != w {
		t.Errorf("\n actual = %v\n expect = %v", a, w)
	}
}

func TestLogConfigFile2(t *testing.T) {
	os.RemoveAll("conftest")
	defer os.RemoveAll("conftest")

	lg := log.NewLog()
	if err := lg.Config("testdata/log-file2.json"); err != nil {
		t.Fatalf(`log.Config("testdata/log-file2.json") = %v`, err)
	}
	lg.Info("This is info.")
	lg.Warn("This is warn.")
	lg.Error("This is error.")

	tl := lg.GetLogger("test")
	tl.Warn("This is WARN.")
	tl.Error("This is ERROR.")
	lg.Close()

	bs, _ := os.ReadFile("conftest/logs/file1.log")
	a := string(bs)
	w := "ERROR - This is error." + eol + "ERROR - This is ERROR." + eol
	if a != w {
		t.Errorf("\n actual = %v\n expect = %v", a, w)
	}

	bs, _ = os.ReadFile("conftest/logs/file2.log")
	a = string(bs)
	w = "WARN - This is WARN." + eol + "ERROR - This is ERROR." + eol
	if a != w {
		t.Errorf("\n actual = %v\n expect = %v", a, w)
	}
}
