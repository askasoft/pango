package log

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pandafw/pango/bye"
	"github.com/pandafw/pango/net/slack"
)

// SlackWriter implements log Writer Interface and send log message to slack.
type SlackWriter struct {
	Webhook  string
	Username string
	Channel  string
	Timeout  time.Duration
	Subfmt   Formatter // subject formatter
	Logfmt   Formatter // log formatter
	Logfil   Filter    // log filter

	sb bytes.Buffer // subject buffer
	mb bytes.Buffer // message buffer
	eb *EventBuffer // error event buffer
}

// SetSubject set the subject formatter
func (sw *SlackWriter) SetSubject(format string) {
	sw.Subfmt = NewLogFormatter(format)
}

// SetFormat set the log formatter
func (sw *SlackWriter) SetFormat(format string) {
	sw.Logfmt = NewLogFormatter(format)
}

// SetFilter set the log filter
func (sw *SlackWriter) SetFilter(filter string) {
	sw.Logfil = NewLogFilter(filter)
}

// SetTimeout set timeout
func (sw *SlackWriter) SetTimeout(timeout string) error {
	tmo, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("SlackWriter - Invalid timeout: %w", err)
	}
	sw.Timeout = tmo
	return nil
}

// SetErrBuffer set the error buffer size
func (sw *SlackWriter) SetErrBuffer(buffer string) error {
	bsz, err := strconv.Atoi(buffer)
	if err != nil {
		return fmt.Errorf("SlackWriter - Invalid error buffer: %w", err)
	}
	if bsz > 0 {
		sw.eb = &EventBuffer{BufSize: bsz}
	}
	return nil
}

// Write send log message to slack
func (sw *SlackWriter) Write(le *Event) {
	if sw.Logfil != nil && sw.Logfil.Reject(le) {
		return
	}

	if sw.eb == nil {
		sw.write(le) //nolint: errcheck
		return
	}

	var err error
	for le1 := sw.eb.Peek(); le1 != nil; sw.eb.Poll() {
		if err = sw.write(le1); err != nil {
			break
		}
	}

	if err == nil {
		err = sw.write(le)
	}

	if err != nil {
		sw.eb.Push(le)
		fmt.Fprintln(os.Stderr, err)
	}
}

// format format log event to (subject, message)
func (sw *SlackWriter) format(le *Event) (sb, mb string) {
	sf := sw.Subfmt
	if sf == nil {
		sf = TextFmtSubject
	}

	lf := sw.Logfmt
	if lf == nil {
		lf = le.Logger().GetFormatter()
		if lf == nil {
			lf = TextFmtDefault
		}
	}

	sw.sb.Reset()
	sf.Write(&sw.sb, le)
	sb = bye.UnsafeString(sw.sb.Bytes())

	sw.mb.Reset()
	lf.Write(&sw.mb, le)
	mb = bye.UnsafeString(sw.mb.Bytes())

	return
}

func (sw *SlackWriter) write(le *Event) (err error) {
	sb, mb := sw.format(le)

	sm := &slack.Message{}
	sm.IconEmoji = getIconEmoji(le.Level())
	sm.Channel = sw.Channel
	sm.Username = sw.Username
	sm.Text = sb

	sa := &slack.Attachment{Text: mb}
	sm.AddAttachment(sa)

	if sw.Timeout.Milliseconds() == 0 {
		sw.Timeout = time.Second * 2
	}

	if err = slack.Post(sw.Webhook, sw.Timeout, sm); err != nil {
		err = fmt.Errorf("SlackWriter(%q) - Post(): %w", sw.Webhook, err)
	}
	return
}

// Flush implementing method. empty.
func (sw *SlackWriter) Flush() {
}

// Close implementing method. empty.
func (sw *SlackWriter) Close() {
}

func getIconEmoji(lvl Level) string {
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

func init() {
	RegisterWriter("slack", func() Writer {
		return &SlackWriter{}
	})
}
