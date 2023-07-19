package freshdesk

import (
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

type TicketSource int
type TicketStatus int
type TicketPriority int
type TicketFilter string

const (
	TicketSourceEmail          TicketSource = 1
	TicketSourcePortal         TicketSource = 2
	TicketSourcePhone          TicketSource = 3
	TicketSourceChat           TicketSource = 7
	TicketSourceFeedbackWidget TicketSource = 9
	TicketSourceOutboundEmail  TicketSource = 10

	TicketStatusOpen     TicketStatus = 2
	TicketStatusPending  TicketStatus = 3
	TicketStatusResolved TicketStatus = 4
	TicketStatusClosed   TicketStatus = 5

	TicketPriorityLow    TicketPriority = 1
	TicketPriorityMedium TicketPriority = 2
	TicketPriorityHigh   TicketPriority = 3
	TicketPriorityUrgent TicketPriority = 4

	TicketFilterNewAndMyOpen TicketFilter = "new_and_my_open"
	TicketFilterWatching     TicketFilter = "watching"
	TicketFilterSpam         TicketFilter = "spam"
	TicketFilterDeleted      TicketFilter = "deleted"
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
	case TicketSourceOutboundEmail:
		return "OutboundEmail"
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
	case "outboundemail":
		return TicketSourceOutboundEmail
	default:
		return 0
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
		return 0
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
		return 0
	}
}

