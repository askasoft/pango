package freshdesk

type Category struct {
	ID int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`

	// List of portal IDs where this category is visible
	VisibleInPortals []int64 `json:"visible_in_portals,omitempty"`

	CreatedAt Time `json:"created_at,omitempty"`

	UpdatedAt Time `json:"updated_at,omitempty"`
}

func (c *Category) String() string {
	return toString(c)
}

type CategoryCreate struct {
	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`

	// List of portal IDs where this category is visible
	VisibleInPortals []int64 `json:"visible_in_portals,omitempty"`
}

func (c *CategoryCreate) String() string {
	return toString(c)
}

type CategoryUpdate = CategoryCreate
