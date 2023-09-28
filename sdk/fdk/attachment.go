package fdk

type Attachment struct {
	ID int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Size int `json:"size,omitempty"`

	ContentType string `json:"content_type,omitempty"`

	AttachmentURL string `json:"attachment_url,omitempty"`

	CanonicalURL string `json:"canonical_url,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`

	// data attachment data
	data []byte

	// file attachment file
	file string
}

func (a *Attachment) String() string {
	return ToJSONIndent(a)
}

func (a *Attachment) Field() string {
	return "attachments[]"
}

func (a *Attachment) File() string {
	return a.file
}

func (a *Attachment) Data() []byte {
	return a.data
}

func NewAttachment(file string, data ...[]byte) *Attachment {
	a := &Attachment{
		file: file,
	}
	if len(data) > 0 {
		a.data = data[0]
	}
	return a
}

type Attachments []*Attachment

func (as Attachments) Files() Files {
	fs := make(Files, len(as))
	for i, a := range as {
		fs[i] = a
	}
	return fs
}
