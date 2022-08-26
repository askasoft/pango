package iox

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pandafw/pango/bye"
	"github.com/pandafw/pango/str"
)

// Discard is an io.Writer on which all Write calls succeed
// without doing anything.
var Discard = io.Discard

// ReadAll reads from r until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.
func ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

// ReadFile reads the file named by filename and returns the contents.
// A successful call returns err == nil, not err == EOF. Because ReadFile
// reads the whole file, it does not treat an EOF from Read as an error
// to be reported.
func ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// ReadString reads the file named by filename and returns the contents as string.
// A successful call returns err == nil, not err == EOF. Because ReadString
// reads the whole file, it does not treat an EOF from Read as an error
// to be reported.
func ReadString(filename string) (string, error) {
	bs, err := os.ReadFile(filename)
	return bye.UnsafeString(bs), err
}

// WriteFile writes data to a file named by filename.
// If the file does not exist, WriteFile creates it with permissions perm
// (before umask); otherwise WriteFile truncates it before writing, without changing permissions.
func WriteFile(filename string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

// WriteString writes string data to a file named by filename.
// If the file does not exist, WriteString creates it with permissions perm
// (before umask); otherwise WriteString truncates it before writing, without changing permissions.
func WriteString(filename string, data string, perm fs.FileMode) error {
	return os.WriteFile(filename, str.UnsafeBytes(data), perm)
}

// ReadDir reads the named directory,
// returning all its directory entries sorted by filename.
// If an error occurs reading the directory,
// ReadDir returns the entries it was able to read before the error,
// along with the error.
func ReadDir(dirname string) ([]os.DirEntry, error) {
	return os.ReadDir(dirname)
}

// NopCloser returns a ReadCloser with a no-op Close method wrapping
// the provided Reader r.
func NopCloser(r io.Reader) io.ReadCloser {
	return io.NopCloser(r)
}

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
	err = os.MkdirAll(dd, ss.Mode().Perm())
	if err != nil {
		return err
	}

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

// SkipBOM skip bom and return a reader
func SkipBOM(r io.Reader) (io.Reader, error) {
	br := bufio.NewReader(r)
	c, _, err := br.ReadRune()
	if err != nil {
		return br, err
	}
	if c != '\uFEFF' {
		// Not a BOM -- put the rune back
		err = br.UnreadRune()
	}
	return br, err
}

// Drain drain the reader
func Drain(r io.Reader) {
	io.Copy(Discard, r) //nolint: errcheck
}

// DrainAndClose drain and close the reader
func DrainAndClose(r io.ReadCloser) {
	Drain(r)
	r.Close()
}
