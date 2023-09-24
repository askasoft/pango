package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sdk"
)

const (
	contentTypeJSON = `application/json; charset="utf-8"`
)

type OpenAI struct {
	Domain string
	Apikey string

	Transport http.RoundTripper
	Timeout   time.Duration
	Logger    log.Logger

	RetryOnRateLimited int
}

func (oai *OpenAI) endpoint(format string, a ...any) string {
	return "https://" + oai.Domain + "/v1" + fmt.Sprintf(format, a...)
}

func (oai *OpenAI) buildJsonRequest(a any) (io.Reader, string, error) {
	body, err := json.Marshal(a)
	if err != nil {
		return nil, "", err
	}

	buf := bytes.NewReader(body)
	return buf, contentTypeJSON, nil
}

func (oai *OpenAI) decodeResponse(res *http.Response, obj any) error {
	defer iox.DrainAndClose(res.Body)

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode == http.StatusOK {
		if obj != nil {
			return decoder.Decode(obj)
		}
		return nil
	}

	er := &ErrorResult{
		Status:     res.Status,
		StatusCode: res.StatusCode,
	}
	_ = decoder.Decode(er)

	return er
}

func (oai *OpenAI) call(req *http.Request) (res *http.Response, err error) {
	client := &http.Client{
		Transport: oai.Transport,
		Timeout:   oai.Timeout,
	}

	if oai.Logger != nil {
		oai.Logger.Debugf("%s %s", req.Method, req.URL)
	}

	rid := log.TraceHttpRequest(oai.Logger, req)

	res, err = client.Do(req)
	if err == nil {
		log.TraceHttpResponse(oai.Logger, res, rid)

		if res.StatusCode == http.StatusTooManyRequests {
			rle := &RateLimitedError{
				Status:     res.Status,
				StatusCode: res.StatusCode,
			}

			er := &ErrorResult{}
			decoder := json.NewDecoder(res.Body)
			_ = decoder.Decode(er)

			if er.Detail != nil {
				rle.Message = er.Detail.String()
			}

			iox.DrainAndClose(res.Body)
			return res, rle
		}
	}

	return res, err
}

func (oai *OpenAI) authAndCall(req *http.Request) (res *http.Response, err error) {
	oai.authenticate(req)
	return oai.call(req)
}

func (oai *OpenAI) authenticate(req *http.Request) {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentTypeJSON)
	}

	req.Header.Set("Authorization", "Bearer "+oai.Apikey)
}

func (oai *OpenAI) doCall(req *http.Request, result any) error {
	res, err := oai.authAndCall(req)
	if err != nil {
		return err
	}

	return oai.decodeResponse(res, result)
}

func (oai *OpenAI) doPostWithRetry(url string, source, result any) error {
	return oai.RetryForRateLimited(func() error {
		return oai.doPost(url, source, result)
	})
}

func (oai *OpenAI) doPost(url string, source, result any) error {
	buf, ct, err := oai.buildJsonRequest(source)
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

	return oai.doCall(req, result)
}

func (oai *OpenAI) RetryForRateLimited(api func() error) (err error) {
	return sdk.RetryForRateLimited(api, oai.RetryOnRateLimited, oai.Logger)
}

// https://platform.openai.com/docs/api-reference/chat/create
func (oai *OpenAI) CreateChatCompletion(req *ChatCompeletionRequest) (*ChatCompeletionResponse, error) {
	url := oai.endpoint("/chat/completions")

	res := &ChatCompeletionResponse{}
	err := oai.doPostWithRetry(url, req, res)
	return res, err
}
