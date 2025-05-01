package freshdesk

import (
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

type AgentTicketScope int

const (
	AgentTicketScopeGlobal     AgentTicketScope = 1
	AgentTicketScopeGroup      AgentTicketScope = 2
	AgentTicketScopeRestricted AgentTicketScope = 3
)

func (ats AgentTicketScope) String() string {
	switch ats {
	case AgentTicketScopeGlobal:
		return "Global"
	case AgentTicketScopeGroup:
		return "Group"
	case AgentTicketScopeRestricted:
		return "Restricted"
	default:
		return num.Itoa(int(ats))
	}
}

func ParseAgentTicketScope(s string) AgentTicketScope {
	switch str.ToLower(s) {
	case "global":
		return AgentTicketScopeGlobal
	case "group":
		return AgentTicketScopeGroup
	case "restricted":
		return AgentTicketScopeRestricted
	default:
		return 0
	}
}

type Agent struct {
	ID int64 `json:"id,omitempty"`

	// If the agent is in a group that has enabled "Automatic Ticket Assignment", this attribute will be set to true if the agent is accepting new tickets
	Available bool `json:"available,omitempty"`

	// Timestamp that denotes when the agent became available/unavailable (depending on the value of the 'available' attribute)
	AvailableSince *Time `json:"available_since,omitempty"`

	// Set to true if this is an occasional agent (true => occasional, false => full-time)
	Occasional bool `json:"occasional,omitempty"`

	// Signature of the agent in HTML format
	Signature string `json:"signature,omitempty"`

	// Ticket permission of the agent (1 -> Global Access, 2 -> Group Access, 3 -> Restricted Access)
	TicketScope AgentTicketScope `json:"ticket_scope,omitempty"`

	// Type of Agent (support_agent -> Support Agent, field_agent -> Field Agent, collaborator -> Collaborator)
	Type string `json:"type,omitempty"`

	// Skill ids associated with the agent
	SkillIDs []int64 `json:"skill_ids,omitempty"`

	// Group IDs associated with the agent
	GroupIDs []int64 `json:"group_ids,omitempty"`

	// Role IDs associated with the agent
	RoleIDs []int64 `json:"role_ids,omitempty"`

	Contact *Contact `json:"contact,omitempty"`

	CreatedAt Time `json:"created_at,omitempty"`

	UpdatedAt Time `json:"updated_at,omitempty"`
}

func (a *Agent) String() string {
	return toString(a)
}

type AgentCreate struct {
	// Set to true if this is an occasional agent (true => occasional, false => full-time)
	Occasional bool `json:"occasional,omitempty"`

	// Signature of the agent in HTML format
	Signature string `json:"signature,omitempty"`

	// Ticket permission of the agent (1 -> Global Access, 2 -> Group Access, 3 -> Restricted Access)
	TicketScope AgentTicketScope `json:"ticket_scope,omitempty"`

	// Type of Agent (support_agent -> Support Agent, field_agent -> Field Agent, collaborator -> Collaborator)
	AgentType string `json:"agent_type,omitempty"`

	// Skill ids associated with the agent
	SkillIDs []int64 `json:"skill_ids,omitempty"`

	// Group IDs associated with the agent
	GroupIDs []int64 `json:"group_ids,omitempty"`

	// Role IDs associated with the agent
	RoleIDs []int64 `json:"role_ids,omitempty"`

	// Name of the contact
	Name string `json:"name,omitempty"`

	// Primary email address of the contact. If you want to associate additional email(s) with this contact, use the other_emails attribute
	Email string `json:"email,omitempty"`

	// Telephone number of the contact
	Phone string `json:"phone,omitempty"`

	// Mobile number of the contact
	Mobile string `json:"mobile,omitempty"`

	// Job title of the contact
	JobTitle string `json:"job_title,omitempty"`

	// Language of the contact
	Language string `json:"language,omitempty"`

	// Time zone in which the contact resides
	TimeZone string `json:"time_zone,omitempty"`

	// Avatar of the contact
	//Avatar *Avatar `json:"avatar,omitempty"`
}

func (a *AgentCreate) String() string {
	return toString(a)
}

// func (aq *AgentRequest) Files() Files {
// 	if aq.Avatar != nil {
// 		return Files{aq.Avatar}
// 	}
// 	return nil
// }

// func (aq *AgentRequest) Values() Values {
// 	vs := Values{}

// 	vs.SetBool("occasional", aq.Occasional)
// 	vs.SetString("signature", aq.Signature)
// 	vs.SetInt("ticket_scope", aq.TicketScope)
// 	vs.SetString("agent_type", aq.AgentType)
// 	vs.SetInt64s("skill_ids", aq.SkillIDs)
// 	vs.SetInt64s("group_ids", aq.GroupIDs)
// 	vs.SetInt64s("role_ids", aq.RoleIDs)

// 	vs.SetString("name", aq.Name)
// 	vs.SetString("email", aq.Email)
// 	vs.SetString("phone", aq.Phone)
// 	vs.SetString("mobile", aq.Mobile)
// 	vs.SetString("job_title", aq.JobTitle)
// 	vs.SetString("language", aq.Language)
// 	vs.SetString("time_zone", aq.TimeZone)
// 	return vs
// }

type AgentUpdate = AgentCreate
