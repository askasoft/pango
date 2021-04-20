package log

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/pandafw/pango/net/email"
	"github.com/pandafw/pango/str"
)

// SMTPWriter implements log Writer Interface and send log message.
type SMTPWriter struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	Tos      []string
	Ccs      []string
	Timeout  time.Duration
	Subfmt   Formatter // subject formatter
	Logfmt   Formatter // log formatter
	Logfil   Filter    // log filter

	email  *email.Email  // email
	sender *email.Sender // email sender
}

// SetSubject set the subject formatter
func (sw *SMTPWriter) SetSubject(format string) {
	sw.Subfmt = NewLogFormatter(format)
}

// SetFormat set the log formatter
func (sw *SMTPWriter) SetFormat(format string) {
	sw.Logfmt = NewLogFormatter(format)
}

// SetFilter set the log filter
func (sw *SMTPWriter) SetFilter(filter string) {
	sw.Logfil = NewLogFilter(filter)
}

// SetTo set To recipients
func (sw *SMTPWriter) SetTo(s string) {
	sw.Tos = str.RemoveEmptys(str.TrimSpaces(str.SplitAnyNoEmpty(s, ",;")))
}

// SetCc set Cc recipients
func (sw *SMTPWriter) SetCc(s string) {
	sw.Ccs = str.RemoveEmptys(str.TrimSpaces(str.SplitAnyNoEmpty(s, ",;")))
}

// SetTimeout set timeout
func (sw *SMTPWriter) SetTimeout(timeout string) error {
	tmo, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("SMTPWriter - Invalid timeout: %v", err)
	}
	sw.Timeout = tmo
	return nil
}

// Write send log message to smtp server.
func (sw *SMTPWriter) Write(le *Event) {
	if sw.Logfil != nil && sw.Logfil.Reject(le) {
		return
	}

	if sw.Subfmt == nil {
		sw.Subfmt = TextFmtSubject
	}
	if sw.Logfmt == nil {
		sw.Logfmt = le.Logger.GetFormatter()
		if sw.Logfmt == nil {
			sw.Logfmt = TextFmtDefault
		}
	}

	if sw.email == nil {
		m := &email.Email{}
		err := m.SetFrom(sw.From)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SMTPWriter(%s:%d) - SetFrom(): %v\n", sw.Host, sw.Port, err)
			return
		}
		for _, a := range sw.Tos {
			err := m.AddTo(a)
			if err != nil {
				fmt.Fprintf(os.Stderr, "SMTPWriter(%s:%d) - AddTo(): %v\n", sw.Host, sw.Port, err)
				return
			}
		}
		for _, a := range sw.Ccs {
			err := m.AddCc(a)
			if err != nil {
				fmt.Fprintf(os.Stderr, "SMTPWriter(%s:%d) - AddCc(): %v\n", sw.Host, sw.Port, err)
				return
			}
		}
		sw.email = m
	}

	if sw.sender == nil {
		sw.sender = &email.Sender{
			Host:      sw.Host,
			Port:      sw.Port,
			Timeout:   sw.Timeout,
			Username:  sw.Username,
			Password:  sw.Password,
			TLSConfig: &tls.Config{ServerName: sw.Host, InsecureSkipVerify: true},
		}
	}
	if !sw.sender.IsDialed() {
		err := sw.sender.Dial()
		if err != nil {
			fmt.Fprintf(os.Stderr, "SMTPWriter(%s:%d) - Dial(): %v\n", sw.Host, sw.Port, err)
			return
		}
	}

	sw.email.Subject = sw.Subfmt.Format(le)
	sw.email.Message = sw.Logfmt.Format(le)

	err := sw.sender.Send(sw.email)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SMTPWriter(%s:%d) - Send(): %v\n", sw.Host, sw.Port, err)
	}
}

// Flush implementing method. empty.
func (sw *SMTPWriter) Flush() {
}

// Close close the mail sender
func (sw *SMTPWriter) Close() {
	if sw.sender != nil {
		err := sw.sender.Close()
		fmt.Fprintf(os.Stderr, "SMTPWriter(%s:%d) - Close(): %v\n", sw.Host, sw.Port, err)
		sw.sender = nil
	}
}

func init() {
	RegisterWriter("smtp", func() Writer {
		return &SMTPWriter{}
	})
}
