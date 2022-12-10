package freshdesk

type Conversation struct {
	// ID of the conversation
	ID int64 `json:"id,omitempty"`

	// Attachments associated with the conversation. The total size of all of a ticket's attachments cannot exceed 20MB.
	Attachments []*Attachment `json:"attachments,omitempty"`

	// Content of the conversation in HTML
	Body string `json:"body,omitempty"`

	// Content of the conversation in plain text
	BodyText string `json:"body_text,omitempty"`

	// Set to true if a particular conversation should appear as being created from outside (i.e., not through web portal)
	Incoming bool `json:"incoming,omitempty"`

	// Email addresses of agents/users who need to be notified about this conversation
	ToEmails []string `json:"to_emails,omitempty"`

	// Set to true if the note is private
	Private bool `json:"private,omitempty"`

	// Denotes the type of the conversation.
	Source int64 `json:"source,omitempty"`

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

	// Email address
	NotifiedTos []string `json:"notified_to,omitempty"`

	// Conversation creation timestamp
	CreatedAt *Time `json:"created_at,omitempty"`

	// Conversation updated timestamp
	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (c *Conversation) AddAttachment(path string, data ...[]byte) {
	a := NewAttachment(path, data...)
	c.Attachments = append(c.Attachments, a)
}

func (c *Conversation) Files() Files {
	fs := make(Files, len(c.Attachments))
	for i, a := range c.Attachments {
		fs[i] = a
	}
	return fs
}

func (c *Conversation) Values() Values {
	vs := Values{}
	vs.SetString("body", c.Body)
	return vs
}

func (c *Conversation) String() string {
	return toString(c)
}
