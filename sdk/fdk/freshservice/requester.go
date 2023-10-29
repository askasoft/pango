package freshservice

type Requester struct {
	ID int64 `json:"id,omitempty"`

	// First name of the requester
	FirstName string `json:"first_name,omitempty"`

	// Last name of the requester.
	LastName string `json:"last_name,omitempty"`

	// Job title of the requester.
	JobTitle string `json:"job_title,omitempty"`

	// Primary email address of the requester.
	PrimaryEmail string `json:"primary_email,omitempty"`

	// Additional/secondary emails associated with the requester.
	SecondaryEmails []string `json:"secondary_emails,omitempty"`

	// Work phone number of the requester.
	WorkPhoneNumber string `json:"work_phone_number,omitempty"`

	// Mobile phone number of the requester.
	MobilePhoneNumber string `json:"mobile_phone_number,omitempty"`

	// Unique IDs of the departments associated with the requester
	DepartmentIDs []int64 `json:"department_ids,omitempty"`

	// Set to true if the requester must be allowed to view tickets filed by other members of the department, and false otherwise
	CanSeeAllTicketsFromAssociatedDepartments bool `json:"can_see_all_tickets_from_associated_departments,omitempty"`

	// User ID of the requester’s reporting manager.
	ReportingManagerID int64 `json:"reporting_manager_id,omitempty"`

	// Address of the requester.
	Address string `json:"address,omitempty"`

	// Time zone of the requester.
	TimeZone string `json:"time_zone,omitempty"`

	// Time format for the requester.Possible values: 12h (12 hour format), 24h (24 hour format)
	TimeFormat string `json:"time_format,omitempty"`

	// Language used by the requester.
	Language string `json:"language,omitempty"`

	// Unique ID of the location associated with the requester.
	LocationID int64 `json:"location_id,omitempty"`

	// Background information of the requester.
	BackgroundInformation string `json:"background_information,omitempty"`

	// Key-value pair containing the names and values of the (custom) requester fields.
	CustomFields map[string]any `json:"custom_fields,omitempty"`

	// Set to true if the user is active, and false if the user account has been deactivated.
	Active bool `json:"active,omitempty"`

	// Set to true if the user has logged in to Freshservice at least once, and false otherwise.
	HasLoggedIn bool `json:"has_logged_in,omitempty"`

	// Set to true if the user is an agent, false if the user is a requester
	IsAgent bool `json:"is_agent,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (r *Requester) String() string {
	return toString(r)
}

type requesterResult struct {
	Requester *Requester `json:"requester,omitempty"`
}

type requestersResult struct {
	Requesters []*Requester `json:"requesters,omitempty"`
}
