package freshdesk

type Forward struct {
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

	// ------------------------------------------------------
	// response

	ID                int64         `json:"id,omitempty"`
	Attachments       []*Attachment `json:"attachments,omitempty"`
	Incoming          bool          `json:"incoming,omitempty"`
	Private           bool          `json:"private,omitempty"`
	SupportEmail      string        `json:"support_email,omitempty"`
	Source            int           `json:"source,omitempty"`
	Category          int           `json:"category,omitempty"`
	FromEmail         string        `json:"from_email,omitempty"`
	CcEmails          []string      `json:"cc_emails,omitempty"`
	BccEmails         []string      `json:"bcc_emails,omitempty"`
	EmailFailureCount int           `json:"email_failure_count,omitempty"`
	OutgoingFailures  string        `json:"outgoing_failures,omitempty"`
	Beleted           bool          `json:"deleted,omitempty"`
	LastEditedAt      *Time         `json:"last_edited_at,omitempty"`
	LastEditedUserID  int64         `json:"last_edited_user_id,omitempty"`
	CloudFiles        []string      `json:"cloud_files,omitempty"`
	HasQuotedText     bool          `json:"has_quoted_text,omitempty"`
	CreatedAt         *Time         `json:"created_at,omitempty"`
	UpdatedAt         *Time         `json:"updated_at,omitempty"`
}

func (f *Forward) String() string {
	return toString(f)
}
