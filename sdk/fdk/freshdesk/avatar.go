package freshdesk

type Avatar struct {
	ID int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Size int64 `json:"size,omitempty"`

	ContentType string `json:"content_type,omitempty"`

	AvatarURL string `json:"avatar_url,omitempty"`

	AttachmentURL string `json:"attachment_url,omitempty"`

	ThumbURL string `json:"thumb_url,omitempty"`

	CreatedAt Time `json:"created_at,omitempty"`

	UpdatedAt Time `json:"updated_at,omitempty"`

	// data attachment data
	data []byte

	// file attachment file
	file string
}

func (a *Avatar) String() string {
	return toString(a)
}

func (a *Avatar) Field() string {
	return "avatar"
}

func (a *Avatar) File() string {
	return a.file
}

func (a *Avatar) Data() []byte {
	return a.data
}

func NewAvatar(file string, data ...[]byte) *Avatar {
	a := &Avatar{
		file: file,
	}
	if len(data) > 0 {
		a.data = data[0]
	}
	return a
}
