package freshdesk

import (
	"testing"
)

func TestTicketFieldsAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	otf := &TicketField{
		Label:                "testfieldlabel",
		LabelForCustomers:    "testfieldlabelforcustomers",
		Type:                 TicketFieldTypeCustomText,
		CustomersCanEdit:     true,
		DisplayedToCustomers: true,
	}

	ctf, err := fd.CreateTicketField(ctxbg, otf)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(ctf)

	tuf := &TicketField{LabelForCustomers: "testfieldlabelforcustomersupd"}

	utf, err := fd.UpdateTicketField(ctxbg, ctf.ID, tuf)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	} else {
		fd.Logger.Debug(utf)
	}

	gtfr, err := fd.GetTicketField(ctxbg, ctf.ID, TicketFieldIncludeSection)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	} else {
		fd.Logger.Debug(gtfr)
	}

	err = fd.DeleteTicketField(ctxbg, ctf.ID)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
}

func TestListTicketFieldsAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	tfs, err := fd.ListTicketFields(ctxbg)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(tfs)
}
