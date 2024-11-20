package freshdesk

import (
	"context"
	"testing"

	"github.com/askasoft/pango/ran"
)

var ctxbg = context.Background()

func TestAgentAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	ac := &AgentRequest{
		Email:       ran.RandNumbers(11) + "@" + ran.RandLetters(8) + ".com",
		TicketScope: AgentTicketScopeGlobal,
	}

	ca, err := fd.CreateAgent(ctxbg, ac)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(ca)

	au := &AgentRequest{
		Occasional: true,
	}
	//au.Avatar = NewAvatar("../../logo.png")

	ua, err := fd.UpdateAgent(ctxbg, ca.ID, au)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(ua)

	ga, err := fd.GetAgent(ctxbg, ua.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(ga)

	err = fd.IterAgents(ctxbg, nil, func(a *Agent) error {
		fd.Logger.Debug(a)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	err = fd.DeleteAgent(ctxbg, ga.ID)
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
	as, _, err := fd.ListAgents(ctxbg, lao)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if len(as) < 1 {
		t.Fatal("ListAgents return empty array")
	}
	fd.Logger.Debug(as)
}
