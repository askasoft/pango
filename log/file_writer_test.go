package log

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pandafw/pango/iox"
)

func TestFileTextFormatSimple(t *testing.T) {
	path := "TestFileTextFormatSimple/filetest"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(&FileWriter{Path: path})
	log.Info("hello")
	log.Close()

	// check lastest file
	bs, _ := ioutil.ReadFile(path + ".log")
	e := `[I] hello` + EOL
	a := string(bs)
	if a != e {
		t.Errorf("TestFileTextFormatSimple\n expect: %q, actual: %q", e, a)
	}
}

func TestFilePropGlobal(t *testing.T) {
	path := "TestFilePropGlobal/filetest"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	SetFormatter(NewTextFormatter("%l - %x{key} - %m%n%T"))
	SetWriter(&FileWriter{Path: path})
	SetProp("key", "val")
	Info("hello")
	Close()

	// check lastest file
	bs, _ := ioutil.ReadFile(path + ".log")
	e := fmt.Sprintf("INFO - %s - %s%s", "val", "hello", EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFilePropGlobal\n expect: %q, actual: %q", e, a)
	}
}

func TestFilePropDefault(t *testing.T) {
	path := "TestFilePropDefault/filetest"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	log1 := Default()
	log1.SetFormatter(NewTextFormatter("%l - %x{key} - %m%n%T"))
	log1.SetWriter(&FileWriter{Path: path})
	log1.SetProp("key", "val")
	log1.Info("hello")
	log1.Close()

	// check lastest file
	bs, _ := ioutil.ReadFile(path + ".log")
	e := fmt.Sprintf("INFO - %s - %s%s", "val", "hello", EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFilePropDefault\n expect: %q, actual: %q", e, a)
	}
}

func TestFilePropNewLog(t *testing.T) {
	path := "TestFilePropNewLog/filetest"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	log1 := NewLog()
	log1.SetFormatter(NewTextFormatter("%l - %X - %m%n%T"))
	log1.SetWriter(&FileWriter{Path: path})
	log1.SetProp("key1", "val1")
	log1.Info("hello")

	log2 := log1.GetLogger("")
	log2.SetProp("key2", "val2")
	log2.Info("hello")
	log1.Close()

	// check lastest file
	bs, _ := ioutil.ReadFile(path + ".log")
	e1 := fmt.Sprintf("INFO - %s=%s - %s%s", "key1", "val1", "hello", EOL) +
		fmt.Sprintf("INFO - %s=%s %s=%s - %s%s", "key1", "val1", "key2", "val2", "hello", EOL)
	e2 := fmt.Sprintf("INFO - %s=%s - %s%s", "key1", "val1", "hello", EOL) +
		fmt.Sprintf("INFO - %s=%s %s=%s - %s%s", "key2", "val2", "key1", "val1", "hello", EOL)
	a := string(bs)
	if a != e1 && a != e2 {
		t.Errorf("TestFilePropNewLog\n expect: %q\n actual: %q", e1, a)
	}
}

func TestFileCallerGlobal(t *testing.T) {
	path := "TestFileCallerGlobal/filetest"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	SetFormatter(NewTextFormatter("%l %S:%L %F() - %m%n%T"))
	SetWriter(&FileWriter{Path: path})
	file, line, ffun := testGetCaller(1)
	Info("hello")
	Close()

	// check lastest file
	bs, _ := ioutil.ReadFile(path + ".log")
	e := fmt.Sprintf("INFO %s:%d %s() - hello%s", file, line, ffun, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileCallerGlobal\n expect: %q, actual: %q", e, a)
	}
}

func TestFileCallerNewLog(t *testing.T) {
	path := "TestFileCallerNewLog/filetest"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	log := NewLog()
	log.SetFormatter(NewTextFormatter("%l %S:%L %F() - %m%n%T"))
	log.SetWriter(&FileWriter{Path: path})
	file, line, ffun := testGetCaller(1)
	log.Info("hello")
	log.Close()

	// check lastest file
	bs, _ := ioutil.ReadFile(path + ".log")
	e := fmt.Sprintf("INFO %s:%d %s() - hello%s", file, line, ffun, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileCallerNewLog\n expect: %q, actual: %q", e, a)
	}
}

