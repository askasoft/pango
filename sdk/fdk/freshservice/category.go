package freshservice

type Category struct {
	ID int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`

	// The rank of the solution category in the category listing
	Position int `json:"position,omitempty"`

	// Set as true if the category is a default one
	DefaultCategory bool `json:"default_category,omitempty"`

	// List of portal IDs where this category is visible
	VisibleInPortals []int64 `json:"visible_in_portals,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (c *Category) String() string {
	return toString(c)
}

type categoryResult struct {
	Category *Category `json:"category,omitempty"`
}

type categoriesResult struct {
	Categories []*Category `json:"categories,omitempty"`
}
