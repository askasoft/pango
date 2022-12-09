package freshdesk

import (
	"os"
	"testing"
	"time"

	"github.com/pandafw/pango/log"
)

func testNewFreshdesk(t *testing.T) *Freshdesk {
	apikey := os.Getenv("FD_API_KEY")
	if apikey == "" {
		t.Skip("FD_API_KEY not set")
		return nil
	}

	domain := os.Getenv("FD_DOMAIN")
	if domain == "" {
		t.Skip("FD_DOMAIN not set")
		return nil
	}

	logs := log.NewLog()
	logs.SetLevel(log.LevelTrace)
	fd := &Freshdesk{
		Domain: domain,
		Apikey: apikey,
		Logger: logs.GetLogger("FD"),
	}

	return fd
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
		Status:      TicketStatusOpen,
		Priority:    TicketPriorityMedium,
	}

	ct, err := fd.CreateTicket(ot)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	gt, err := fd.GetTicket(ct.ID, IncludeConversations)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if gt.ID != ct.ID {
		t.Fatalf("TicketID: %v, want %v", gt.ID, ct.ID)
	}

	tu := &Ticket{}
	tu.Description = `<div>
<div>test05 - 非公開メモ</div>
<div>問い合わせです。</div>
<p> 外部 image</p><img src="https://github.com/pandafw/pango/raw/master/logo.png"><br/><br/><br/>
</div>`
	tu.AddAttachment("./avatar.go")

	ut, err := fd.UpdateTicket(ct.ID, tu)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if ut.ID != ct.ID {
		t.Fatalf("TicketID: %v, want %v", ut.ID, ct.ID)
	}

	on := &Note{
		Body:    "create note " + time.Now().String(),
		Private: true,
	}
	nn, err := fd.CreateNote(ct.ID, on)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if nn.TicketID != ct.ID {
		t.Fatalf("TicketID: %v, want %v", nn.TicketID, ct.ID)
	}

	uc := &Conversation{
		Body: "update note " + time.Now().String(),
	}
	uc.AddAttachment("./conversation.go")
	nc, err := fd.UpdateConversation(nn.ID, uc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if nc.TicketID != ct.ID {
		t.Fatalf("TicketID: %v, want %v", nc.TicketID, ct.ID)
	}

	ltp := &ListTicketsOption{}
	ts, err := fd.ListTickets(ltp)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(ts) < 1 {
		t.Fatal("ListTickets return empty array")
	}

	err = fd.DeleteTicket(ct.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}
