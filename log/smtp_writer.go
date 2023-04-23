package log

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/net/email"
	"github.com/askasoft/pango/str"
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

	email  *email.Email      // email
	sender *email.SMTPSender // email sender

	sb bytes.Buffer // subject buffer
	mb bytes.Buffer // message buffer
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
	sw.Tos = str.RemoveEmptys(str.TrimSpaces(str.FieldsAny(s, ",;")))
}

// SetCc set Cc recipients
func (sw *SMTPWriter) SetCc(s string) {
	sw.Ccs = str.RemoveEmptys(str.TrimSpaces(str.FieldsAny(s, ",;")))
}

// SetTimeout set timeout
func (sw *SMTPWriter) SetTimeout(timeout string) error {
	tmo, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("SMTPWriter - Invalid timeout: %w", err)
	}
	sw.Timeout = tmo
	return nil
}

// Write send log message to smtp server.
func (sw *SMTPWriter) Write(le *Event) (err error) {
	if sw.Logfil != nil && sw.Logfil.Reject(le) {
		return
	}

	if sw.email == nil {
		if err = sw.initEmail(); err != nil {
			return
		}
	}

	if sw.sender == nil {
		sw.initSender()
	}

	if !sw.sender.IsDialed() {
		if err = sw.sender.Dial(); err != nil {
			err = fmt.Errorf("SMTPWriter(%s:%d) - Dial(): %w", sw.Host, sw.Port, err)
			return
		}

		if err = sw.sender.Login(); err != nil {
			err = fmt.Errorf("SMTPWriter(%s:%d) - Login(%s, %s): %w", sw.Host, sw.Port, sw.Username, sw.Password, err)
			sw.sender.Close()
			return
		}
	}

	sub, msg := sw.format(le)
	sw.email.Subject = sub
	sw.email.Message = msg

	if err = sw.sender.Send(sw.email); err != nil {
		err = fmt.Errorf("SMTPWriter(%s:%d) - Send(): %w", sw.Host, sw.Port, err)
	}
	return
}

// format format log event to (subject, message)
func (sw *SMTPWriter) format(le *Event) (sub, msg string) {
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
	sub = bye.UnsafeString(sw.sb.Bytes())

	sw.mb.Reset()
	lf.Write(&sw.mb, le)
	msg = bye.UnsafeString(sw.mb.Bytes())

	return
}

func (sw *SMTPWriter) initSender() {
	if sw.Timeout.Milliseconds() == 0 {
		sw.Timeout = time.Second * 2
	}

	sw.sender = &email.SMTPSender{
		Host:     sw.Host,
		Port:     sw.Port,
		Username: sw.Username,
		Password: sw.Password,
	}
	sw.sender.Helo = "localhost"
	sw.sender.Timeout = sw.Timeout
	sw.sender.TLSConfig = &tls.Config{ServerName: sw.Host, InsecureSkipVerify: true} //nolint: gosec
}

func (sw *SMTPWriter) initEmail() (err error) {
	m := &email.Email{}

	if err = m.SetFrom(sw.From); err != nil {
		err = fmt.Errorf("SMTPWriter(%s:%d) - SetFrom(): %w", sw.Host, sw.Port, err)
		return
	}

	for _, a := range sw.Tos {
		if err = m.AddTo(a); err != nil {
			err = fmt.Errorf("SMTPWriter(%s:%d) - AddTo(): %w", sw.Host, sw.Port, err)
			return
		}
	}

	for _, a := range sw.Ccs {
		if err = m.AddCc(a); err != nil {
			err = fmt.Errorf("SMTPWriter(%s:%d) - AddCc(): %w", sw.Host, sw.Port, err)
			return
		}
	}

	sw.email = m
	return
}

// Flush implementing method. empty.
func (sw *SMTPWriter) Flush() {
}

// Close close the mail sender
func (sw *SMTPWriter) Close() {
	if sw.sender != nil {
		if err := sw.sender.Close(); err != nil {
			perrorf("SMTPWriter(%s:%d) - Close(): %v", sw.Host, sw.Port, err)
		}
		sw.sender = nil
	}
}

func init() {
	RegisterWriter("smtp", func() Writer {
		return &SMTPWriter{}
	})
}
