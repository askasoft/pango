package freshservice

import (
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

type FolderVisibility int

const (
	FolderVisibilityAllUsers      FolderVisibility = 1
	FolderVisibilityLoggedInUsers FolderVisibility = 2
	FolderVisibilityAgentsOnly    FolderVisibility = 3
	FolderVisibilityDepartments   FolderVisibility = 4
	FolderVisibilityAgentGroups   FolderVisibility = 5
	FolderVisibilityContactGroups FolderVisibility = 6
)

func (fv FolderVisibility) String() string {
	switch fv {
	case FolderVisibilityAllUsers:
		return "AllUsers"
	case FolderVisibilityLoggedInUsers:
		return "LoggedInUsers"
	case FolderVisibilityAgentsOnly:
		return "AgentsOnly"
	case FolderVisibilityDepartments:
		return "Departments"
	case FolderVisibilityAgentGroups:
		return "AgentGroups"
	case FolderVisibilityContactGroups:
		return "ContactGroups"
	default:
		return num.Itoa(int(fv))
	}
}

func ParseFolderVisibility(s string) FolderVisibility {
	switch str.ToLower(s) {
	case "allusers":
		return FolderVisibilityAllUsers
	case "loggedinusers":
		return FolderVisibilityLoggedInUsers
	case "agentsonly":
		return FolderVisibilityAgentsOnly
	case "departments":
		return FolderVisibilityDepartments
	case "agentgroups":
		return FolderVisibilityAgentGroups
	case "contactgroups":
		return FolderVisibilityContactGroups
	default:
		return 0
	}
}

type Folder struct {
	ID int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`

	// Describes the position in which the folder is listed
	Position int `json:"position,omitempty"`

	// Set as true is it is a default folder
	DefaultFolder bool `json:"default_folder,omitempty"`

	CategoryID int64 `json:"category_id,omitempty"`

	// Accessibility of this folder. Please refer to Folder Properties table.
	Visibility FolderVisibility `json:"visibility,omitempty"`

	// Approval settings that have been associated with the folder. Key-value pair containing the approval_type, approval_ids and its values.
	ApprovalSettings map[string]string `json:"approval_settings,omitempty"`

	// ID of the department to which this solution folder is visible. ( Mandatory if visibility is set to '4')
	DepartmentIDs []int64 `json:"department_ids,omitempty"`

	// ID of the Agent Groups to which this solution folder is visible. ( Mandatory if visibility is set to '5')
	GroupIDs []int64 `json:"group_ids,omitempty"`

	// ID of the Contact Groups to which this solution folder is visible. ( Mandatory if visibility is set to '6')
	RequesterGroupIDs []int64 `json:"requester_group_ids,omitempty"`

	ManageByGroupIDs []int64 `json:"manage_by_group_ids,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (f *Folder) String() string {
	return toString(f)
}

type folderResult struct {
	Foler   *Folder   `json:"folder,omitempty"`
	Folders []*Folder `json:"folders,omitempty"`
}
