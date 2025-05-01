package freshservice

import (
	"context"
	"fmt"
	"testing"

	"github.com/askasoft/pango/str"
)

var ctxbg = context.Background()

func TestAgentGroups(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	var ags []*AgentGroup
	defer func() {
		for _, ag := range ags {
			err := fs.DeleteAgentGroup(ctxbg, ag.ID)
			if err != nil {
				fs.Logger.Errorf("Failed to delete group [%d] %s: %v", ag.ID, ag.Name, err)
			}
		}
	}()

	count := 100
	for i := 0; i < count; i++ {
		agc := &AgentGroupCreate{
			Name: fmt.Sprintf("ApiTestGroup%03d", i+1),
		}

		cg, err := fs.CreateAgentGroup(ctxbg, agc)
		if err != nil {
			t.Fatalf("ERROR: %v", err)
		}
		ags = append(ags, cg)
	}

	itcnt := 0
	err := fs.IterAgentGroups(ctxbg, nil, func(ag *AgentGroup) error {
		fs.Logger.Debugf("Iterate group [%d] %s", ag.ID, ag.Name)
		if str.Contains(ag.Name, "ApiTestGroup") {
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
	fs.Logger.Infof("Iterate %d groups", itcnt)
}
