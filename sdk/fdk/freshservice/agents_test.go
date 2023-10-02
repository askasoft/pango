package freshservice

import (
	"testing"
)

func TestIterAgents(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	itcnt := 0
	err := fs.IterAgents(nil, func(a *Agent) error {
		itcnt++
		fs.Logger.Debugf("Iterate Agent [%d] %s %s <%s>", a.ID, a.FirstName, a.LastName, a.Email)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Infof("Iterate %d agents", itcnt)
}

func TestIterFilterAgents(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	fao := &FilterAgentsOption{
		Query: `last_name:'ou'`,
	}

	itcnt := 0
	err := fs.IterFilterAgents(fao, func(a *Agent) error {
		itcnt++
		fs.Logger.Debugf("Iterate Agent [%d] %s %s <%s>", a.ID, a.FirstName, a.LastName, a.Email)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Infof("Iterate %d agents", itcnt)
}
