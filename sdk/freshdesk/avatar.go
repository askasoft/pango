package freshdesk

type Avatar struct {
	ID string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Size int64 `json:"size,omitempty"`

	ContentType string `json:"content_type,omitempty"`

	AvatarURL string `json:"avatar_url,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (a *Avatar) String() string {
	return toString(a)
}
