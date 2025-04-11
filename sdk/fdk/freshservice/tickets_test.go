package freshservice

import (
	"reflect"
	"testing"
	"time"

	"github.com/askasoft/pango/fsu"
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

	tm1, _ := time.ParseInLocation("2006-1-2 15:04:05", "2000-01-02 10:20:30", time.Local)
	ot := &Ticket{
		Name:        "test",
		Phone:       "09012345678",
		Subject:     "test " + time.Now().String(),
		Description: "description " + time.Now().String(),
		Status:      TicketStatusPending,
		Priority:    TicketPriorityMedium,
		CustomFields: map[string]any{
			"rand75874209": true,
		},
		CreatedAt: &Time{Time: tm1},
	}

	ct, err := fs.CreateTicket(ctxbg, ot)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Debug(ct)

	tu := &Ticket{}
	tu.Description = `<div>
<div>test05 - 非公開メモ</div>
<div>問い合わせです。</div>
<p> 外部 image</p><img src="https://github.com/askasoft/pango/raw/master/logo.png"><br/><br/><br/>
</div>`
	tu.AddAttachment("./ticket.go")

	ut, err := fs.UpdateTicket(ctxbg, ct.ID, tu)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	} else {
		fs.Logger.Debug(ut)
	}

	tm2, _ := time.ParseInLocation("2006-1-2 15:04:05", "2000-01-02 10:20:30", time.Local)
	nc := &Note{
		Body:      "create note " + time.Now().String(),
		Private:   true,
		CreatedAt: &Time{Time: tm2},
	}

	cn, err := fs.CreateNote(ctxbg, ct.ID, nc)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	} else {
		fs.Logger.Debug(cn)
	}

	cu := &Conversation{
		Body: "update note " + time.Now().String(),
	}
	cu.AddAttachment("./conversation.go")
	uc, err := fs.UpdateConversation(ctxbg, cn.ID, cu)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	} else {
		fs.Logger.Debug(uc)
	}

	gtr, err := fs.GetTicket(ctxbg, ct.ID, TicketIncludeRequester)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	} else {
		fs.Logger.Debug(gtr)
	}

	if len(gtr.Attachments) != 1 {
		t.Errorf("Attachments: %d", len(gtr.Attachments))
	} else {
		bs, err := fs.DownloadNoAuth(ctxbg, gtr.Attachments[0].AttachmentURL)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		} else {
			ws, _ := fsu.ReadFile("./ticket.go")
			if !reflect.DeepEqual(bs, ws) {
				t.Error("./ticket.go different")
			}
		}
	}

	gtc, err := fs.GetTicket(ctxbg, ct.ID, TicketIncludeConversations)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	} else {
		fs.Logger.Debug(gtc)
	}

	if len(gtc.Conversations) != 1 {
		t.Errorf("Conversations: %d", len(gtc.Conversations))
	} else {
		c := gtc.Conversations[0]
		if len(c.Attachments) != 1 {
			t.Errorf("Attachments: %d", len(c.Attachments))
		} else {
			bs, err := fs.DownloadNoAuth(ctxbg, c.Attachments[0].AttachmentURL)
			if err != nil {
				t.Errorf("ERROR: %v", err)
			} else {
				ws, _ := fsu.ReadFile("./conversation.go")
				if !reflect.DeepEqual(bs, ws) {
					t.Error("./conversation.go different")
				}
			}
		}
	}

	err = fs.DeleteTicketAttachment(ctxbg, ut.ID, ut.Attachments[0].ID)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	err = fs.DeleteTicket(ctxbg, ct.ID)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
}

func TestListTicket(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	ltp := &ListTicketsOption{PerPage: 1}
	ts, _, err := fs.ListTickets(ctxbg, ltp)
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
	err := fs.IterTickets(ctxbg, ltp, func(t *Ticket) error {
		fs.Logger.Debug(t)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}

func TestListTicketFields(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	tfs, err := fs.ListTicketFields(ctxbg)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	fs.Logger.Debug(tfs)
}
