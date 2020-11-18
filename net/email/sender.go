package email

import (
	"crypto/tls"
	"encoding/base64"
	"io"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pandafw/pango/enc"
	"github.com/pandafw/pango/strutil"
)

// Sender email sender
type Sender struct {
	// LocalName is the hostname sent to the SMTP server with the HELO command.
	// By default, "localhost" is sent.
	Helo string
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
	// SSL defines whether an SSL connection is used. It should be false in
	// most cases since the authentication mechanism should use the STARTTLS
	// extension instead.
	SSL bool
	// TSLConfig represents the TLS configuration used for the TLS (when the
	// STARTTLS extension is used) or SSL connection.
	TLSConfig *tls.Config
	StartTLS  bool
	// Timeout timeout when connect to the SMTP server
	Timeout time.Duration

	client *smtp.Client
}

// Dial dials and authenticates to an SMTP server.
// Should call Close() when done.
func (s *Sender) Dial() error {
	if s.Port <= 0 {
		if s.SSL {
			s.Port = 25
		} else {
			s.Port = 465
		}
	}

	conn, err := net.DialTimeout("tcp", s.Host+":"+strconv.Itoa(s.Port), s.Timeout)
	if err != nil {
		return err
	}

	if s.SSL {
		conn = tls.Client(conn, s.tlsConfig())
	}

	c, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		return err
	}

	if s.Helo != "" {
		if err := c.Hello(s.Helo); err != nil {
			c.Close()
			return err
		}
	}

	if !s.SSL {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(s.tlsConfig()); err != nil {
				c.Close()
				return err
			}
		}
	}

	if s.Auth == nil && s.Username != "" {
		if ok, auths := c.Extension("AUTH"); ok {
			if strings.Contains(auths, "CRAM-MD5") {
				s.Auth = smtp.CRAMMD5Auth(s.Username, s.Password)
			} else if strings.Contains(auths, "LOGIN") && !strings.Contains(auths, "PLAIN") {
				s.Auth = &loginAuth{
					username: s.Username,
					password: s.Password,
					host:     s.Host,
				}
			} else {
				s.Auth = smtp.PlainAuth("", s.Username, s.Password, s.Host)
			}
		}
	}

	if s.Auth != nil {
		if err = c.Auth(s.Auth); err != nil {
			c.Close()
			return err
		}
	}

	s.client = c
	return nil
}

func (s *Sender) tlsConfig() *tls.Config {
	if s.TLSConfig == nil {
		s.TLSConfig = &tls.Config{ServerName: s.Host}
	}
	return s.TLSConfig
}

