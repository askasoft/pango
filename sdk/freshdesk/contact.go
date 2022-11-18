package freshdesk

type Contact struct {
	ID int64 `json:"id,omitempty"`

	// Set to true if the contact has been verified
	Active bool `json:"active,omitempty"`

	// Address of the contact
	Address string `json:"address,omitempty"`

	// Avatar of the contact
	Avatar *Avatar `json:"avatar,omitempty"`

	// ID of the primary company to which this contact belongs
	CompanyID int64 `json:"company_id,omitempty"`

	// Set to true if the contact can see all tickets that are associated with the company to which he belong
	ViewAllTickets bool `json:"view_all_tickets,omitempty"`

	// Key value pair containing the name and value of the custom fields. Read more here
	CustomFields map[string]any `json:"custom_fields,omitempty"`

	// Set to true if the contact has been deleted. Note that this attribute will only be present for deleted contacts
	Deleted bool `json:"deleted,omitempty"`

	// A short description of the contact
	Description string `json:"description,omitempty"`

	// Primary email address of the contact. If you want to associate additional email(s) with this contact, use the other_emails attribute
	Email string `json:"email,omitempty"`

	// Job title of the contact
	JobTitle string `json:"job_title,omitempty"`

	// Language of the contact
	Language string `json:"language,omitempty"`

	// Mobile number of the contact
	Mobile string `json:"mobile,omitempty"`

	// Name of the contact
	Name string `json:"name,omitempty"`

	// Additional emails associated with the contact
	OtherEmails []string `json:"other_emails,omitempty"`

	// Telephone number of the contact
	Phone string `json:"phone,omitempty"`

	// Tags associated with this contact
	Tags []string `json:"tags,omitempty"`

	// Time zone in which the contact resides
	TimeZone string `json:"time_zone,omitempty"`

	// Twitter handle of the contact
	TwitterID string `json:"twitter_id,omitempty"`

	// Additional companies associated with the contact
	OtherCompanies []string `json:"other_companies,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (c *Contact) String() string {
	return toString(c)
}