package freshdesk

type TicketForward struct {
	// Content of the note in HTML format
	Body string `json:"body,omitempty"`

	// Content of the note in plain text
	BodyText string `json:"body_text,omitempty"`

	// ID of the forwarded note
	ID int64 `json:"id,omitempty"`

	// Set to true if a particular conversation should appear as being created from outside (i.e., not through web portal)
	Incoming bool `json:"incoming,omitempty"`

	// Set to true if the note is private
	Private bool `json:"private,omitempty"`

	// ID of the agent/user who is forwarding the ticket
	UserID int64 `json:"user_id,omitempty"`

	// Email address from which the reply is sent. For notes, this value will be null.
	SupportEmail string `json:"support_email,omitempty"`

	// Denotes the type of the conversation.
	Source int `json:"source,omitempty"`

	// ID of the ticket to which this conversation is being added
	TicketID int64 `json:"ticket_id,omitempty"`

	// Include the quoted text conversations in the forwarded email. The default value is True.
	IncludeQuotedText bool `json:"include_quoted_text"`

	// Include the ticket attachments in the forwarded email. The default value is True.
	IncludeOriginalAttachments bool `json:"include_original_attachments"`

	// The email address from which the forward is sent. By default the global support email will be used.
	FromEmail string `json:"from_email,omitempty"`

	// Emails to which the ticket gets forwarded
	ToEmails []string `json:"to_emails,omitempty"`

	// Email address added in the 'cc' field of the outgoing forward email.
	CcEmails []string `json:"cc_emails,omitempty"`

	// Email address added in the 'bcc' field of the outgoing forward email.
	BccEmails []string `json:"bcc_emails,omitempty"`
}

func (tf *TicketForward) String() string {
	return toString(tf)
}

type ReplyForward struct {
	// Content of the conversation in HTML
	Body string `json:"body,omitempty"`

	// Content of the conversation in plain text
	BodyText string `json:"body_text,omitempty"`

	// ID of the agent/user who is adding the conversation
	UserID int64 `json:"user_id,omitempty"`

	// ID of the ticket to which this conversation is being added
	TicketID int64 `json:"ticket_id,omitempty"`

	// Emails to which the reply is addressed
	ToEmails []string `json:"to_emails,omitempty"`
}

func (rf *ReplyForward) String() string {
	return toString(rf)
}

type ForwardResult struct {
	ID                int64         `json:"id,omitempty"`
	Body              string        `json:"body,omitempty"`
	BodyText          string        `json:"body_text,omitempty"`
	UserID            int64         `json:"user_id,omitempty"`
	Incoming          bool          `json:"incoming,omitempty"`
	Private           bool          `json:"private,omitempty"`
	SupportEmail      string        `json:"support_email,omitempty"`
	Source            int           `json:"source,omitempty"`
	Category          int           `json:"category,omitempty"`
	TicketID          int64         `json:"ticket_id,omitempty"`
	FromEmail         string        `json:"from_email,omitempty"`
	ToEmails          []string      `json:"to_emails,omitempty"`
	CcEmails          []string      `json:"cc_emails,omitempty"`
	BccEmails         []string      `json:"bcc_emails,omitempty"`
	EmailFailureCount int           `json:"email_failure_count,omitempty"`
	OutgoingFailures  string        `json:"outgoing_failures,omitempty"`
	Deleted           bool          `json:"deleted,omitempty"`
	LastEditedAt      *Time         `json:"last_edited_at,omitempty"`
	LastEditedUserID  int64         `json:"last_edited_user_id,omitempty"`
	Attachments       []*Attachment `json:"attachments,omitempty"`
	CloudFiles        []string      `json:"cloud_files,omitempty"`
	HasQuotedText     bool          `json:"has_quoted_text,omitempty"`
	CreatedAt         Time          `json:"created_at,omitempty"`
	UpdatedAt         Time          `json:"updated_at,omitempty"`
}

func (fr *ForwardResult) String() string {
	return toString(fr)
}
