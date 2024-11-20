package freshdesk

import (
	"testing"
)

func TestListCompanies(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	lco := &ListCompaniesOption{PerPage: 10}
	cs, _, err := fd.ListCompanies(ctxbg, lco)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(cs) < 1 {
		t.Fatal("ListCompanies return empty array")
	}
	fd.Logger.Debug(cs)
}

func TestExportCompany(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	id, err := fd.ExportCompanies(ctxbg, []string{"name", "description"}, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	job, err := fd.GetExportedCompaniesURL(ctxbg, id)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(job)
}
