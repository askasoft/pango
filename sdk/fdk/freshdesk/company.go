package freshdesk

type Company struct {
	ID int64 `json:"id,omitempty"`

	// Name of the company
	Name string `json:"name,omitempty"`

	// Description of the company
	Description string `json:"description,omitempty"`

	// Any specific note about the company
	Note string `json:"note,omitempty"`

	// Domains of the company. Email addresses of contacts that contain this domain will be associated with that company automatically.
	Domains []string `json:"domains,omitempty"`

	// The strength of your relationship with the company
	HealthScore string `json:"health_score,omitempty"`

	// Classification based on how much value the company brings to your business
	AccountTier string `json:"account_tier,omitempty"`

	// Date when your contract or relationship with the company is due for renewal
	RenewalDate *Time `json:"renewal_date,omitempty"`

	//The industry the company serves in
	Industry string `json:"industry,omitempty"`

	// Key value pair containing the name and value of the custom fields.
	CustomFields map[string]any `json:"custom_fields,omitempty"`

	CreatedAt *Time `json:"created_at,omitempty"`

	UpdatedAt *Time `json:"updated_at,omitempty"`
}

func (c *Company) String() string {
	return toString(c)
}

type companyResult struct {
	Companies []*Company `json:"companies,omitempty"`
}
