package freshservice

type AgentField struct {
	ID int64 `json:"id,omitempty"`

	// Name of the agent field.
	Name string `json:"name,omitempty"`

	// Field label shown to admins while editing an agent’s profile
	LabelForAdmins string `json:"label_for_admins,omitempty"`

	// Signifies whether admins can edit the field value.
	EditableByAdmins bool `json:"editable_by_admins,omitempty"`

	// Signifies whether admins have to necessarily provide a value for this field.
	MandatoryForAdmins bool `json:"mandatory_for_admins,omitempty"`

	// Field label shown to agents while viewing their own profile
	LabelForAgents string `json:"label_for_agents,omitempty"`

	// Signifies whether the field is shown to agents while viewing their own profile
	VisibleTo_Agents bool `json:"visible_to_agents,omitempty"`

	// Signifies whether the field is editable by agents while viewing their own profile
	EditableByAgents bool `json:"editable_by_agents,omitempty"`

	// Signifies whether agents have to necessarily provide a value for this field while viewing their own profile.
	MandatoryForAgents bool `json:"mandatory_for_agents,omitempty"`

	// Type of data stored in this field
	Type string `json:"type,omitempty"`

	// Indicates whether it is a default field or user-defined custom field
	DefaultField bool `json:"default_field,omitempty"`

	CreatedAt Time `json:"created_at,omitempty"`

	UpdatedAt Time `json:"updated_at,omitempty"`
}

func (af *AgentField) String() string {
	return toString(af)
}

type agentFieldsResult struct {
	AgentFields []*AgentField `json:"agent_fields,omitempty"`
}
