package freshdesk

import (
	"testing"
	"time"

	"github.com/askasoft/pango/ran"
)

var (
	_ WithFiles = &ContactCreate{}
	_ WithFiles = &ContactUpdate{}
)

func TestContactAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	cc := &ContactCreate{
		Mobile:      ran.RandNumbers(11),
		Description: "create description " + time.Now().String(),
	}
	cc.Name = "panda " + cc.Mobile

	c1, err := fd.CreateContact(ctxbg, cc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Info(c1)

	defer func() {
		err = fd.DeleteContact(ctxbg, c1.ID)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
	}()

	cu := &ContactUpdate{}
	cu.Description = "update description " + time.Now().String()
	cu.Avatar = NewAvatar("../../../logo.png")

	uc, err := fd.UpdateContact(ctxbg, c1.ID, cu)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Info(uc)

	gc, err := fd.GetContact(ctxbg, c1.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Info(gc)

	err = fd.IterContacts(ctxbg, nil, func(c *Contact) error {
		fd.Logger.Info(c)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}
func TestExportContacts(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	id, err := fd.ExportContacts(ctxbg, []string{"name", "email"}, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	job, err := fd.GetExportedContactsURL(ctxbg, id)
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

	c, err := fd.GetContact(ctxbg, 2043042684973)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	fd.Logger.Info(c)
}
