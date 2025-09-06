package httplog

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
)

// HTTPWriter implements log Writer Interface and batch send log messages to webhook.
type HTTPWriter struct {
	log.BatchSupport
	log.RetrySupport
	log.FilterSupport
	log.FormatSupport

	URL         string // request URL
	Method      string // http method
	Insecure    bool
	Username    string // basic auth username
	Password    string // basic auth password
	ContentType string
	Timeout     time.Duration

	client *http.Client
}

// SetUrl set the request url
func (hw *HTTPWriter) SetUrl(u string) error {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return fmt.Errorf("HTTPWriter: invalid URL %q: %w", u, err)
	}
	hw.URL = u
	return nil
}

// SetTimeout set timeout
func (hw *HTTPWriter) SetTimeout(timeout string) error {
	td, err := tmu.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("httplog: invalid timeout %q: %w", timeout, err)
	}
	hw.Timeout = td
	return nil
}

// Write cache log message, flush if needed
func (hw *HTTPWriter) Write(le *log.Event) {
	if hw.Reject(le) {
		le = nil
	}

	if hw.Retries > 0 {
		hw.RetryWrite(le, hw.write)
	} else {
		hw.BatchWrite(le, hw.flush)
	}
}

// Flush flush cached events
func (hw *HTTPWriter) Flush() {
	if hw.Retries > 0 {
		hw.RetryFlush(hw.write)
	} else {
		hw.BatchFlush(hw.flush)
	}
}

// Close flush and close the writer
func (hw *HTTPWriter) Close() {
	hw.Flush()
}

func (hw *HTTPWriter) write(le *log.Event) error {
	hw.initClient()

	hw.Buffer.Reset()
	lf := hw.GetFormatter(le, log.JSONFmtDefault)
	lf.Write(&hw.Buffer, le)

	return hw.send()
}

func (hw *HTTPWriter) flush(eb *log.EventBuffer) error {
	hw.initClient()

	hw.Buffer.Reset()
	for it := eb.Iterator(); it.Next(); {
		le := it.Value()
		lf := hw.GetFormatter(le, log.JSONFmtDefault)
		lf.Write(&hw.Buffer, le)
	}

	return hw.send()
}

func (hw *HTTPWriter) initClient() {
	if hw.client == nil {
		if hw.Method == "" {
			hw.Method = http.MethodPost
		}
		if hw.Timeout.Milliseconds() == 0 {
			bc := 2
			if hw.BatchCount > 2 {
				bc = hw.BatchCount
			}
			hw.Timeout = time.Second * time.Duration(bc)
		}

		hw.client = &http.Client{Timeout: hw.Timeout}
		if hw.Insecure {
			hw.client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint: gosec
			}
		}
	}
}

func (hw *HTTPWriter) send() error {
	req, err := http.NewRequest(hw.Method, hw.URL, &hw.Buffer)
	if err != nil {
		err = fmt.Errorf("httplog: NewRequest(%q, %q): %w", hw.URL, hw.Method, err)
		return err
	}
	if hw.ContentType != "" {
		req.Header.Set("Content-Type", hw.ContentType)
	}
	if hw.Username != "" {
		req.SetBasicAuth(hw.Username, hw.Password)
	}

	res, err := hw.client.Do(req)
	if err != nil {
		err = fmt.Errorf("httplog: Send(%q): %w", hw.URL, err)
		return err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		buf, _ := iox.ReadAll(res.Body)
		err = fmt.Errorf("httplog: Read(%q): %s: %s", hw.URL, res.Status, str.UnsafeString(buf))
	}

	iox.DrainAndClose(res.Body)
	return err
}

func init() {
	log.RegisterWriter("http", func() log.Writer {
		return &HTTPWriter{}
	})
}
