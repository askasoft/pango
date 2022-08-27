package log

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/pandafw/pango/bye"
	"github.com/pandafw/pango/net/teams"
)

// TeamsWriter implements log Writer Interface and send log message to teams.
type TeamsWriter struct {
	Webhook string
	Timeout time.Duration
	Subfmt  Formatter // subject formatter
	Logfmt  Formatter // log formatter
	Logfil  Filter    // log filter

	sb bytes.Buffer // subject buffer
	mb bytes.Buffer // message buffer
	eb *EventBuffer // event buffer
}

// SetWebhook set the webhook URL
func (tw *TeamsWriter) SetWebhook(webhook string) error {
	_, err := url.ParseRequestURI(webhook)
	if err != nil {
		return fmt.Errorf("TeamsWriter - Invalid webhook: %w", err)
	}
	tw.Webhook = webhook
	return nil
}

// SetSubject set the subject formatter
func (tw *TeamsWriter) SetSubject(format string) {
	tw.Subfmt = NewLogFormatter(format)
}

// SetFormat set the log formatter
func (tw *TeamsWriter) SetFormat(format string) {
	tw.Logfmt = NewLogFormatter(format)
}

// SetFilter set the log filter
func (tw *TeamsWriter) SetFilter(filter string) {
	tw.Logfil = NewLogFilter(filter)
}

// SetTimeout set timeout
func (tw *TeamsWriter) SetTimeout(timeout string) error {
	tmo, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("TeamsWriter - Invalid timeout: %w", err)
	}
	tw.Timeout = tmo
	return nil
}

// SetBuffer set the event buffer size
func (tw *TeamsWriter) SetBuffer(buffer string) error {
	bsz, err := strconv.Atoi(buffer)
	if err != nil {
		return fmt.Errorf("TeamsWriter - Invalid buffer: %w", err)
	}
	if bsz > 0 {
		tw.eb = &EventBuffer{BufSize: bsz}
	}
	return nil
}

// Write send log message to teams
func (tw *TeamsWriter) Write(le *Event) (err error) {
	if tw.Logfil != nil && tw.Logfil.Reject(le) {
		return
	}

	sub, msg := tw.format(le)

	tm := &teams.Message{}
	tm.Title = sub
	tm.Text = msg

	if tw.Timeout.Milliseconds() == 0 {
		tw.Timeout = time.Second * 2
	}

	if err = teams.Post(tw.Webhook, tw.Timeout, tm); err != nil {
		err = fmt.Errorf("TeamsWriter(%q) - Post(): %w", tw.Webhook, err)
	}
	return
}

// format format log event to (message)
func (tw *TeamsWriter) format(le *Event) (sub, msg string) {
	sf := tw.Subfmt
	if sf == nil {
		sf = TextFmtSubject
	}

	lf := tw.Logfmt
	if lf == nil {
		lf = le.Logger().GetFormatter()
		if lf == nil {
			lf = TextFmtDefault
		}
	}

	tw.sb.Reset()
	sf.Write(&tw.sb, le)
	sub = bye.UnsafeString(tw.sb.Bytes())

	tw.mb.Reset()
	lf.Write(&tw.mb, le)
	msg = bye.UnsafeString(tw.mb.Bytes())

	return
}

// Flush implementing method. empty.
func (tw *TeamsWriter) Flush() {
}

// Close implementing method. empty.
func (tw *TeamsWriter) Close() {
}

func init() {
	RegisterWriter("teams", func() Writer {
		return &TeamsWriter{}
	})
}
