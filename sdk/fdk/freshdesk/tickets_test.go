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

	// tm1, _ := time.ParseInLocation("2006-1-2 15:04:05", "2000-01-02 10:20:30", time.Local)
	ot := &Ticket{
		Name:        "test",
		Phone:       "09012345678",
		Subject:     "test " + time.Now().String(),
		Description: "description " + time.Now().String(),
		Status:      TicketStatusClosed,
		Priority:    TicketPriorityMedium,
		// CreatedAt:   &Time{Time: tm1}, // unsupport
	}

	ct, err := fd.CreateTicket(ot)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(ct)

	tu := &Ticket{}
	tu.Description = `<div>
<div>test05 - 非公開メモ</div>
<div>問い合わせです。</div>
<p> 外部 image</p><img src="https://github.com/askasoft/pango/raw/master/logo.png"><br/><br/><br/>
</div>`
	tu.AddAttachment("./agent.go")

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
	// tm2, _ := time.ParseInLocation("2006-1-2 15:04:05", "2000-01-02 10:20:30", time.Local)
	nuc := &Note{
		Body:   "public user note " + time.Now().String(),
		UserID: cs[0].ID,
		// CreatedAt: &Time{Time: tm2}, // unsupport
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
	cu.AddAttachment("./agent.go")
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

func TestListTickets(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	ltp := &ListTicketsOption{PerPage: 1}
	ts, _, err := fd.ListTickets(ltp)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	for i, t := range ts {
		fd.Logger.Debugf("%d: #%d [%s] %s", i+1, t.ID, t.CreatedAt.String(), t.Subject)
	}
}

func TestIterTickets(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	ltp := &ListTicketsOption{PerPage: 10}

	i := 0
	err := fd.IterTickets(ltp, func(t *Ticket) error {
		i++
		fd.Logger.Debugf("%d: #%d [%s] %s", i, t.ID, t.CreatedAt.String(), t.Subject)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}

func TestFilterTickets(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	ftp := &FilterTicketsOption{
		Query: `"created_at:>'2023-10-01'"`,
	}

	ts, total, err := fd.FilterTickets(ftp)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	fd.Logger.Debugf("Total: %d", total)
	for i, t := range ts {
		fd.Logger.Debugf("%d: #%d [%s] %s", i+1, t.ID, t.CreatedAt.String(), t.Subject)
	}
}

func TestIterFilterTickets(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	ftp := &FilterTicketsOption{
		Query: `"created_at:>'2023-10-01'"`,
	}

	i := 0
	err := fd.IterFilterTickets(ftp, func(t *Ticket) error {
		i++
		fd.Logger.Debugf("%d: #%d [%s] %s", i, t.ID, t.CreatedAt.String(), t.Subject)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}
