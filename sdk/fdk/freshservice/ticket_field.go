package freshservice

type TicketField struct {
	ID int64 `json:"id,omitempty"`

	// ID of the workspace to which this ticket belongs. The attribute is applicable only for accounts with the 'Workspaces' feature enabled.
	WorkspaceID int64 `json:"workspace_id,omitempty"`

	// Name of the ticket field.
	Name string `json:"name,omitempty"`

	// Label of the field for display
	Label string `json:"label,omitempty"`

	// Description of the field
	Description string `json:"description,omitempty"`

	// Indicates if the field is a checkbox, dropdown, text field
	FieldType string `json:"field_type,omitempty"`

	// True if the field is marked mandatory
	Required bool `json:"required,omitempty"`

	// True if the field is marked mandatory while closing the Release item
	RequiredForClosure bool `json:"required_for_closure,omitempty"`

	// True if the field is a default field. False if customm
	DefaultField bool `json:"default_field,omitempty"`

	// List of values supported by the field.
	Choices any `json:"choices,omitempty"`

	// contain details of nested fields
	NestedFields any `json:"nested_fields,omitempty"`

	// Set to true if the field is mandatory for Agents
	RequiredForAgents bool `json:"required_for_agents,omitempty"`

	// Set to true if the field is mandatory in the customer portal
	RequiredForCustomers bool `json:"required_for_customer,omitempty"`

	// Display name for the field (as seen in the customer portal)
	LabelForCustomers bool `json:"label_for_customers,omitempty"`

	// Set to true if the field can be updated by customers
	CustomersCanEdit bool `json:"customer_can_edit,omitempty"`

	// Set to true if the field is displayed in the customer portal
	DisplayedToCustomers bool `json:"displayed_to_customers,omitempty"`

	// Applicable only for the requester field. Set to true if customer can add additional requesters to a ticket
	PortalCc bool `json:"portal_cc,omitempty"`

	// Applicable only if portal_cc is set to true. Value will be all when a customer can add any requester to the CC list and company when a customer can add only company contacts to the CC list
	PortalCcTo bool `json:"portal_cc_to,omitempty"`

	// Applicable only for custom_date field. When ‘Request time information’ is selected, the field date_only returns false and if unchecked, date_only returns true
	DateOnly bool `json:"date_only,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (tf *TicketField) String() string {
	return toString(tf)
}

type ticketFieldResult struct {
	TicketField  *TicketField   `json:"ticket_field,omitempty"`
	TicketFields []*TicketField `json:"ticket_fields,omitempty"`
}
