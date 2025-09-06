package smtplog

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/net/email"
	"github.com/askasoft/pango/str"
	"github.com/askasoft/pango/tmu"
)

// SMTPWriter implements log Writer Interface and send log message.
type SMTPWriter struct {
	log.RetrySupport
	log.FilterSupport
	log.FormatSupport
	log.SubjectSuport

	Host     string
	Port     int
	Insecure bool
	Username string
	Password string
	Timeout  time.Duration

	email  email.Email       // email
	sender *email.SMTPSender // email sender
}

// SetFrom set From recipient
func (sw *SMTPWriter) SetFrom(s string) error {
	return sw.email.SetFrom(s)
}

// SetTo set To recipients
func (sw *SMTPWriter) SetTo(s string) error {
	return sw.email.SetTo(s)
}

// SetCc set Cc recipients
func (sw *SMTPWriter) SetCc(s string) error {
	return sw.email.SetCc(s)
}

// SetTimeout set timeout
func (sw *SMTPWriter) SetTimeout(timeout string) error {
	td, err := tmu.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("smtplog: invalid timeout %q: %w", timeout, err)
	}
	sw.Timeout = td
	return nil
}

// Write send log message to smtp server.
func (sw *SMTPWriter) Write(le *log.Event) {
	if sw.Reject(le) {
		sw.Flush()
		return
	}

	sw.RetryWrite(le, sw.write)
}

// Flush retry send failed events.
func (sw *SMTPWriter) Flush() {
	sw.RetryFlush(sw.write)
}

// Close flush and close the mail sender
func (sw *SMTPWriter) Close() {
	sw.Flush()

	if sw.sender != nil {
		if err := sw.sender.Close(); err != nil {
			log.Perrorf("smtplog: Close(%s:%d): %v", sw.Host, sw.Port, err)
		}
		sw.sender = nil
	}
}

func (sw *SMTPWriter) write(le *log.Event) (err error) {
	if sw.sender == nil {
		sw.initSender()
	}

	if !sw.sender.IsDialed() {
		if err = sw.sender.Dial(); err != nil {
			err = fmt.Errorf("smtplog: Dial(%s:%d): %w", sw.Host, sw.Port, err)
			return
		}

		if err = sw.sender.Login(); err != nil {
			err = fmt.Errorf("smtplog: (%s:%d) Login(%s, %s): %w", sw.Host, sw.Port, sw.Username, sw.Password, err)
			sw.sender.Close()
			return
		}
	}

	sub, msg := sw.format(le)
	sw.email.Subject = sub
	sw.email.Message = msg
	sw.email.Date = time.Now()

	if err = sw.sender.Send(&sw.email); err != nil {
		err = fmt.Errorf("smtplog: (%s:%d) Send(): %w", sw.Host, sw.Port, err)
	}
	return
}

// format format log event to (subject, message)
func (sw *SMTPWriter) format(le *log.Event) (sub, msg string) {
	sbs := sw.SubFormat(le)
	sub = str.UnsafeString(sbs)

	mbs := sw.Format(le)
	msg = str.UnsafeString(mbs)
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
	if sw.Insecure {
		sw.sender.TLSConfig = &tls.Config{ServerName: sw.Host, InsecureSkipVerify: true} //nolint: gosec
	}
}

func init() {
	log.RegisterWriter("smtp", func() log.Writer {
		return &SMTPWriter{}
	})
}
