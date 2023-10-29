package freshservice

type AgentGroupUnassignedFor string

const (
	AgentGroupUnassignedFor30m AgentGroupUnassignedFor = "30m"
	AgentGroupUnassignedFor1h  AgentGroupUnassignedFor = "1h"
	AgentGroupUnassignedFor2h  AgentGroupUnassignedFor = "2h"
	AgentGroupUnassignedFor4h  AgentGroupUnassignedFor = "4h"
	AgentGroupUnassignedFor8h  AgentGroupUnassignedFor = "8h"
	AgentGroupUnassignedFor12h AgentGroupUnassignedFor = "12h"
	AgentGroupUnassignedFor1d  AgentGroupUnassignedFor = "1d"
	AgentGroupUnassignedFor2d  AgentGroupUnassignedFor = "2d"
	AgentGroupUnassignedFor3d  AgentGroupUnassignedFor = "3d"
)

type AgentGroup struct {
	ID int64 `json:"id,omitempty"`

	// Name of the group
	Name string `json:"name,omitempty"`

	// Description of the group.
	Description string `json:"description,omitempty"`

	// The time after which an escalation email is sent if a ticket in the group remains unassigned.
	// The accepted values are “30m” for 30 minutes, “1h” for 1 hour, “2h” for 2 hours, “4h” for 4 hours, “8h” for 8 hours, “12h” for 12 hours, “1d” for 1 day, “2d” for 2 days, and “3d” for 3 days.
	UnassignedFor AgentGroupUnassignedFor `json:"unassigned_for,omitempty"`

	// Unique ID of the business hours configuration associated with the group.
	BusinessHoursID int64 `json:"business_hours_id,omitempty"`

	// The Unique ID of the user to whom an escalation email is sent if a ticket in this group is unassigned. To create/update a group with an escalate_to value of ‘none’, please set the value of this parameter to ‘null’.
	EscalateTo int64 `json:"escalate_to,omitempty"`

	// A comma separated array of user IDs of agents who are members of this group.
	// If the group is restricted and approvals-enabled, the input value for this field should also include the user IDs of agents whose member access to the group is pending approval by a group leader.
	// The response value for this field would only contain the list of approved members.
	// The members_pending_approval read-only attribute in the response will include the list of members whose approval is pending.
	Members []int64 `json:"members,omitempty"`

	// A comma separated array of user IDs of agents who are observers of this group.
	// If the group is restricted and approvals-enabled, the input value for this field should also include the user IDs of agents whose observer access to the group is pending approval by a group leader.
	// The response value for this field would only contain the list of approved observers.
	// The observers_pending_approval read-only attribute in the response will include the list of observers whose approval is pending.
	// This attribute is only applicable for accounts which have the “Access Controls Pro” feature enabled.
	Observers []int64 `json:"ovservers,omitempty"`

	// Signifies whether a group is marked as restricted.
	// This attribute won't be supported if the "Access Controls Pro" feature is unavailable for the account. Set to true if the group is restricted, and false otherwise.
	Restricted bool `json:"restricted,omitempty"`

	// A comma separated array of user IDs of agents who are leaders of this group.
	// If the group is restricted and approvals-enabled, the input value for this field should also include the user IDs of agents whose leader access to the group is pending approval by another group leader. The response value for this field would only contain the list of approved leaders.
	// The leaders_pending_approval read-only attribute in the response will include the list of leaders whose approval is pending.
	// This attribute is only applicable for accounts which have the “Access Controls Pro” feature enabled.
	Leaders []int64 `json:"leaders,omitempty"`

	// A comma-separated array of user IDs of agents whose member access to the group is pending approval by a group leader.
	// This attribute is only applicable for accounts which have the “Access Controls Pro” feature enabled.
	// READ ONLY
	MembersPendingApproval []int64 `json:"members_pending_approval,omitempty"`

	// A comma-separated array of user IDs of agents whose observer access to the group is pending approval by a group leader.
	// This attribute is only applicable for accounts which have the “Access Controls Pro” feature enabled.
	// READ ONLY
	ObserversPendingApproval []int64 `json:"observers_pending_approval,omitempty"`

	// A comma-separated array of user IDs of agents whose leader access to the group is pending approval by an existing group leader.
	// This attribute is only applicable for accounts which have the “Access Controls Pro” feature enabled.
	// READ ONLY
	LeadersPendingApproval []int64 `json:"leaders_pending_approval,omitempty"`

	// Signifies whether the restricted group requires approvals for membership changes.
	// This attribute is only applicable for accounts which have the “Access Controls Pro” feature enabled.
	ApprovalRequired bool `json:"approval_required,omitempty"`

	// Describes the automatic ticket assignment type.
	// Will not be supported if the "Round Robin" feature is disabled for the account.
	// Set to true if automatic ticket assignment is enabled, and false otherwise.
	AutoTicketAssign bool `json:"auto_ticket_assign,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (ag *AgentGroup) String() string {
	return toString(ag)
}

type agentGroupResult struct {
	Group *AgentGroup `json:"group,omitempty"`
}

type agentGroupsResult struct {
	Groups []*AgentGroup `json:"groups,omitempty"`
}