func TestFileCallerNewLog2(t *testing.T) {
	path := "TestFileCallerNewLog2/filetest"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	log := NewLog()
	log.SetFormatter(NewTextFormatter("%l %S:%L %F() - %m%n%T"))
	log.SetWriter(&FileWriter{Path: path})
	file, line, ffun := testGetCaller(1)
	log.Log(LevelInfo, "hello")
	log.Close()

	// check lastest file
	bs, _ := ioutil.ReadFile(path + ".log")
	e := fmt.Sprintf("INFO %s:%d %s() - hello%s", file, line, ffun, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileCallerNewLog2\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateMaxSize(t *testing.T) {
	path := "TestFileRotateMaxSize/filetest.log"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(&FileWriter{Path: path, MaxSize: 10})
	for i := 1; i < 10; i++ {
		log.Info("hello test ", i)
	}
	log.Close()

	// check existing files
	for i := 1; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%03d.log", i))
		bs, _ := ioutil.ReadFile(sp)
		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateMaxSize\n expect: %q, actual: %q", e, a)
		}
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)
	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateMaxSize\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateMaxSizeGzip(t *testing.T) {
	path := "TestFileRotateMaxSizeGzip/filetest.log"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(&FileWriter{Path: path, MaxSize: 10, Gzip: true})
	for i := 1; i < 10; i++ {
		log.Info("hello test ", i)
	}
	time.Sleep(time.Millisecond * 10)
	log.Close()

	// check existing files
	for i := 1; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%03d.log.gz", i))
		bs, _ := ioutil.ReadFile(sp)
		gr, _ := gzip.NewReader(bytes.NewReader(bs))
		bs, _ = ioutil.ReadAll(gr)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateMaxSizeGzip\n expect: %q, actual: %q", e, a)
		}
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)

	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateMaxSizeGzip\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateMaxSizeDaily(t *testing.T) {
	path := "TestFileRotateMaxSizeDaily/filetest.log"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(&FileWriter{Path: path, MaxDays: 100, MaxSize: 10})

	now := time.Now()
	for i := 1; i < 10; i++ {
		log.Info("hello test ", i)
	}
	log.Close()

	// check existing files
	for i := 1; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s-%03d.log", now.Format("20060102"), i))
		bs, _ := ioutil.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateMaxSizeDaily\n expect: %q, actual: %q", e, a)
		}
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)

	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateMaxSizeDaily\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateMaxSplit(t *testing.T) {
	path := "TestFileRotateMaxSplit/filetest.log"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(&FileWriter{Path: path, MaxSplit: 3, MaxSize: 10})
	for i := 1; i < 10; i++ {
		log.Info("hello test ", i)
	}
	log.Close()

	// check deleted files
	for i := 1; i < 6; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%03d.log", i))

		err := iox.FileExists(sp)
		if err == nil {
			t.Errorf("TestFileRotateMaxSplit file %q exists: %v", sp, err)
		}
	}

	// check existing files
	for i := 6; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%03d.log", i))
		bs, _ := ioutil.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateMaxSplit\n expect: %q, actual: %q", e, a)
		}
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)

	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateMaxSplit\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateMaxFilesHourly(t *testing.T) {
	path := "TestFileRotateMaxFilesHourly/filetest.log"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(&FileWriter{Path: path, MaxSplit: 3, MaxSize: 10, MaxHours: 100})
	now := time.Now()
	for i := 1; i < 10; i++ {
		log.Info("hello test ", i)
	}
	log.Close()

	// check deleted files
	for i := 1; i < 6; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s-%03d.log", now.Format("2006010215"), i))

		err := iox.FileExists(sp)
		if err == nil {
			t.Errorf("TestFileRotateMaxFilesHourly file %q exists: %v", sp, err)
		}
	}

	// check existing files
	for i := 6; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s-%03d.log", now.Format("2006010215"), i))
		bs, _ := ioutil.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateMaxFilesHourly\n expect: %q, actual: %q", e, a)
		}
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)

	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateMaxFilesHourly\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateDaily(t *testing.T) {
	path := "TestFileRotateDaily/filetest.log"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	fw := &FileWriter{Path: path, MaxDays: 100}
	lg := NewLog()
	lg.SetFormatter(TextFmtSimple)

	now := time.Now()
	tm := now
	for i := 1; i < 10; i++ {
		le := newEvent(lg, LevelInfo, "hello test "+strconv.Itoa(i))
		le.when = tm
		fw.Write(le)
		fw.openTime = tm
		fw.openDay = tm.Day()
		fw.openHour = tm.Hour()
		tm = tm.Add(time.Hour * 24)
	}
	time.Sleep(time.Millisecond * 100)
	fw.Close()

	// check existing files
	tm = now
	for i := 1; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("20060102")))
		bs, _ := ioutil.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateDaily\n expect: %q, actual: %q", e, a)
		}

		tm = tm.Add(time.Hour * 24)
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)

	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateDaily\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateDailyOutdated(t *testing.T) {
	path := "TestFileRotateDailyOutdated/filetest.log"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	fw := &FileWriter{Path: path, MaxDays: 3}
	lg := NewLog()
	lg.SetFormatter(TextFmtSimple)

	now := time.Now().Add(time.Hour * 24 * -8)
	tm := now
	for i := 1; i < 10; i++ {
		le := newEvent(lg, LevelInfo, "hello test "+strconv.Itoa(i))
		le.when = tm
		fw.Write(le)
		fw.openTime = tm
		fw.openDay = tm.Day()
		fw.openHour = tm.Hour()
		if i > 1 && i < 10 {
			tm0 := tm.Add(time.Hour * -24)
			sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm0.Format("20060102")))
			os.Chtimes(sp, tm0, tm0)
		}
		tm = tm.Add(time.Hour * 24)

		// let deleteOutdatedFiles goroutine finish
		time.Sleep(time.Millisecond * 100)
	}
	fw.Close()

	// check outdated files
	tm = now
	for i := 1; i < 7; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("20060102")))

		err := iox.FileExists(sp)
		if err == nil {
			t.Errorf("TestFileRotateDailyOutdated file %q exists: %v", sp, err)
		}

		tm = tm.Add(time.Hour * 24)
	}

	// check existing files
	for i := 7; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("20060102")))
		bs, _ := ioutil.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateDailyOutdated\n expect: %q, actual: %q", e, a)
		}

		tm = tm.Add(time.Hour * 24)
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)
	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateDailyOutdated\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateHourly(t *testing.T) {
	path := "TestFileRotateHourly/filetest.log"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	fw := &FileWriter{Path: path, MaxHours: 100}
	lg := NewLog()
	lg.SetFormatter(TextFmtSimple)

	now := time.Now()
	tm := now
	for i := 1; i < 10; i++ {
		le := newEvent(lg, LevelInfo, "hello test "+strconv.Itoa(i))
		le.when = tm
		fw.Write(le)
		fw.openTime = tm
		fw.openDay = tm.Day()
		fw.openHour = tm.Hour()
		tm = tm.Add(time.Hour)
	}
	time.Sleep(time.Millisecond * 100)
	fw.Close()

	// check existing files
	tm = now
	for i := 1; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("2006010215")))
		bs, _ := ioutil.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateHourly\n expect: %q, actual: %q", e, a)
		}

		tm = tm.Add(time.Hour)
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)

	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateHourly\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateHourlyOutdated(t *testing.T) {
	path := "TestFileRotateHourlyOutdated/filetest.log"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	fw := &FileWriter{Path: path, MaxHours: 3}
	lg := NewLog()
	lg.SetFormatter(TextFmtSimple)

	now := time.Now().Add(time.Hour * -8)
	tm := now
	for i := 1; i < 10; i++ {
		le := newEvent(lg, LevelInfo, "hello test "+strconv.Itoa(i))
		le.when = tm
		fw.Write(le)
		fw.openTime = tm
		fw.openDay = tm.Day()
		fw.openHour = tm.Hour()
		if i > 1 && i < 10 {
			tm0 := tm.Add(time.Hour * -1)
			sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm0.Format("2006010215")))
			os.Chtimes(sp, tm0, tm0)
		}
		tm = tm.Add(time.Hour)

		// let deleteOutdatedFiles goroutine finish
		time.Sleep(time.Millisecond * 100)
	}
	fw.Close()

	// check outdated files
	tm = now
	for i := 1; i < 7; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("2006010215")))

		err := iox.FileExists(sp)
		if err == nil {
			t.Errorf("TestFileRotateHourlyOutdated file %q exists: %v", sp, err)
		}

		tm = tm.Add(time.Hour)
	}

	// check existing files
	for i := 7; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("2006010215")))
		bs, _ := ioutil.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateHourlyOutdated\n expect: %q, actual: %q", e, a)
		}

		tm = tm.Add(time.Hour)
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)
	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateHourlyOutdated\n expect: %q, actual: %q", e, a)
	}
}
