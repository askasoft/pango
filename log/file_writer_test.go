package log

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/askasoft/pango/fsu"
)

func TestFileTextFormatSimple(t *testing.T) {
	testdir := "TestFileTextFormatSimple-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(&FileWriter{Path: path})
	log.Info("hello")
	log.Close()

	// check lastest file
	bs, _ := os.ReadFile(path + ".log")
	e := `[I] hello` + EOL
	a := string(bs)
	if a != e {
		t.Errorf("TestFileTextFormatSimple\n expect: %q, actual: %q", e, a)
	}
}

func TestFilePropGlobal(t *testing.T) {
	testdir := "TestFilePropGlobal-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	SetFormatter(NewTextFormatter("%l - %x{key} - %m%n%T"))
	SetWriter(&FileWriter{Path: path})
	SetProp("key", "val")
	Info("hello")
	Close()

	// check lastest file
	bs, _ := os.ReadFile(path + ".log")
	e := fmt.Sprintf("INFO - %s - %s%s", "val", "hello", EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFilePropGlobal\n expect: %q, actual: %q", e, a)
	}
}

func TestFilePropDefault(t *testing.T) {
	testdir := "TestFilePropDefault-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	log1 := Default()
	log1.SetFormatter(NewTextFormatter("%l - %x{key} - %m%n%T"))
	log1.SetWriter(&FileWriter{Path: path})
	log1.SetProp("key", "val")
	log1.Info("hello")
	log1.Close()

	// check lastest file
	bs, _ := os.ReadFile(path + ".log")
	e := fmt.Sprintf("INFO - %s - %s%s", "val", "hello", EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFilePropDefault\n expect: %q, actual: %q", e, a)
	}
}

func TestFilePropNewLog(t *testing.T) {
	testdir := "TestFilePropNewLog-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

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
	bs, _ := os.ReadFile(path + ".log")
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
	testdir := "TestFileCallerGlobal-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	SetFormatter(NewTextFormatter("%l %S:%L %F() - %m%n%T"))
	SetWriter(&FileWriter{Path: path})
	file, line, ffun := testGetCaller(1)
	Info("hello")
	Close()

	// check lastest file
	bs, _ := os.ReadFile(path + ".log")
	e := fmt.Sprintf("INFO %s:%d %s() - hello%s", file, line, ffun, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileCallerGlobal\n expect: %q, actual: %q", e, a)
	}
}

func TestFileCallerNewLog(t *testing.T) {
	testdir := "TestFileCallerNewLog-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	log := NewLog()
	log.SetFormatter(NewTextFormatter("%l %S:%L %F() - %m%n%T"))
	log.SetWriter(&FileWriter{Path: path})
	file, line, ffun := testGetCaller(1)
	log.Info("hello")
	log.Close()

	// check lastest file
	bs, _ := os.ReadFile(path + ".log")
	e := fmt.Sprintf("INFO %s:%d %s() - hello%s", file, line, ffun, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileCallerNewLog\n expect: %q, actual: %q", e, a)
	}
}

