//go:build !go1.18
// +build !go1.18

package log

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/iox"
)

// WebhookBatchWriter implements log Writer Interface and batch send log messages to webhook.
type WebhookBatchWriter struct {
	LogFilter
	LogFormatter

	Webhook     string // webhook URL
	Method      string // http method
	Insecure    bool
	Username    string // basic auth username
	Password    string // basic auth password
	ContentType string
	Timeout     time.Duration
	CacheCount  int           // max cacheable event count
	BatchCount  int           // messages send batch count
	FlushLevel  Level         // flush events if event <= FlushLevel
	FlushDelta  time.Duration // flush events if [current log event time] - [first log event time] >= FlushDelta

	evtbuf *EventBuffer
	hc     *http.Client
}

// SetWebhook set the webhook URL
func (wbw *WebhookBatchWriter) SetWebhook(webhook string) error {
	_, err := url.ParseRequestURI(webhook)
	if err != nil {
		return fmt.Errorf("WebhookBatchWriter - Invalid webhook '%s': %w", webhook, err)
	}
	wbw.Webhook = webhook
	return nil
}

// SetTimeout set timeout
func (wbw *WebhookBatchWriter) SetTimeout(timeout string) error {
	td, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("WebhookBatchWriter - Invalid timeout '%s': %w", timeout, err)
	}
	wbw.Timeout = td
	return nil
}

// SetFlushLevel set the flush level
func (wbw *WebhookBatchWriter) SetFlushLevel(lvl string) {
	wbw.FlushLevel = ParseLevel(lvl)
}

func (wbw *WebhookBatchWriter) init() {
	if wbw.BatchCount < 1 {
		wbw.BatchCount = 10
	}
	if wbw.CacheCount < wbw.BatchCount {
		wbw.CacheCount = wbw.BatchCount * 2
	}

	if wbw.evtbuf == nil {
		wbw.evtbuf = NewEventBuffer(wbw.CacheCount)
	}
}

func (wbw *WebhookBatchWriter) shouldFlush(le *Event) bool {
	if wbw.evtbuf.Len() >= wbw.BatchCount {
		return true
	}
	if le.Level <= wbw.FlushLevel {
		return true
	}
	if wbw.FlushDelta > 0 && wbw.evtbuf.Len() > 1 {
		if fle, ok := wbw.evtbuf.Peek(); ok {
			if le.When.Sub(fle.When) >= wbw.FlushDelta {
				return true
			}
		}
	}
	return false
}

// Write cache log message, flush if needed
func (wbw *WebhookBatchWriter) Write(le *Event) error {
	if wbw.Reject(le) {
		return nil
	}

	wbw.init()
	wbw.evtbuf.Push(le)

	if wbw.shouldFlush(le) {
		wbw.Flush()
	}

	return nil
}

// Flush flush cached events
func (wbw *WebhookBatchWriter) Flush() {
	if err := wbw.flush(); err == nil {
		wbw.evtbuf.Clear()
	} else {
		perror(err)
	}
}

func (wbw *WebhookBatchWriter) flush() error {
	if wbw.hc == nil {
		if wbw.Method == "" {
			wbw.Method = http.MethodPost
		}
		if wbw.Timeout.Milliseconds() == 0 {
			wbw.Timeout = time.Second * 2
		}

		wbw.hc = &http.Client{Timeout: wbw.Timeout}
		if wbw.Insecure {
			wbw.hc.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint: gosec
			}
		}
	}

	wbw.Buffer.Reset()
	for _, le := range wbw.evtbuf.Values() {
		lf := wbw.GetFormatter(le, JSONFmtDefault)
		lf.Write(&wbw.Buffer, le.(*Event))
	}

	req, err := http.NewRequest(wbw.Method, wbw.Webhook, &wbw.Buffer)
	if err != nil {
		err = fmt.Errorf("WebhookBatchWriter(%q) - NewRequest(%v): %w", wbw.Webhook, wbw.Method, err)
		return err
	}
	if wbw.ContentType != "" {
		req.Header.Set("Content-Type", wbw.ContentType)
	}
	if wbw.Username != "" {
		req.SetBasicAuth(wbw.Username, wbw.Password)
	}

	res, err := wbw.hc.Do(req)
	if err != nil {
		err = fmt.Errorf("WebhookBatchWriter(%q) - Send(): %w", wbw.Webhook, err)
		return err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		buf, _ := iox.ReadAll(res.Body)
		err = fmt.Errorf("WebhookBatchWriter(%q) - %s: %s", wbw.Webhook, res.Status, bye.UnsafeString(buf))
	}

	iox.DrainAndClose(res.Body)
	return err
}

// Close flush and close the writer
func (wbw *WebhookBatchWriter) Close() {
	wbw.Flush()
}

func init() {
	RegisterWriter("bathook", func() Writer {
		return &WebhookBatchWriter{}
	})
}
