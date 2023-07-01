package fdk

import (
	"bytes"
	"net/http"
	"time"

	"github.com/askasoft/pango/fsu"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
)

type FDK struct {
	Domain   string
	Apikey   string
	Username string
	Password string

	Transport http.RoundTripper
	Timeout   time.Duration
	Logger    log.Logger

	RetryOnRateLimited int
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
	fdk.authenticate(req)

	client := &http.Client{
		Transport: fdk.Transport,
		Timeout:   fdk.Timeout,
	}

	return Call(client, req, fdk.RetryOnRateLimited, fdk.Logger)
}

func (fdk *FDK) doCall(req *http.Request, result any) error {
	res, err := fdk.call(req)
	if err != nil {
		return err
	}

	return DecodeResponse(res, result)
}

func (fdk *FDK) DoDownload(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := fdk.call(req)
	if err != nil {
		return nil, err
	}

	defer iox.DrainAndClose(res.Body)
	buf := &bytes.Buffer{}
	_, err = iox.Copy(buf, res.Body)
	return buf.Bytes(), err
}

func (fdk *FDK) DoSave(url string, filename string) error {
	bs, err := fdk.DoDownload(url)
	if err != nil {
		return err
	}
	return fsu.WriteFile(filename, bs, fsu.FileMode(0660))
}

func (fdk *FDK) DoGet(url string, result any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	return fdk.doCall(req, result)
}

func (fdk *FDK) DoList(url string, lo ListOption, ap any) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	if !lo.IsNil() {
		q := lo.Values()
		req.URL.RawQuery = q.Encode()
	}

	res, err := fdk.call(req)
	if err != nil {
		return false, err
	}

	if err := DecodeResponse(res, ap); err != nil {
		return false, err
	}

	next := res.Header.Get("Link") != ""
	return next, nil
}

func (fdk *FDK) DoPost(url string, source, result any) error {
	buf, ct, err := BuildRequest(source)
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
	buf, ct, err := BuildRequest(source)
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
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return fdk.doCall(req, nil)
}

// SleepForRetry if err is RateLimitedError, sleep Retry-After and return true
func (fdk *FDK) SleepForRetry(err error) bool {
	return SleepForRetry(err, fdk.Logger)
}

func (fdk *FDK) SleepAndRetry(api func() error, maxRetry int) (err error) {
	return SleepAndRetry(api, maxRetry, fdk.Logger)
}
