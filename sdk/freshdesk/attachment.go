package freshdesk

type Attachment struct {
	ID int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Size int `json:"size,omitempty"`

	ContentType string `json:"content_type,omitempty"`

	AttachmentURL string `json:"attachment_url,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`

	// data attachment data
	data []byte

	// path attachment file path
	path string
}

func NewAttachment(path string, data ...[]byte) *Attachment {
	a := &Attachment{
		path: path,
	}
	if len(data) > 0 {
		a.data = data[0]
	}
	return a
}

func (a *Attachment) String() string {
	return toString(a)
}
