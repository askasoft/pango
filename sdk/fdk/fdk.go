package fdk

import (
	"context"
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

	MaxRetries  int
	RetryAfter  time.Duration
	ShouldRetry func(error) bool // default retry on not canceled error or (status = 429 || (status >= 500 && status <= 599))
}

// Endpoint formats endpoint url
func (fdk *FDK) Endpoint(format string, a ...any) string {
	return "https://" + fdk.Domain + "/api/v2" + fmt.Sprintf(format, a...)
}

func (fdk *FDK) RetryForError(ctx context.Context, api func() error) (err error) {
	return sdk.RetryForError(ctx, api, fdk.MaxRetries, fdk.Logger)
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
		fsr := fdk.ShouldRetry
		if fsr == nil {
			fsr = shouldRetry
		}
		if fsr(err) {
			err = sdk.NewRetryError(err, fdk.RetryAfter)
		}
		return res, err
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

func (fdk *FDK) decodeResponse(res *http.Response, result any) error {
	defer iox.DrainAndClose(res.Body)

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated || res.StatusCode == http.StatusNoContent {
		if result != nil {
			return decoder.Decode(result)
		}
		return nil
	}

	re := newResultError(res)
	if res.StatusCode != http.StatusNotFound {
		_ = decoder.Decode(re)
	}

	fsr := fdk.ShouldRetry
	if fsr == nil {
		fsr = shouldRetry
	}

	if fsr(re) {
		s := res.Header.Get("Retry-After")
		n := num.Atoi(s)
		if n > 0 {
			re.RetryAfter = time.Second * time.Duration(n)
		} else {
			re.RetryAfter = fdk.RetryAfter
		}
	}

	return re
}

func (fdk *FDK) DoGet(ctx context.Context, url string, result any) error {
	return fdk.RetryForError(ctx, func() error {
		return fdk.doGet(ctx, url, result)
	})
}

func (fdk *FDK) doGet(ctx context.Context, url string, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	return fdk.doCall(req, result)
}

func (fdk *FDK) DoList(ctx context.Context, url string, lo ListOption, ap any) (next bool, err error) {
	err = fdk.RetryForError(ctx, func() error {
		next, err = fdk.doList(ctx, url, lo, ap)
		return err
	})
	return
}

func (fdk *FDK) doList(ctx context.Context, url string, lo ListOption, result any) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

func (fdk *FDK) DoPost(ctx context.Context, url string, source, result any) error {
	return fdk.RetryForError(ctx, func() error {
		return fdk.doPost(ctx, url, source, result)
	})
}

func (fdk *FDK) doPost(ctx context.Context, url string, source, result any) error {
	buf, ct, err := buildRequest(source)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return err
	}
	if ct != "" {
		req.Header.Set("Content-Type", ct)
	}

	return fdk.doCall(req, result)
}

func (fdk *FDK) DoPut(ctx context.Context, url string, source, result any) error {
	return fdk.RetryForError(ctx, func() error {
		return fdk.doPut(ctx, url, source, result)
	})
}

func (fdk *FDK) doPut(ctx context.Context, url string, source, result any) error {
	buf, ct, err := buildRequest(source)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, buf)
	if err != nil {
		return err
	}
	if ct != "" {
		req.Header.Set("Content-Type", ct)
	}

	return fdk.doCall(req, result)
}

func (fdk *FDK) DoDelete(ctx context.Context, url string) error {
	return fdk.RetryForError(ctx, func() error {
		return fdk.doDelete(ctx, url)
	})
}

func (fdk *FDK) doDelete(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	return fdk.doCall(req, nil)
}

func (fdk *FDK) DoDownload(ctx context.Context, url string) (buf []byte, err error) {
	err = fdk.RetryForError(ctx, func() error {
		buf, err = fdk.doDownload(ctx, url)
		return err
	})
	return
}

func (fdk *FDK) doDownload(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := fdk.authAndCall(req)
	if err != nil {
		return nil, err
	}

	return copyResponse(res)
}

func (fdk *FDK) DoSaveFile(ctx context.Context, url string, path string) error {
	return fdk.RetryForError(ctx, func() error {
		return fdk.doSaveFile(ctx, url, path)
	})
}

func (fdk *FDK) doSaveFile(ctx context.Context, url string, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := fdk.authAndCall(req)
	if err != nil {
		return err
	}

	return saveResponse(res, path)
}

func (fdk *FDK) DoDownloadNoAuth(ctx context.Context, url string) (buf []byte, err error) {
	err = fdk.RetryForError(ctx, func() error {
		buf, err = fdk.doDownloadNoAuth(ctx, url)
		return err
	})
	return
}

func (fdk *FDK) doDownloadNoAuth(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := fdk.call(req)
	if err != nil {
		return nil, err
	}

	return copyResponse(res)
}

func (fdk *FDK) DoSaveFileNoAuth(ctx context.Context, url string, path string) error {
	return fdk.RetryForError(ctx, func() error {
		return fdk.doSaveFileNoAuth(ctx, url, path)
	})
}

func (fdk *FDK) doSaveFileNoAuth(ctx context.Context, url string, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := fdk.call(req)
	if err != nil {
		return err
	}

	return saveResponse(res, path)
}
