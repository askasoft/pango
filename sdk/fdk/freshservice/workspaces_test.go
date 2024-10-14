package freshservice

import (
	"testing"
)

func TestWorkspaces(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	itcnt := 0
	err := fs.IterWorkspaces(nil, func(w *Workspace) error {
		itcnt++
		fs.Logger.Debugf("Iterate workspace #%d: %s", w.ID, w.Name)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Infof("Iterate %d workspaces", itcnt)
}
