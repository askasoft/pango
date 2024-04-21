package freshservice

import (
	"testing"
)

func TestIterApprovals(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	itcnt := 0
	err := fs.IterApprovals(&ListApprovalsOption{
		Parent: "ticket",
		Status: ApprovalStatusRequested.String(),
	}, func(a *Approval) error {
		itcnt++
		fs.Logger.Debugf("Iterate Approval #%d: %s", itcnt, a)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Infof("Iterate %d approvals", itcnt)
}
