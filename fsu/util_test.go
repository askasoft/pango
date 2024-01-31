package fsu

import (
	"errors"
	"path"
	"reflect"
	"strings"
	"testing"
)

const testdir = "_testdir"

func TestFileExists(t *testing.T) {
	defer RemoveAll(testdir)
	MkdirAll(testdir, FileMode(0777))
	WriteFile(path.Join(testdir, "test.txt"), []byte("a"), FileMode(0666))

	cs := []struct {
		f string
		e error
	}{
		{"_testdir/test.txt", nil},
		{"_testdir/notexist.txt", ErrNotExist},
		{"_testdir", ErrIsDir},
	}

	for i, c := range cs {
		err := FileExists(c.f)
		if !errors.Is(err, c.e) {
			t.Errorf("[%d] %s - %v, want %v", i, c.f, err, c.e)
		}
	}
}

func TestFileSize(t *testing.T) {
	defer RemoveAll(testdir)
	MkdirAll(testdir, FileMode(0777))
	WriteFile(path.Join(testdir, "test.txt"), []byte("a"), FileMode(0666))

	cs := []struct {
		f string
		z int64
		e error
	}{
		{"_testdir/test.txt", 1, nil},
		{"_testdir/notexist.txt", 0, ErrNotExist},
		{"_testdir", 0, ErrIsDir},
	}

	for i, c := range cs {
		err := FileExists(c.f)
		if !errors.Is(err, c.e) {
			t.Errorf("[%d] %s - %v, want %v", i, c.f, err, c.e)
		}
	}
}

func TestDirExists(t *testing.T) {
	defer RemoveAll(testdir)
	MkdirAll(testdir, FileMode(0777))
	WriteFile(path.Join(testdir, "test.txt"), []byte("a"), FileMode(0666))

	cs := []struct {
		f string
		e error
	}{
		{"_testdir/test.txt", ErrNotDir},
		{"_testdir/notexist", ErrNotExist},
		{"_testdir", nil},
	}

	for i, c := range cs {
		err := DirExists(c.f)
		if !errors.Is(err, c.e) {
			t.Errorf("[%d] %s - %v, want %v", i, c.f, err, c.e)
		}
	}
}

func TestDirIsEmpty(t *testing.T) {
	defer RemoveAll(testdir)
	MkdirAll(testdir, FileMode(0777))
	MkdirAll(path.Join(testdir, "empty"), FileMode(0777))
	WriteFile(path.Join(testdir, "test.txt"), []byte("a"), FileMode(0666))

	cs := []struct {
		f string
		e error
	}{
		{"_testdir/test.txt", ErrNotDir},
		{"_testdir/notexist", ErrNotExist},
		{"_testdir/empty", nil},
	}

	for i, c := range cs {
		err := DirIsEmpty(c.f)
		if !errors.Is(err, c.e) {
			t.Errorf("[%d] %s - %v, want %v", i, c.f, err, c.e)
		}
	}
}

func TestCopyFile(t *testing.T) {
	defer RemoveAll(testdir)
	MkdirAll(testdir, FileMode(0777))

	cs := []struct {
		s string
		d string
	}{
		{"1.txt", "_testdir/1.txt"},
	}

	for i, c := range cs {
		defer Remove(c.s)
		defer Remove(c.d)

		sbs := []byte(strings.Repeat("a", (i+1)*10))
		err := WriteFile(c.s, sbs, FileMode(0600))
		if err != nil {
			t.Fatalf("#%d: %v", i, err)
		}
		err = CopyFile(c.s, c.d)
		if err != nil {
			t.Fatalf("#%d: %v", i, err)
		}
		dbs, err := ReadFile(c.d)
		if err != nil {
			t.Fatalf("#%d: %v", i, err)
		}
		if !reflect.DeepEqual(sbs, dbs) {
			t.Fatalf("#%d: \n GOT: %v\nWANT: %v\n", i, dbs, sbs)
		}
	}
}
