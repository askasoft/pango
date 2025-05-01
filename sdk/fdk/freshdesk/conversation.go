package freshdesk

import (
	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

type ConversationSource int

const (
	ConversationSourceReply          ConversationSource = 0
	ConversationSourceNote           ConversationSource = 2
	ConversationSourceTweets         ConversationSource = 5
	ConversationSourceSurveyFeedback ConversationSource = 6
	ConversationSourceFacebookPost   ConversationSource = 7
	ConversationSourceForwardedEmail ConversationSource = 8
	ConversationSourcePhone          ConversationSource = 9
	ConversationSourceECommerce      ConversationSource = 11
)

func (cs ConversationSource) String() string {
	switch cs {
	case ConversationSourceReply:
		return "Reply"
	case ConversationSourceNote:
		return "Note"
	case ConversationSourceTweets:
		return "Tweets"
	case ConversationSourceSurveyFeedback:
		return "SurveyFeedback"
	case ConversationSourceFacebookPost:
		return "FacebookPost"
	case ConversationSourceForwardedEmail:
		return "ForwardedEmail"
	case ConversationSourcePhone:
		return "Phone"
	case ConversationSourceECommerce:
		return "E-Commerce"
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
	case "surveyfeekback":
		return ConversationSourceSurveyFeedback
	case "facebookpost":
		return ConversationSourceFacebookPost
	case "forwardedemail":
		return ConversationSourceForwardedEmail
	case "phone":
		return ConversationSourcePhone
	case "e-commerce":
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

	// Thread ID undocument
	ThreadID int64 `json:"thread_id,omitempty"`

	// Thread Message ID undocument
	ThreadMessageID int64 `json:"thread_message_id,omitempty"`

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

	NotifiedTos []string `json:"notified_to,omitempty"`

	// Conversation creation timestamp
	CreatedAt Time `json:"created_at,omitempty"`

	// Conversation updated timestamp
	UpdatedAt Time `json:"updated_at,omitempty"`
}

func (c *Conversation) String() string {
	return toString(c)
}
