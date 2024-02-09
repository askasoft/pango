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
