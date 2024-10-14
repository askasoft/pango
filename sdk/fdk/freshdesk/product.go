package freshdesk

type Product struct {
	ID int64 `json:"id,omitempty"`

	// Name of the product
	Name string `json:"name,omitempty"`

	// Description of the product
	Description string `json:"description,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (p *Product) String() string {
	return toString(p)
}
