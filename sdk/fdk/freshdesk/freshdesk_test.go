package freshdesk

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/str"
)

func TestWithFiles(t *testing.T) {
	var (
		tt WithFiles = &Ticket{}
		tc WithFiles = &Conversation{}
		at WithFiles = &Article{}
		ac WithFiles = &Contact{}
		//ag WithFiles = &Agent{}
	)
	fmt.Println(tt, tc, at, ac)
}

func testNewFreshdesk(t *testing.T) *FreshDesk {
	apikey := os.Getenv("FDK_APIKEY")
	if apikey == "" {
		t.Skip("FDK_APIKEY not set")
		return nil
	}

	domain := os.Getenv("FDK_DOMAIN")
	if domain == "" {
		t.Skip("FDK_DOMAIN not set")
		return nil
	}

	logs := log.NewLog()
	//logs.SetLevel(log.LevelDebug)
	fd := &FreshDesk{
		Domain:             domain,
		Apikey:             apikey,
		Logger:             logs.GetLogger("FDK"),
		RetryOnRateLimited: 1,
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

	nc := &Note{
		Body:    "create note " + time.Now().String(),
		Private: true,
	}
	cn, err := fd.CreateNote(ct.ID, nc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(cn)

	cu := &Conversation{
		Body: "update note " + time.Now().String(),
	}
	cu.AddAttachment("./any.go")
	uc, err := fd.UpdateConversation(cn.ID, cu)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(uc)

	gtc, err := fd.GetTicket(ct.ID, IncludeConversations)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(gtc)

	gtr, err := fd.GetTicket(ct.ID, IncludeRequester)
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
	err := fd.IterTickets(ltp, func(t *Ticket) bool {
		fd.Logger.Debug(t)
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

	err = fd.IterContacts(nil, func(c *Contact) bool {
		fd.Logger.Debug(c)
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
	fd.Logger.Debug(ca)

	au := &AgentRequest{
		Occasional: true,
	}
	//au.Avatar = NewAvatar("../../logo.png")

	ua, err := fd.UpdateAgent(ca.ID, au)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(ua)

	ga, err := fd.GetAgent(ua.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(ga)

	err = fd.IterAgents(nil, func(a *Agent) bool {
		fd.Logger.Debug(a)
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
	fd.Logger.Debug(as)
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

func TestListRoles(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	roles, err := fd.ListRoles()
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(roles)
}

func TestListGroups(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	groups, err := fd.ListGroups()
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(groups)
}

func TestAutomationAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	rules, err := fd.ListAutomationRules(AutomationTypeTicketCreation)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(rules)

	for _, rule := range rules {
		rule, err := fd.GetAutomationRule(AutomationTypeTicketCreation, rule.ID)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
		fd.Logger.Debug(rule)
		break
	}
}

func TestSolutionAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	cc := &Category{
		Name:        "Test Category",
		Description: "Test Category For API Test",
	}
	cat, err := fd.CreateCategory(cc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteCategory(cat.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	cf := &Folder{
		Name:        "Test Folder",
		Description: "Test Folder For API Test",
		Visibility:  FolderVisibilityAgents,
	}
	fol, err := fd.CreateFolder(cat.ID, cf)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteFolder(fol.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	ca := &Article{
		Title:       "Test Article",
		Description: "Test Article for API Test",
		Status:      ArticleStatusDraft,
	}
	art, err := fd.CreateArticle(fol.ID, ca)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer func() {
		err = fd.DeleteArticle(art.ID)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}
	}()

	art.AddAttachment("./any.go")
	_, err = fd.UpdateArticle(art.ID, art)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	arts, err := fd.ListFolderArticles(fol.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(arts) != 1 {
		t.Fatalf("ERROR: articles=%d", len(arts))
	}

	fols, err := fd.ListCategoryFolders(cat.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(fols) != 1 {
		t.Fatalf("ERROR: folders=%d", len(fols))
	}
}
