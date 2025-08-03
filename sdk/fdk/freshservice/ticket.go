package freshservice

import (
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

type TicketSource int
type TicketStatus int
type TicketPriority int
type TicketUrgency int
type TicketImpact int

const (
	TicketSourceEmail          TicketSource = 1
	TicketSourcePortal         TicketSource = 2
	TicketSourcePhone          TicketSource = 3
	TicketSourceChat           TicketSource = 4
	TicketSourceFeedbackWidget TicketSource = 5
	TicketSourceYammer         TicketSource = 6
	TicketSourceAWSCloudwatch  TicketSource = 7
	TicketSourcePagerduty      TicketSource = 8
	TicketSourceWalkup         TicketSource = 9
	TicketSourceSlack          TicketSource = 10

	TicketStatusOpen     TicketStatus = 2
	TicketStatusPending  TicketStatus = 3
	TicketStatusResolved TicketStatus = 4
	TicketStatusClosed   TicketStatus = 5

	TicketPriorityLow    TicketPriority = 1
	TicketPriorityMedium TicketPriority = 2
	TicketPriorityHigh   TicketPriority = 3
	TicketPriorityUrgent TicketPriority = 4

	TicketUrgencyLow    TicketUrgency = 1
	TicketUrgencyMedium TicketUrgency = 2
	TicketUrgencyHigh   TicketUrgency = 3

	TicketImpactLow    TicketImpact = 1
	TicketImpactMedium TicketImpact = 2
	TicketImpactHigh   TicketImpact = 3

	TicketFilterNewAndMyOpen = "new_and_my_open"
	TicketFilterWatching     = "watching"
	TicketFilterSpam         = "spam"
	TicketFilterDeleted      = "deleted"
)

func (ts TicketSource) String() string {
	switch ts {
	case TicketSourceEmail:
		return "Email"
	case TicketSourcePortal:
		return "Portal"
	case TicketSourcePhone:
		return "Phone"
	case TicketSourceChat:
		return "Chat"
	case TicketSourceFeedbackWidget:
		return "FeedbackWidget"
	case TicketSourceYammer:
		return "Yammer"
	case TicketSourceAWSCloudwatch:
		return "AWSCloudwatch"
	case TicketSourcePagerduty:
		return "Pagerduty"
	case TicketSourceWalkup:
		return "Walkup"
	case TicketSourceSlack:
		return "Slack"
	default:
		return num.Itoa(int(ts))
	}
}

func ParseTicketSource(s string) TicketSource {
	switch str.ToLower(s) {
	case "email":
		return TicketSourceEmail
	case "portal":
		return TicketSourcePortal
	case "phone":
		return TicketSourcePhone
	case "chat":
		return TicketSourceChat
	case "feedbackwidget":
		return TicketSourceFeedbackWidget
	case "yammer":
		return TicketSourceYammer
	case "awscloudwatch":
		return TicketSourceAWSCloudwatch
	case "pagerduty":
		return TicketSourcePagerduty
	case "walkup":
		return TicketSourceWalkup
	case "slack":
		return TicketSourceSlack
	default:
		return TicketSource(num.Atoi(s))
	}
}

func (ts TicketStatus) String() string {
	switch ts {
	case TicketStatusOpen:
		return "Open"
	case TicketStatusPending:
		return "Pending"
	case TicketStatusResolved:
		return "Resolved"
	case TicketStatusClosed:
		return "Closed"
	default:
		return num.Itoa(int(ts))
	}
}

func ParseTicketStatus(s string) TicketStatus {
	switch str.ToLower(s) {
	case "open":
		return TicketStatusOpen
	case "pending":
		return TicketStatusPending
	case "resolved":
		return TicketStatusResolved
	case "closed":
		return TicketStatusClosed
	default:
		return TicketStatus(num.Atoi(s))
	}
}

func (tp TicketPriority) String() string {
	switch tp {
	case TicketPriorityLow:
		return "Low"
	case TicketPriorityMedium:
		return "Medium"
	case TicketPriorityHigh:
		return "High"
	case TicketPriorityUrgent:
		return "Urgent"
	default:
		return num.Itoa(int(tp))
	}
}

func ParseTicketPriority(s string) TicketPriority {
	switch str.ToLower(s) {
	case "low":
		return TicketPriorityLow
	case "medium":
		return TicketPriorityMedium
	case "high":
		return TicketPriorityHigh
	case "urgent":
		return TicketPriorityUrgent
	default:
		return TicketPriority(num.Atoi(s))
	}
}

func (tu TicketUrgency) String() string {
	switch tu {
	case TicketUrgencyLow:
		return "Low"
	case TicketUrgencyMedium:
		return "Medium"
	case TicketUrgencyHigh:
		return "High"
	default:
		return num.Itoa(int(tu))
	}
}

func ParseTicketUrgency(s string) TicketUrgency {
	switch str.ToLower(s) {
	case "low":
		return TicketUrgencyLow
	case "medium":
		return TicketUrgencyMedium
	case "high":
		return TicketUrgencyHigh
	default:
		return TicketUrgency(num.Atoi(s))
	}
}

func (ti TicketImpact) String() string {
	switch ti {
	case TicketImpactLow:
		return "Low"
	case TicketImpactMedium:
		return "Medium"
	case TicketImpactHigh:
		return "High"
	default:
		return num.Itoa(int(ti))
	}
}

func ParseTicketImpact(s string) TicketImpact {
	switch str.ToLower(s) {
	case "low":
		return TicketImpactLow
	case "medium":
		return TicketImpactMedium
	case "high":
		return TicketImpactHigh
	default:
		return TicketImpact(num.Atoi(s))
	}
}

type Ticket struct {
	// Unique ID of the ticket
	ID int64 `json:"id,omitempty"`

	// ID of the workspace to which the ticket belongs.
	// The attribute is applicable only for accounts with the 'Workspaces' feature enabled.
	WorkspaceID int64 `json:"workspace_id,omitempty"`

	// ID of the department to which this ticket belongs.
	DepartmentID int64 `json:"department_id,omitempty"`

	// Ticket attachments. The total size of these attachments cannot exceed 40MB.
	Attachments []*Attachment `json:"attachments,omitempty"`

	// Name of the requester
	Name string `json:"name,omitempty"`

	// Email address of the requester. If no contact exists with this email address in Freshdesk, it will be added as a new contact.
	Email string `json:"email,omitempty"`

	// Phone number of the requester. If no contact exists with this phone number in Freshdesk, it will be added as a new contact. If the phone number is set and the email address is not, then the name attribute is mandatory.
	Phone string `json:"phone,omitempty"`

	// User ID of the requester. For existing contacts, the requester_id can be passed instead of the requester's email.
	RequesterID int64 `json:"requester_id,omitempty"`

	// ID of the agent to whom the ticket has been assigned.
	ResponderID int64 `json:"responder_id,omitempty"`

	// Helps categorize the ticket according to the different kinds of issues your support team deals with.
	Type string `json:"type,omitempty"`

	// Status of the ticket
	Status TicketStatus `json:"status,omitempty"`

	// Priority of the ticket
	Priority TicketPriority `json:"priority,omitempty"`

	// The channel through which the ticket was created
	Source TicketSource `json:"source,omitempty"`

	// Set to true if the ticket has been deleted/trashed. Deleted tickets will not be displayed in any views except the "deleted" filter
	Deleted bool `json:"deleted,omitempty"`

	// Set to true if the ticket has been marked as spam
	Spam bool `json:"spam,omitempty"`

	// Timestamp that denotes when the ticket is due to be resolved
	DueBy *Time `json:"due_by,omitempty"`

	// Timestamp that denotes when the first response is due
	FrDueBy *Time `json:"fr_due_by,omitempty"`

	// Set to true if the ticket has been escalated for any reason
	IsEscalated bool `json:"is_escalated,omitempty"`

	// Set to true if the ticket has been escalated as the result of first response time being breached
	FrEscalated bool `json:"fr_escalated,omitempty"`

	// Email addresses to which the ticket was originally sent
	ToEmails []string `json:"to_emails,omitempty"`

	// Email address added in the 'cc' field of the incoming ticket email
	CcEmails []string `json:"cc_emails,omitempty"`

	// Email address(e)s added while forwarding a ticket
	FwdEmails []string `json:"fwd_emails,omitempty"`

	// Email address added while replying to a ticket
	ReplyCcEmails []string `json:"reply_cc_emails,omitempty"`

	// ID of email config which is used for this ticket. (i.e., support@yourcompany.com/sales@yourcompany.com)
	EmailConfigID int64 `json:"email_config_id,omitempty"`

	// ID of the group to which the ticket has been assigned
	GroupID int64 `json:"group_id,omitempty"`

	// Subject of the ticket
	Subject string `json:"subject,omitempty"`

	// HTML content of the ticket
	Description string `json:"description,omitempty"`

	// Content of the ticket in plain text
	DescriptionText string `json:"description_text,omitempty"`

	// Tags that have been associated with the ticket
	Tags []string `json:"tags,omitempty"`

	// Key value pairs containing the names and values of custom fields.
	CustomFields map[string]any `json:"custom_fields,omitempty"`

	// include=requester
	Requester *Requester `json:"requester,omitempty"`

	// include=conversations
	Conversations []*Conversation `json:"conversations,omitempty"`

	// Ticket Category.
	Category string `json:"category,omitempty"`

	// Ticket sub category.
	SubCategory string `json:"sub_category,omitempty"`

	// Ticket item category.
	ItemCategory string `json:"item_category,omitempty"`

	// Ticket urgency.
	Urgency TicketUrgency `json:"urgency,omitempty"`

	// Ticket impact.
	Impact TicketImpact `json:"impact,omitempty"`

	// Content of the ticket resolution note in plain text
	ResolutionNotes string `json:"resolution_notes,omitempty"`

	// HTML content of the ticket resolution note
	ResolutionNotesHTML string `json:"resolution_notes_html,omitempty"`

	// Ticket creation timestamp
	CreatedAt Time `json:"created_at,omitempty"`

	// Ticket updated timestamp
	UpdatedAt Time `json:"updated_at,omitempty"`
}

func (t *Ticket) String() string {
	return toString(t)
}

type ticketResult struct {
	Ticket  *Ticket   `json:"ticket,omitempty"`
	Tickets []*Ticket `json:"tickets,omitempty"`
}

type TicketCreate struct {
	// Unique ID of the ticket
	ID int64 `json:"id,omitempty"`

	// ID of the workspace to which the ticket belongs.
	// The attribute is applicable only for accounts with the 'Workspaces' feature enabled.
	WorkspaceID int64 `json:"workspace_id,omitempty"`

	// ID of the department to which this ticket belongs.
	DepartmentID int64 `json:"department_id,omitempty"`

	// Ticket attachments. The total size of these attachments cannot exceed 40MB.
	Attachments []*Attachment `json:"attachments,omitempty"`

	// Name of the requester
	Name string `json:"name,omitempty"`

	// Email address of the requester. If no contact exists with this email address in Freshdesk, it will be added as a new contact.
	Email string `json:"email,omitempty"`

	// Phone number of the requester. If no contact exists with this phone number in Freshdesk, it will be added as a new contact. If the phone number is set and the email address is not, then the name attribute is mandatory.
	Phone string `json:"phone,omitempty"`

	// User ID of the requester. For existing contacts, the requester_id can be passed instead of the requester's email.
	RequesterID int64 `json:"requester_id,omitempty"`

	// ID of the agent to whom the ticket has been assigned.
	ResponderID int64 `json:"responder_id,omitempty"`

	// Helps categorize the ticket according to the different kinds of issues your support team deals with.
	Type string `json:"type,omitempty"`

	// Status of the ticket
	Status TicketStatus `json:"status,omitempty"`

	// Priority of the ticket
	Priority TicketPriority `json:"priority,omitempty"`

	// The channel through which the ticket was created
	Source TicketSource `json:"source,omitempty"`

	// Set to true if the ticket has been deleted/trashed. Deleted tickets will not be displayed in any views except the "deleted" filter
	Deleted bool `json:"deleted,omitempty"`

	// Set to true if the ticket has been marked as spam
	Spam bool `json:"spam,omitempty"`

	// Timestamp that denotes when the ticket is due to be resolved
	DueBy *Time `json:"due_by,omitempty"`

	// Timestamp that denotes when the first response is due
	FrDueBy *Time `json:"fr_due_by,omitempty"`

	// Set to true if the ticket has been escalated for any reason
	IsEscalated bool `json:"is_escalated,omitempty"`

	// Set to true if the ticket has been escalated as the result of first response time being breached
	FrEscalated bool `json:"fr_escalated,omitempty"`

	// Email addresses to which the ticket was originally sent
	ToEmails []string `json:"to_emails,omitempty"`

	// Email address added in the 'cc' field of the incoming ticket email
	CcEmails []string `json:"cc_emails,omitempty"`

	// Email address(e)s added while forwarding a ticket
	FwdEmails []string `json:"fwd_emails,omitempty"`

	// Email address added while replying to a ticket
	ReplyCcEmails []string `json:"reply_cc_emails,omitempty"`

	// ID of email config which is used for this ticket. (i.e., support@yourcompany.com/sales@yourcompany.com)
	EmailConfigID int64 `json:"email_config_id,omitempty"`

	// ID of the group to which the ticket has been assigned
	GroupID int64 `json:"group_id,omitempty"`

	// Subject of the ticket
	Subject string `json:"subject,omitempty"`

	// HTML content of the ticket
	Description string `json:"description,omitempty"`

	// Content of the ticket in plain text
	DescriptionText string `json:"description_text,omitempty"`

	// Tags that have been associated with the ticket
	Tags *[]string `json:"tags,omitempty"`

	// Key value pairs containing the names and values of custom fields.
	CustomFields map[string]any `json:"custom_fields,omitempty"`

	// Ticket Category.
	Category string `json:"category,omitempty"`

	// Ticket sub category.
	SubCategory string `json:"sub_category,omitempty"`

	// Ticket item category.
	ItemCategory string `json:"item_category,omitempty"`

	// Ticket urgency.
	Urgency TicketUrgency `json:"urgency,omitempty"`

	// Ticket impact.
	Impact TicketImpact `json:"impact,omitempty"`

	// Content of the ticket resolution note in plain text
	ResolutionNotes string `json:"resolution_notes,omitempty"`

	// HTML content of the ticket resolution note
	ResolutionNotesHTML string `json:"resolution_notes_html,omitempty"`

	// Ticket creation timestamp
	CreatedAt *Time `json:"created_at,omitempty"`

	// Ticket updated timestamp
	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (t *TicketCreate) AddAttachment(path string, data ...[]byte) {
	a := NewAttachment(path, data...)
	t.Attachments = append(t.Attachments, a)
}

func (t *TicketCreate) Files() Files {
	return ((Attachments)(t.Attachments)).Files()
}

func (t *TicketCreate) Values() Values {
	vs := Values{}

	vs.SetString("name", t.Name)
	vs.SetInt64("requester_id", t.RequesterID)
	vs.SetString("email", t.Email)
	vs.SetString("phone", t.Phone)
	vs.SetInt64("email_config_id", t.EmailConfigID)
	vs.SetString("subject", t.Subject)
	vs.SetString("type", t.Type)
	vs.SetInt("status", (int)(t.Status))
	vs.SetInt("priority", (int)(t.Priority))
	vs.SetString("description", t.Description)
	vs.SetInt64("responder_id", t.ResponderID)
	vs.SetStrings("cc_emails", t.CcEmails)
	vs.SetTimePtr("due_by", t.DueBy)
	vs.SetTimePtr("fr_due_by", t.FrDueBy)
	vs.SetInt64("group_id", t.GroupID)
	vs.SetInt64("workspace_id", t.WorkspaceID)
	vs.SetInt64("department_id", t.DepartmentID)
	vs.SetInt("source", (int)(t.Source))
	vs.SetStringsPtr("tags", t.Tags)
	vs.SetString("category", t.Category)
	vs.SetString("sub_category", t.SubCategory)
	vs.SetString("item_category", t.ItemCategory)
	vs.SetMap("custom_fields", t.CustomFields)
	vs.SetTimePtr("created_at", t.CreatedAt)
	vs.SetTimePtr("updated_at", t.UpdatedAt)

	return vs
}

func (t *TicketCreate) String() string {
	return toString(t)
}

type TicketUpdate = TicketCreate
