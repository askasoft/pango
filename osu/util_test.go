package osu

import (
	"errors"
	"path"
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