type Ticket struct {
	// Parent Unique ID of the ticket
	ParentID int64 `json:"parent_id,omitempty"`

	// Unique ID of the ticket
	ID int64 `json:"id,omitempty"`

	// Name of the requester
	Name string `json:"name,omitempty"`

	// Email address of the requester. If no contact exists with this email address in Freshdesk, it will be added as a new contact.
	Email string `json:"email,omitempty"`

	// Phone number of the requester. If no contact exists with this phone number in Freshdesk, it will be added as a new contact. If the phone number is set and the email address is not, then the name attribute is mandatory.
	Phone string `json:"phone,omitempty"`

	// Facebook ID of the requester. A contact should exist with this facebook_id in Freshdesk.
	FacebookID string `json:"facebook_id,omitempty"`

	// Twitter handle of the requester. If no contact exists with this handle in Freshdesk, it will be added as a new contact.
	TwitterID string `json:"twitter_id,omitempty"`

	// External ID of the requester. If no contact exists with this external ID in Freshdesk, they will be added as a new contact.
	UniqueExternalID string `json:"unique_external_id,omitempty"`

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

	// ID of the product to which the ticket is associated
	ProductID int64 `json:"product_id,omitempty"`

	// ID of the group to which the ticket has been assigned
	GroupID int64 `json:"group_id,omitempty"`

	// ID of the company to which this ticket belongs
	CompanyID int64 `json:"company_id,omitempty"`

	// Subject of the ticket
	Subject string `json:"subject,omitempty"`

	// HTML content of the ticket
	Description string `json:"description,omitempty"`

	// Content of the ticket in plain text
	DescriptionText string `json:"description_text,omitempty"`

	// Tags that have been associated with the ticket
	Tags []string `json:"tags,omitempty"`

	// Ticket attachments. The total size of these attachments cannot exceed 15MB.
	Attachments []*Attachment `json:"attachments,omitempty"`

	// Key value pairs containing the names and values of custom fields.
	CustomFields map[string]any `json:"custom_fields,omitempty"`

	// include=requester
	Requester *Contact `json:"requester,omitempty"`

	// include=conversations
	Conversations []*Conversation `json:"conversations,omitempty"`

	// ID of the internal agent which the ticket should be assigned with
	InternalAgentID int64 `json:"internal_agent_id,omitempty"`

	// ID of the internal group to which the ticket should be assigned with
	InternalGroupID int64 `json:"internal_group_id,omitempty"`

	// Ticket creation timestamp
	CreatedAt *Time `json:"created_at,omitempty"`

	// Ticket updated timestamp
	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (t *Ticket) AddAttachment(path string, data ...[]byte) {
	a := NewAttachment(path, data...)
	t.Attachments = append(t.Attachments, a)
}

func (t *Ticket) Files() Files {
	return ((Attachments)(t.Attachments)).Files()
}

func (t *Ticket) Values() Values {
	vs := Values{}

	vs.SetInt64("parent_id", t.ParentID)
	vs.SetString("name", t.Name)
	vs.SetInt64("requester_id", t.RequesterID)
	vs.SetString("email", t.Email)
	vs.SetString("facebook_id", t.FacebookID)
	vs.SetString("phone", t.Phone)
	vs.SetString("twitter_id", t.TwitterID)
	vs.SetString("unique_external_id", t.UniqueExternalID)
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
	vs.SetInt64("product_id", t.ProductID)
	vs.SetInt("source", (int)(t.Source))
	vs.SetStrings("tags", t.Tags)
	vs.SetInt64("company_id", t.CompanyID)
	vs.SetInt64("internal_agent_id", t.InternalAgentID)
	vs.SetInt64("internal_group_id", t.InternalGroupID)
	vs.SetMap("custom_fields", t.CustomFields)

	return vs
}

func (t *Ticket) String() string {
	return toString(t)
}

type TicketProperties struct {
	// Support email from which the reply should be sent
	FromEmail string `json:"from_email,omitempty"`

	// Used to skip email notifications sent to requesters on closing a ticket
	SkipCloseNotification bool `json:"skip_close_notification"`

	// ID of email config which is used for this ticket. (i.e., support@yourcompany.com/sales@yourcompany.com)
	EmailConfigID int64 `json:"email_config_id,omitempty"`

	// ID of the group to which the ticket has been assigned
	GroupID int64 `json:"group_id,omitempty"`

	// Priority of the ticket
	Priority TicketPriority `json:"priority,omitempty"`

	// Parent Unique ID of the ticket
	ParentID int64 `json:"parent_id,omitempty"`

	// User ID of the requester. For existing contacts, the requester_id can be passed instead of the requester's email.
	RequesterID int64 `json:"requester_id,omitempty"`

	// ID of the agent to whom the ticket has been assigned.
	ResponderID int64 `json:"responder_id,omitempty"`

	// Status of the ticket
	Status TicketStatus `json:"status,omitempty"`

	// The channel through which the ticket was created
	Source TicketSource `json:"source,omitempty"`

	// Helps categorize the ticket according to the different kinds of issues your support team deals with.
	Type string `json:"type,omitempty"`

	// ID of the product to which the ticket is associated
	ProductID int64 `json:"product_id,omitempty"`

	// Timestamp that denotes when the ticket is due to be resolved
	DueBy *Time `json:"due_by,omitempty"`

	// Timestamp that denotes when the first response is due
	FrDueBy *Time `json:"fr_due_by,omitempty"`

	// Key value pairs containing the names and values of custom fields.
	CustomFields map[string]any `json:"custom_fields,omitempty"`

	// Tags that have been associated with the ticket
	Tags []string `json:"tags,omitempty"`

	// ID of the internal agent which the ticket should be assigned with
	InternalAgentID int64 `json:"internal_agent_id,omitempty"`

	// ID of the internal group to which the ticket should be assigned with
	InternalGroupID int64 `json:"internal_group_id,omitempty"`
}

func (tp *TicketProperties) String() string {
	return toString(tp)
}

type TicketsMerge struct {
	// Ticket to which conversations from secondary tickets will be merged
	PrimaryID int64 `json:"primary_id,omitempty"`

	// Array of numbers	IDs of tickets to be merged
	TicketIDs []int64 `json:"ticket_ids,omitempty"`

	// This contains the note added to the primary ticket along with the type of the note (public/private). A default note gets added if this is not specified in the request
	NoteInPrimary *Note `json:"note_in_primary,omitempty"`

	// This contains the note added to the secondary tickets along with the type of the note (public/private). A default note gets added if this is not specified in the request
	NoteInSecondary *Note `json:"note_in_secondary,omitempty"`

	// This will add recipients from secondary tickets in CC of the primary ticket
	ConvertRecepientsToCc bool `json:"convert_recepients_to_cc,omitempty"` //nolint: misspell
}

func (tm *TicketsMerge) String() string {
	return toString(tm)
}

type TicketWatchers struct {
	WatcherIDs []int64 `json:"watcher_ids,omitempty"`
}

func (tw *TicketWatchers) String() string {
	return toString(tw)
}
