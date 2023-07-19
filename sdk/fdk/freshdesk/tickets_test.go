package freshdesk

import (
	"testing"
	"time"

	"github.com/askasoft/pango/log"
)

func TestTicketTypes(t *testing.T) {
	var tt WithFiles = &Ticket{}
	var n WithFiles = &Note{}
	var r WithFiles = &Reply{}
	var c WithFiles = &Conversation{}

	log.Trace(tt, n, r, c)
}

func TestTicketAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	ot := &Ticket{
		Name:        "test",
		Phone:       "09012345678",
		Subject:     "test " + time.Now().String(),
		Description: "description " + time.Now().String(),
		Status:      TicketStatusClosed,
		Priority:    TicketPriorityMedium,
	}

	ct, err := fd.CreateTicket(ot)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	tu := &Ticket{}
	tu.Description = `<div>
<div>test05 - 非公開メモ</div>
<div>問い合わせです。</div>
<p> 外部 image</p><img src="https://github.com/askasoft/pango/raw/master/logo.png"><br/><br/><br/>
</div>`
	tu.AddAttachment("./any.go")

	ut, err := fd.UpdateTicket(ct.ID, tu)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	err = fd.DeleteAttachment(ut.Attachments[0].ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	// find contact
	cs, _, err := fd.ListContacts(&ListContactsOption{
		Phone: "09012345678",
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(cs) == 0 {
		t.Fatalf("ERROR: %v", "missing contact")
	}

	// public note
	nuc := &Note{
		Body:   "public user note " + time.Now().String(),
		UserID: cs[0].ID,
	}
	cnu, err := fd.CreateNote(ct.ID, nuc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(cnu)

	// private note
	nc := &Note{
		Body:    "private agent note " + time.Now().String(),
		Private: true,
	}
	cn, err := fd.CreateNote(ct.ID, nc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(cn)

	cu := &Conversation{
		Body: "private agent update note " + time.Now().String(),
	}
	cu.AddAttachment("./any.go")
	uc, err := fd.UpdateConversation(cn.ID, cu)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(uc)

	gtc, err := fd.GetTicket(ct.ID, TicketIncludeConversations)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(gtc)

	gtr, err := fd.GetTicket(ct.ID, TicketIncludeRequester)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(gtr)

	err = fd.DeleteTicket(ct.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}

func TestListTicket(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	ltp := &ListTicketsOption{PerPage: 1}
	ts, _, err := fd.ListTickets(ltp)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(ts)
}

func TestIterTicket(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	ltp := &ListTicketsOption{PerPage: 2}
	err := fd.IterTickets(ltp, func(t *Ticket) error {
		fd.Logger.Debug(t)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}
