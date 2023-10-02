package freshservice

import (
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

type AgentScoreboardLevel int

const (
	AgentScoreboardLevelBeginner     AgentScoreboardLevel = 1
	AgentScoreboardLevelIntermediate AgentScoreboardLevel = 2
	AgentScoreboardLevelProfessional AgentScoreboardLevel = 3
	AgentScoreboardLevelExpert       AgentScoreboardLevel = 4
	AgentScoreboardLevelMaster       AgentScoreboardLevel = 5
	AgentScoreboardLevelGuru         AgentScoreboardLevel = 6
)

func (asl AgentScoreboardLevel) String() string {
	switch asl {
	case AgentScoreboardLevelBeginner:
		return "Beginner"
	case AgentScoreboardLevelIntermediate:
		return "Intermediate"
	case AgentScoreboardLevelProfessional:
		return "Professional"
	case AgentScoreboardLevelExpert:
		return "Expert"
	case AgentScoreboardLevelMaster:
		return "Master"
	case AgentScoreboardLevelGuru:
		return "Guru"
	default:
		return num.Itoa(int(asl))
	}
}

func ParseAgentScoreboardLevel(s string) AgentScoreboardLevel {
	switch str.ToLower(s) {
	case "beginner":
		return AgentScoreboardLevelBeginner
	case "intermediate":
		return AgentScoreboardLevelIntermediate
	case "professional":
		return AgentScoreboardLevelProfessional
	case "expert":
		return AgentScoreboardLevelExpert
	case "master":
		return AgentScoreboardLevelMaster
	case "guru":
		return AgentScoreboardLevelGuru
	default:
		return AgentScoreboardLevel(num.Atoi(s))
	}
}

type AssignRole struct {
	// Unique ID of the role assigned
	RoleID int64 `json:"role_id,omitempty"`

	// The scope in which the agent can use the permissions granted by this role.
	// Possible values include entire_helpdesk (all plans),
	// member_groups (all plans; in the Pro and Enterprise plans, this also includes groups that the agent is an observer of),
	// specified_groups (Pro and Enterprise only), and assigned_items (all plans)
	AssignmentScope string `json:"assignment_scope,omitempty"`

	// Unique IDs of Groups in which the permissions granted by the role applies.
	// Mandatory only when the assignment_scope is specified_groups, and should be ignored otherwise.
	Groups []int64 `json:"groups,omitempty"`
}

func (ar *AssignRole) String() string {
	return toString(ar)
}

type Agent struct {
	ID int64 `json:"id,omitempty"`

	// First name of the agent
	FirstName string `json:"first_name,omitempty"`

	// Last name of the agent.
	LastName string `json:"last_name,omitempty"`

	// Set to true if this is an occasional agent (true => occasional, false => full-time)
	Occasional bool `json:"occasional,omitempty"`

	// Job title of the agent.
	JobTitle string `json:"job_title,omitempty"`

	// Email address of the agent.
	Email string `json:"email,omitempty"`

	// Work phone number of the agent.
	WorkPhoneNumber string `json:"work_phone_number,omitempty"`

	// Mobile phone number of the agent.
	MobilePhoneNumber string `json:"mobile_phone_number,omitempty"`

	// Unique IDs of the departments associated with the requester
	DepartmentIDs []int64 `json:"department_ids,omitempty"`

	// Set to true if the requester must be allowed to view tickets filed by other members of the department, and false otherwise
	CanSeeAllTicketsFromAssociatedDepartments bool `json:"can_see_all_tickets_from_associated_departments,omitempty"`

	// User ID of the requester’s reporting manager.
	ReportingManagerID int64 `json:"reporting_manager_id,omitempty"`

	// Address of the requester.
	Address string `json:"address,omitempty"`

	// Time zone of the requester.
	TimeZone string `json:"time_zone,omitempty"`

	// Time format for the requester.Possible values: 12h (12 hour format), 24h (24 hour format)
	TimeFormat string `json:"time_format,omitempty"`

	// Language used by the requester.
	Language string `json:"language,omitempty"`

	// Unique ID of the location associated with the requester.
	LocationID int64 `json:"location_id,omitempty"`

	// Background information of the requester.
	BackgroundInformation string `json:"background_information,omitempty"`

	// Unique ID of the level of the agent in the Arcade.
	ScoreboardLevelID AgentScoreboardLevel `json:"scoreboard_level_id,omitempty"`

	// Unique IDs of the groups that the agent is a member of.
	// The input value for this field should also include the IDs of approval-enabled restricted groups for which the agent’s member access is pending approval by a group leader.
	// The response value for this field would only contain the list of groups that the agent is an approved member of.
	// The member_of_pending_approval read-only attribute in the response will include the list of groups for which the agent’s member access is pending approval by a group leader.
	MemberOf []int64 `json:"member_of,omitempty"`

	// Unique IDs of the groups that the agent is an observer of.
	// The input value for this field should also include the IDs of approval-enabled restricted groups for which the agent’s observer access is pending approval by a group leader.
	// The response value for this field would only contain the list of groups that the agent is an approved observer of.
	// The observer_of_pending_approval read-only attribute in the response will include the list of groups for which the agent’s observer access is pending approval by a group leader.
	ObserverOf []int64 `json:"observer_of,omitempty"`

	// Unique IDs of the restricted groups to which the agent’s addition as a member is pending approval.
	MemberOfPendingApproval []int64 `json:"member_of_pending_approval,omitempty"`

	// Unique IDs of the restricted groups to which the agent’s addition as an observer is pending approval.
	ObserverOfPendingApproval []int64 `json:"observer_of_pending_approval,omitempty"`

	Roles []*AssignRole `json:"roles,omitempty"`

	// Timestamp of the agent's last successful login.
	LastLoginAt *Time `json:"last_login_at,omitempty"`

	// Timestamp of the agent's recent activity.
	LastActiveAt *Time `json:"last_active_at,omitempty"`

	// Key-value pair containing the names and values of the (custom) requester fields.
	CustomFields map[string]any `json:"custom_fields,omitempty"`

	// Set to true if the user has logged in to Freshservice at least once, and false otherwise.
	HasLoggedIn bool `json:"has_logged_in,omitempty"`

	// True if the agent is active, false if the agent has been deactivated.
	Active bool `json:"active,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (a *Agent) String() string {
	return toString(a)
}

type agentResult struct {
	Agent  *Agent   `json:"agent,omitempty"`
	Agents []*Agent `json:"agents,omitempty"`
}
