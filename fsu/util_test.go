package fsu

import (
	"errors"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestFileExists(t *testing.T) {
	testdir := "TestFileExists-" + strconv.Itoa(rand.Int())
	defer RemoveAll(testdir)
	MkdirAll(testdir, FileMode(0777))
	WriteFile(filepath.Join(testdir, "test.txt"), []byte("a"), FileMode(0666))

	cs := []struct {
		f string
		e error
	}{
		{"/test.txt", nil},
		{"/notexist.txt", ErrNotExist},
		{"", ErrIsDir},
	}

	for i, c := range cs {
		f := testdir + c.f
		err := FileExists(f)
		if !errors.Is(err, c.e) {
			t.Errorf("#%d FileExists(%q) = %v, want %v", i, c.f, err, c.e)
		}
	}
}

func TestFileSize(t *testing.T) {
	testdir := "TestFileSize-" + strconv.Itoa(rand.Int())
	defer RemoveAll(testdir)
	MkdirAll(testdir, FileMode(0777))
	WriteFile(filepath.Join(testdir, "test.txt"), []byte("a"), FileMode(0666))

	cs := []struct {
		f string
		z int64
		e error
	}{
		{"/test.txt", 1, nil},
		{"/notexist.txt", 0, ErrNotExist},
		{"", 0, ErrIsDir},
	}

	for i, c := range cs {
		f := testdir + c.f
		err := FileExists(f)
		if !errors.Is(err, c.e) {
			t.Errorf("#%d FileSize(%q) = %v, want %v", i, c.f, err, c.e)
		}
	}
}

func TestDirExists(t *testing.T) {
	testdir := "TestDirExists-" + strconv.Itoa(rand.Int())
	defer RemoveAll(testdir)
	MkdirAll(testdir, FileMode(0777))
	WriteFile(filepath.Join(testdir, "test.txt"), []byte("a"), FileMode(0666))

	cs := []struct {
		f string
		e error
	}{
		{"/test.txt", ErrNotDir},
		{"/notexist", ErrNotExist},
		{"", nil},
	}

	for i, c := range cs {
		f := testdir + c.f
		err := DirExists(f)
		if !errors.Is(err, c.e) {
			t.Errorf("#%d DirExists(%q) = %v, want %v", i, c.f, err, c.e)
		}
	}
}

func TestDirIsEmpty(t *testing.T) {
	testdir := "TestDirIsEmpty-" + strconv.Itoa(rand.Int())
	defer RemoveAll(testdir)
	MkdirAll(testdir, FileMode(0777))
	WriteFile(filepath.Join(testdir, "test.txt"), []byte("a"), FileMode(0666))
	MkdirAll(filepath.Join(testdir, "empty"), FileMode(0777))
	MkdirAll(filepath.Join(testdir, "hasdir", "empty"), FileMode(0777))
	MkdirAll(filepath.Join(testdir, "hasfile"), FileMode(0777))
	WriteFile(filepath.Join(testdir, "hasfile", "test.txt"), []byte("a"), FileMode(0666))

	cs := []struct {
		f string
		e error
	}{
		{"/test.txt", ErrNotDir},
		{"/notexist", ErrNotExist},
		{"/hasdir", ErrNotEmpty},
		{"/hasfile", ErrNotEmpty},
		{"/empty", nil},
	}

	for i, c := range cs {
		f := testdir + c.f
		err := DirIsEmpty(f)
		if !errors.Is(err, c.e) {
			t.Errorf("#%d DirIsEmpty(%q) = %v, want %v", i, c.f, err, c.e)
		}
	}
}

func TestCopyFile(t *testing.T) {
	srcdir := "TestCopyFile-" + strconv.Itoa(rand.Int())
	dstdir := "TestCopyFile-" + strconv.Itoa(rand.Int())
	defer RemoveAll(srcdir)
	defer RemoveAll(dstdir)

	cs := []struct {
		s, d string
	}{
		{"1.txt", "1.txt"},
	}

	for i, c := range cs {
		sf := filepath.Join(srcdir, c.s)
		df := filepath.Join(dstdir, c.d)

		MkdirAll(srcdir, FileMode(0777))

		ss := strings.Repeat("a", (i+1)*10)
		err := WriteString(sf, ss, FileMode(0666))
		if err != nil {
			t.Fatalf("#%d: %v", i, err)
		}

		err = Chmod(sf, 0400)
		if err != nil {
			t.Fatalf("#%d: %v", i, err)
		}

		err = CopyFile(sf, df)
		if err != nil {
			t.Fatalf("#%d: %v", i, err)
		}

		ds, err := ReadString(df)
		if err != nil {
			t.Fatalf("#%d: %v", i, err)
		}
		if ss != ds {
			t.Fatalf("#%d: \n GOT: %v\nWANT: %v\n", i, ds, ss)
		}
	}
}
