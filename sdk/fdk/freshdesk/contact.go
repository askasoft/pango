package freshdesk

import (
	"net/url"

	"github.com/askasoft/pango/bol"
	"github.com/askasoft/pango/num"
)

type OtherCompany struct {
	// ID of the primary company to which this contact belongs
	CompanyID int64 `json:"company_id,omitempty"`

	// Set to true if the contact can see all tickets that are associated with the company to which he belong
	ViewAllTickets bool `json:"view_all_tickets,omitempty"`
}

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

	// Key value pair containing the name and value of the custom fields.
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

	// External ID of the contact
	UniqueExternalID string `json:"unique_external_id,omitempty"`

	// Additional companies associated with the contact
	OtherCompanies []any `json:"other_companies,omitempty"`

	// IDs of the companies associated with the contact (only used by MergeContact)
	CompanyIDs int64 `json:"company_ids,omitempty"`

	// Return by MakeAgent()
	Agent *Agent `json:"agent,omitempty"`

	CreatedAt Time `json:"created_at,omitempty"`

	UpdatedAt Time `json:"updated_at,omitempty"`
}

func (c *Contact) String() string {
	return toString(c)
}

type ContactCreate struct {
	// Name of the contact
	Name string `json:"name,omitempty"`

	// Primary email address of the contact. If you want to associate additional email(s) with this contact, use the other_emails attribute
	Email string `json:"email,omitempty"`

	// Telephone number of the contact
	Phone string `json:"phone,omitempty"`

	// Mobile number of the contact
	Mobile string `json:"mobile,omitempty"`

	// Twitter handle of the contact
	TwitterID string `json:"twitter_id,omitempty"`

	// External ID of the contact
	UniqueExternalID string `json:"unique_external_id,omitempty"`

	// Additional emails associated with the contact
	OtherEmails []string `json:"other_emails,omitempty"`

	// ID of the primary company to which this contact belongs
	CompanyID int64 `json:"company_id,omitempty"`

	// Set to true if the contact can see all tickets that are associated with the company to which he belong
	ViewAllTickets bool `json:"view_all_tickets,omitempty"`

	// Additional companies associated with the contact
	OtherCompanies []any `json:"other_companies,omitempty"`

	// Address of the contact
	Address string `json:"address,omitempty"`

	// Avatar of the contact
	Avatar *Avatar `json:"avatar,omitempty"`

	// Key value pair containing the name and value of the custom fields.
	CustomFields map[string]any `json:"custom_fields,omitempty"`

	// A short description of the contact
	Description string `json:"description,omitempty"`

	// Job title of the contact
	JobTitle string `json:"job_title,omitempty"`

	// Language of the contact
	Language string `json:"language,omitempty"`

	// Tags associated with this contact
	Tags *[]string `json:"tags,omitempty"`

	// Time zone in which the contact resides
	TimeZone string `json:"time_zone,omitempty"`

	// This attribute for contacts can only be set if the Custom Objects feature is enabled. The value can either be in the form of the display_id (record id) or primary_field_value (user defined record value). The default value is display_id.
	LookupParameter string `json:"lookup_parameter,omitempty"`
}

func (c *ContactCreate) Files() Files {
	if c.Avatar != nil {
		return Files{c.Avatar}
	}
	return nil
}

func (c *ContactCreate) Values() Values {
	vs := Values{}
	vs.SetString("name", c.Name)
	vs.SetString("email", c.Email)
	vs.SetString("phone", c.Phone)
	vs.SetString("mobile", c.Mobile)
	vs.SetString("twitter_id", c.TwitterID)
	vs.SetString("unique_external_id", c.UniqueExternalID)
	vs.SetStrings("other_emails", c.OtherEmails)
	vs.SetInt64("company_id", c.CompanyID)
	vs.SetBool("view_all_tickets", c.ViewAllTickets)
	vs.SetString("address", c.Address)
	vs.SetMap("custom_fields", c.CustomFields)
	vs.SetString("description", c.Description)
	vs.SetString("job_title", c.JobTitle)
	vs.SetString("language", c.Language)
	vs.SetStringsPtr("tags", c.Tags)
	vs.SetString("time_zone", c.TimeZone)
	vs.SetString("lookup_parameter", c.LookupParameter)
	if len(c.OtherCompanies) > 0 {
		for _, o := range c.OtherCompanies {
			if oc, ok := o.(*OtherCompany); ok {
				(url.Values)(vs).Add("other_companies[company_id]", num.Ltoa(oc.CompanyID))
				(url.Values)(vs).Add("other_companies[view_all_tickets]", bol.Btoa(oc.ViewAllTickets))
			}
		}
	}
	return vs
}

func (c *ContactCreate) String() string {
	return toString(c)
}

type ContactUpdate = ContactCreate

type ContactsMerge struct {
	// ID of the primary contact
	PrimaryContactID int64 `json:"primary_contact_id,omitempty"`

	// Array of numbers	IDs of contacts to be merged
	SecondaryContactIDs []int64 `json:"secondary_contact_ids,omitempty"`

	// Contains attributes that need to be updated in the primary contact during merge (optional)
	// email, phone, mobile, twitter_id, unique_external_id, other_emails, company_ids
	Contact *Contact `json:"contact,omitempty"`
}

func (cm *ContactsMerge) String() string {
	return toString(cm)
}
