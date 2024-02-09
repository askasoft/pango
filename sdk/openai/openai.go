package openai

import (
	"fmt"
	"net/http"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sdk"
)

type OpenAI struct {
	Domain string
	Apikey string

	Transport http.RoundTripper
	Timeout   time.Duration
	Logger    log.Logger

	MaxRetries    int
	RetryAfter    time.Duration
	AbortOnRetry  func() bool
	AbortInterval time.Duration
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

	rid := log.TraceHttpRequest(oai.Logger, req)

	res, err = client.Do(req)
	if err != nil {
		return res, sdk.NewNetError(err, oai.RetryAfter)
	}

	log.TraceHttpResponse(oai.Logger, res, rid)
	return res, nil
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

	return decodeResponse(res, result)
}

func (oai *OpenAI) doPostWithRetry(url string, source, result any) error {
	return sdk.RetryForError(func() error {
		return oai.doPost(url, source, result)
	}, oai.MaxRetries, oai.AbortOnRetry, oai.AbortInterval, oai.Logger)
}

func (oai *OpenAI) doPost(url string, source, result any) error {
	buf, ct, err := buildJsonRequest(source)
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

// https://platform.openai.com/docs/api-reference/chat/create
func (oai *OpenAI) CreateChatCompletion(req *ChatCompeletionRequest) (*ChatCompeletionResponse, error) {
	url := oai.endpoint("/chat/completions")

	res := &ChatCompeletionResponse{}
	err := oai.doPostWithRetry(url, req, res)
	return res, err
}

// https://platform.openai.com/docs/api-reference/embeddings/create
func (oai *OpenAI) CreateTextEmbeddings(req *TextEmbeddingsRequest) (*TextEmbeddingsResponse, error) {
	url := oai.endpoint("/embeddings")

	res := &TextEmbeddingsResponse{}
	err := oai.doPostWithRetry(url, req, res)
	return res, err
}
