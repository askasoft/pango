package freshdesk

type Note struct {
	// ID of the conversation
	ID int64 `json:"id,omitempty"`

	// Attachments associated with the conversation. The total size of all of a ticket's attachments cannot exceed 20MB.
	Attachments []*Attachment `json:"attachments,omitempty"`

	// Content of the conversation in HTML
	Body string `json:"body,omitempty"`

	// Content of the conversation in plain text
	BodyText string `json:"body_text,omitempty"`

	// Set to true if a particular conversation should appear as being created from outside (i.e., not through web portal)
	Incoming bool `json:"incoming"`

	// Set to true if the note is private. The default value is true.
	Private bool `json:"private"`

	// The email address from which the reply is sent. By default the global support email will be used.
	FromEmail string `json:"from_email,omitempty"`

	// ID of the agent/user who is adding the conversation
	UserID int64 `json:"user_id,omitempty"`

	// Email addresses of agents/users who need to be notified about this note
	NotifyEmails []string `json:"notify_emails,omitempty"`

	// ID of the ticket to which this conversation is being added
	TicketID int64 `json:"ticket_id,omitempty"`

	// Email address from which the reply is sent. For notes, this value will be null.
	SupportEmail string `json:"support_email,omitempty"`

	// Email address
	NotifiedTos []string `json:"notified_to,omitempty"`

	// Conversation creation timestamp
	CreatedAt *Time `json:"created_at,omitempty"`

	// Conversation updated timestamp
	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (n *Note) AddAttachment(path string, data ...[]byte) {
	a := NewAttachment(path, data...)
	n.Attachments = append(n.Attachments, a)
}

func (n *Note) GetAttachments() []*Attachment {
	return n.Attachments
}

func (n *Note) Values() Values {
	vs := Values{}

	vs.SetBool("incoming", n.Incoming)
	vs.SetBool("private", n.Private)
	vs.SetStrings("notify_emails", n.NotifyEmails)
	vs.SetInt64("user_id", n.UserID)
	vs.SetString("body", n.Body)

	return vs
}

func (n *Note) String() string {
	return toString(n)
}
