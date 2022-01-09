package log

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
	eb *EventBuffer // error event buffer
}

// SetFormat set the log formatter
func (ww *WebhookWriter) SetFormat(format string) {
	ww.Logfmt = NewJSONFormatter(format)
}

// SetFilter set the log filter
func (ww *WebhookWriter) SetFilter(filter string) {
	ww.Logfil = NewLogFilter(filter)
}

// SetTimeout set timeout
func (ww *WebhookWriter) SetTimeout(timeout string) error {
	tmo, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("WebhookWriter - Invalid timeout: %v", err)
	}
	ww.Timeout = tmo
	return nil
}

// SetErrBuffer set the error buffer size
func (ww *WebhookWriter) SetErrBuffer(buffer string) error {
	bsz, err := strconv.Atoi(buffer)
	if err != nil {
		return fmt.Errorf("SlackWriter - Invalid error buffer: %v", err)
	}
	if bsz > 0 {
		ww.eb = &EventBuffer{BufSize: bsz}
	}
	return nil
}

// Write send log message to webhook
func (ww *WebhookWriter) Write(le *Event) {
	if ww.Logfil != nil && ww.Logfil.Reject(le) {
		return
	}

	if ww.eb == nil {
		ww.write(le)
		return
	}

	var err error
	for le1 := ww.eb.Peek(); le1 != nil; ww.eb.Poll() {
		if err = ww.write(le1); err != nil {
			break
		}
	}

	if err == nil {
		err = ww.write(le)
	}

	if err != nil {
		ww.eb.Push(le)
		fmt.Fprintln(os.Stderr, err)
	}
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
	lf := ww.Logfmt
	if lf == nil {
		lf = le.Logger().GetFormatter()
		if lf == nil {
			lf = JSONFmtDefault
		}
	}

	ww.bb.Reset()
	lf.Write(&ww.bb, le)
}

// write send log message to webhook
func (ww *WebhookWriter) write(le *Event) error {
	ww.initClient()

	// format msg
	ww.format(le)

	req, err := http.NewRequest(ww.Method, ww.Webhook, &ww.bb)
	if err != nil {
		err = fmt.Errorf("WebhookWriter(%q) - NewRequest(%v): %v", ww.Webhook, ww.Method, err)
		return err
	}
	if ww.ContentType != "" {
		req.Header.Set("Content-Type", ww.ContentType)
	}

	res, err := ww.hc.Do(req)
	if err != nil {
		err = fmt.Errorf("WebhookWriter(%q) - Send(): %v", ww.Webhook, err)
		return err
	}
	io.Copy(ioutil.Discard, res.Body)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		err = fmt.Errorf("WebhookWriter(%q) - Status: %s", ww.Webhook, res.Status)
	}

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
