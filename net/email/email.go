package email

import (
	"fmt"
	"net/mail"
	"path/filepath"
	"strings"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/iox/fsu"
	"github.com/askasoft/pango/str"
)

// Email email struct
type Email struct {
	MsgID       string
	Date        time.Time
	Subject     string
	Message     string
	HTML        bool
	Attachments []*Attachment

	Sender   *mail.Address
	From     *mail.Address
	Tos      []*mail.Address
	Ccs      []*mail.Address
	Bccs     []*mail.Address
	ReplyTos []*mail.Address

	DkimDomain     string
	DkimSelector   string
	DkimPrivateKey string
}

// ParseAddress parse email address
func ParseAddress(s string) (*mail.Address, error) {
	a, err := mail.ParseAddress(s)
	if err != nil {
		return nil, fmt.Errorf("%w - %s", err, s)
	}
	return a, nil
}

// ParseAddressList parse email address list (comma separate)
func ParseAddressList(s string) ([]*mail.Address, error) {
	as, err := mail.ParseAddressList(s)
	if err != nil {
		return nil, fmt.Errorf("%w - %s", err, s)
	}
	return as, nil
}

// ParseAddresses parse email addresses
func ParseAddresses(ss ...string) ([]*mail.Address, error) {
	as := make([]*mail.Address, len(ss))
	for i, s := range ss {
		a, err := ParseAddress(s)
		if err != nil {
			return nil, err
		}
		as[i] = a
	}
	return as, nil
}

// GetDate return the email date
func (m *Email) GetDate() time.Time {
	if m.Date.IsZero() {
		return time.Now()
	}
	return m.Date
}

// SenderAddress get sender address
func (m *Email) SenderAddress() string {
	s := m.Sender
	if s == nil {
		s = m.From
	}
	if s == nil {
		return ""
	}
	return s.Address
}

// SetSender set sender
func (m *Email) SetSender(s string) error {
	a, err := ParseAddress(s)
	if err != nil {
		m.Sender = a
	}
	return err
}

// SetFrom set the From address
func (m *Email) SetFrom(s string) error {
	a, err := ParseAddress(s)
	if err == nil {
		m.From = a
	}
	return err
}

// AddTo add To address
func (m *Email) AddTo(tos ...string) error {
	as, err := ParseAddresses(tos...)
	if err == nil {
		m.Tos = append(m.Tos, as...)
	}
	return err
}

// SetTo set To address list (comma separate)
func (m *Email) SetTo(to string) error {
	as, err := ParseAddressList(to)
	if err == nil {
		m.Tos = as
	}
	return err
}

// AddCc add Cc address
func (m *Email) AddCc(ccs ...string) error {
	as, err := ParseAddresses(ccs...)
	if err == nil {
		m.Ccs = append(m.Ccs, as...)
	}
	return err
}

// SetCc set Cc address list (comma separate)
func (m *Email) SetCc(cc string) error {
	as, err := ParseAddressList(cc)
	if err == nil {
		m.Ccs = as
	}
	return err
}

// AddBcc add Bcc address
func (m *Email) AddBcc(bccs ...string) error {
	as, err := ParseAddresses(bccs...)
	if err == nil {
		m.Bccs = append(m.Bccs, as...)
	}
	return err
}

// SetBcc set Bcc address list (comma separate)
func (m *Email) SetBcc(bcc string) error {
	as, err := ParseAddressList(bcc)
	if err == nil {
		m.Bccs = as
	}
	return err
}

// AddReplyTo add Reply-To address
func (m *Email) AddReplyTo(rs ...string) error {
	as, err := ParseAddresses(rs...)
	if err == nil {
		m.ReplyTos = append(m.ReplyTos, as...)
	}
	return err
}

// SetReplyTo set Reply-To address list (comma separate)
func (m *Email) SetReply(rs string) error {
	as, err := ParseAddressList(rs)
	if err == nil {
		m.ReplyTos = as
	}
	return err
}

// AddAttachment add a attachment
func (m *Email) AddAttachment(as ...*Attachment) {
	m.Attachments = append(m.Attachments, as...)
}

// AttachFile attach a file
func (m *Email) AttachFile(paths ...string) error {
	for _, path := range paths {
		if err := fsu.FileExists(path); err != nil {
			return err
		}
		m.AddAttachment(&Attachment{Name: filepath.Base(path), Data: &iox.FileReader{Path: path}})
	}
	return nil
}

// SetTextMsg set the text message
func (m *Email) SetTextMsg(msg string) {
	m.HTML = false
	m.Message = msg
}

// SetHTMLMsg set the HTML message
func (m *Email) SetHTMLMsg(msg string) {
	m.HTML = true
	m.Message = msg
}

// AttachString attach a string date file
func (m *Email) AttachString(name string, data string) {
	m.AddAttachment(&Attachment{Name: name, Data: strings.NewReader(data)})
}

// EmbedFile embed a file
func (m *Email) EmbedFile(cid string, path string) error {
	if err := fsu.FileExists(path); err != nil {
		return err
	}
	m.AddAttachment(&Attachment{Cid: cid, Name: filepath.Base(path), Data: &iox.FileReader{Path: path}})
	return nil
}

// EmbedString embed a string date file
func (m *Email) EmbedString(cid string, name string, data string) {
	m.AddAttachment(&Attachment{Cid: cid, Name: name, Data: strings.NewReader(data)})
}

// GetRecipients get all recipients
func (m *Email) GetRecipients() []string {
	rs := make(map[string]struct{}, len(m.Tos)+len(m.Ccs)+len(m.Bccs))
	for _, a := range m.Tos {
		rs[a.Address] = struct{}{}
	}
	for _, a := range m.Ccs {
		rs[a.Address] = struct{}{}
	}
	for _, a := range m.Bccs {
		rs[a.Address] = struct{}{}
	}

	as := make([]string, 0, len(rs))
	for a := range rs {
		as = append(as, a)
	}
	return as
}

// GetRecipientsByDomain return a (domain => address string array) map
func (m *Email) GetRecipientsByDomain() map[string][]string {
	rcpts := m.GetRecipients()
	rs := make(map[string][]string, len(rcpts))
	for _, r := range rcpts {
		d := str.SubstrAfterByte(r, '@')
		if d == "" {
			continue
		}

		rs[d] = append(rs[d], r)
	}
	return rs
}
