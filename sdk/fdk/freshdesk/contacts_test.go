package freshdesk

import (
	"testing"
	"time"

	"github.com/askasoft/pango/ran"
)

func TestContactAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	cn := &Contact{
		Mobile:      ran.RandNumbers(11),
		Description: "create description " + time.Now().String(),
	}
	cn.Name = "panda " + cn.Mobile

	cc, err := fd.CreateContact(cn)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Info(cc)

	cu := &Contact{}
	cu.Description = "update description " + time.Now().String()
	cu.Avatar = NewAvatar("../../../logo.png")

	uc, err := fd.UpdateContact(cc.ID, cu)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Info(uc)

	gc, err := fd.GetContact(cc.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Info(gc)

	err = fd.IterContacts(nil, func(c *Contact) error {
		fd.Logger.Info(c)
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
	fd.Logger.Info(job)
}

func TestGetContacts(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	c, err := fd.GetContact(2043042684973)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	fd.Logger.Info(c)
}
