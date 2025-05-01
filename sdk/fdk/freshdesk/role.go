package freshdesk

type Role struct {
	ID int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`

	// Set to true if this is the default role
	Default bool `json:"default,omitempty"`

	CreatedAt Time `json:"created_at,omitempty"`

	UpdatedAt Time `json:"updated_at,omitempty"`
}

func (r *Role) String() string {
	return toString(r)
}
