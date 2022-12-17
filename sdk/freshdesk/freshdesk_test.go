package freshdesk

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/pandafw/pango/log"
	"github.com/pandafw/pango/str"
)

func TestWithFiles(t *testing.T) {
	var (
		tt WithFiles = &Ticket{}
		tc WithFiles = &Conversation{}
		ac WithFiles = &Contact{}
	)
	fmt.Println(tt, tc, ac)
}

func testNewFreshdesk(t *testing.T) *Freshdesk {
	apikey := os.Getenv("FD_APIKEY")
	if apikey == "" {
		t.Skip("FD_APIKEY not set")
		return nil
	}

	domain := os.Getenv("FD_DOMAIN")
	if domain == "" {
		t.Skip("FD_DOMAIN not set")
		return nil
	}

	logs := log.NewLog()
	//logs.SetLevel(log.LevelInfo)
	fd := &Freshdesk{
		Domain: domain,
		Apikey: apikey,
		Logger: logs.GetLogger("FDK"),
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

	err = fd.IterTickets(nil, func(t *Ticket) bool {
		fmt.Println(t)
		return true
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

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
	if len(ts) < 1 {
		t.Fatal("ListTickets return empty array")
	}
}

func TestIterTicket(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	ltp := &ListTicketsOption{PerPage: 2}
	err := fd.IterTickets(ltp, func(t *Ticket) bool {
		fmt.Println(t)
		return true
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}

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
	fmt.Println(cc)

	cu := &Contact{}
	cu.Description = "update description " + time.Now().String()
	cu.Avatar = NewAvatar("../../logo.png")

	uc, err := fd.UpdateContact(cc.ID, cu)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fmt.Println(uc)

	gc, err := fd.GetContact(cc.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fmt.Println(gc)

	err = fd.IterContacts(nil, func(c *Contact) bool {
		fmt.Println(c)
		return true
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	err = fd.DeleteContact(cc.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}

func TestAgentAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	ac := &AgentRequest{
		Email:       str.RandNumbers(11) + "@" + str.RandLetters(8) + ".com",
		TicketScope: AgentTicketScopeGlobal,
	}

	ca, err := fd.CreateAgent(ac)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fmt.Println(ca)

	au := &AgentRequest{
		Occasional: true,
	}
	//au.Avatar = NewAvatar("../../logo.png")

	ua, err := fd.UpdateAgent(ca.ID, au)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fmt.Println(ua)

	ga, err := fd.GetAgent(ua.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fmt.Println(ga)

	err = fd.IterAgents(nil, func(a *Agent) bool {
		fmt.Println(a)
		return true
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	err = fd.DeleteAgent(ga.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}

func TestListAgents(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	lao := &ListAgentsOption{PerPage: 10}
	as, _, err := fd.ListAgents(lao)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(as) < 1 {
		t.Fatal("ListAgents return empty array")
	}
	fmt.Println(as)
}
