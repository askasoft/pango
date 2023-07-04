package freshservice

type Reply struct {
	// Content of the conversation in HTML
	Body string `json:"body,omitempty"`

	// Attachments associated with the conversation. The total size of all of a ticket's attachments cannot exceed 20MB.
	Attachments []*Attachment `json:"attachments,omitempty"`

	// The email address from which the reply is sent. By default the global support email will be used.
	FromEmail string `json:"from_email,omitempty"`

	// ID of the agent/user who is adding the conversation
	UserID int64 `json:"user_id,omitempty"`

	// Email address added in the 'cc' field of the outgoing ticket email.
	CcEmails []string `json:"cc_emails,omitempty"`

	// Email address added in the 'bcc' field of the outgoing ticket email.
	BccEmails []string `json:"bcc_emails,omitempty"`
}

func (r *Reply) AddAttachment(path string, data ...[]byte) {
	a := NewAttachment(path, data...)
	r.Attachments = append(r.Attachments, a)
}

func (r *Reply) Files() Files {
	return ((Attachments)(r.Attachments)).Files()
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
