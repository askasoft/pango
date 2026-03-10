package xmlrpc

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/askasoft/pango/iox"
)

type Client struct {
	Endpoint     string
	Timeout      time.Duration
	Transport    http.RoundTripper
	Authenticate func(req *http.Request)
}

func (c *Client) Call(ctx context.Context, method string, result any, params ...any) error {
	buf := &bytes.Buffer{}
	if err := EncodeMethodCall(buf, method, params...); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Endpoint, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", `text/xml; charset="utf-8"`)

	if a := c.Authenticate; a != nil {
		a(req)
	}

	res, err := c.call(req)
	if err != nil {
		return err
	}
	defer iox.DrainAndClose(res.Body)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return &HTTPError{
			Method:     method,
			StatusCode: res.StatusCode,
			Status:     res.Status,
		}
	}

	if result == nil {
		result = new(any)
	}

	err = Decode(res.Body, result)
	if err != nil {
		if fe, ok := AsFaultError(err); ok {
			fe.Method = method
		}
	}
	return err
}

func (c *Client) call(req *http.Request) (*http.Response, error) {
	hc := http.Client{
		Transport: c.Transport,
		Timeout:   c.Timeout,
	}
	return hc.Do(req)
}
