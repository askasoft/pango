package email

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/pandafw/pango/enc"
	"github.com/pandafw/pango/str"
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

	// Timeout timeout when connect to the SMTP server
	Timeout time.Duration

	// SSL defines whether an SSL connection is used. It should be false in
	// most cases since the authentication mechanism should use the STARTTLS
	// extension instead.
	SSL bool

	// SkipTLS Skip StartTLS when the STARTTLS extension is used
	SkipTLS bool

	// TSLConfig represents the TLS configuration used for the TLS (when the
	// STARTTLS extension is used) or SSL connection.
	TLSConfig *tls.Config

	// ConnDebug  a conn wrap func
	ConnDebug func(conn net.Conn) net.Conn

	// DataDebug  a data writer wrap func
	DataDebug func(w io.Writer) io.Writer

	client *smtp.Client
}

// DialAndSend opens a connection to the SMTP server, sends the given emails and
// closes the connection.
func (s *Sender) DialAndSend(ms ...*Email) error {
	err := s.Dial()
	if err != nil {
		return err
	}
	defer s.Close()

	return s.Send(ms...)
}

// IsDialed return true if the sender is dialed
func (s *Sender) IsDialed() bool {
	return s.client != nil
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
	return s.dial()
}

// Send send mail to SMTP server.
func (s *Sender) Send(ms ...*Email) error {
	for i, m := range ms {
		if err := s.send(m); err != nil {
			return fmt.Errorf("Failed to send email %d: %v", i+1, err)
		}
	}

	return nil
}

// Close close the SMTP client
func (s *Sender) Close() error {
	if s.client == nil {
		return nil
	}

	err := s.client.Quit()
	s.client = nil
	return err
}

func (s *Sender) setConn(i interface{}, c net.Conn) {
	v := reflect.ValueOf(i).Elem()
	f := v.FieldByName("conn")
	pc := (*net.Conn)(unsafe.Pointer(f.UnsafeAddr()))
	*pc = c
}

func (s *Sender) getConn(i interface{}) net.Conn {
	v := reflect.ValueOf(i).Elem()
	f := v.FieldByName("conn")
	pc := (*net.Conn)(unsafe.Pointer(f.UnsafeAddr()))
	return *pc
}

func (s *Sender) dial() error {
	addr := s.Host + ":" + strconv.Itoa(s.Port)
	conn, err := net.DialTimeout("tcp", addr, s.Timeout)
	if err != nil {
		return fmt.Errorf("Failed to dial %s - %v", addr, err)
	}

	if s.SSL {
		conn = tls.Client(conn, s.tlsConfig())
	}

	corg := conn
	if s.ConnDebug != nil {
		conn = s.ConnDebug(conn)
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

	if !s.SSL && !s.SkipTLS {
		if s.ConnDebug != nil {
			s.setConn(c, corg)
			c.Text = textproto.NewConn(corg)
		}
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(s.tlsConfig()); err != nil {
				c.Close()
				return err
			}
		}
		if s.ConnDebug != nil {
			ctls := s.getConn(c)
			conn = s.ConnDebug(ctls)
			c.Text = textproto.NewConn(conn)
		}
	}

	if s.Auth == nil && s.Username != "" {
		if ok, auths := c.Extension("AUTH"); ok {
			if strings.Contains(auths, "CRAM-MD5") {
				s.Auth = smtp.CRAMMD5Auth(s.Username, s.Password)
			} else if strings.Contains(auths, "LOGIN") && !strings.Contains(auths, "PLAIN") {
				s.Auth = &loginAuth{host: s.Host, username: s.Username, password: s.Password}
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

func (s *Sender) send(mail *Email) error {
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
		w.Close()
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Sender) tlsConfig() *tls.Config {
	if s.TLSConfig == nil {
		s.TLSConfig = &tls.Config{ServerName: s.Host}
	}
	return s.TLSConfig
}

// header SMTP message header
type header map[string]string

func formatDate(t time.Time) string {
	return t.Format(time.RFC1123Z)
}

// Encode a RFC 822 "word" token into mail-safe form as per RFC 2047.
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

// http://www.faqs.org/rfcs/rfc2822.html
func writeHeader(w io.Writer, h header) error {
	var err error
	for k, v := range h {
		err = writeString(w, k)
		err = writeString(w, ": ")
		if len(k)+len(v) > 74 {
			// folding header
			v = strings.ReplaceAll(v, " ", "\r\n ")
		}
		err = writeString(w, v)
		err = writeString(w, "\r\n")
	}
	return err
}

func closeAttach(r io.Reader) {
	if c, ok := r.(io.Closer); ok {
		c.Close()
	}
}

func copyAttach(w io.Writer, r io.Reader) error {
	defer closeAttach(r)

	b := base64.NewEncoder(base64.StdEncoding, enc.NewBase64LineWriter(w))
	_, err := io.Copy(b, r)
	if err != nil {
		return err
	}

	return b.Close()
}

func writeMsg(w io.Writer, encoding string, msg string) error {
	if encoding == "Base64" {
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

	header := header{}
	if m.HTML {
		header["Content-Type"] = "text/html; charset=UTF-8"
	} else {
		header["Content-Type"] = "text/plain; charset=UTF-8"
	}
	header["Content-Disposition"] = "inline"
	header["Content-Transfer-Encoding"] = enc

	// Write the message part
	err = writeString(w, "--"+boundary+"\n")
	err = writeHeader(w, header)
	err = writeString(w, "\r\n")
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

		header := header{}
		header["Content-Type"] = mt + "; name=\"" + encodeWord(a.Name) + "\""
		if a.Cid != "" {
			header["Content-Disposition"] = "inline; filename=\"" + encodeWord(a.Name) + "\""
		} else {
			header["Content-Disposition"] = "attachment; filename=\"" + encodeWord(a.Name) + "\""
		}
		if a.Cid != "" {
			header["Content-ID"] = a.Cid
		}
		header["Content-Transfer-Encoding"] = "Base64"

		err = writeString(w, "--"+boundary+"\r\n")
		err = writeHeader(w, header)
		err = writeString(w, "\r\n")

		err = copyAttach(w, a.Data)
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

	header["MIME-Version"] = "1.0"
	if m.MsgID != "" {
		header["Message-ID"] = m.MsgID
	}
	header["Date"] = formatDate(m.GetDate())
	header["From"] = m.GetFrom().String()
	header["To"] = encodeAddress(m.GetTos()...)
	if len(m.GetCcs()) > 0 {
		header["Cc"] = encodeAddress(m.GetCcs()...)
	}
	if len(m.GetBccs()) > 0 {
		header["Bcc"] = encodeAddress(m.GetBccs()...)
	}
	if len(m.GetReplys()) > 0 {
		header["Reply-To"] = encodeAddress(m.GetReplys()...)
	}
	header["Subject"] = encodeWord(m.Subject)

	var boundary string
	if m.HTML || len(m.Attachments) > 0 {
		boundary = str.RandDigitLetters(28)
		header["Content-Type"] = "multipart/mixed; boundary=" + boundary
	} else {
		header["Content-Type"] = "text/plain; charset=UTF-8"
	}

	enc := "7bit"
	if !str.IsASCII(m.Message) {
		enc = "Base64"
	}
	header["Content-Transfer-Encoding"] = enc

	if s.DataDebug != nil {
		w = s.DataDebug(w)
	}

	err := writeHeader(w, header)
	if err != nil {
		return err
	}

	return writeBody(w, enc, m, boundary)
}