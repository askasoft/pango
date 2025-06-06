package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/log/httplog"
	"github.com/askasoft/pango/sdk"
)

type OpenAI struct {
	Domain string
	Apikey string

	Transport http.RoundTripper
	Timeout   time.Duration
	Logger    log.Logger

	MaxRetries  int
	RetryAfter  time.Duration
	ShouldRetry func(error) bool // default retry on not canceled error or (status = 429 || (status >= 500 && status <= 599))
}

func (oai *OpenAI) endpoint(format string, a ...any) string {
	return "https://" + oai.Domain + "/v1" + fmt.Sprintf(format, a...)
}

func (oai *OpenAI) call(req *http.Request) (res *http.Response, err error) {
	client := &http.Client{
		Transport: oai.Transport,
		Timeout:   oai.Timeout,
	}

	if log := oai.Logger; log != nil {
		log.Debugf("%s %s", req.Method, req.URL)
	}

	rid := httplog.TraceHttpRequest(oai.Logger, req)

	res, err = client.Do(req)
	if err != nil {
		return res, err
	}

	httplog.TraceHttpResponse(oai.Logger, res, rid)
	return res, nil
}

func (oai *OpenAI) authenticate(req *http.Request) {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentTypeJSON)
	}

	req.Header.Set("Authorization", "Bearer "+oai.Apikey)
}

func (oai *OpenAI) doCall(req *http.Request, result any) error {
	sr := oai.ShouldRetry
	if sr == nil {
		sr = shouldRetry
	}

	oai.authenticate(req)

	res, err := oai.call(req)
	if err != nil {
		if sr(err) {
			return sdk.NewRetryError(err, oai.RetryAfter)
		}
		return err
	}
	defer iox.DrainAndClose(res.Body)

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode == http.StatusOK {
		if result != nil {
			return decoder.Decode(result)
		}
		return nil
	}

	re := &ResultError{
		Status:     res.Status,
		StatusCode: res.StatusCode,
	}
	_ = decoder.Decode(re)

	if sr(re) {
		re.RetryAfter = oai.RetryAfter
	}
	return re
}

func (oai *OpenAI) doPostWithRetry(ctx context.Context, url string, source, result any) error {
	return sdk.RetryForError(ctx, func() error {
		return oai.doPost(ctx, url, source, result)
	}, oai.MaxRetries, oai.Logger)
}

func (oai *OpenAI) doPost(ctx context.Context, url string, source, result any) error {
	buf, ct, err := buildJsonRequest(source)
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

	return oai.doCall(req, result)
}

// https://platform.openai.com/docs/api-reference/chat/create
func (oai *OpenAI) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	url := oai.endpoint("/chat/completions")

	res := &ChatCompletionResponse{}
	err := oai.doPostWithRetry(ctx, url, req, res)
	return res, err
}

// https://platform.openai.com/docs/api-reference/embeddings/create
func (oai *OpenAI) CreateTextEmbeddings(ctx context.Context, req *TextEmbeddingsRequest) (*TextEmbeddingsResponse, error) {
	url := oai.endpoint("/embeddings")

	res := &TextEmbeddingsResponse{}
	err := oai.doPostWithRetry(ctx, url, req, res)
	return res, err
}
