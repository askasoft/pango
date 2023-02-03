package fdk

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/pandafw/pango/bye"
	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/log"
	"github.com/pandafw/pango/net/httpx"
)

const (
	contentTypeJSON = `application/json; charset="utf-8"`
	logTimeFormat   = "2006-01-02T15:04:05.000"
)

func LogRequest(log log.Logger, req *http.Request) (rid uint64) {
	if log != nil && log.IsTraceEnabled() {
		rid = rand.Uint64() //nolint: gosec
		bs, _ := httputil.DumpRequestOut(req, true)
		log.Tracef(">>>>>>>> %s %016x >>>>>>>>", time.Now().Format(logTimeFormat), rid)
		log.Trace(bye.UnsafeString(bs))
	}
	return
}

func LogResponse(log log.Logger, res *http.Response, rid uint64) {
	if log != nil && log.IsTraceEnabled() {
		bs, _ := httputil.DumpResponse(res, true)
		log.Tracef("<<<<<<<< %s %016x <<<<<<<<", time.Now().Format(logTimeFormat), rid)
		log.Trace(bye.UnsafeString(bs))
	}
}

func ToString(o any) string {
	bs, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(bs)
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

// BuildRequest build a fdk request, returns buffer, contentType, error
func BuildRequest(a any) (io.Reader, string, error) {
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

func Call(client *http.Client, req *http.Request, retryOnRateLimited int, log log.Logger) (res *http.Response, err error) {
	err = SleepAndRetry(func() error {
		if log != nil {
			log.Infof("%s %s", req.Method, req.URL)
		}

		rid := LogRequest(log, req)

		res, err = client.Do(req)
		if err != nil {
			return err
		}
		LogResponse(log, res, rid)

		if res.StatusCode == http.StatusTooManyRequests {
			s := res.Header.Get("Retry-After")
			n, _ := strconv.Atoi(s)
			if n <= 0 {
				n = 60 // invalid number, default to 60s
			}
			iox.DrainAndClose(res.Body)
			return &RateLimitedError{StatusCode: res.StatusCode, RetryAfter: n}
		}

		return err
	}, retryOnRateLimited, log)

	return
}

func DecodeResponse(res *http.Response, obj any) error {
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
		if err := decoder.Decode(er); err != nil {
			return err
		}
	}
	return er
}

// SleepForRetry if err is RateLimitedError, sleep Retry-After and return true
func SleepForRetry(err error, log log.Logger) bool {
	if err != nil {
		if rle, ok := err.(*RateLimitedError); ok { //nolint: errorlint
			if log != nil {
				log.Warnf("Sleep %d seconds for API Rate Limited", rle.RetryAfter)
			}
			time.Sleep(time.Duration(rle.RetryAfter) * time.Second)
			return true
		}
	}
	return false
}

func SleepAndRetry(api func() error, maxRetry int, log log.Logger) (err error) {
	for i := 0; ; i++ {
		err = api()
		if i >= maxRetry {
			break
		}
		if !SleepForRetry(err, log) {
			break
		}
	}
	return err
}