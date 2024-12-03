package fdk

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/net/httpx"
	"github.com/askasoft/pango/str"
)

const (
	contentTypeJSON = `application/json; charset="utf-8"`
)

func ToJSONIndent(o any) string {
	if o == nil {
		return ""
	}

	bs, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err.Error()
	}
	return str.UnsafeString(bs)
}

func ToJSON(o any) string {
	if o == nil {
		return ""
	}

	bs, err := json.Marshal(o)
	if err != nil {
		return err.Error()
	}
	return str.UnsafeString(bs)
}

func addMultipartValues(mw *httpx.MultipartWriter, vs Values) error {
	return mw.WriteFields(url.Values(vs))
}

func addMultipartFiles(mw *httpx.MultipartWriter, fs Files) (err error) {
	for _, f := range fs {
		if len(f.Data()) > 0 {
			err = mw.WriteFileData(f.Field(), f.File(), f.Data())
		} else {
			err = mw.WriteFile(f.Field(), f.File())
		}
		if err != nil {
			return
		}
	}
	return
}

// buildRequest build a fdk request, returns buffer, contentType, error
func buildRequest(a any) (io.Reader, string, error) {
	if a == nil {
		return nil, "", nil
	}

	if wf, ok := a.(WithFiles); ok {
		fs := wf.Files()
		if len(fs) > 0 {
			vs := wf.Values()
			return buildFileRequest(vs, fs)
		}
	}
	return buildJsonRequest(a)
}

func buildFileRequest(vs Values, fs Files) (io.Reader, string, error) {
	buf := &bytes.Buffer{}
	mw := httpx.NewMultipartWriter(buf)

	contentType := mw.FormDataContentType()

	if err := addMultipartValues(mw, vs); err != nil {
		return nil, "", err
	}
	if err := addMultipartFiles(mw, fs); err != nil {
		return nil, "", err
	}
	if err := mw.Close(); err != nil {
		return nil, "", err
	}

	return buf, contentType, nil
}

func buildJsonRequest(a any) (io.Reader, string, error) {
	body, err := json.Marshal(a)
	if err != nil {
		return nil, "", err
	}

	buf := bytes.NewReader(body)
	return buf, contentTypeJSON, nil
}

func copyResponse(res *http.Response) ([]byte, error) {
	defer iox.DrainAndClose(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, &ResultError{StatusCode: res.StatusCode, Status: res.Status}
	}

	buf := &bytes.Buffer{}
	_, err := iox.Copy(buf, res.Body)
	return buf.Bytes(), err
}

func saveResponse(res *http.Response, path string) error {
	defer iox.DrainAndClose(res.Body)

	if res.StatusCode != http.StatusOK {
		return &ResultError{StatusCode: res.StatusCode, Status: res.Status}
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.FileMode(0770)); err != nil {
		return err
	}

	return fsu.WriteReader(path, res.Body, fsu.FileMode(0660))
}