func TestFileCallerNewLog2(t *testing.T) {
	testdir := "TestFileCallerNewLog2-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	log := NewLog()
	log.SetFormatter(NewTextFormatter("%l %S:%L %F() - %m%n%T"))
	log.SetWriter(&FileWriter{Path: path})
	file, line, ffun := testGetCaller(1)
	log.Log(LevelInfo, "hello")
	log.Close()

	// check lastest file
	bs, _ := os.ReadFile(path + ".log")
	e := fmt.Sprintf("INFO %s:%d %s() - hello%s", file, line, ffun, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileCallerNewLog2\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateMaxSize(t *testing.T) {
	testdir := "TestFileRotateMaxSize-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest.log"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

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
		bs, _ := os.ReadFile(sp)
		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateMaxSize\n expect: %q, actual: %q", e, a)
		}
	}

	// check lastest file
	bs, _ := os.ReadFile(path)
	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateMaxSize\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateMaxSizeGzip(t *testing.T) {
	testdir := "TestFileRotateMaxSizeGzip-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest.log"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(&FileWriter{Path: path, MaxSize: 10, Gzip: true})
	for i := 1; i < 10; i++ {
		log.Info("hello test ", i)
	}
	log.Close()

	// sleep for gzip compress
	time.Sleep(time.Second * 3)

	// check existing files
	for i := 1; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%03d.log.gz", i))
		bs, err := os.ReadFile(sp)
		if err != nil {
			t.Fatalf("TestFileRotateMaxSizeGzip\n failed to read file %q, %v", sp, err)
		}

		gr, err := gzip.NewReader(bytes.NewReader(bs))
		if err != nil {
			t.Fatalf("TestFileRotateMaxSizeGzip\n failed to read gzip %q, %v", sp, err)
		}

		bs, err = io.ReadAll(gr)
		if err != nil {
			t.Fatalf("TestFileRotateMaxSizeGzip\n failed to read gzip %q, %v", sp, err)
		}

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateMaxSizeGzip\n expect: %q, actual: %q", e, a)
		}
	}

	// check lastest file
	bs, _ := os.ReadFile(path)

	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateMaxSizeGzip\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateMaxSizeDaily(t *testing.T) {
	testdir := "TestFileRotateMaxSizeDaily-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest.log"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

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
		bs, _ := os.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateMaxSizeDaily\n expect: %q, actual: %q", e, a)
		}
	}

	// check lastest file
	bs, _ := os.ReadFile(path)

	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateMaxSizeDaily\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateMaxSplit(t *testing.T) {
	testdir := "TestFileRotateMaxSplit-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest.log"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

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

		err := fsu.FileExists(sp)
		if err == nil {
			t.Errorf("TestFileRotateMaxSplit file %q exists: %v", sp, err)
		}
	}

	// check existing files
	for i := 6; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%03d.log", i))
		bs, _ := os.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateMaxSplit\n expect: %q, actual: %q", e, a)
		}
	}

	// check lastest file
	bs, _ := os.ReadFile(path)

	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateMaxSplit\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateMaxFilesHourly(t *testing.T) {
	testdir := "TestFileRotateMaxFilesHourly-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest.log"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

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

		err := fsu.FileExists(sp)
		if err == nil {
			t.Errorf("TestFileRotateMaxFilesHourly file %q exists: %v", sp, err)
		}
	}

	// check existing files
	for i := 6; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s-%03d.log", now.Format("2006010215"), i))
		bs, _ := os.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateMaxFilesHourly\n expect: %q, actual: %q", e, a)
		}
	}

	// check lastest file
	bs, _ := os.ReadFile(path)

	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateMaxFilesHourly\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateDaily(t *testing.T) {
	testdir := "TestFileRotateDaily-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest.log"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	fw := &FileWriter{Path: path, MaxDays: 100}
	lg := NewLog()
	lg.SetFormatter(TextFmtSimple)

	now := time.Now()
	tm := now
	for i := 1; i < 10; i++ {
		le := newEvent(lg, LevelInfo, "hello test "+strconv.Itoa(i))
		le.Time = tm
		fw.Write(le)
		fw.openTime = tm
		tm = tm.Add(time.Hour * 24)
	}
	time.Sleep(time.Millisecond * 100)
	fw.Close()

	// check existing files
	tm = now
	for i := 1; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("20060102")))
		bs, _ := os.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateDaily\n expect: %q, actual: %q", e, a)
		}

		tm = tm.Add(time.Hour * 24)
	}

	// check lastest file
	bs, _ := os.ReadFile(path)

	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateDaily\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateDailyInit(t *testing.T) {
	testdir := "TestFileRotateDailyInit-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest.log"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	// create yesterday file
	os.MkdirAll(testdir, os.FileMode(0777))
	os.WriteFile(path, []byte("init"), os.FileMode(0666))
	yes := time.Now().Add(time.Hour * -24)
	os.Chtimes(path, yes, yes)

	fw := &FileWriter{Path: path, MaxDays: 100}
	lg := NewLog()
	lg.SetFormatter(TextFmtSimple)

	now := time.Now()
	tm := now
	for i := 1; i < 10; i++ {
		le := newEvent(lg, LevelInfo, "hello test "+strconv.Itoa(i))
		le.Time = tm
		fw.Write(le)
		fw.openTime = tm
		tm = tm.Add(time.Hour * 24)
	}
	time.Sleep(time.Millisecond * 100)
	fw.Close()

	// check init rotated files
	tm = yes
	for {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("20060102")))
		bs, _ := os.ReadFile(sp)

		e := "init"
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateDailyInit\n expect: %q, actual: %q", e, a)
		}

		break
	}

	// check rotated files
	tm = now
	for i := 1; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("20060102")))
		bs, _ := os.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateDailyInit\n expect: %q, actual: %q", e, a)
		}

		tm = tm.Add(time.Hour * 24)
	}

	// check lastest file
	bs, _ := os.ReadFile(path)

	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateDaily\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateDailyOutdated(t *testing.T) {
	testdir := "TestFileRotateDailyOutdated-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest.log"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	fw := &FileWriter{Path: path, MaxDays: 3}
	lg := NewLog()
	lg.SetFormatter(TextFmtSimple)

	now := time.Now().Add(time.Hour * 24 * -8)
	tm := now
	for i := 1; i < 10; i++ {
		le := newEvent(lg, LevelInfo, "hello test "+strconv.Itoa(i))
		le.Time = tm
		fw.Write(le)
		fw.openTime = tm
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

		err := fsu.FileExists(sp)
		if err == nil {
			t.Errorf("TestFileRotateDailyOutdated file %q exists: %v", sp, err)
		}

		tm = tm.Add(time.Hour * 24)
	}

	// check existing files
	for i := 7; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("20060102")))
		bs, _ := os.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateDailyOutdated\n expect: %q, actual: %q", e, a)
		}

		tm = tm.Add(time.Hour * 24)
	}

	// check lastest file
	bs, _ := os.ReadFile(path)
	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateDailyOutdated\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateHourly(t *testing.T) {
	testdir := "TestFileRotateHourly-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest.log"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	fw := &FileWriter{Path: path, MaxHours: 100}
	lg := NewLog()
	lg.SetFormatter(TextFmtSimple)

	now := time.Now()
	tm := now
	for i := 1; i < 10; i++ {
		le := newEvent(lg, LevelInfo, "hello test "+strconv.Itoa(i))
		le.Time = tm
		fw.Write(le)
		fw.openTime = tm
		tm = tm.Add(time.Hour)
	}
	time.Sleep(time.Millisecond * 100)
	fw.Close()

	// check existing files
	tm = now
	for i := 1; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("2006010215")))
		bs, _ := os.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateHourly\n expect: %q, actual: %q", e, a)
		}

		tm = tm.Add(time.Hour)
	}

	// check lastest file
	bs, _ := os.ReadFile(path)

	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateHourly\n expect: %q, actual: %q", e, a)
	}
}

