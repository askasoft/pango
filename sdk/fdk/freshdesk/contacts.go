package freshdesk

import (
	"net/url"
)

// ---------------------------------------------------
// Contact

type ListContactsOption struct {
	Email            string
	Mobile           string
	Phone            string
	UniqueExternalID string
	CompanyID        int64
	UpdatedSince     Time
	State            string // [blocked/deleted/unverified/verified]
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
	q.SetString("state", lco.State)
	q.SetTime("updated_since", lco.UpdatedSince)
	q.SetInt("page", lco.Page)
	q.SetInt("per_page", lco.PerPage)
	return q
}

func (fd *Freshdesk) CreateContact(contact *Contact) (*Contact, error) {
	url := fd.endpoint("/contacts")
	result := &Contact{}
	err := fd.doPost(url, contact, result)
	return result, err
}

func (fd *Freshdesk) UpdateContact(cid int64, contact *Contact) (*Contact, error) {
	url := fd.endpoint("/contacts/%d", cid)
	result := &Contact{}
	err := fd.doPut(url, contact, result)
	return result, err
}

func (fd *Freshdesk) GetContact(cid int64) (*Contact, error) {
	url := fd.endpoint("/contacts/%d", cid)
	contact := &Contact{}
	err := fd.doGet(url, contact)
	return contact, err
}

func (fd *Freshdesk) DeleteContact(cid int64) error {
	url := fd.endpoint("/contacts/%d", cid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) HardDeleteContact(cid int64, force ...bool) error {
	url := fd.endpoint("/contacts/%d/hard_delete", cid)
	if len(force) > 0 && force[0] {
		url += "?force=true"
	}
	return fd.doDelete(url)
}

func (fd *Freshdesk) ListContacts(lco *ListContactsOption) ([]*Contact, bool, error) {
	url := fd.endpoint("/contacts")
	contacts := []*Contact{}
	next, err := fd.doList(url, lco, &contacts)
	return contacts, next, err
}

func (fd *Freshdesk) IterContacts(lco *ListContactsOption, itf func(*Contact) error) error {
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
		contacts, next, err := fd.ListContacts(lco)
		if err != nil {
			return err
		}
		for _, c := range contacts {
			if err = itf(c); err != nil {
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

func (fd *Freshdesk) SearchContacts(keyword string) ([]*Contact, error) {
	url := fd.endpoint("/contacts/autocomplete?term=%s", url.QueryEscape(keyword))
	contacts := []*Contact{}
	err := fd.doGet(url, &contacts)
	return contacts, err
}

func (fd *Freshdesk) RestoreContact(cid int64) error {
	url := fd.endpoint("/contacts/%d/restore", cid)
	return fd.doPut(url, nil, nil)
}

func (fd *Freshdesk) InviteContact(cid int64) error {
	url := fd.endpoint("/contacts/%d/send_invite", cid)
	return fd.doPut(url, nil, nil)
}

func (fd *Freshdesk) MergeContacts(cm *ContactsMerge) error {
	url := fd.endpoint("/contacts/merge")
	return fd.doPost(url, nil, nil)
}

// ExportContacts return a job id, call GetExportedContactsURL() to get the job detail
func (fd *Freshdesk) ExportContacts(defaultFields, customFields []string) (string, error) {
	url := fd.endpoint("/contacts/export")
	opt := &ExportOption{
		Fields: &ExportFields{
			DefaultFields: defaultFields,
			CustomFields:  customFields,
		},
	}
	job := &Job{}
	err := fd.doPost(url, opt, job)
	return job.ID, err
}

// GetExportedContactsURL get the exported contacts url
func (fd *Freshdesk) GetExportedContactsURL(jid string) (*Job, error) {
	url := fd.endpoint("/contacts/export/%s", jid)
	job := &Job{}
	err := fd.doGet(url, job)
	return job, err
}

func (fd *Freshdesk) MakeAgent(cid int64, agent *Agent) (*Contact, error) {
	url := fd.endpoint("/contacts/%d/make_agent", cid)
	result := &Contact{}
	err := fd.doPut(url, agent, result)
	return result, err
}
