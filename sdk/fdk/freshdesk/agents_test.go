package freshdesk

import (
	"testing"

	"github.com/askasoft/pango/str"
)

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

	err = fd.IterAgents(nil, func(a *Agent) error {
		fd.Logger.Debug(a)
		return nil
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
