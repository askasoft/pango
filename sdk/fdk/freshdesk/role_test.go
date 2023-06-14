package freshdesk

import (
	"testing"
)

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
