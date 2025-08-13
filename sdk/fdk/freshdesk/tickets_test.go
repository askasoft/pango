package freshdesk

import (
	"reflect"
	"testing"
	"time"

	"github.com/askasoft/pango/iox/fsu"
)

var (
	_ WithFiles = &TicketCreate{}
	_ WithFiles = &TicketUpdate{}
	_ WithFiles = &NoteCreate{}
	_ WithFiles = &NoteUpdate{}
	_ WithFiles = &ReplyCreate{}
)

func TestTicketAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	// tm1, _ := time.ParseInLocation("2006-1-2 15:04:05", "2000-01-02 10:20:30", time.Local)
	tc := &TicketCreate{
		Name:        "test",
		Phone:       "09012345678",
		Subject:     "test " + time.Now().String(),
		Description: "description " + time.Now().String(),
		Status:      TicketStatusPending,
		Priority:    TicketPriorityMedium,
		// CreatedAt:   &Time{Time: tm1}, // unsupport
	}

	ct, err := fd.CreateTicket(ctxbg, tc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(ct)

	tu := &TicketUpdate{}
	tu.Description = `<div>
<div>test05 - 非公開メモ</div>
<div>問い合わせです。</div>
<p> 外部 image</p><img src="https://github.com/askasoft/pango/raw/master/logo.png"><br/><br/><br/>
</div>`
	tu.AddAttachment("./ticket.go")

	ut, err := fd.UpdateTicket(ctxbg, ct.ID, tu)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	// find contact
	cs, _, err := fd.ListContacts(ctxbg, &ListContactsOption{
		Phone: "09012345678",
	})
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
	if len(cs) == 0 {
		t.Errorf("ERROR: %v", "missing contact")
	} else {
		fd.Logger.Debug(cs)
	}

	// public note
	// tm2, _ := time.ParseInLocation("2006-1-2 15:04:05", "2000-01-02 10:20:30", time.Local)
	nc1 := &NoteCreate{
		Body:   "public user note " + time.Now().String(),
		UserID: cs[0].ID,
		// CreatedAt: &Time{Time: tm2}, // unsupport
	}
	cnu, err := fd.CreateNote(ctxbg, ct.ID, nc1)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	} else {
		fd.Logger.Debug(cnu)
	}

	// private note
	nc2 := &NoteCreate{
		Body:    "private agent note " + time.Now().String(),
		Private: true,
	}
	cn, err := fd.CreateNote(ctxbg, ct.ID, nc2)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	} else {
		fd.Logger.Debug(cn)
	}

	nu := &NoteUpdate{
		Body: "private agent update note " + time.Now().String(),
	}
	nu.AddAttachment("./conversation.go")
	uc, err := fd.UpdateConversation(ctxbg, cn.ID, nu)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	} else {
		fd.Logger.Debug(uc)
	}

	gtr, err := fd.GetTicket(ctxbg, ct.ID, TicketIncludeRequester)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	} else {
		fd.Logger.Debug(gtr)
	}

	if len(gtr.Attachments) != 1 {
		t.Errorf("Attachments: %d", len(gtr.Attachments))
	} else {
		bs, err := fd.DownloadNoAuth(ctxbg, gtr.Attachments[0].AttachmentURL)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		} else {
			ws, _ := fsu.ReadFile("./ticket.go")
			if !reflect.DeepEqual(bs, ws) {
				t.Error("./ticket.go different")
			}
		}
	}

	gtc, err := fd.GetTicket(ctxbg, ct.ID, TicketIncludeConversations)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	} else {
		fd.Logger.Debug(gtc)
	}

	if len(gtc.Conversations) != 2 {
		t.Errorf("Conversations: %d", len(gtc.Conversations))
	} else {
		for _, c := range gtc.Conversations {
			if c.Private {
				if len(c.Attachments) != 1 {
					t.Errorf("Attachments: %d", len(c.Attachments))
				} else {
					bs, err := fd.DownloadNoAuth(ctxbg, c.Attachments[0].AttachmentURL)
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
		}
	}

	err = fd.DeleteAttachment(ctxbg, ut.Attachments[0].ID)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	err = fd.DeleteTicket(ctxbg, ct.ID)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
}

func TestListTickets(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	ltp := &ListTicketsOption{PerPage: 1}
	ts, _, err := fd.ListTickets(ctxbg, ltp)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	for i, t := range ts {
		fd.Logger.Infof("%d: #%d [%s] %s", i+1, t.ID, t.CreatedAt.String(), t.Subject)
	}
}

func TestIterTicketsAndConversations(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	ltp := &ListTicketsOption{PerPage: 10}

	i, j := 0, 0
	err := fd.IterTickets(ctxbg, ltp, func(t *Ticket) error {
		i++
		fd.Logger.Infof("%d: #%d [%s] %s", i, t.ID, t.CreatedAt.String(), t.Subject)

		j = 0
		fd.IterTicketConversations(ctxbg, t.ID, nil, func(c *Conversation) error {
			j++
			fd.Logger.Infof("#%d [%d]: Source: %v, Income: %v, UID: %d, THID: %d, FROM: %s, TO: %v, CC: %v",
				t.ID, j, c.Source, c.Incoming, c.UserID, c.ThreadID, c.FromEmail, c.ToEmails, c.CcEmails)
			return nil
		})
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

	ts, total, err := fd.FilterTickets(ctxbg, ftp)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	fd.Logger.Infof("Total: %d", total)
	for i, t := range ts {
		fd.Logger.Infof("%d: #%d [%s] %s", i+1, t.ID, t.CreatedAt.String(), t.Subject)
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
	err := fd.IterFilterTickets(ctxbg, ftp, func(t *Ticket) error {
		i++
		fd.Logger.Infof("%d: #%d [%s] %s", i, t.ID, t.CreatedAt.String(), t.Subject)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}