// Send send message to SMTP server.
func (s *Sender) Send(mail *Email) error {
	c := s.client
	if err := c.Mail(mail.GetSender()); err != nil {
		if err == io.EOF {
			// This is probably due to a timeout, so reconnect and try again.
			derr := s.Dial()
			if derr != nil {
				return derr
			}
		} else {
			return err
		}
	}

	for _, addr := range mail.GetTos() {
		if err := c.Rcpt(addr.Address); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	err = s.writeMail(w, mail)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

// header SMTP message header
type header map[string]string

func formatDate(t time.Time) string {
	return t.Format(time.RFC1123Z)
}

func encodeWord(s string) string {
	return mime.QEncoding.Encode("UTF-8", s)
}

func encodeBody(s string) string {
	sb := strings.Builder{}
	b := base64.NewEncoder(base64.StdEncoding, enc.NewBase64LineWriter(&sb))
	b.Write([]byte(s))
	b.Close()
	return sb.String()
}

func encodeAddress(as ...*mail.Address) string {
	sb := strings.Builder{}
	for i, a := range as {
		sb.WriteString(a.String())
		if i < len(as)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

func writeHead(w io.Writer, h header) error {
	var err error
	for k, v := range h {
		err = writeString(w, k)
		err = writeString(w, ": ")
		err = writeString(w, v)
		err = writeString(w, "\r\n")
	}
	return err
}

func copyFile(w io.Writer, r io.Reader) error {
	b := base64.NewEncoder(base64.StdEncoding, enc.NewBase64LineWriter(w))
	_, err := io.Copy(w, r)
	if err != nil {
		return err
	}
	return b.Close()
}

func writeMsg(w io.Writer, enc string, msg string) error {
	if enc == "Base64" {
		b := base64.NewEncoder(base64.StdEncoding, enc.NewBase64LineWriter(w))
		err := writeString(b, msg)
		if err != nil {
			return err
		}
		return b.Close()
	}
	return writeString(w, msg)
}

func writeString(w io.Writer, s string) error {
	_, err := w.Write([]byte(s))
	return err
}

func writeBody(w io.Writer, enc string, m *Email, boundary string) error {
	if boundary == "" {
		return writeMsg(w, enc, m.Message)
	}

	var err error

	// Write the message part
	err = writeString(w, "--"+boundary+"\n")
	err = writeString(w, HEADER_CONTENT_TYPE)
	err = writeString(w, ": ")
	if m.HTML {
		err = writeString(w, "text/html")
	} else {
		err = writeString(w, "text/plain")
	}
	err = writeString(w, "; charset=UTF-8")
	err = writeString(w, "\r\n")

	err = writeString(w, HEADER_CONTENT_DISPOSITION)
	err = writeString(w, ": ")
	err = writeString(w, HEADER_CONTENT_DISPOSITION_INLIE)
	err = writeString(w, "\r\n")

	err = writeString(w, HEADER_CONTENT_TRANSFER_ENCODING)
	err = writeString(w, ": ")
	err = writeString(w, enc)
	err = writeString(w, "\r\n\r\n")
	if err != nil {
		return err
	}

	err = writeMsg(w, enc, m.Message)
	err = writeString(w, "\r\n")
	if err != nil {
		return err
	}

	// Append attachments
	err = writeAttachments(w, m.Attachments, boundary)
	if err != nil {
		return err
	}

	err = writeString(w, "--"+boundary+"--\r\n\r\n")
	return err
}

func writeAttachments(w io.Writer, as []*Attachment, boundary string) error {
	var err error
	for _, a := range as {
		mt := mime.TypeByExtension(filepath.Ext(a.Name))
		if mt == "" {
			mt = "application/octet-stream"
		}

		err = writeString(w, "--")
		err = writeString(w, boundary)
		err = writeString(w, "\r\n")

		err = writeString(w, HEADER_CONTENT_TYPE)
		err = writeString(w, ": ")
		err = writeString(w, encodeWord(mt))
		err = writeString(w, "; name=\"")
		err = writeString(w, encodeWord(a.Name))
		err = writeString(w, "\"\r\n")

		err = writeString(w, HEADER_CONTENT_DISPOSITION)
		err = writeString(w, ": ")
		if a.Cid != "" {
			err = writeString(w, HEADER_CONTENT_DISPOSITION_INLIE)
		} else {
			err = writeString(w, HEADER_CONTENT_DISPOSITION_ATTACHMENT)
		}
		err = writeString(w, "; filename=\"")
		err = writeString(w, encodeWord(a.Name))
		err = writeString(w, "\"\r\n")

		if a.Cid != "" {
			err = writeString(w, HEADER_CONTENT_ID)
			err = writeString(w, ": ")
			err = writeString(w, a.Cid)
			err = writeString(w, "\r\n")
		}
		err = writeString(w, HEADER_CONTENT_TRANSFER_ENCODING)
		err = writeString(w, ": Base64\r\n\r\n")

		err = copyFile(w, a.Data)
		if err != nil {
			return err
		}

		err = writeString(w, "\r\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sender) writeMail(w io.Writer, m *Email) error {
	header := header{}

	header[HEADER_DATE] = formatDate(m.GetDate())
	header[HEADER_FROM] = m.GetFrom().String()
	header[HEADER_TO] = encodeAddress(m.GetTos()...)
	if len(m.GetCcs()) > 0 {
		header[HEADER_CC] = encodeAddress(m.GetCcs()...)
	}
	if len(m.GetBccs()) > 0 {
		header[HEADER_BCC] = encodeAddress(m.GetBccs()...)
	}
	if len(m.GetReplys()) > 0 {
		header[HEADER_REPLY_TO] = encodeAddress(m.GetReplys()...)
	}

	if m.MsgID != "" {
		header[HEADER_MESSAGE_ID] = m.MsgID
	}
	header[HEADER_SUBJECT] = encodeWord(m.Subject)
	header[HEADER_MIME_VERSION] = HEADER_MIME_VERSION_10

	var boundary string
	if m.HTML || len(m.Attachments) > 0 {
		boundary = strutil.RandDigitLetters(28)
		header[HEADER_CONTENT_TYPE] = "multipart/mixed; boundary=" + boundary
	} else {
		header[HEADER_CONTENT_TYPE] = "text/plain; charset=UTF-8"
	}

	enc := "7bit"
	if !strutil.IsASCII(m.Message) {
		enc = "Base64"
	}

	header[HEADER_CONTENT_TRANSFER_ENCODING] = enc

	writeHead(w, header)
	writeBody(w, enc, m, boundary)
	return nil
}
