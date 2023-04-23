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
	"unicode"
	"unsafe"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/str"
)

// Sender email sender
type Sender struct {
	// LocalName is the hostname sent to the SMTP server with the HELO command.
	// By default, "localhost" is sent.
	Helo string

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

// IsDialed return true if the sender is dialed
func (s *Sender) IsDialed() bool {
	return s.client != nil
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

func (s *Sender) setConn(i any, c net.Conn) {
	v := reflect.ValueOf(i).Elem()
	f := v.FieldByName("conn")
	pc := (*net.Conn)(unsafe.Pointer(f.UnsafeAddr()))
	*pc = c
}

func (s *Sender) getConn(i any) net.Conn {
	v := reflect.ValueOf(i).Elem()
	f := v.FieldByName("conn")
	pc := (*net.Conn)(unsafe.Pointer(f.UnsafeAddr()))
	return *pc
}

func (s *Sender) tlsConfig(host string) *tls.Config {
	if s.TLSConfig == nil {
		s.TLSConfig = &tls.Config{ServerName: host} //nolint: gosec
	}
	return s.TLSConfig
}

func (s *Sender) dial(host string, port int) error {
	addr := host + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", addr, s.Timeout)
	if err != nil {
		return fmt.Errorf("Failed to dial %s - %w", addr, err)
	}

	if s.SSL {
		conn = tls.Client(conn, s.tlsConfig(host))
	}

	corg := conn
	if s.ConnDebug != nil {
		conn = s.ConnDebug(conn)
	}

	c, err := smtp.NewClient(conn, host)
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
			if err := c.StartTLS(s.tlsConfig(host)); err != nil {
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

	s.client = c
	return nil
}

func (s *Sender) send(recipients []string, email *Email) error {
	c := s.client
	if err := c.Mail(email.GetSender()); err != nil {
		return err
	}

	for _, addr := range recipients {
		if err := c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	err = s.writeMail(w, email)
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

func formatDate(t time.Time) string {
	return t.Format(time.RFC1123Z)
}

// Encode a RFC 822 "word" token into mail-safe form as per RFC 2047.
func encodeWord(s string) string {
	return mime.QEncoding.Encode("utf-8", s)
}

// http://www.faqs.org/rfcs/rfc2822.html
// if string's length > 75, it will be splitted with ' ' by mime encoding
func encodeString(s string) string {
	// Text in an encoded-word in a display-name must not contain certain
	// characters like quotes or parentheses (see RFC 2047 section 5.3).
	// When this is the case encode the name using base64 encoding.
	if strings.ContainsAny(s, "\"#$%&'(),.:;<>@[]^`{|}~") {
		return mime.BEncoding.Encode("UTF-8", s)
	}
	return mime.QEncoding.Encode("UTF-8", s)
}

func writeStrings(w io.Writer, bs ...string) (err error) {
	for _, b := range bs {
		if err = writeString(w, b); err != nil {
			return
		}
	}
	if err = writeEOL(w); err != nil {
		return
	}
	return
}

func writeFolding(w io.Writer, h string, v string) (err error) {
	// http://www.faqs.org/rfcs/rfc2822.html
	if len(h)+len(v) < 75 {
		err = writeString(w, v)
		return
	}

	ll := len(h) + 2 // line length

	// folding header
	for v != "" {
		i := strings.IndexByte(v, ' ')

		// skip empty
		if i == 0 {
			v = v[1:]
			continue
		}

		s := v
		if i > 0 {
			s = v[0:i]
			v = v[i:]
		} else {
			v = ""
		}

		if ll+len(s) > 74 {
			if _, err = w.Write([]byte{'\r', '\n', ' '}); err != nil {
				return
			}
			ll = 1
		}
		if err = writeString(w, s); err != nil {
			return
		}
		ll += len(s)
	}
	return
}

func writeHeader(w io.Writer, h string, v string) (err error) {
	if err = writeString(w, h); err != nil {
		return
	}
	if _, err = w.Write([]byte{':', ' '}); err != nil {
		return
	}
	if err = writeFolding(w, h, v); err != nil {
		return
	}
	if err = writeEOL(w); err != nil {
		return
	}
	return
}

func writeAddress(w io.Writer, h string, as ...*mail.Address) (err error) {
	if len(as) < 1 {
		return
	}

	if err = writeString(w, h); err != nil {
		return
	}
	if _, err = w.Write([]byte{':', ' '}); err != nil {
		return
	}

	for i, a := range as {
		if i > 0 {
			if _, err = w.Write([]byte{'\r', '\n', ' ', ',', ' '}); err != nil {
				return
			}
		}
		if err = writeFolding(w, h, a.String()); err != nil {
			return
		}
	}

	err = writeEOL(w)
	return
}

func writeEOL(w io.Writer) (err error) {
	_, err = w.Write([]byte{'\r', '\n'})
	return
}

func writeString(w io.Writer, s string) (err error) {
	_, err = w.Write(str.UnsafeBytes(s))
	return
}

func needEncoding(s string) bool {
	// http://www.faqs.org/rfcs/rfc2822.html
	// 2.1.1. Line Length Limits
	if s == "" {
		return false
	}

	ll := 0 // line length
	sl := len(s)
	for i := 0; i < sl; i++ {
		ll++

		ch := s[i]
		if ch > unicode.MaxASCII {
			return true
		}
		if ll > 1000 {
			return true
		}
		if ch == '\n' {
			ll = 0
		}
	}
	return false
}

func (s *Sender) writeMail(w io.Writer, m *Email) (err error) {
	if s.DataDebug != nil {
		w = s.DataDebug(w)
	}

	if err = writeHeader(w, "MIME-Version", "1.0"); err != nil {
		return
	}
	if m.MsgID != "" {
		if err = writeHeader(w, "Message-ID", m.MsgID); err != nil {
			return
		}
	}
	if err = writeHeader(w, "Date", formatDate(m.GetDate())); err != nil {
		return
	}
	if err = writeAddress(w, "From", m.GetFrom()); err != nil {
		return
	}
	if err = writeAddress(w, "To", m.GetTos()...); err != nil {
		return
	}
	if err = writeAddress(w, "Cc", m.GetCcs()...); err != nil {
		return
	}
	if err = writeAddress(w, "Bcc", m.GetBccs()...); err != nil {
		return
	}
	if err = writeAddress(w, "Reply-To", m.GetReplys()...); err != nil {
		return
	}
	if err = writeHeader(w, "Subject", encodeString(m.Subject)); err != nil {
		return
	}

	var boundary string
	if m.HTML || len(m.Attachments) > 0 {
		boundary = str.RandLetterNumbers(28)
		if err = writeHeader(w, "Content-Type", "multipart/mixed; boundary="+boundary); err != nil {
			return
		}
	} else {
		if err = writeHeader(w, "Content-Type", "text/plain; charset=UTF-8"); err != nil {
			return
		}
	}

	enc := "7bit"
	if needEncoding(m.Message) {
		enc = "Base64"
	}
	if err = writeHeader(w, "Content-Transfer-Encoding", enc); err != nil {
		return
	}

	if err = writeEOL(w); err != nil {
		return
	}

	return writeBody(w, enc, m, boundary)
}

func closeAttach(r io.Reader) {
	if c, ok := r.(io.Closer); ok {
		c.Close()
	}
}

func copyAttach(w io.Writer, r io.Reader) error {
	defer closeAttach(r)

	b := base64.NewEncoder(base64.StdEncoding, iox.NewMimeChunkWriter(w))
	_, err := io.Copy(b, r)
	if err != nil {
		return err
	}

	return b.Close()
}

func writeMsg(w io.Writer, encoding string, msg string) error {
	if encoding == "Base64" {
		b := base64.NewEncoder(base64.StdEncoding, iox.NewMimeChunkWriter(w))
		err := writeString(b, msg)
		if err != nil {
			return err
		}
		return b.Close()
	}
	return writeString(w, msg)
}

func writeBody(w io.Writer, enc string, m *Email, boundary string) (err error) {
	if boundary == "" {
		return writeMsg(w, enc, m.Message)
	}

	// Write the message part
	if err = writeStrings(w, "--", boundary); err != nil {
		return
	}

	cot := "plain"
	if m.HTML {
		cot = "html"
	}

	if err = writeHeader(w, "Content-Type", fmt.Sprintf("text/%s; charset=UTF-8", cot)); err != nil {
		return
	}
	if err = writeHeader(w, "Content-Disposition", "inline"); err != nil {
		return
	}
	if err = writeHeader(w, "Content-Transfer-Encoding", enc); err != nil {
		return
	}
	if err = writeEOL(w); err != nil {
		return
	}

	if err = writeMsg(w, enc, m.Message); err != nil {
		return
	}
	if err = writeEOL(w); err != nil {
		return
	}

	// Append attachments
	if err = writeAttachments(w, m.Attachments, boundary); err != nil {
		return
	}

	err = writeStrings(w, "--", boundary, "--")
	return
}

func writeAttachments(w io.Writer, as []*Attachment, boundary string) (err error) {
	for _, a := range as {
		mt := mime.TypeByExtension(filepath.Ext(a.Name))
		if mt == "" {
			mt = "application/octet-stream"
		}

		if err = writeStrings(w, "--", boundary); err != nil {
			return
		}
		if err = writeHeader(w, "Content-Type", mt+"; name=\""+encodeWord(a.Name)+"\""); err != nil {
			return
		}
		cod := "attachment"
		if a.Cid != "" {
			cod = "inline"
			if err = writeHeader(w, "Content-ID", a.Cid); err != nil {
				return
			}
		}
		if err = writeHeader(w, "Content-Disposition", cod+"; filename=\""+encodeWord(a.Name)+"\""); err != nil {
			return
		}
		if err = writeHeader(w, "Content-Transfer-Encoding", "Base64"); err != nil {
			return
		}
		if err = writeEOL(w); err != nil {
			return
		}

		if err = copyAttach(w, a.Data); err != nil {
			return
		}
		if err = writeEOL(w); err != nil {
			return
		}
	}
	return
}
