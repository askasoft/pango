package email

import (
	"errors"
	"io"
)

// Attachment email attachment
type Attachment struct {
	Cid  string
	Name string
	Data io.Reader
}

// NewAttachment new a email attachment
func NewAttachment(cid string, name string, data io.Reader) (*Attachment, error) {
	if len(name) == 0 {
		return nil, errors.New("Empty name of email attachment")
	}
	if data == nil {
		return nil, errors.New("Null data of email attachment")
	}

	return &Attachment{Cid: cid, Name: name, Data: data}, nil
}

// IsInline return true if cid is not empty
func (a *Attachment) IsInline() bool {
	return len(a.Cid) > 0
}
