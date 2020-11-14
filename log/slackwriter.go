package log

import (
	"fmt"
	"os"
	"time"

	"github.com/pandafw/pango/slack"
)

// SlackWriter implements LogWriter Interface and writes messages to slack.
type SlackWriter struct {
	Level    int    `json:"level"`
	Webhook  string `json:"webhook"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
	Timeout  string `json:"timeout"`
	Subfmt   Formatter
	Logfmt   Formatter
}

// SetLevel set the log level
func (sw *SlackWriter) SetLevel(level string) {
	sw.Level = ParseLevel(level)
}

// SetFormat set a log formatter
func (sw *SlackWriter) SetFormat(format string) {
	sw.Logfmt = NewFormatter(format)
}

// SetSubject set a subject formatter
func (sw *SlackWriter) SetSubject(format string) {
	sw.Subfmt = NewFormatter(format)
}

// Write write message in smtp writer.
// it will send an email with subject and only this message.
func (sw *SlackWriter) Write(le *Event) {
	if sw.Level < le.Level {
		return
	}
	if sw.Subfmt == nil {
		sw.Subfmt = FormatterSimple
	}
	if sw.Logfmt == nil {
		sw.Logfmt = le.Logger.GetFormatter()
	}

	sm := slack.Message{}
	sm.IconEmoji = getIconEmoji(le.Level)
	sm.Channel = sw.Channel
	sm.Username = sw.Username
	sm.Text = sw.Subfmt.Format(le)

	sa := &slack.Attachment{Text: sw.Logfmt.Format(le)}
	sm.AddAttachment(sa)

	timeout, _ := time.ParseDuration(sw.Timeout)
	err := sm.Post(sw.Webhook, timeout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// Flush implementing method. empty.
func (sw *SlackWriter) Flush() {
}

// Close implementing method. empty.
func (sw *SlackWriter) Close() {
}

func getIconEmoji(lvl int) string {
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
	}
	return ":ghost:"
}
