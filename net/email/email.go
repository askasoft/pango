package email

import (
	"fmt"
	"net/mail"
	"path/filepath"
	"strings"
	"time"

	"github.com/pandafw/pango/iox"
	"github.com/pandafw/pango/osu"
	"github.com/pandafw/pango/str"
)

// Email email struct
type Email struct {
	MsgID       string
	Date        time.Time
	Subject     string
	Message     string
	HTML        bool
	Attachments []*Attachment

	sender *mail.Address
	from   *mail.Address
	tos    []*mail.Address
	ccs    []*mail.Address
	bccs   []*mail.Address
	replys []*mail.Address

	DkimDomain     string
	DkimSelector   string
	DkimPrivateKey string
}

// ParseAddress parse email address
func ParseAddress(s string) (*mail.Address, error) {
	a, err := mail.ParseAddress(s)
	if err != nil {
		return nil, fmt.Errorf("Invalid address %s - %w", s, err)
	}
	return a, nil
}

// GetDate return the email date
func (m *Email) GetDate() time.Time {
	if m.Date.IsZero() {
		m.Date = time.Now()
	}
	return m.Date
}

// GetSender get sender
func (m *Email) GetSender() string {
	if m.sender == nil {
		m.sender = m.from
	}
	if m.sender == nil {
		return ""
	}
	return m.sender.Address
}

// SetSender set sender
func (m *Email) SetSender(s string) error {
	a, err := ParseAddress(s)
	if err != nil {
		m.sender = a
	}
	return err
}

// GetFrom get from
func (m *Email) GetFrom() *mail.Address {
	return m.from
}

// SetFrom set from
func (m *Email) SetFrom(s string) error {
	a, err := ParseAddress(s)
	if err == nil {
		m.from = a
	}
	return err
}

// GetTos get to address array
func (m *Email) GetTos() []*mail.Address {
	return m.tos
}

// AddTo add to address
func (m *Email) AddTo(tos ...string) error {
	for _, s := range tos {
		a, err := ParseAddress(s)
		if err != nil {
			return err
		}
		m.tos = append(m.tos, a)
	}
	return nil
}

// GetCcs get cc address array
func (m *Email) GetCcs() []*mail.Address {
	return m.ccs
}

// AddCc add cc address
func (m *Email) AddCc(ccs ...string) error {
	for _, s := range ccs {
		a, err := ParseAddress(s)
		if err != nil {
			return err
		}
		m.ccs = append(m.ccs, a)
	}
	return nil
}

// GetBccs get bcc address array
func (m *Email) GetBccs() []*mail.Address {
	return m.bccs
}

// AddBcc add bcc address
func (m *Email) AddBcc(bccs ...string) error {
	for _, s := range bccs {
		a, err := ParseAddress(s)
		if err != nil {
			return err
		}
		m.bccs = append(m.bccs, a)
	}
	return nil
}

// GetReplys get reply address array
func (m *Email) GetReplys() []*mail.Address {
	return m.replys
}

// AddReply add reply address
func (m *Email) AddReply(rs ...string) error {
	for _, s := range rs {
		a, err := ParseAddress(s)
		if err != nil {
			return err
		}
		m.replys = append(m.replys, a)
	}
	return nil
}

// AddAttachment add a attachment
func (m *Email) AddAttachment(as ...*Attachment) {
	m.Attachments = append(m.Attachments, as...)
}

// AttachFile attach a file
func (m *Email) AttachFile(paths ...string) error {
	for _, path := range paths {
		if err := osu.FileExists(path); err != nil {
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
	if err := osu.FileExists(path); err != nil {
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
	rs := make(map[string]bool, len(m.tos)+len(m.ccs)+len(m.bccs))
	for _, a := range m.tos {
		rs[a.Address] = true
	}
	for _, a := range m.ccs {
		rs[a.Address] = true
	}
	for _, a := range m.bccs {
		rs[a.Address] = true
	}

	as := make([]string, 0, len(rs))
	for a := range rs {
		as = append(as, a)
	}

	return as
}

// GetRcptsByDomain return a (domain => address string array) map
func (m *Email) GetRcptsByDomain() map[string][]string {
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
