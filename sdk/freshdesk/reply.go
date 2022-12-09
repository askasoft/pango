package freshdesk

type Reply struct {
	// ID of the reply
	ID int64 `json:"id,omitempty"`

	// Attachments associated with the conversation. The total size of all of a ticket's attachments cannot exceed 20MB.
	Attachments []*Attachment `json:"attachments,omitempty"`

	// Content of the conversation in HTML
	Body string `json:"body,omitempty"`

	// Content of the conversation in plain text
	BodyText string `json:"body_text,omitempty"`

	// The email address from which the reply is sent. By default the global support email will be used.
	FromEmail string `json:"from_email,omitempty"`

	// ID of the agent/user who is adding the conversation
	UserID int64 `json:"user_id,omitempty"`

	// Email address added in the 'cc' field of the outgoing ticket email.
	CcEmails []string `json:"cc_emails,omitempty"`

	// Email address added in the 'bcc' field of the outgoing ticket email.
	BccEmails []string `json:"bcc_emails,omitempty"`

	// ID of the ticket to which this conversation is being added
	TicketID int64 `json:"ticket_id,omitempty"`

	// Email address
	RepliedTos []string `json:"replied_to,omitempty"`

	// Conversation creation timestamp
	CreatedAt *Time `json:"created_at,omitempty"`

	// Conversation updated timestamp
	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (r *Reply) AddAttachment(path string, data ...[]byte) {
	a := NewAttachment(path, data...)
	r.Attachments = append(r.Attachments, a)
}

func (r *Reply) GetAttachments() []*Attachment {
	return r.Attachments
}

func (r *Reply) Values() Values {
	vs := Values{}

	vs.SetString("from_email", r.FromEmail)
	vs.SetInt64("user_id", r.UserID)
	vs.SetStrings("cc_emails", r.CcEmails)
	vs.SetStrings("bcc_emails", r.BccEmails)
	vs.SetString("body", r.Body)

	return vs
}

func (r *Reply) String() string {
	return toString(r)
}
