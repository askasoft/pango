package freshservice

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
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	ot := &Ticket{
		Name:        "test",
		Phone:       "09012345678",
		Subject:     "test " + time.Now().String(),
		Description: "description " + time.Now().String(),
		Status:      TicketStatusOpen,
		Priority:    TicketPriorityMedium,
	}

	ct, err := fs.CreateTicket(ot)
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

	ut, err := fs.UpdateTicket(ct.ID, tu)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	err = fs.DeleteTicketAttachment(ut.ID, ut.Attachments[0].ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	nc := &Note{
		Body:    "create note " + time.Now().String(),
		Private: true,
	}
	cn, err := fs.CreateNote(ct.ID, nc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Debug(cn)

	cu := &Conversation{
		Body: "update note " + time.Now().String(),
	}
	cu.AddAttachment("./any.go")
	uc, err := fs.UpdateConversation(cn.ID, cu)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Debug(uc)

	gtc, err := fs.GetTicket(ct.ID, TicketIncludeConversations)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Debug(gtc)

	gtr, err := fs.GetTicket(ct.ID, TicketIncludeRequester)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Debug(gtr)

	err = fs.DeleteTicket(ct.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}

func TestListTicket(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	ltp := &ListTicketsOption{PerPage: 1}
	ts, _, err := fs.ListTickets(ltp)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Debug(ts)
}

func TestIterTicket(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	ltp := &ListTicketsOption{PerPage: 2}
	err := fs.IterTickets(ltp, func(t *Ticket) error {
		fs.Logger.Debug(t)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}
