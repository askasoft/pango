package httpx

import (
	"io"
	"mime/multipart"
	"os"
	"path"
)

// SaveMultipartFile save multipart file to the specific local file 'dst'.
func SaveMultipartFile(file *multipart.FileHeader, dst string) error {
	dir := path.Dir(dst)
	if err := os.MkdirAll(dir, os.FileMode(0770)); err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
