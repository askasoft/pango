package freshdesk

import (
	"context"
	"net/url"
)

// ---------------------------------------------------
// Contact

type ContactState string

const (
	ContactStateBlocked    ContactState = "blocked"
	ContactStateDeleted    ContactState = "deleted"
	ContactStateUnverified ContactState = "unverified"
	ContactStateVerified   ContactState = "verified"
)

type ListContactsOption struct {
	Email            string
	Mobile           string
	Phone            string
	UniqueExternalID string
	CompanyID        int64
	UpdatedSince     Time
	State            ContactState // [blocked/deleted/unverified/verified]
	Page             int
	PerPage          int
}

func (lco *ListContactsOption) IsNil() bool {
	return lco == nil
}

func (lco *ListContactsOption) Values() Values {
	q := Values{}
	q.SetString("email", lco.Email)
	q.SetString("mobile", lco.Mobile)
	q.SetString("phone", lco.Phone)
	q.SetString("unique_external_id", lco.UniqueExternalID)
	q.SetInt64("company_id", lco.CompanyID)
	q.SetString("state", (string)(lco.State))
	q.SetTime("updated_since", lco.UpdatedSince)
	q.SetInt("page", lco.Page)
	q.SetInt("per_page", lco.PerPage)
	return q
}

func (fd *Freshdesk) CreateContact(ctx context.Context, contact *ContactCreate) (*Contact, error) {
	url := fd.Endpoint("/contacts")
	result := &Contact{}
	if err := fd.DoPost(ctx, url, contact, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateContact(ctx context.Context, cid int64, contact *ContactUpdate) (*Contact, error) {
	url := fd.Endpoint("/contacts/%d", cid)
	result := &Contact{}
	if err := fd.DoPut(ctx, url, contact, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) GetContact(ctx context.Context, cid int64) (*Contact, error) {
	url := fd.Endpoint("/contacts/%d", cid)
	contact := &Contact{}
	err := fd.DoGet(ctx, url, contact)
	return contact, err
}

func (fd *Freshdesk) DeleteContact(ctx context.Context, cid int64) error {
	url := fd.Endpoint("/contacts/%d", cid)
	return fd.DoDelete(ctx, url)
}

func (fd *Freshdesk) HardDeleteContact(ctx context.Context, cid int64, force ...bool) error {
	url := fd.Endpoint("/contacts/%d/hard_delete", cid)
	if len(force) > 0 && force[0] {
		url += "?force=true"
	}
	return fd.DoDelete(ctx, url)
}

func (fd *Freshdesk) ListContacts(ctx context.Context, lco *ListContactsOption) ([]*Contact, bool, error) {
	url := fd.Endpoint("/contacts")
	contacts := []*Contact{}
	next, err := fd.DoList(ctx, url, lco, &contacts)
	return contacts, next, err
}

func (fd *Freshdesk) IterContacts(ctx context.Context, lco *ListContactsOption, icf func(*Contact) error) error {
	if lco == nil {
		lco = &ListContactsOption{}
	}
	if lco.Page < 1 {
		lco.Page = 1
	}
	if lco.PerPage < 1 {
		lco.PerPage = 100
	}

	for {
		contacts, next, err := fd.ListContacts(ctx, lco)
		if err != nil {
			return err
		}
		for _, c := range contacts {
			if err = icf(c); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lco.Page++
	}
	return nil
}

func (fd *Freshdesk) SearchContacts(ctx context.Context, keyword string) ([]*Contact, error) {
	url := fd.Endpoint("/contacts/autocomplete?term=%s", url.QueryEscape(keyword))
	contacts := []*Contact{}
	err := fd.DoGet(ctx, url, &contacts)
	return contacts, err
}

func (fd *Freshdesk) RestoreContact(ctx context.Context, cid int64) error {
	url := fd.Endpoint("/contacts/%d/restore", cid)
	return fd.DoPut(ctx, url, nil, nil)
}

func (fd *Freshdesk) InviteContact(ctx context.Context, cid int64) error {
	url := fd.Endpoint("/contacts/%d/send_invite", cid)
	return fd.DoPut(ctx, url, nil, nil)
}

func (fd *Freshdesk) MergeContacts(ctx context.Context, cm *ContactsMerge) error {
	url := fd.Endpoint("/contacts/merge")
	return fd.DoPost(ctx, url, nil, nil)
}

// ExportContacts return a job id, call GetExportedContactsURL() to get the job detail
func (fd *Freshdesk) ExportContacts(ctx context.Context, defaultFields, customFields []string) (string, error) {
	url := fd.Endpoint("/contacts/export")
	opt := &ExportOption{
		Fields: &ExportFields{
			DefaultFields: defaultFields,
			CustomFields:  customFields,
		},
	}
	job := &Job{}
	err := fd.DoPost(ctx, url, opt, job)
	return job.ID, err
}

// GetExportedContactsURL get the exported contacts url
func (fd *Freshdesk) GetExportedContactsURL(ctx context.Context, jid string) (*Job, error) {
	url := fd.Endpoint("/contacts/export/%s", jid)
	job := &Job{}
	err := fd.DoGet(ctx, url, job)
	return job, err
}

func (fd *Freshdesk) MakeAgent(ctx context.Context, cid int64, agent *Agent) (*Contact, error) {
	url := fd.Endpoint("/contacts/%d/make_agent", cid)
	result := &Contact{}
	if err := fd.DoPut(ctx, url, agent, result); err != nil {
		return nil, err
	}
	return result, nil
}
