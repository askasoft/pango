package fdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/log/httplog"
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/sdk"
)

type FDK struct {
	Domain   string
	Apikey   string
	Username string
	Password string

	Transport http.RoundTripper
	Timeout   time.Duration
	Logger    log.Logger

	MaxRetries    int
	RetryAfter    time.Duration
	AbortOnRetry  func() error
	AbortInterval time.Duration
}

// Endpoint formats endpoint url
func (fdk *FDK) Endpoint(format string, a ...any) string {
	return "https://" + fdk.Domain + "/api/v2" + fmt.Sprintf(format, a...)
}

func (fdk *FDK) RetryForError(api func() error) (err error) {
	return sdk.RetryForError(api, fdk.MaxRetries, fdk.AbortOnRetry, fdk.AbortInterval, fdk.Logger)
}

func (fdk *FDK) authenticate(req *http.Request) {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentTypeJSON)
	}

	if fdk.Apikey != "" {
		req.SetBasicAuth(fdk.Apikey, "X")
	} else {
		req.SetBasicAuth(fdk.Username, fdk.Password)
	}
}

func (fdk *FDK) call(req *http.Request) (res *http.Response, err error) {
	client := &http.Client{
		Transport: fdk.Transport,
		Timeout:   fdk.Timeout,
	}

	if log := fdk.Logger; log != nil {
		log.Debugf("%s %s", req.Method, req.URL)
	}

	rid := httplog.TraceHttpRequest(fdk.Logger, req)

	res, err = client.Do(req)
	if err != nil {
		return res, sdk.NewNetError(err, fdk.RetryAfter)
	}

	httplog.TraceHttpResponse(fdk.Logger, res, rid)
	return res, nil
}

func (fdk *FDK) authAndCall(req *http.Request) (res *http.Response, err error) {
	fdk.authenticate(req)
	return fdk.call(req)
}

func (fdk *FDK) doCall(req *http.Request, result any) error {
	res, err := fdk.authAndCall(req)
	if err != nil {
		return err
	}

	return fdk.decodeResponse(res, result)
}

func (fdk *FDK) decodeResponse(res *http.Response, obj any) error {
	defer iox.DrainAndClose(res.Body)

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated || res.StatusCode == http.StatusNoContent {
		if obj != nil {
			return decoder.Decode(obj)
		}
		return nil
	}

	re := &ResultError{StatusCode: res.StatusCode, Status: res.Status}
	if res.StatusCode != http.StatusNotFound {
		_ = decoder.Decode(re)
	}

	switch {
	case res.StatusCode == http.StatusTooManyRequests:
		s := res.Header.Get("Retry-After")
		n := num.Atoi(s)
		if n > 0 {
			re.RetryAfter = time.Second * time.Duration(n)
		} else {
			re.RetryAfter = fdk.RetryAfter
		}
	case res.StatusCode >= 500 && res.StatusCode <= 599:
		re.RetryAfter = fdk.RetryAfter
	}

	return re
}

func (fdk *FDK) DoGet(url string, result any) error {
	return fdk.RetryForError(func() error {
		return fdk.doGet(url, result)
	})
}

func (fdk *FDK) doGet(url string, result any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	return fdk.doCall(req, result)
}

func (fdk *FDK) DoList(url string, lo ListOption, ap any) (next bool, err error) {
	err = fdk.RetryForError(func() error {
		next, err = fdk.doList(url, lo, ap)
		return err
	})
	return
}

func (fdk *FDK) doList(url string, lo ListOption, result any) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	if lo != nil && !lo.IsNil() {
		q := lo.Values()
		req.URL.RawQuery = q.Encode()
	}

	res, err := fdk.authAndCall(req)
	if err != nil {
		return false, err
	}

	err = fdk.decodeResponse(res, result)
	if err != nil {
		return false, err
	}

	next := res.Header.Get("Link") != ""
	return next, nil
}

func (fdk *FDK) DoPost(url string, source, result any) error {
	return fdk.RetryForError(func() error {
		return fdk.doPost(url, source, result)
	})
}

func (fdk *FDK) doPost(url string, source, result any) error {
	buf, ct, err := buildRequest(source)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return err
	}
	if ct != "" {
		req.Header.Set("Content-Type", ct)
	}

	return fdk.doCall(req, result)
}

func (fdk *FDK) DoPut(url string, source, result any) error {
	return fdk.RetryForError(func() error {
		return fdk.doPut(url, source, result)
	})
}

func (fdk *FDK) doPut(url string, source, result any) error {
	buf, ct, err := buildRequest(source)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, buf)
	if err != nil {
		return err
	}
	if ct != "" {
		req.Header.Set("Content-Type", ct)
	}

	return fdk.doCall(req, result)
}

func (fdk *FDK) DoDelete(url string) error {
	return fdk.RetryForError(func() error {
		return fdk.doDelete(url)
	})
}

func (fdk *FDK) doDelete(url string) error {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return fdk.doCall(req, nil)
}

func (fdk *FDK) DoDownload(url string) (buf []byte, err error) {
	err = fdk.RetryForError(func() error {
		buf, err = fdk.doDownload(url)
		return err
	})
	return
}

func (fdk *FDK) doDownload(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := fdk.authAndCall(req)
	if err != nil {
		return nil, err
	}

	return copyResponse(res)
}

func (fdk *FDK) DoSaveFile(url string, path string) error {
	return fdk.RetryForError(func() error {
		return fdk.doSaveFile(url, path)
	})
}

func (fdk *FDK) doSaveFile(url string, path string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := fdk.authAndCall(req)
	if err != nil {
		return err
	}

	return saveResponse(res, path)
}

func (fdk *FDK) DoDownloadNoAuth(url string) (buf []byte, err error) {
	err = fdk.RetryForError(func() error {
		buf, err = fdk.doDownloadNoAuth(url)
		return err
	})
	return
}

func (fdk *FDK) doDownloadNoAuth(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := fdk.call(req)
	if err != nil {
		return nil, err
	}

	return copyResponse(res)
}

func (fdk *FDK) DoSaveFileNoAuth(url string, path string) error {
	return fdk.RetryForError(func() error {
		return fdk.doSaveFileNoAuth(url, path)
	})
}

func (fdk *FDK) doSaveFileNoAuth(url string, path string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := fdk.call(req)
	if err != nil {
		return err
	}

	return saveResponse(res, path)
}
