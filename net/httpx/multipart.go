package httpx

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

type MultipartWriter struct {
	*multipart.Writer
}

func NewMultipartWriter(w io.Writer) *MultipartWriter {
	return &MultipartWriter{multipart.NewWriter(w)}
}

func (mw *MultipartWriter) CreateFormFile(fieldname, filename string) (io.Writer, error) {
	mh := make(textproto.MIMEHeader)

	cd := fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(fieldname), escapeQuotes(filepath.Base(filename)))
	mh.Set("Content-Disposition", cd)

	ct := mime.TypeByExtension(filepath.Ext(filename))
	if ct == "" {
		ct = "application/octet-stream"
	}
	mh.Set("Content-Type", ct)

	return mw.CreatePart(mh)
}

// WriteFields calls WriteField and then writes the given fields values.
func (mw *MultipartWriter) WriteFields(fields url.Values) error {
	for k, vs := range fields {
		for _, v := range vs {
			err := mw.WriteField(k, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteFile calls CreateFormFile and then writes the given file.
func (mw *MultipartWriter) WriteFile(fieldname, filename string) error {
	fw, err := mw.CreateFormFile(fieldname, filename)
	if err != nil {
		return err
	}

	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = io.Copy(fw, fp)
	return err
}

// WriteFileData calls CreateFormFile and then writes the given file data.
func (mw *MultipartWriter) WriteFileData(fieldname, filename string, data []byte) error {
	fw, err := mw.CreateFormFile(fieldname, filename)
	if err != nil {
		return err
	}

	_, err = fw.Write(data)
	return err
}

//----------------------------------------------

// SaveMultipartFile save multipart file to the specific local file 'dst'.
func SaveMultipartFile(file *multipart.FileHeader, dst string) error {
	dir := filepath.Dir(dst)
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

// CopyMultipartFile read multipart file to the specific buffer 'dst'.
func CopyMultipartFile(file *multipart.FileHeader, dst io.Writer) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = io.Copy(dst, src)
	return err
}

// ReadMultipartFile read multipart file and return it's content []byte.
func ReadMultipartFile(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	data := make([]byte, file.Size)
	_, err = src.Read(data)
	if err != nil {
		return nil, err
	}

	return data, err
}