func TestFileRotateHourlyOutdated(t *testing.T) {
	testdir := "TestFileRotateHourlyOutdated-" + strconv.Itoa(rand.Int())
	path := testdir + "/filetest.log"

	os.RemoveAll(testdir)
	defer os.RemoveAll(testdir)

	fw := &FileWriter{Path: path, MaxHours: 3}
	lg := NewLog()
	lg.SetFormatter(TextFmtSimple)

	now := time.Now().Add(time.Hour * -8)
	tm := now
	for i := 1; i < 10; i++ {
		le := newEvent(lg, LevelInfo, "hello test "+strconv.Itoa(i))
		le.Time = tm
		fw.Write(le)
		fw.openTime = tm
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

		err := fsu.FileExists(sp)
		if err == nil {
			t.Errorf("TestFileRotateHourlyOutdated file %q exists: %v", sp, err)
		}

		tm = tm.Add(time.Hour)
	}

	// check existing files
	for i := 7; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("2006010215")))
		bs, _ := os.ReadFile(sp)

		e := fmt.Sprintf(`[I] hello test %d%s`, i, EOL)
		a := string(bs)
		if a != e {
			t.Errorf("TestFileRotateHourlyOutdated\n expect: %q, actual: %q", e, a)
		}

		tm = tm.Add(time.Hour)
	}

	// check lastest file
	bs, _ := os.ReadFile(path)
	e := fmt.Sprintf(`[I] hello test %d%s`, 9, EOL)
	a := string(bs)
	if a != e {
		t.Errorf("TestFileRotateHourlyOutdated\n expect: %q, actual: %q", e, a)
	}
}
