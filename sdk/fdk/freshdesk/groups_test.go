package freshdesk

import (
	"testing"
)

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
