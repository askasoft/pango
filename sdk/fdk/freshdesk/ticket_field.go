package freshdesk

type TicketField struct {
	ID int64 `json:"id,omitempty"`

	// Name of the ticket field.
	Name string `json:"name,omitempty"`

	// Label of the field for display
	Label string `json:"label,omitempty"`

	// Description of the field
	Description string `json:"description,omitempty"`

	// True if the field is a default field. False if customm
	Default bool `json:"default,omitempty"`

	// Position in which the ticket field is displayed in the form
	Position int `json:"field_type,omitempty"`

	// For custom ticket fields, The type of value associated with the field will be given (Examples custom_date, custom_text...)
	Type string `json:"type,omitempty"`

	// True if the field is marked mandatory while closing the Release item
	RequiredForClosure bool `json:"required_for_closure,omitempty"`

	// Set to true if the field is mandatory for Agents
	RequiredForAgents bool `json:"required_for_agents,omitempty"`

	// Set to true if the field is mandatory in the customer portal
	RequiredForCustomers bool `json:"required_for_customer,omitempty"`

	// Display name for the field (as seen in the customer portal)
	LabelForCustomers string `json:"label_for_customers,omitempty"`

	// Set to true if the field can be updated by customers
	CustomersCanEdit bool `json:"customers_can_edit,omitempty"`

	// Set to true if the field is displayed in the customer portal
	DisplayedToCustomers bool `json:"displayed_to_customers,omitempty"`

	// Applicable only for the requester field. Set to true if customer can add additional requesters to a ticket
	PortalCc bool `json:"portal_cc,omitempty"`

	// Applicable only if portal_cc is set to true. Value will be all when a customer can add any requester to the CC list and company when a customer can add only company contacts to the CC list
	PortalCcTo string `json:"portal_cc_to,omitempty"`

	// List of values supported by the field.
	Choices any `json:"choices,omitempty"`

	// Applicable only for dependent fields, this contains details of nested fields
	DependentFields any `json:"dependent_fields,omitempty"`

	// Applicable only if the field is part of a section. This contains the details of a section (ID, position) for which it is been a part of
	SectionMappings any `json:"SectionMappings,omitempty"`

	// True if the Ticket field is inside FSM section (Applicable only if FSM is enabled)
	IsFsm bool `json:"is_fsm,omitempty"`

	// True if the choice update is in progress (Applicable for the fields which has 100+ choices)
	FieldUpdateInProgress bool `json:"field_update_in_progress,omitempty"`

	CreatedAt Time `json:"created_at,omitempty"`

	UpdatedAt Time `json:"updated_at,omitempty"`
}

func (tf *TicketField) String() string {
	return toString(tf)
}

type TicketFieldCreate struct {
	// Label of the field for display
	Label string `json:"label,omitempty"`

	// True if the field is a default field. False if customm
	Default bool `json:"default,omitempty"`

	// Position in which the ticket field is displayed in the form
	Position int `json:"field_type,omitempty"`

	// For custom ticket fields, The type of value associated with the field will be given (Examples custom_date, custom_text...)
	Type string `json:"type,omitempty"`

	// True if the field is marked mandatory while closing the Release item
	RequiredForClosure bool `json:"required_for_closure,omitempty"`

	// Set to true if the field is mandatory for Agents
	RequiredForAgents bool `json:"required_for_agents,omitempty"`

	// Set to true if the field is mandatory in the customer portal
	RequiredForCustomers bool `json:"required_for_customer,omitempty"`

	// Display name for the field (as seen in the customer portal)
	LabelForCustomers string `json:"label_for_customers,omitempty"`

	// Set to true if the field can be updated by customers
	CustomersCanEdit bool `json:"customers_can_edit,omitempty"`

	// Set to true if the field is displayed in the customer portal
	DisplayedToCustomers bool `json:"displayed_to_customers,omitempty"`

	// List of values supported by the field.
	Choices any `json:"choices,omitempty"`

	// Applicable only for dependent fields, this contains details of nested fields
	DependentFields any `json:"dependent_fields,omitempty"`

	// Applicable only if the field is part of a section. This contains the details of a section (ID, position) for which it is been a part of
	SectionMappings any `json:"SectionMappings,omitempty"`
}

func (tf *TicketFieldCreate) String() string {
	return toString(tf)
}

type TicketFieldUpdate = TicketFieldCreate
