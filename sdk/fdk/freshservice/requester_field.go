package freshservice

type RequesterField struct {
	ID int64 `json:"id,omitempty"`

	// Label of the requester field.
	Label string `json:"label,omitempty"`

	// Name of the requester field.
	Name string `json:"name,omitempty"`

	// Position of the requester field.
	Position int `json:"position,omitempty"`

	// For custom requester fields, type of value associated with the field will be given (Examples custom_date, custom_text...).
	Type string `json:"type,omitempty"`

	// ID of the user field group.
	UserFieldGroupID int64 `json:"user_field_group_id,omitempty"`

	// Set to true if the field can be updated by requesters during signup.
	EditableInSignup bool `json:"editable_in_signup,omitempty"`

	// Requesters can edit the field in the support portal.
	RequestersCanEdit bool `json:"requesters_can_edit,omitempty"`

	// Set to true if the field is mandatory in the support portal.
	RequiredForRequesters bool `json:"required_for_requesters,omitempty"`

	// Display name for the field (as seen in the support portal).
	LabelForRequesters string `json:"label_for_requesters,omitempty"`

	// Requesters can see the field in the support portal.
	DisplayedForRequesters bool `json:"displayed_for_requesters,omitempty"`

	// Set to true if the field is mandatory for agents.
	RequiredForAgents bool `json:"required_for_agents,omitempty"`

	// List of values supported by the field.
	Choices any `json:"choices,omitempty"`

	CreatedAt Time `json:"created_at,omitempty"`

	UpdatedAt Time `json:"updated_at,omitempty"`
}

func (rf *RequesterField) String() string {
	return toString(rf)
}

type requesterFieldsResult struct {
	RequesterFields []*RequesterField `json:"requester_fields,omitempty"`
}
