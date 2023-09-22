package fdk

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/net/httpx"
	"github.com/askasoft/pango/num"
)

const (
	contentTypeJSON = `application/json; charset="utf-8"`
)

func ToString(o any) string {
	bs, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err.Error()
	}
	return bye.UnsafeString(bs)
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

func call(client *http.Client, req *http.Request, logger log.Logger) (res *http.Response, err error) {
	if logger != nil {
		logger.Debugf("%s %s", req.Method, req.URL)
	}

	rid := log.TraceHttpRequest(logger, req)

	res, err = client.Do(req)
	if err == nil {
		log.TraceHttpResponse(logger, res, rid)

		if res.StatusCode == http.StatusTooManyRequests {
			s := res.Header.Get("Retry-After")
			n := num.Atoi(s)
			if n <= 0 {
				n = 30 // invalid number, default to 30s
			}
			iox.DrainAndClose(res.Body)
			return res, &RateLimitedError{StatusCode: res.StatusCode, RetryAfter: time.Second * time.Duration(n)}
		}
	}

	return res, err
}

func decodeResponse(res *http.Response, obj any) error {
	defer iox.DrainAndClose(res.Body)

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated || res.StatusCode == http.StatusNoContent {
		if obj != nil {
			return decoder.Decode(obj)
		}
		return nil
	}

	er := &ErrorResult{StatusCode: res.StatusCode, Status: res.Status}
	if res.StatusCode != http.StatusNotFound {
		_ = decoder.Decode(er)
	}
	return er
}

func copyResponse(res *http.Response) ([]byte, error) {
	defer iox.DrainAndClose(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, &ErrorResult{StatusCode: res.StatusCode, Status: res.Status}
	}

	buf := &bytes.Buffer{}
	_, err := iox.Copy(buf, res.Body)
	return buf.Bytes(), err
}

func saveResponse(res *http.Response, path string) error {
	defer iox.DrainAndClose(res.Body)

	if res.StatusCode != http.StatusOK {
		return &ErrorResult{StatusCode: res.StatusCode, Status: res.Status}
	}

	return fsu.WriteReader(path, res.Body, fsu.FileMode(0660))
}
