package httplog

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/log/internal"
	"github.com/askasoft/pango/str"
)

// HTTPWriter implements log Writer Interface and batch send log messages to webhook.
type HTTPWriter struct {
	log.LogFilter
	log.LogFormatter
	log.BatchWriter

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
	td, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("HTTPWriter: invalid timeout %q: %w", timeout, err)
	}
	hw.Timeout = td
	return nil
}

// Write cache log message, flush if needed
func (hw *HTTPWriter) Write(le *log.Event) error {
	if hw.Reject(le) {
		return nil
	}

	if hw.BatchCount > 1 {
		hw.InitBuffer()
		hw.EventBuffer.Push(le)

		if hw.ShouldFlush(le) {
			if err := hw.flush(); err != nil {
				return err
			}
			hw.EventBuffer.Clear()
		}
		return nil
	}

	return hw.write(le)
}

func (hw *HTTPWriter) write(le *log.Event) error {
	hw.initClient()
	hw.Format(le, log.JSONFmtDefault)
	return hw.send()
}

func (hw *HTTPWriter) flush() error {
	hw.initClient()

	hw.Buffer.Reset()
	for it := hw.EventBuffer.Iterator(); it.Next(); {
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
		err = fmt.Errorf("HTTPWriter(%q): NewRequest(%v): %w", hw.URL, hw.Method, err)
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
		err = fmt.Errorf("HTTPWriter(%q): Send(): %w", hw.URL, err)
		return err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		buf, _ := iox.ReadAll(res.Body)
		err = fmt.Errorf("HTTPWriter(%q): Send(): %s: %s", hw.URL, res.Status, str.UnsafeString(buf))
	}

	iox.DrainAndClose(res.Body)
	return err
}

// Flush flush cached events
func (hw *HTTPWriter) Flush() {
	if hw.EventBuffer == nil || hw.EventBuffer.IsEmpty() {
		return
	}

	if err := hw.flush(); err == nil {
		hw.EventBuffer.Clear()
	} else {
		internal.Perror(err)
	}
}

// Close flush and close the writer
func (hw *HTTPWriter) Close() {
	hw.Flush()
}

func init() {
	log.RegisterWriter("http", func() log.Writer {
		return &HTTPWriter{}
	})
}
