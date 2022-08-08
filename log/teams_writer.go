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
	Logfmt  Formatter // log formatter
	Logfil  Filter    // log filter

	mb bytes.Buffer // message buffer
	eb *EventBuffer // error event buffer
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

// SetErrBuffer set the error buffer size
func (tw *TeamsWriter) SetErrBuffer(buffer string) error {
	bsz, err := strconv.Atoi(buffer)
	if err != nil {
		return fmt.Errorf("TeamsWriter - Invalid error buffer: %w", err)
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

	var err error
	for le1 := tw.eb.Peek(); le1 != nil; tw.eb.Poll() {
		if err = tw.write(le1); err != nil {
			break
		}
	}

	if err == nil {
		err = tw.write(le)
	}

	if err != nil {
		tw.eb.Push(le)
		fmt.Fprintln(os.Stderr, err)
	}
}

// format format log event to (message)
func (tw *TeamsWriter) format(le *Event) (mb string) {
	lf := tw.Logfmt
	if lf == nil {
		lf = le.Logger().GetFormatter()
		if lf == nil {
			lf = TextFmtDefault
		}
	}

	tw.mb.Reset()
	lf.Write(&tw.mb, le)
	mb = bye.UnsafeString(tw.mb.Bytes())

	return
}

func (tw *TeamsWriter) write(le *Event) (err error) {
	mb := tw.format(le)

	tm := &teams.Message{}
	tm.Text = mb

	if tw.Timeout.Milliseconds() == 0 {
		tw.Timeout = time.Second * 2
	}

	if err = teams.Post(tw.Webhook, tw.Timeout, tm); err != nil {
		err = fmt.Errorf("TeamsWriter(%q) - Post(): %w", tw.Webhook, err)
	}
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