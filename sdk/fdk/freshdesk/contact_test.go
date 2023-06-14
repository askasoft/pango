package freshdesk

import (
	"testing"
	"time"

	"github.com/askasoft/pango/str"
)

func TestContactAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	cn := &Contact{
		Mobile:      str.RandNumbers(11),
		Description: "create description " + time.Now().String(),
	}
	cn.Name = "panda " + cn.Mobile

	cc, err := fd.CreateContact(cn)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(cc)

	cu := &Contact{}
	cu.Description = "update description " + time.Now().String()
	cu.Avatar = NewAvatar("../../../logo.png")

	uc, err := fd.UpdateContact(cc.ID, cu)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(uc)

	gc, err := fd.GetContact(cc.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(gc)

	err = fd.IterContacts(nil, func(c *Contact) error {
		fd.Logger.Debug(c)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	err = fd.DeleteContact(cc.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}
func TestExportContacts(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	id, err := fd.ExportContacts([]string{"name", "email"}, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	job, err := fd.GetExportedContactsURL(id)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(job)
}
