package log

import (
	"fmt"
	"net/url"
	"time"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/sdk/teams"
)

// TeamsWriter implements log Writer Interface and send log message to teams.
type TeamsWriter struct {
	LogFilter
	LogFormatter
	SubFormatter

	Webhook string
	Timeout time.Duration
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

// SetTimeout set timeout
func (tw *TeamsWriter) SetTimeout(timeout string) error {
	tmo, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("TeamsWriter - Invalid timeout: %w", err)
	}
	tw.Timeout = tmo
	return nil
}

// Write send log message to teams
func (tw *TeamsWriter) Write(le *Event) (err error) {
	if tw.Reject(le) {
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
	sbs := tw.SubFormat(le)
	sub = bye.UnsafeString(sbs)

	mbs := tw.Format(le)
	msg = bye.UnsafeString(mbs)

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
