package teamslog

import (
	"fmt"
	"net/url"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sdk/teams"
	"github.com/askasoft/pango/str"
)

// TeamsWriter implements log Writer Interface and send log message to teams.
type TeamsWriter struct {
	log.RetrySupport
	log.FilterSupport
	log.FormatSupport
	log.SubjectSuport

	Webhook string
	Timeout time.Duration

	message teams.Message
}

// SetWebhook set the webhook URL
func (tw *TeamsWriter) SetWebhook(webhook string) error {
	_, err := url.ParseRequestURI(webhook)
	if err != nil {
		return fmt.Errorf("TeamsWriter: invalid webhook %q: %w", webhook, err)
	}
	tw.Webhook = webhook
	return nil
}

// SetTimeout set timeout
func (tw *TeamsWriter) SetTimeout(timeout string) error {
	td, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("TeamsWriter: invalid timeout %q: %w", timeout, err)
	}
	tw.Timeout = td
	return nil
}

// Write send log message to teams
func (tw *TeamsWriter) Write(le *log.Event) {
	if tw.Reject(le) {
		return
	}

	tw.RetryWrite(le, tw.write)
}

// Flush retry send failed events.
func (tw *TeamsWriter) Flush() {
	tw.RetryFlush(tw.write)
}

// Close flush and close.
func (tw *TeamsWriter) Close() {
	tw.Flush()
}

func (tw *TeamsWriter) write(le *log.Event) (err error) {
	tw.message.Title, tw.message.Text = tw.format(le)

	if tw.Timeout.Milliseconds() == 0 {
		tw.Timeout = time.Second * 2
	}

	if err = teams.Post(tw.Webhook, tw.Timeout, &tw.message); err != nil {
		err = fmt.Errorf("TeamsWriter(%s): Post(): %w", tw.Webhook, err)
	}
	return
}

// format format log event to (message)
func (tw *TeamsWriter) format(le *log.Event) (sub, msg string) {
	sbs := tw.SubFormat(le)
	sub = str.UnsafeString(sbs)

	mbs := tw.Format(le)
	msg = str.UnsafeString(mbs)
	return
}

func init() {
	log.RegisterWriter("teams", func() log.Writer {
		return &TeamsWriter{}
	})
}
