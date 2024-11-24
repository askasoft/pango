package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/log/httplog"
	"github.com/askasoft/pango/sdk"
)

type AzureOpenAI struct {
	Domain     string
	Apikey     string
	Apiver     string
	Deployment string

	Transport http.RoundTripper
	Timeout   time.Duration
	Logger    log.Logger

	MaxRetries int
	RetryAfter time.Duration
}

func (aoai *AzureOpenAI) endpoint(format string, a ...any) string {
	return "https://" + aoai.Domain + "/openai/deployments/" + aoai.Deployment + fmt.Sprintf(format, a...) + "?api-version=" + aoai.Apiver
}

func (aoai *AzureOpenAI) call(req *http.Request) (res *http.Response, err error) {
	client := &http.Client{
		Transport: aoai.Transport,
		Timeout:   aoai.Timeout,
	}

	if log := aoai.Logger; log != nil {
		log.Debugf("%s %s", req.Method, req.URL)
	}

	rid := httplog.TraceHttpRequest(aoai.Logger, req)

	res, err = client.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return res, err
		}
		return res, sdk.NewNetError(err, aoai.RetryAfter)
	}

	httplog.TraceHttpResponse(aoai.Logger, res, rid)
	return res, nil
}

func (aoai *AzureOpenAI) authenticate(req *http.Request) {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentTypeJSON)
	}

	req.Header.Set("API-KEY", aoai.Apikey)
}

func (aoai *AzureOpenAI) doCall(req *http.Request, result any) error {
	aoai.authenticate(req)
	res, err := aoai.call(req)
	if err != nil {
		return err
	}

	return decodeResponse(res, result, aoai.RetryAfter)
}

func (aoai *AzureOpenAI) doPostWithRetry(ctx context.Context, url string, source, result any) error {
	return sdk.RetryForError(ctx, func() error {
		return aoai.doPost(ctx, url, source, result)
	}, aoai.MaxRetries, aoai.Logger)
}

func (aoai *AzureOpenAI) doPost(ctx context.Context, url string, source, result any) error {
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

	return aoai.doCall(req, result)
}

// https://platform.openai.com/docs/api-reference/chat/create
func (aoai *AzureOpenAI) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	url := aoai.endpoint("/chat/completions")

	res := &ChatCompletionResponse{}
	err := aoai.doPostWithRetry(ctx, url, req, res)
	return res, err
}

// https://platform.openai.com/docs/api-reference/embeddings/create
func (aoai *AzureOpenAI) CreateTextEmbeddings(ctx context.Context, req *TextEmbeddingsRequest) (*TextEmbeddingsResponse, error) {
	url := aoai.endpoint("/embeddings")

	res := &TextEmbeddingsResponse{}
	err := aoai.doPostWithRetry(ctx, url, req, res)
	return res, err
}
