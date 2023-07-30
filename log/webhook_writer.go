package log

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/askasoft/pango/iox"
)

// WebhookWriter implements log Writer Interface and send log message to webhook.
type WebhookWriter struct {
	LogFilter
	LogFormatter

	Webhook     string // webhook URL
	Method      string // http method
	ContentType string
	Timeout     time.Duration

	hc *http.Client
}

// SetWebhook set the webhook URL
func (ww *WebhookWriter) SetWebhook(webhook string) error {
	_, err := url.ParseRequestURI(webhook)
	if err != nil {
		return fmt.Errorf("WebhookWriter - Invalid webhook: %w", err)
	}
	ww.Webhook = webhook
	return nil
}

// SetTimeout set timeout
func (ww *WebhookWriter) SetTimeout(timeout string) error {
	tmo, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("WebhookWriter - Invalid timeout: %w", err)
	}
	ww.Timeout = tmo
	return nil
}

// Write send log message to webhook
func (ww *WebhookWriter) Write(le *Event) error {
	if ww.Reject(le) {
		return nil
	}

	ww.initClient()

	// format msg
	ww.format(le)

	req, err := http.NewRequest(ww.Method, ww.Webhook, &ww.bb)
	if err != nil {
		err = fmt.Errorf("WebhookWriter(%q) - NewRequest(%v): %w", ww.Webhook, ww.Method, err)
		return err
	}
	if ww.ContentType != "" {
		req.Header.Set("Content-Type", ww.ContentType)
	}

	res, err := ww.hc.Do(req)
	if err != nil {
		err = fmt.Errorf("WebhookWriter(%q) - Send(): %w", ww.Webhook, err)
		return err
	}

	defer iox.DrainAndClose(res.Body)

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		err = fmt.Errorf("WebhookWriter(%q) - Status: %s", ww.Webhook, res.Status)
	}

	return err
}

func (ww *WebhookWriter) initClient() {
	if ww.Method == "" {
		ww.Method = "POST"
	}

	if ww.Timeout.Milliseconds() == 0 {
		ww.Timeout = time.Second * 2
	}

	if ww.hc == nil {
		ww.hc = &http.Client{Timeout: ww.Timeout}
	}
}

func (ww *WebhookWriter) format(le *Event) {
	lf := ww.Formatter
	if lf == nil {
		lf = JSONFmtDefault
	}

	ww.bb.Reset()
	lf.Write(&ww.bb, le)
}

// Flush implementing method. empty.
func (ww *WebhookWriter) Flush() {
}

// Close implementing method. empty.
func (ww *WebhookWriter) Close() {
}

func init() {
	RegisterWriter("webhook", func() Writer {
		return &WebhookWriter{}
	})
}
