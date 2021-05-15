package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

// SMTPSender SMTP email sender
type SMTPSender struct {
	// Host represents the host of the SMTP server.
	Host string

	// Port represents the port of the SMTP server.
	Port int

	// Username is the username to use to authenticate to the SMTP server.
	Username string

	// Password is the password to use to authenticate to the SMTP server.
	Password string

	// Auth represents the authentication mechanism used to authenticate to the
	// SMTP server.
	Auth smtp.Auth

	Sender
}

// DialAndSend opens a connection to the SMTP server, sends the given emails and
// closes the connection.
func (ss *SMTPSender) DialAndSend(ms ...*Email) error {
	if err := ss.Dial(); err != nil {
		return err
	}
	defer ss.Close()

	if err := ss.Login(); err != nil {
		return err
	}
	return ss.Send(ms...)
}

// Dial dials and authenticates to an SMTP server.
// Should call Close() when done.
func (ss *SMTPSender) Dial() error {
	if ss.Port <= 0 {
		if ss.SSL {
			ss.Port = 25
		} else {
			ss.Port = 465
		}
	}
	return ss.dial(ss.Host, ss.Port)
}

// Send send mail to SMTP server.
func (ss *SMTPSender) Send(ms ...*Email) error {
	for i, m := range ms {
		if err := ss.send(m); err != nil {
			return fmt.Errorf("Failed to send email %d: %v", i+1, err)
		}
	}

	return nil
}

// Login login smtp server
func (ss *SMTPSender) Login() error {
	c := ss.client
	if ss.Auth == nil && ss.Username != "" {
		if ok, auths := c.Extension("AUTH"); ok {
			if strings.Contains(auths, "CRAM-MD5") {
				ss.Auth = smtp.CRAMMD5Auth(ss.Username, ss.Password)
			} else if strings.Contains(auths, "LOGIN") && !strings.Contains(auths, "PLAIN") {
				ss.Auth = &loginAuth{host: ss.Host, username: ss.Username, password: ss.Password}
			} else {
				ss.Auth = smtp.PlainAuth("", ss.Username, ss.Password, ss.Host)
			}
		}
	}

	if ss.Auth != nil {
		if err := c.Auth(ss.Auth); err != nil {
			return err
		}
	}

	return nil
}
