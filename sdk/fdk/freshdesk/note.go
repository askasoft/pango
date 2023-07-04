package freshdesk

type Note struct {
	// Attachments associated with the conversation. The total size of all of a ticket's attachments cannot exceed 20MB.
	Attachments []*Attachment `json:"attachments,omitempty"`

	// Content of the conversation in HTML
	Body string `json:"body,omitempty"`

	// Set to true if a particular conversation should appear as being created from outside (i.e., not through web portal)
	Incoming bool `json:"incoming"`

	// Email addresses of agents/users who need to be notified about this note
	NotifyEmails []string `json:"notify_emails,omitempty"`

	// Set to true if the note is private. The default value is true.
	Private bool `json:"private"`

	// ID of the agent/user who is adding the conversation
	UserID int64 `json:"user_id,omitempty"`

	// ------------------------------------------------------
	// response

	ID           int64    `json:"id,omitempty"`
	BodyText     string   `json:"body_text,omitempty"`
	SupportEmail string   `json:"support_email,omitempty"`
	TicketID     int64    `json:"ticket_id,omitempty"`
	NotifiedTos  []string `json:"notified_to,omitempty"`
	CreatedAt    *Time    `json:"created_at,omitempty"`
	UpdatedAt    *Time    `json:"updated_at,omitempty"`
}

func (n *Note) AddAttachment(path string, data ...[]byte) {
	a := NewAttachment(path, data...)
	n.Attachments = append(n.Attachments, a)
}

func (n *Note) Files() Files {
	return ((Attachments)(n.Attachments)).Files()
}

func (n *Note) Values() Values {
	vs := Values{}

	vs.SetString("body", n.Body)
	vs.SetBool("incoming", n.Incoming)
	vs.SetStrings("notify_emails", n.NotifyEmails)
	vs.SetBool("private", n.Private)
	vs.SetInt64("user_id", n.UserID)

	return vs
}

func (n *Note) String() string {
	return toString(n)
}
