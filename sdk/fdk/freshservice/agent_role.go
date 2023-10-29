package freshservice

type AgentRole struct {
	ID int64 `json:"id,omitempty"`

	// Name of the role
	Name string `json:"name,omitempty"`

	// Description of the role.
	Description string `json:"description,omitempty"`

	// default	boolean	Set to true if it is a default role, and false otherwise
	Default bool `json:"default,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (ar *AgentRole) String() string {
	return toString(ar)
}

type agentRoleResult struct {
	Role *AgentRole `json:"role,omitempty"`
}

type agentRolesResult struct {
	Roles []*AgentRole `json:"roles,omitempty"`
}
