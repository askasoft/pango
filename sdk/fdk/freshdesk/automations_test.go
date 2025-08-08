package freshdesk

import (
	"testing"
)

func TestAutomationAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	rules, _, err := fd.ListAutomationRules(ctxbg, AutomationTypeTicketCreation, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(rules)

	for _, rule := range rules {
		rule, err := fd.GetAutomationRule(ctxbg, AutomationTypeTicketCreation, rule.ID)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
		fd.Logger.Info(rule)
	}
}
