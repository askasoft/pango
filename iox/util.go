package iox

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// DirExists check if the directory dir exists
func DirExists(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("%q is not a directory", dir)
	}
	return nil
}

// FileExists check if the file exists
func FileExists(file string) error {
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return fmt.Errorf("%q is directory", file)
	}
	return nil
}

// CopyFile copy src file to des file
func CopyFile(src string, dst string) error {
	ss, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !ss.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	dd := filepath.Dir(dst)
	os.MkdirAll(dd, ss.Mode().Perm())

	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, ss.Mode().Perm())
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	return err
}
