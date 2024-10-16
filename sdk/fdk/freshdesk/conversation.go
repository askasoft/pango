package freshdesk

import (
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

type ConversationSource int

const (
	ConversationSourceReply     ConversationSource = 0
	ConversationSourceNote      ConversationSource = 2
	ConversationSourceTweets    ConversationSource = 5
	ConversationSourceSurvey    ConversationSource = 6
	ConversationSourceFacebook  ConversationSource = 7
	ConversationSourceForwarded ConversationSource = 8
	ConversationSourcePhone     ConversationSource = 9
	ConversationSourceECommerce ConversationSource = 11
)

func (cs ConversationSource) String() string {
	switch cs {
	case ConversationSourceReply:
		return "Reply"
	case ConversationSourceNote:
		return "Note"
	case ConversationSourceTweets:
		return "Tweets"
	case ConversationSourceSurvey:
		return "Survey"
	case ConversationSourceFacebook:
		return "Facebook"
	case ConversationSourceForwarded:
		return "Forwarded"
	case ConversationSourcePhone:
		return "Phone"
	case ConversationSourceECommerce:
		return "ECommerce"
	default:
		return num.Itoa(int(cs))
	}
}

func ParseConversationSource(s string) ConversationSource {
	switch str.ToLower(s) {
	case "reply":
		return ConversationSourceReply
	case "note":
		return ConversationSourceNote
	case "tweets":
		return ConversationSourceTweets
	case "survey":
		return ConversationSourceSurvey
	case "facebook":
		return ConversationSourceFacebook
	case "forwarded":
		return ConversationSourceForwarded
	case "phone":
		return ConversationSourcePhone
	case "ecommerce":
		return ConversationSourceECommerce
	default:
		return -1
	}
}

type Conversation struct {
	// ID of the conversation
	ID int64 `json:"id,omitempty"`

	// Attachments (Updatable) associated with the conversation. The total size of all of a ticket's attachments cannot exceed 20MB.
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

	// Timestamp when the conversation last edited
	LastEditedAt *Time `json:"last_edited_at,omitempty"`

	// ID of the agent who has last edited the conversation
	LastEditedUserID int64 `json:"last_edited_user_id,omitempty"`

	// The email address from which the reply is sent.
	FromEmail string `json:"from_email,omitempty"`

	// Email addresses of agents/users who need to be notified about this conversation
	ToEmails []string `json:"to_emails,omitempty"`

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
