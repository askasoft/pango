package freshdesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/net/httpx"
	"github.com/pandafw/pango/num"
)

func toString(o any) string {
	bs, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

type ListOption interface {
	Values() Values
}

type File interface {
	Field() string
	File() string
	Data() []byte
}

type Files []File

type WithFiles interface {
	Files() Files
	Values() Values
}

type Values url.Values

func (vs Values) Map() map[string][]string {
	return (map[string][]string)(vs)
}

func (vs Values) SetBool(name string, value bool) {
	s := "false"
	if value {
		s = "true"
	}
	(url.Values)(vs).Set(name, s)
}

func (vs Values) SetString(name string, value string) {
	if value != "" {
		(url.Values)(vs).Set(name, value)
	}
}

func (vs Values) SetStrings(name string, value []string) {
	name += "[]"
	if len(value) > 0 {
		for _, s := range value {
			(url.Values)(vs).Add(name, s)
		}
	}
}

func (vs Values) SetMap(name string, value map[string]any) {
	if len(value) > 0 {
		for k, v := range value {
			(url.Values)(vs).Add(fmt.Sprintf("%s[%s]", name, k), fmt.Sprint(v))
		}
	}
}

func (vs Values) SetTime(name string, value Time) {
	if !value.IsZero() {
		(url.Values)(vs).Set(name, value.String())
	}
}

func (vs Values) SetTimePtr(name string, value *Time) {
	if value != nil && !value.IsZero() {
		(url.Values)(vs).Set(name, value.String())
	}
}

func (vs Values) SetInt(name string, value int) {
	if value != 0 {
		(url.Values)(vs).Set(name, num.Itoa(value))
	}
}

func (vs Values) SetInt64(name string, value int64) {
	if value != 0 {
		(url.Values)(vs).Set(name, num.Ltoa(value))
	}
}

func (vs Values) Encode() string {
	return (url.Values)(vs).Encode()
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

func buildRequest(a any) (io.Reader, string, error) {
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

func decodeResponse(res *http.Response, status int, obj any) error {
	defer iox.DrainAndClose(res.Body)

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode == status {
		if obj != nil {
			return decoder.Decode(obj)
		}
		return nil
	}

	er := &ErrorResult{StatusCode: res.StatusCode}
	if err := decoder.Decode(er); err != nil {
		return err
	}
	return er
}
