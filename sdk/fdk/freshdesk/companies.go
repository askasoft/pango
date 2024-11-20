package freshdesk

import (
	"context"
	"net/url"
)

// ---------------------------------------------------
// Company

type ListCompaniesOption = PageOption

func (fd *Freshdesk) CreateCompany(ctx context.Context, company *Company) (*Company, error) {
	url := fd.endpoint("/companies")
	result := &Company{}
	if err := fd.doPost(ctx, url, company, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) GetCompany(ctx context.Context, cid int64) (*Company, error) {
	url := fd.endpoint("/companies/%d", cid)
	result := &Company{}
	err := fd.doGet(ctx, url, result)
	return result, err
}

func (fd *Freshdesk) ListCompanies(ctx context.Context, lco *ListCompaniesOption) ([]*Company, bool, error) {
	url := fd.endpoint("/companies")
	result := []*Company{}
	next, err := fd.doList(ctx, url, lco, &result)
	return result, next, err
}

func (fd *Freshdesk) IterCompanies(ctx context.Context, lco *ListCompaniesOption, icf func(*Company) error) error {
	if lco == nil {
		lco = &ListCompaniesOption{}
	}
	if lco.Page < 1 {
		lco.Page = 1
	}
	if lco.PerPage < 1 {
		lco.PerPage = 100
	}

	for {
		companies, next, err := fd.ListCompanies(ctx, lco)
		if err != nil {
			return err
		}
		for _, c := range companies {
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

// Search Companies
// Search for a company using its name.
// Note:
// 1. The search is case insensitive.
// 2. You cannot search with a substring. For example, a company "Acme Corporation" can be looked up using "acme", "Ac", "Corporation" and "Co". But it will not be returned when you search for "cme" or "orporation".
func (fd *Freshdesk) SearchCompanies(ctx context.Context, name string) ([]*Company, error) {
	url := fd.endpoint("/companies/autocomplete?name=%s", url.QueryEscape(name))
	result := &companyResult{}
	err := fd.doGet(ctx, url, result)
	return result.Companies, err
}

func (fd *Freshdesk) UpdateCompany(ctx context.Context, cid int64, company *Company) (*Company, error) {
	url := fd.endpoint("/companies/%d", cid)
	result := &Company{}
	if err := fd.doPut(ctx, url, company, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) DeleteCompany(ctx context.Context, cid int64) error {
	url := fd.endpoint("/companies/%d", cid)
	return fd.doDelete(ctx, url)
}

// ExportCompanies return a job id, call GetExportedCompaniesURL() to get the job detail
func (fd *Freshdesk) ExportCompanies(ctx context.Context, defaultFields, customFields []string) (string, error) {
	url := fd.endpoint("/companies/export")
	opt := &ExportOption{
		Fields: &ExportFields{
			DefaultFields: defaultFields,
			CustomFields:  customFields,
		},
	}
	job := &Job{}
	err := fd.doPost(ctx, url, opt, &job)
	return job.ID, err
}

// GetExportedCompaniesURL get the exported companies url
func (fd *Freshdesk) GetExportedCompaniesURL(ctx context.Context, jid string) (*Job, error) {
	url := fd.endpoint("/companies/export/%s", jid)
	job := &Job{}
	err := fd.doGet(ctx, url, job)
	return job, err
}
