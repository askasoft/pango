package fsu

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/str"
)

// A DirEntry is an entry read from a directory
// (using the ReadDir function or a File's ReadDir method).
type DirEntry = fs.DirEntry

// A FileInfo describes a file and is returned by Stat and Lstat.
type FileInfo = fs.FileInfo

// A FileMode represents a file's mode and permission bits.
// The bits have the same definition on all systems, so that
// information about files can be moved from one system
// to another portably. Not all bits apply to all systems.
// The only required bit is ModeDir for directories.
type FileMode = fs.FileMode

var (
	ErrInvalid    = fs.ErrInvalid    // "invalid argument"
	ErrPermission = fs.ErrPermission // "permission denied"
	ErrExist      = fs.ErrExist      // "file already exists"
	ErrNotExist   = fs.ErrNotExist   // "file does not exist"
	ErrClosed     = fs.ErrClosed     // "file already closed"
	ErrNotDir     = syscall.ENOTDIR
	ErrIsDir      = errors.New("file is directory")
	ErrNotEmpty   = errors.New("directory not empty")
)

// CopyFile copy src file to des file
func CopyFile(src string, dst string) error {
	ss, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !ss.Mode().IsRegular() {
		return fmt.Errorf("'%s' is not a regular file", src)
	}

	dd := filepath.Dir(dst)
	err = MkdirAll(dd, FileMode(0770))
	if err != nil {
		return err
	}

	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, FileMode(0660))
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	return err
}

// DirIsEmpty check if the directory dir contains sub folders or files
func DirIsEmpty(dir string) error {
	if err := DirExists(dir); err != nil {
		return err
	}

	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err != nil {
		return err
	}
	return ErrNotEmpty
}

// DirExists check if the directory dir exists
// return ErrIsNotDir if dir is not directory
func DirExists(dir string) error {
	fi, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return &fs.PathError{Op: "DirExists", Path: dir, Err: ErrNotDir}
	}
	return nil
}

// FileExists check if the file exists
// return ErrIsDir if file is directory
func FileExists(file string) error {
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return &fs.PathError{Op: "FileExists", Path: file, Err: ErrIsDir}
	}
	return nil
}

// FileSize get the file size
func FileSize(file string) (int64, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	if fi.IsDir() {
		return 0, &fs.PathError{Op: "FileSize", Path: file, Err: ErrIsDir}
	}
	return fi.Size(), nil
}

// Mkdir creates a new directory with the specified name and permission
// bits (before umask).
// If there is an error, it will be of type *PathError.
func Mkdir(name string, perm FileMode) error {
	return os.Mkdir(name, perm)
}

// MkdirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission bits perm (before umask) are used for all
// directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing
// and returns nil.
func MkdirAll(path string, perm FileMode) error {
	return os.MkdirAll(path, perm)
}

// Remove removes the named file or directory.
// If there is an error, it will be of type *PathError.
func Remove(name string) error {
	return os.Remove(name)
}

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters. If the path does not exist, RemoveAll
// returns nil (no error).
// If there is an error, it will be of type *PathError.
func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Remove removes all files return the by filepath.Glob().
func RemoveGlob(path string) error {
	fns, err := filepath.Glob(path)
	if err != nil {
		return err
	}
	for _, fn := range fns {
		err := os.Remove(fn)
		if err != nil {
			return err
		}
	}
	return nil
}

// Remove removes all files and children return the by filepath.Glob().
func RemoveGlobAll(path string) error {
	fns, err := filepath.Glob(path)
	if err != nil {
		return err
	}
	for _, fn := range fns {
		err := os.RemoveAll(fn)
		if err != nil {
			return err
		}
	}
	return nil
}

// Chdir changes the current working directory to the named directory.
// If there is an error, it will be of type *PathError.
func Chdir(dir string) error {
	return os.Chdir(dir)
}

// Chmod changes the mode of the named file to mode.
// If the file is a symbolic link, it changes the mode of the link's target.
// If there is an error, it will be of type *PathError.
//
// A different subset of the mode bits are used, depending on the
// operating system.
//
// On Unix, the mode's permission bits, ModeSetuid, ModeSetgid, and
// ModeSticky are used.
//
// On Windows, only the 0200 bit (owner writable) of mode is used; it
// controls whether the file's read-only attribute is set or cleared.
// The other bits are currently unused. For compatibility with Go 1.12
// and earlier, use a non-zero mode. Use mode 0400 for a read-only
// file and 0600 for a readable+writable file.
//
// On Plan 9, the mode's permission bits, ModeAppend, ModeExclusive,
// and ModeTemporary are used.
func Chmod(name string, mode FileMode) error {
	return os.Chmod(name, mode)
}

// ReadDir reads the named directory,
// returning all its directory entries sorted by filename.
// If an error occurs reading the directory,
// ReadDir returns the entries it was able to read before the error,
// along with the error.
func ReadDir(dirname string) ([]DirEntry, error) {
	return os.ReadDir(dirname)
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
func WriteFile(filename string, data []byte, perm FileMode) error {
	return os.WriteFile(filename, data, perm)
}

// WriteReader writes reader data to a file named by filename.
// If the file does not exist, WriteReader creates it with permissions perm
// (before umask); otherwise WriteReader truncates it before writing, without changing permissions.
func WriteReader(filename string, src io.Reader, perm FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, src)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

// WriteString writes string data to a file named by filename.
// If the file does not exist, WriteString creates it with permissions perm
// (before umask); otherwise WriteString truncates it before writing, without changing permissions.
func WriteString(filename string, data string, perm FileMode) error {
	return os.WriteFile(filename, str.UnsafeBytes(data), perm)
}
