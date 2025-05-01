package freshdesk

import (
	"fmt"
	"testing"

	"github.com/askasoft/pango/str"
)

func TestGroups(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	var gs []*Group
	defer func() {
		for _, g := range gs {
			err := fd.DeleteGroup(ctxbg, g.ID)
			if err != nil {
				fd.Logger.Errorf("Failed to delete group [%d] %s: %v", g.ID, g.Name, err)
			}
		}
	}()

	count := 100
	for i := 0; i < count; i++ {
		gc := &GroupCreate{
			Name: fmt.Sprintf("ApiTestGroup%03d", i+1),
		}

		cg, err := fd.CreateGroup(ctxbg, gc)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
		gs = append(gs, cg)
	}

	itcnt := 0
	err := fd.IterGroups(ctxbg, nil, func(g *Group) error {
		fd.Logger.Debugf("Iterate group [%d] %s", g.ID, g.Name)
		if str.Contains(g.Name, "ApiTestGroup") {
			itcnt++
		}
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	if count != itcnt {
		t.Fatalf("%d created, but %d iterated", count, itcnt)
	}
}
