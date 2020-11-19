package ios

import (
	"fmt"
	"os"
)

// FileExists check if file exists
func FileExists(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return fmt.Errorf("%q is directory", path)
	}
	return nil
}

// FileReader a file reader
type FileReader struct {
	Path string
	file *os.File
}

// Read implements io.Reader
func (fr *FileReader) Read(p []byte) (n int, err error) {
	if fr.file == nil {
		file, err := os.Open(fr.Path)
		if err != nil {
			return 0, err
		}
		fr.file = file
	}
	return fr.file.Read(p)
}

// Close implements io.Close
func (fr *FileReader) Close() error {
	if fr.file == nil {
		return nil
	}

	err := fr.file.Close()
	fr.file = nil
	return err
}
