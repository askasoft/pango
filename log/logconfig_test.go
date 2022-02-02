package log

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
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

func assertLogEqual(t *testing.T, msg string, want interface{}, val interface{}) {
	if want != val {
		t.Errorf("msg = %v, want %v", val, want)
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

	assertLogEqual(t, `len(mw.Writers)`, 7, len(mw.Writers))

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
		assertLogEqual(t, `lf.Level`, LevelDebug, lf.Level)
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
		assertLogEqual(t, `f.Level`, LevelError, f.Level)
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
		assertLogEqual(t, `w.SyncLevel`, LevelError, w.SyncLevel)

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

		w, ok := aw.writer.(*SlackWriter)
		if !ok {
			t.Fatalf("Not SlackWriter")
		}
		assertLogEqual(t, `w.Channel`, "develop", w.Channel)
		assertLogEqual(t, `w.Username`, "gotest", w.Username)
		assertLogEqual(t, `w.Webhook`, "https://hooks.slack.com/services/...", w.Webhook)
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

		w, ok := aw.writer.(*SMTPWriter)
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
		assertLogEqual(t, `f.Level`, LevelError, f.Level)
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

	bs, _ := ioutil.ReadFile("conftest/logs/file1.log")
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

	bs, _ := ioutil.ReadFile("conftest/logs/file1.log")
	a := string(bs)
	w := "ERROR - This is error." + EOL + "ERROR - This is ERROR." + EOL
	if a != w {
		t.Errorf("\n actual = %v\n expect = %v", a, w)
	}

	bs, _ = ioutil.ReadFile("conftest/logs/file2.log")
	a = string(bs)
	w = "WARN - This is WARN." + EOL + "ERROR - This is ERROR." + EOL
	if a != w {
		t.Errorf("\n actual = %v\n expect = %v", a, w)
	}
}
