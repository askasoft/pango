package log

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, `[I] hello`+eol, string(bs))
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
	assert.Equal(t, fmt.Sprintf("INFO %s:%d %s() - hello%s", file, line, ffun, eol), string(bs))
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
	assert.Equal(t, fmt.Sprintf("INFO %s:%d %s() - hello%s", file, line, ffun, eol), string(bs))
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
	assert.Equal(t, fmt.Sprintf("INFO %s:%d %s() - hello%s", file, line, ffun, eol), string(bs))
}

func TestFileSyncWrite(t *testing.T) {
	path := "TestFileSyncWrite/filetest.log"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	log := NewLog()
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(&FileWriter{Path: path})

	wg := sync.WaitGroup{}
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			time.Sleep(time.Microsecond * 10)
			for i := 1; i < 10; i++ {
				log.Info(n, i)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Close()

	// read actual log
	bs, _ := ioutil.ReadFile(path)
	as := strings.Split(strings.TrimSuffix(string(bs), eol), eol)
	sort.Strings(as)

	// expected data
	es := []string{}
	for n := 1; n < 10; n++ {
		for i := 1; i < 10; i++ {
			es = append(es, fmt.Sprint("[I] ", n, i))
		}
	}

	assert.Equal(t, es, as)
}

func TestFileAsyncWrite(t *testing.T) {
	path := "TestFileAsyncWrite/filetest.log"
	dir := filepath.Dir(path)
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	log := NewLog()
	log.Async(10)
	log.SetFormatter(TextFmtSimple)
	log.SetWriter(&FileWriter{Path: path})

	wg := sync.WaitGroup{}
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			time.Sleep(time.Microsecond * 10)
			for i := 1; i < 10; i++ {
				log.Info(n, i)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Close()

	// read actual log
	bs, _ := ioutil.ReadFile(path)
	as := strings.Split(strings.TrimSuffix(string(bs), eol), eol)
	sort.Strings(as)

	// expected data
	es := []string{}
	for n := 1; n < 10; n++ {
		for i := 1; i < 10; i++ {
			es = append(es, fmt.Sprint("[I] ", n, i))
		}
	}

	assert.Equal(t, es, as)
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
		assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, i, eol), string(bs))
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)
	assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, 9, eol), string(bs))
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
		assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, i, eol), string(bs))
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)
	assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, 9, eol), string(bs))
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
		assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, i, eol), string(bs))
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)
	assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, 9, eol), string(bs))
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
		assert.NoFileExists(t, sp)
	}

	// check existing files
	for i := 6; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%03d.log", i))
		bs, _ := ioutil.ReadFile(sp)
		assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, i, eol), string(bs))
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)
	assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, 9, eol), string(bs))
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
		assert.NoFileExists(t, sp)
	}

	// check existing files
	for i := 6; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s-%03d.log", now.Format("2006010215"), i))
		bs, _ := ioutil.ReadFile(sp)
		assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, i, eol), string(bs))
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)
	assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, 9, eol), string(bs))
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
		le.When = tm
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
		assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, i, eol), string(bs))
		tm = tm.Add(time.Hour * 24)
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)
	assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, 9, eol), string(bs))
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
		le.When = tm
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
		assert.NoFileExists(t, sp)
		tm = tm.Add(time.Hour * 24)
	}

	// check existing files
	for i := 7; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("20060102")))
		bs, _ := ioutil.ReadFile(sp)
		assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, i, eol), string(bs))
		tm = tm.Add(time.Hour * 24)
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)
	assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, 9, eol), string(bs))
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
		le.When = tm
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
		assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, i, eol), string(bs))
		tm = tm.Add(time.Hour)
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)
	assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, 9, eol), string(bs))
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
		le.When = tm
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
		assert.NoFileExists(t, sp)
		tm = tm.Add(time.Hour)
	}

	// check existing files
	for i := 7; i < 9; i++ {
		sp := strings.ReplaceAll(path, ".log", fmt.Sprintf("-%s.log", tm.Format("2006010215")))
		bs, _ := ioutil.ReadFile(sp)
		assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, i, eol), string(bs))
		tm = tm.Add(time.Hour)
	}

	// check lastest file
	bs, _ := ioutil.ReadFile(path)
	assert.Equal(t, fmt.Sprintf(`[I] hello test %d%s`, 9, eol), string(bs))
}
