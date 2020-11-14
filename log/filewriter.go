package log

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// FileWriter implements Writer.
// It writes messages and rotate by file size limit, daily, hourly.
type FileWriter struct {
	Level      int    `json:"level"`      // Level threshold
	File       string `json:"file"`       // Log file name
	Perm       uint32 `json:"perm"`       // Log file permission
	Rotate     bool   `json:"rotate"`     // Rotate log files
	MaxFiles   int    `json:"maxfiles"`   // Max split files
	MaxSize    int64  `json:"maxsize"`    // Rotate at size
	Daily      bool   `json:"daily"`      // Rotate daily
	MaxDays    int    `json:"maxdays"`    // Max daily files
	Hourly     bool   `json:"hourly"`     // Rotate hourly
	MaxHours   int    `json:"maxhours"`   // Max hourly files
	Gzip       bool   `json:"gzip"`       // Compress rotated log files
	FlushLevel int    `json:"flushlevel"` // Flush by log level
	Logfmt     Formatter

	dir      string
	prefix   string
	suffix   string
	file     *os.File
	fileSize int64
	fileNum  int
	openTime time.Time
	openDay  int
	openHour int
	sync.RWMutex
}

// SetLevel set the log level
func (fw *FileWriter) SetLevel(level string) {
	fw.Level = ParseLevel(level)
}

// SetFormat set a log formatter
func (fw *FileWriter) SetFormat(format string) {
	fw.Logfmt = NewTextFormatter(format)
}

// Write write logger message into file.
func (fw *FileWriter) Write(le *Event) {
	if fw.Level < le.Level {
		return
	}

	fw.init()
	if fw.file == nil {
		return
	}

	if fw.Rotate {
		fw.lock(le)
		d := le.When.Day()
		h := le.When.Hour()
		if fw.needRotate(d, h) {
			fw.runlock(le)
			fw.lock(le)
			if fw.needRotate(d, h) {
				fw.rotate(le.When)
			}
			fw.unlock(le)
		} else {
			fw.runlock(le)
		}
	}

	// format msg
	if fw.Logfmt == nil {
		fw.Logfmt = le.Logger.GetFormatter()
	}
	msg := fw.Logfmt.Format(le)

	// write log
	fw.lock(le)
	defer fw.unlock(le)

	n, _ := fw.file.WriteString(msg)
	fw.fileSize += int64(n)

	if fw.FlushLevel <= le.Level {
		fw.file.Sync()
	}
}

// Flush flush file logger.
// there are no buffering messages in file logger in memory.
// flush file means sync file from disk.
func (fw *FileWriter) Flush() {
	if fw.file != nil {
		fw.file.Sync()
	}
}

// Close close the file description, close file writer.
func (fw *FileWriter) Close() {
	if fw.file != nil {
		fw.file.Close()
	}
}

func (fw *FileWriter) init() {
	if fw.file != nil {
		return
	}

	// init dir, prefix, suffix
	if fw.prefix == "" {
		fw.dir = path.Dir(fw.File)
		name := filepath.Base(fw.File)
		fw.suffix = filepath.Ext(name)
		if fw.suffix == "" {
			fw.suffix = ".log"
		}
		fw.prefix = strings.TrimSuffix(name, fw.suffix)
	}

	// init perm
	if fw.Perm == 0 {
		fw.Perm = 0660
	}

	// create dirs
	os.MkdirAll(fw.dir, os.FileMode(fw.Perm))

	// Open the log file
	file, err := os.OpenFile(fw.File, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(fw.Perm))
	if err != nil {
		fmt.Fprintf(os.Stderr, "FileWriter(%q) - OpenFile(): %s\n", fw.File, err)
		return
	}

	// Make sure file perm is user set perm cause of `os.OpenFile` will obey umask
	os.Chmod(fw.File, os.FileMode(fw.Perm))

	fi, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FileWriter(%q) - Stat(): %s\n", fw.File, err)
		return
	}

	// init file info
	fw.fileSize = fi.Size()
	fw.openTime = time.Now()
	fw.openDay = fw.openTime.Day()
	fw.openHour = fw.openTime.Hour()

	fw.file = file
}

func (fw *FileWriter) lock(le *Event) {
	if !le.Logger.IsAsync() {
		fw.Lock()
	}
}

func (fw *FileWriter) unlock(le *Event) {
	if !le.Logger.IsAsync() {
		fw.Unlock()
	}
}

func (fw *FileWriter) rlock(le *Event) {
	if !le.Logger.IsAsync() {
		fw.RLock()
	}
}

func (fw *FileWriter) runlock(le *Event) {
	if !le.Logger.IsAsync() {
		fw.RUnlock()
	}
}

func (fw *FileWriter) needRotate(day int, hour int) bool {
	return (fw.MaxSize > 0 && fw.fileSize >= fw.MaxSize) || (fw.Hourly && hour != fw.openHour) || (fw.Daily && day != fw.openDay)
}

