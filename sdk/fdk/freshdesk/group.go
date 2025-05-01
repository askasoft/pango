package freshdesk

type AutoTicketAssign int
type GroupUnassignedFor string

const (
	AutoTicketAssignDisabled             AutoTicketAssign = 0
	AutoTicketAssignRoundRobin           AutoTicketAssign = 1
	AutoTicketAssignSkillBasedRoundRobin AutoTicketAssign = 2
	AutoTicketAssignLoadBasedRoundRobin  AutoTicketAssign = 3
	AutoTicketAssignOmniroute            AutoTicketAssign = 4

	GroupUnassignedFor30m GroupUnassignedFor = "30m"
	GroupUnassignedFor1h  GroupUnassignedFor = "1h"
	GroupUnassignedFor2h  GroupUnassignedFor = "2h"
	GroupUnassignedFor4h  GroupUnassignedFor = "4h"
	GroupUnassignedFor8h  GroupUnassignedFor = "8h"
	GroupUnassignedFor12h GroupUnassignedFor = "12h"
	GroupUnassignedFor1d  GroupUnassignedFor = "1d"
	GroupUnassignedFor2d  GroupUnassignedFor = "2d"
	GroupUnassignedFor3d  GroupUnassignedFor = "3d"
)

type Group struct {
	ID int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`

	// Array of agent user IDs separated by commas. Instructions on finding an agent's user ID can be found here.
	AgentIDs []int64 `json:"agent_ids,omitempty"`

	// Describes the type of automatic ticket assignment set for the group. Automatic ticket assignment is only available on certain plans.
	AutoTicketAssign AutoTicketAssign `json:"auto_ticket_assign,omitempty"`

	// Unique ID of the business hour associated with the group
	BusinessHourID int64 `json:"business_hour_id,omitempty"`

	// The ID of the user to whom an escalation email is sent if a ticket is unassigned. To create/update a group with an escalate_to value of 'none', please set the value of this parameter to 'null'
	EscalateTo int64 `json:"escalate_to,omitempty"`

	// The time after which an escalation email is sent if a ticket remains unassigned. The accepted values are "30m" for 30 minutes, "1h" for 1 hour, "2h" for 2 hours, "4h" for 4 hours, "8h" for 8 hours, "12h" for 12 hours, "1d" for 1 day, "2d" for 2 days, and "3d" for 3 days
	UnassignedFor GroupUnassignedFor `json:"unassigned_for,omitempty"`

	CreatedAt Time `json:"created_at,omitempty"`

	UpdatedAt Time `json:"updated_at,omitempty"`
}

func (g *Group) String() string {
	return toString(g)
}

type GroupCreate struct {
	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`

	// Array of agent user IDs separated by commas. Instructions on finding an agent's user ID can be found here.
	AgentIDs []int64 `json:"agent_ids,omitempty"`

	// Describes the type of automatic ticket assignment set for the group. Automatic ticket assignment is only available on certain plans.
	AutoTicketAssign AutoTicketAssign `json:"auto_ticket_assign,omitempty"`

	// The ID of the user to whom an escalation email is sent if a ticket is unassigned. To create/update a group with an escalate_to value of 'none', please set the value of this parameter to 'null'
	EscalateTo int64 `json:"escalate_to,omitempty"`

	// The time after which an escalation email is sent if a ticket remains unassigned. The accepted values are "30m" for 30 minutes, "1h" for 1 hour, "2h" for 2 hours, "4h" for 4 hours, "8h" for 8 hours, "12h" for 12 hours, "1d" for 1 day, "2d" for 2 days, and "3d" for 3 days
	UnassignedFor GroupUnassignedFor `json:"unassigned_for,omitempty"`
}

func (g *GroupCreate) String() string {
	return toString(g)
}

type GroupUpdate = GroupCreate
