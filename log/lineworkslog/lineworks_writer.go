package lineworkslog

import (
	"fmt"
	"net/url"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
	"github.com/askasoft/pango/whk/lineworks"
)

const (
	defaultSubLength = 500
	defaultMsgLength = 2500
)

// LineWorksWriter implements log Writer Interface and send log message to lineworks.
type LineWorksWriter struct {
	log.RetrySupport
	log.FilterSupport
	log.FormatSupport
	log.SubjectSuport

	Webhook      string
	Timeout      time.Duration
	MaxSubLength int
	MaxMsgLength int

	message lineworks.Message
}

// SetWebhook set the webhook URL
func (sw *LineWorksWriter) SetWebhook(webhook string) error {
	_, err := url.ParseRequestURI(webhook)
	if err != nil {
		return fmt.Errorf("lineworkslog: invalid webhook %q: %w", webhook, err)
	}
	sw.Webhook = webhook
	return nil
}

// SetTimeout set timeout
func (sw *LineWorksWriter) SetTimeout(timeout string) error {
	td, err := tmu.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("lineworkslog: invalid timeout %q: %w", timeout, err)
	}
	sw.Timeout = td
	return nil
}

// Write send log message to lineworks
func (sw *LineWorksWriter) Write(le *log.Event) {
	if sw.Reject(le) {
		sw.Flush()
		return
	}

	sw.RetryWrite(le, sw.write)
}

// Flush retry send failed events.
func (sw *LineWorksWriter) Flush() {
	sw.RetryFlush(sw.write)
}

// Close flush and close.
func (sw *LineWorksWriter) Close() {
	sw.Flush()
}

func (sw *LineWorksWriter) write(le *log.Event) (err error) {
	if sw.Timeout.Milliseconds() == 0 {
		sw.Timeout = time.Second * 5
	}

	sub, msg := sw.format(le)
	sw.message.Title = sub
	sw.message.Body.Text = msg

	if err = lineworks.Post(sw.Webhook, sw.Timeout, &sw.message); err != nil {
		err = fmt.Errorf("lineworkslog: Post(%q): %w", sw.Webhook, err)
	}
	return
}

// format format log event to (subject, message)
func (sw *LineWorksWriter) format(le *log.Event) (sub, msg string) {
	sbs := sw.SubFormat(le)
	sub = str.UnsafeString(sbs)
	msl := sw.MaxSubLength
	if msl <= 0 {
		msl = defaultSubLength
	}
	sub = str.Ellipsis(sub, msl)

	mbs := sw.Format(le)
	msg = str.UnsafeString(mbs)
	mml := sw.MaxMsgLength
	if mml <= 0 {
		mml = defaultMsgLength
	}
	msg = str.Ellipsis(msg, mml)
	return
}

func init() {
	log.RegisterWriter("lineworks", func() log.Writer {
		return &LineWorksWriter{}
	})
}
