package log

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// WebhookWriter implements log Writer Interface and send log message to webhook.
type WebhookWriter struct {
	Webhook     string // webhook URL
	Method      string // http method
	ContentType string
	Timeout     time.Duration
	Logfmt      Formatter // log formatter
	Logfil      Filter    // log filter

	hc *http.Client
	bb bytes.Buffer
}

// SetFormat set the log formatter
func (ew *WebhookWriter) SetFormat(format string) {
	ew.Logfmt = NewJSONFormatter(format)
}

// SetFilter set the log filter
func (ew *WebhookWriter) SetFilter(filter string) {
	ew.Logfil = NewLogFilter(filter)
}

// SetTimeout set timeout
func (ew *WebhookWriter) SetTimeout(timeout string) error {
	tmo, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("WebhookWriter - Invalid timeout: %v", err)
	}
	ew.Timeout = tmo
	return nil
}

// Write send log message to webhook
func (ew *WebhookWriter) Write(le *Event) {
	if ew.Logfil != nil && ew.Logfil.Reject(le) {
		return
	}

	lf := ew.Logfmt
	if lf == nil {
		lf = le.Logger().GetFormatter()
		if lf == nil {
			lf = JSONFmtDefault
		}
	}

	if ew.Timeout.Milliseconds() == 0 {
		ew.Timeout = time.Second * 2
	}

	if ew.hc == nil {
		ew.hc = &http.Client{Timeout: ew.Timeout}
	}

	if len(ew.Method) == 0 {
		ew.Method = "POST"
	}

	// format msg
	ew.bb.Reset()
	lf.Write(&ew.bb, le)

	req, err := http.NewRequest(ew.Method, ew.Webhook, &ew.bb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WebhookWriter(%q) - NewRequest(%v): %v\n", ew.Webhook, ew.Method, err)
		return
	}
	if len(ew.ContentType) > 0 {
		req.Header.Set("Content-Type", ew.ContentType)
	}

	res, err := ew.hc.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WebhookWriter(%q) - Send(): %v\n", ew.Webhook, err)
		return
	}
	io.Copy(ioutil.Discard, res.Body)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		fmt.Fprintf(os.Stderr, "WebhookWriter(%q) - Status: %s\n", ew.Webhook, res.Status)
	}
}

// Flush implementing method. empty.
func (ew *WebhookWriter) Flush() {
}

// Close implementing method. empty.
func (ew *WebhookWriter) Close() {
}

func init() {
	RegisterWriter("webhook", func() Writer {
		return &WebhookWriter{}
	})
}
