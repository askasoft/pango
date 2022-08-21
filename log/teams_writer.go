package log

import (
	"bytes"
	"fmt"
	"os"
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
func (tw *TeamsWriter) Write(le *Event) {
	if tw.Logfil != nil && tw.Logfil.Reject(le) {
		return
	}

	if tw.eb == nil {
		tw.write(le) //nolint: errcheck
		return
	}

	err := tw.flush()
	if err == nil {
		err = tw.write(le)
	}

	if err != nil {
		tw.eb.Push(le)
		fmt.Fprintln(os.Stderr, err)
	}
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

func (tw *TeamsWriter) write(le *Event) (err error) {
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

// flush flush buffered event
func (tw *TeamsWriter) flush() error {
	if tw.eb != nil {
		for le := tw.eb.Peek(); le != nil; tw.eb.Poll() {
			if err := tw.write(le); err != nil {
				return err
			}
		}
	}
	return nil
}

// Flush implementing method. empty.
func (tw *TeamsWriter) Flush() {
	tw.flush()
}

// Close implementing method. empty.
func (tw *TeamsWriter) Close() {
	tw.flush()
}

func init() {
	RegisterWriter("teams", func() Writer {
		return &TeamsWriter{}
	})
}
