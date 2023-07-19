package freshservice

import (
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

type ConversationSource int

const (
	ConversationSourceEmail          ConversationSource = 0
	ConversationSourceForm           ConversationSource = 1
	ConversationSourceNote           ConversationSource = 2
	ConversationSourceStatus         ConversationSource = 3
	ConversationSourceMeta           ConversationSource = 4
	ConversationSourceFeedback       ConversationSource = 5
	ConversationSourceForwardedEmail ConversationSource = 6
)

func (cs ConversationSource) String() string {
	switch cs {
	case ConversationSourceEmail:
		return "Email"
	case ConversationSourceForm:
		return "Form"
	case ConversationSourceNote:
		return "Note"
	case ConversationSourceStatus:
		return "Status"
	case ConversationSourceMeta:
		return "Meta"
	case ConversationSourceFeedback:
		return "Feedback"
	case ConversationSourceForwardedEmail:
		return "ForwardedEmail"
	default:
		return num.Itoa(int(cs))
	}
}

func ParseConversationSource(s string) ConversationSource {
	switch str.ToLower(s) {
	case "email":
		return ConversationSourceEmail
	case "form":
		return ConversationSourceForm
	case "note":
		return ConversationSourceNote
	case "status":
		return ConversationSourceStatus
	case "meta":
		return ConversationSourceMeta
	case "feedback":
		return ConversationSourceFeedback
	case "forwardedemail":
		return ConversationSourceForwardedEmail
	default:
		return -1
	}
}

type Conversation struct {
	// ID of the conversation
	ID int64 `json:"id,omitempty"`

	// Attachments (Updatable) associated with the conversation. The total size of all of a ticket's attachments cannot exceed 40MB.
	Attachments []*Attachment `json:"attachments,omitempty"`

	// Body (Updatable) Content of the conversation in HTML
	Body string `json:"body,omitempty"`

	// Content of the conversation in plain text
	BodyText string `json:"body_text,omitempty"`

	// Set to true if a particular conversation should appear as being created from outside (i.e., not through web portal)
	Incoming bool `json:"incoming,omitempty"`

	// Set to true if the note is private
	Private bool `json:"private,omitempty"`

	// Denotes the type of the conversation.
	Source ConversationSource `json:"source,omitempty"`

	// Email address from which the reply is sent. For notes, this value will be null.
	SupportEmail string `json:"support_email,omitempty"`

	// ID of the ticket to which this conversation is being added
	TicketID int64 `json:"ticket_id,omitempty"`

	// ID of the agent/user who is adding the conversation
	UserID int64 `json:"user_id,omitempty"`

	// Email addresses of agents/users who need to be notified about this conversation
	ToEmails []string `json:"to_emails,omitempty"`

	// The email address from which the reply is sent. By default the global support email will be used.
	FromEmail string `json:"from_email,omitempty"`

	// Email address added in the 'cc' field of the outgoing ticket email.
	CcEmails []string `json:"cc_emails,omitempty"`

	// Email address added in the 'bcc' field of the outgoing ticket email.
	BccEmails []string `json:"bcc_emails,omitempty"`

	// Conversation creation timestamp
	CreatedAt *Time `json:"created_at,omitempty"`

	// Conversation updated timestamp
	UpdatedAt *Time `json:"updated_at,omitempty"`

	// ------------------------------------------------------
	// response

	NotifiedTos []string `json:"notified_to,omitempty"`
}

func (c *Conversation) AddAttachment(path string, data ...[]byte) {
	a := NewAttachment(path, data...)
	c.Attachments = append(c.Attachments, a)
}

func (c *Conversation) Files() Files {
	return ((Attachments)(c.Attachments)).Files()
}

func (c *Conversation) Values() Values {
	vs := Values{}
	vs.SetString("body", c.Body)
	return vs
}

func (c *Conversation) String() string {
	return toString(c)
}

type conversationResult struct {
	Conversation  *Conversation   `json:"conversation,omitempty"`
	Conversations []*Conversation `json:"conversations,omitempty"`
}
