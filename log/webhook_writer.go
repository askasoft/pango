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

// WebhookWriter implements log Writer Interface and send log message to webhook.
type WebhookWriter struct {
	LogFilter
	LogFormatter

	Webhook     string // webhook URL
	Method      string // http method
	Insecure    bool
	Username    string // basic auth username
	Password    string // basic auth password
	ContentType string
	Timeout     time.Duration

	hc *http.Client
}

// SetWebhook set the webhook URL
func (ww *WebhookWriter) SetWebhook(webhook string) error {
	_, err := url.ParseRequestURI(webhook)
	if err != nil {
		return fmt.Errorf("WebhookWriter - Invalid webhook '%s': %w", webhook, err)
	}
	ww.Webhook = webhook
	return nil
}

// SetTimeout set timeout
func (ww *WebhookWriter) SetTimeout(timeout string) error {
	td, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("WebhookWriter - Invalid timeout '%s': %w", timeout, err)
	}
	ww.Timeout = td
	return nil
}

// Write send log message to webhook
func (ww *WebhookWriter) Write(le *Event) error {
	if ww.Reject(le) {
		return nil
	}

	if ww.hc == nil {
		if ww.Method == "" {
			ww.Method = http.MethodPost
		}
		if ww.Timeout.Milliseconds() == 0 {
			ww.Timeout = time.Second * 2
		}

		ww.hc = &http.Client{Timeout: ww.Timeout}
		if ww.Insecure {
			ww.hc.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint: gosec
			}
		}
	}

	ww.Format(le, JSONFmtDefault)

	req, err := http.NewRequest(ww.Method, ww.Webhook, &ww.Buffer)
	if err != nil {
		err = fmt.Errorf("WebhookWriter(%q) - NewRequest(%v): %w", ww.Webhook, ww.Method, err)
		return err
	}
	if ww.ContentType != "" {
		req.Header.Set("Content-Type", ww.ContentType)
	}
	if ww.Username != "" {
		req.SetBasicAuth(ww.Username, ww.Password)
	}

	res, err := ww.hc.Do(req)
	if err != nil {
		err = fmt.Errorf("WebhookWriter(%q) - Send(): %w", ww.Webhook, err)
		return err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		buf, _ := iox.ReadAll(res.Body)
		err = fmt.Errorf("WebhookWriter(%q) - %s: %s", ww.Webhook, res.Status, bye.UnsafeString(buf))
	}

	iox.DrainAndClose(res.Body)
	return err
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
