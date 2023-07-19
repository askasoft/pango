package freshdesk

import (
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

type FolderVisibility int

const (
	FolderVisibilityAllUsers                FolderVisibility = 1
	FolderVisibilityLoggedInUsers           FolderVisibility = 2
	FolderVisibilityAgents                  FolderVisibility = 3
	FolderVisibilitySelectedCompanies       FolderVisibility = 4
	FolderVisibilityBots                    FolderVisibility = 5
	FolderVisibilitySelectedContactSegments FolderVisibility = 6
	FolderVisibilitySelectedCompanySegments FolderVisibility = 7
)

func (fv FolderVisibility) String() string {
	switch fv {
	case FolderVisibilityAllUsers:
		return "AllUsers"
	case FolderVisibilityLoggedInUsers:
		return "LoggedInUsers"
	case FolderVisibilityAgents:
		return "Agents"
	case FolderVisibilitySelectedCompanies:
		return "SelectedCompanies"
	case FolderVisibilityBots:
		return "Bots"
	case FolderVisibilitySelectedContactSegments:
		return "SelectedContactSegments"
	case FolderVisibilitySelectedCompanySegments:
		return "SelectedCompanySegments"
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
	case "agents":
		return FolderVisibilityAgents
	case "selectedcompanies":
		return FolderVisibilitySelectedCompanies
	case "bots":
		return FolderVisibilityBots
	case "selectedcontactsegments":
		return FolderVisibilitySelectedContactSegments
	case "selectedcompanysegments":
		return FolderVisibilitySelectedCompanySegments
	default:
		return 0
	}
}

type Folder struct {
	ID int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`

	ParentFolderID int64 `json:"parent_folder_id,omitempty"`

	// Parent category and folders in which the folder is placed
	Hierarchy []map[string]any `json:"hierarchy,omitempty"`

	// Number of articles present inside a folder
	ArticlesCount int `json:"articles_count,omitempty"`

	// Number of folders present inside a folder
	SubFoldersCount int `json:"sub_folders_count,omitempty"`

	// Accessibility of this folder. Please refer to Folder Properties table.
	Visibility FolderVisibility `json:"visibility,omitempty"`

	// IDs of the companies to whom this solution folder is visible
	CompanyIDs []int64 `json:"company_ids,omitempty"`

	// IDs of the contact segments to whom this solution folder is visible
	ContactSegmentIDs []int64 `json:"contact_segment_ids,omitempty"`

	// IDs of the company segments to whom this solution folder is visible
	CompanySegmentIDs []int64 `json:"company_segment_ids,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (f *Folder) String() string {
	return toString(f)
}
