package xmlrpc

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/log/httplog"
	"github.com/askasoft/pango/ret"
)

type Client struct {
	Endpoint string

	Transport http.RoundTripper
	Timeout   time.Duration
	Logger    log.Logger

	MaxRetries  int
	RetryAfter  time.Duration
	ShouldRetry func(error) bool // default retry on not canceled error or (status = 429 || (status >= 500 && status <= 599))

	Authenticate func(req *http.Request)
}

func (c *Client) Call(ctx context.Context, method string, result any, params ...any) error {
	return c.doPostWithRetry(ctx, method, result, params)
}

func (c *Client) doPostWithRetry(ctx context.Context, method string, result any, params []any) error {
	return ret.RetryForError(ctx, func() error {
		return c.doPost(ctx, method, result, params)
	}, c.MaxRetries, c.Logger)
}

func (c *Client) doPost(ctx context.Context, method string, result any, params []any) error {
	buf := &bytes.Buffer{}
	if err := EncodeMethodCall(buf, method, params...); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Endpoint, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", `text/xml; charset="utf-8"`)

	return c.doCall(req, method, result)
}

func (c *Client) doCall(req *http.Request, method string, result any) error {
	sr := c.ShouldRetry
	if sr == nil {
		sr = shouldRetry
	}

	if a := c.Authenticate; a != nil {
		c.Authenticate(req)
	}

	res, err := c.call(req)
	if err != nil {
		if sr(err) {
			return ret.NewRetryError(err, c.RetryAfter)
		}
		return err
	}
	defer iox.DrainAndClose(res.Body)

	if result == nil {
		result = new(any)
	}

	decoder := NewDecoder(res.Body)
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		err = decoder.Decode(result)
		if err != nil {
			if fe, ok := AsFaultError(err); ok {
				fe.Method = method
				fe.StatusCode = res.StatusCode
				if sr(fe) {
					fe.RetryAfter = c.RetryAfter
				}
			}
		}
		return err
	}

	fe := &FaultError{
		Method:     method,
		StatusCode: res.StatusCode,
	}

	err = decoder.Decode(fe)
	if fe2, ok := AsFaultError(err); ok {
		fe.FaultCode = fe2.FaultCode
		fe.FaultString = fe2.FaultString
	}

	if sr(fe) {
		fe.RetryAfter = c.RetryAfter
	}
	return fe
}

func (c *Client) call(req *http.Request) (res *http.Response, err error) {
	client := &http.Client{
		Transport: c.Transport,
		Timeout:   c.Timeout,
	}

	if log := c.Logger; log != nil {
		log.Debugf("%s %s", req.Method, req.URL)
	}

	rid := httplog.TraceHttpRequest(c.Logger, req)

	res, err = client.Do(req)
	if err != nil {
		return res, err
	}

	httplog.TraceHttpResponse(c.Logger, res, rid)
	return res, nil
}
