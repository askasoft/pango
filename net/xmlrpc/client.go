package xmlrpc

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/log/httplog"
)

type Client struct {
	Endpoint string

	Logger       log.Logger
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
		c.Authenticate(req)
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
