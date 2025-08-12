package fdk

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/askasoft/pango/doc/jsonx"
	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/net/httpx"
)

const (
	contentTypeJSON = `application/json; charset="utf-8"`
)

type CustomRequest interface {
	RequestBody() (io.Reader, string, error)
}

func toString(o any) string {
	return jsonx.Prettify(o)
}

func addMultipartValues(mw *httpx.MultipartWriter, vs Values) error {
	return mw.WriteFields(url.Values(vs))
}

func addMultipartFiles(mw *httpx.MultipartWriter, fs Files) (err error) {
	for _, f := range fs {
		if f.Data() == nil {
			err = mw.WriteFile(f.Field(), f.File())
		} else {
			err = mw.WriteFileData(f.Field(), f.File(), f.Data())
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

	if bb, ok := a.(CustomRequest); ok {
		return bb.RequestBody()
	}

	if wf, ok := a.(WithFiles); ok {
		fs := wf.Files()
		if len(fs) > 0 {
			vs := wf.Values()
			return BuildMultipartRequest(vs, fs)
		}
	}

	return BuildJSONRequest(a)
}

func BuildMultipartRequest(vs Values, fs Files) (io.Reader, string, error) {
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

func BuildJSONRequest(a any) (io.Reader, string, error) {
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
		return nil, newResultError(res)
	}

	buf := &bytes.Buffer{}
	_, err := iox.Copy(buf, res.Body)
	return buf.Bytes(), err
}

func saveResponse(res *http.Response, path string) error {
	defer iox.DrainAndClose(res.Body)

	if res.StatusCode != http.StatusOK {
		return newResultError(res)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.FileMode(0770)); err != nil {
		return err
	}

	return fsu.WriteReader(path, res.Body, fsu.FileMode(0660))
}