// DoRotate means it need to write file in new file.
// new file name like xx-20130101.log (daily) or xx-001.log (by line or size)
func (fw *FileWriter) rotate(logTime time.Time) {
	path := "" // rotate file name

	date := ""
	if fw.Hourly {
		date = logTime.Format("-2006010215")
		if fw.openHour != logTime.Hour() {
			fw.fileNum = 0
		}
	} else if fw.Daily {
		date = logTime.Format("-20060102")
		if fw.openDay != logTime.Day() {
			fw.fileNum = 0
		}
	}

	var files []string

	pre := filepath.Join(fw.dir, fw.prefix) + date

	// only when one of them be setted, then the file would be splited
	if fw.MaxSize > 0 {
		path, files = fw.nextFile(pre)
	} else {
		path = pre + fw.suffix
		_, err := os.Stat(path)
		if err == nil {
			// timely rotate file exists (normally impossible)
			// find next split file name
			path, files = fw.nextFile(pre)
		}
	}

	// remove old split files
	if len(files) > fw.MaxFiles {
		files := files[:len(files)-fw.MaxFiles]
		go fw.deleteFiles(files)
	}

	// close file before rename
	err := fw.file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FileWriter(%q) - Close(): %s\n", fw.File, err)
		return
	}
	fw.file = nil

	// Rename the file to its new found name
	// even if occurs error,we MUST guarantee to  restart new logger
	err = os.Rename(fw.File, path)
	if err == nil && fw.Gzip {
		go fw.compressFile(path)
	}

	// Open file again
	fw.init()

	// delete outdated rotated files
	if fw.Hourly || fw.Daily {
		go fw.deleteOutdatedFiles()
	}
}

func (fw *FileWriter) nextFile(pre string) (string, []string) {
	var path string
	var files []string
	for fw.fileNum++; ; fw.fileNum++ {
		path = pre + fmt.Sprintf("%03d", fw.fileNum) + fw.suffix
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			if fw.Gzip {
				p := path + ".gz"
				_, err = os.Stat(p)
				if os.IsNotExist(err) {
					break
				}
			} else {
				break
			}
		}
		if fw.MaxFiles > 0 {
			files = append(files, path)
		}
	}
	return path, files
}

func (fw *FileWriter) compressFile(src string) {
	dst := src + ".gz"

	f, err := os.Open(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FileWriter(%q) - Open(%q): %s\n", fw.File, src, err)
		return
	}
	defer f.Close()

	// If this file already exists, we presume it was created by
	// a previous attempt to compress the log file.
	gzf, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(fw.Perm))
	if err != nil {
		fmt.Fprintf(os.Stderr, "FileWriter(%q) - OpenFile(%q): %s\n", fw.File, dst, err)
		return
	}
	defer gzf.Close()

	gz := gzip.NewWriter(gzf)

	if _, err := io.Copy(gz, f); err != nil {
		fmt.Fprintf(os.Stderr, "FileWriter(%q) - gzip(%q): %s\n", fw.File, dst, err)
		return
	}
	if err := gz.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "FileWriter(%q) - gzip.Close(%q): %s\n", fw.File, dst, err)
		return
	}
	if err := gzf.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "FileWriter(%q) - Close(%q): %s\n", fw.File, dst, err)
		return
	}

	f.Close()
	if err := os.Remove(src); err != nil {
		fmt.Fprintf(os.Stderr, "FileWriter(%q) - Remove(%q): %s\n", fw.File, src, err)
	}
}

func (fw *FileWriter) deleteFiles(files []string) {
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			fmt.Fprintf(os.Stderr, "FileWriter(%q) - Remove(%q): %s\n", fw.File, f, err)
		}
	}
}

func (fw *FileWriter) deleteOutdatedFiles() {
	var due time.Time
	if fw.Hourly {
		due = time.Now().Add(-1 * time.Hour * time.Duration(fw.MaxHours))
	} else {
		due = time.Now().Add(-24 * time.Hour * time.Duration(fw.MaxDays))
	}

	dir := filepath.Dir(fw.File)
	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FileWriter(%q) - Open(%q): %s\n", fw.File, dir, err)
		return
	}
	defer f.Close()

	fis, err := f.Readdir(-1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FileWriter(%q) - Readdir(%q): %s\n", fw.File, dir, err)
		return
	}

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}

		if fi.ModTime().Before(due) {
			name := filepath.Base(fi.Name())
			if strings.HasPrefix(name, fw.prefix) {
				if err := os.Remove(fi.Name()); err != nil {
					fmt.Fprintf(os.Stderr, "FileWriter(%q) - Remove(%q): %s\n", fw.File, fi.Name(), err)
				}
			}
		}
	}
}
