package freshservice

const (
	WorkspaceStateActive = "active"
	WorkspaceStateDraft  = "draft"
)

type Workspace struct {
	ID int64 `json:"id,omitempty"`

	// Name of the workspace
	Name string `json:"name,omitempty"`

	// Description of the workspace
	Description string `json:"description,omitempty"`

	// Signifies if the workspace is assigned as the primary workspace of the account
	Primary bool `json:"primary,omitempty"`

	// Signifies if the workspace is marked as restricted
	Restricted bool `json:"restricted,omitempty"`

	// Status of the workspace. Possible values:
	// active: The workspace is in active state
	// draft: The workspace is in draft state
	State string `json:"state,omitempty"`

	// The template from which the workspace was created
	TemplateName string `json:"template_name,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (w *Workspace) String() string {
	return toString(w)
}

type workspaceResult struct {
	Workspace *Workspace `json:"workspace,omitempty"`
}

type workspacesResult struct {
	Workspaces []*Workspace `json:"workspaces,omitempty"`
}
