package slacklog

import (
	"fmt"
	"net/url"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sdk/slack"
	"github.com/askasoft/pango/str"
)

// SlackWriter implements log Writer Interface and send log message to slack.
type SlackWriter struct {
	log.LogFilter
	log.LogFormatter
	log.SubFormatter

	Webhook string
	Timeout time.Duration

	message slack.Message
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
func (sw *SlackWriter) Write(le *log.Event) (err error) {
	if sw.Reject(le) {
		return
	}

	if len(sw.message.Attachments) == 0 {
		sw.message.AddAttachment(&slack.Attachment{})
	}

	sub, msg := sw.format(le)
	sw.message.IconEmoji = sw.getIconEmoji(le.Level)
	sw.message.Text = sub
	sw.message.Attachments[0].Text = msg

	if sw.Timeout.Milliseconds() == 0 {
		sw.Timeout = time.Second * 2
	}

	if err = slack.Post(sw.Webhook, sw.Timeout, &sw.message); err != nil {
		err = fmt.Errorf("SlackWriter(%s): Post(): %w", sw.Webhook, err)
	}
	return
}

// Flush implementing method. empty.
func (sw *SlackWriter) Flush() {
}

// Close implementing method. empty.
func (sw *SlackWriter) Close() {
}

// format format log event to (subject, message)
func (sw *SlackWriter) format(le *log.Event) (sub, msg string) {
	sbs := sw.SubFormat(le)
	sub = str.UnsafeString(sbs)
	sub = slack.EscapeString(sub)

	mbs := sw.Format(le)
	msg = str.UnsafeString(mbs)
	return
}

func (sw *SlackWriter) getIconEmoji(lvl log.Level) string {
	switch lvl {
	case log.LevelFatal:
		return ":boom:"
	case log.LevelError:
		return ":fire:"
	case log.LevelWarn:
		return ":warning:"
	case log.LevelInfo:
		return ":droplet:"
	case log.LevelDebug:
		return ":bug:"
	case log.LevelTrace:
		return ":ant:"
	default:
		return ":ghost:"
	}
}

func init() {
	log.RegisterWriter("slack", func() log.Writer {
		return &SlackWriter{}
	})
}
