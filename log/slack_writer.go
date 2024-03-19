package log

import (
	"fmt"
	"net/url"
	"time"

	"github.com/askasoft/pango/sdk/slack"
	"github.com/askasoft/pango/str"
)

// SlackWriter implements log Writer Interface and send log message to slack.
type SlackWriter struct {
	LogFilter
	LogFormatter
	SubFormatter

	Webhook string
	Timeout time.Duration
}

// SetWebhook set the webhook URL
func (sw *SlackWriter) SetWebhook(webhook string) error {
	_, err := url.ParseRequestURI(webhook)
	if err != nil {
		return fmt.Errorf("SlackWriter: invalid webhook %q: %w", webhook, err)
	}
	sw.Webhook = webhook
	return nil
}

// SetTimeout set timeout
func (sw *SlackWriter) SetTimeout(timeout string) error {
	td, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("SlackWriter: invalid timeout %q: %w", timeout, err)
	}
	sw.Timeout = td
	return nil
}

// Write send log message to slack
func (sw *SlackWriter) Write(le *Event) (err error) {
	if sw.Reject(le) {
		return
	}

	sub, msg := sw.format(le)

	sm := &slack.Message{}
	sm.IconEmoji = sw.getIconEmoji(le.Level)
	sm.Text = sub

	sa := &slack.Attachment{Text: msg}
	sm.AddAttachment(sa)

	if sw.Timeout.Milliseconds() == 0 {
		sw.Timeout = time.Second * 2
	}

	if err = slack.Post(sw.Webhook, sw.Timeout, sm); err != nil {
		err = fmt.Errorf("SlackWriter(%s): Post(): %w", sw.Webhook, err)
	}
	return
}

// format format log event to (subject, message)
func (sw *SlackWriter) format(le *Event) (sub, msg string) {
	sbs := sw.SubFormat(le)
	sub = str.UnsafeString(sbs)
	sub = slack.EscapeString(sub)

	mbs := sw.Format(le)
	msg = str.UnsafeString(mbs)

	return
}

func (sw *SlackWriter) getIconEmoji(lvl Level) string {
	switch lvl {
	case LevelFatal:
		return ":boom:"
	case LevelError:
		return ":fire:"
	case LevelWarn:
		return ":warning:"
	case LevelInfo:
		return ":droplet:"
	case LevelDebug:
		return ":bug:"
	case LevelTrace:
		return ":ant:"
	default:
		return ":ghost:"
	}
}

// Flush implementing method. empty.
func (sw *SlackWriter) Flush() {
}

// Close implementing method. empty.
func (sw *SlackWriter) Close() {
}

func init() {
	RegisterWriter("slack", func() Writer {
		return &SlackWriter{}
	})
}
