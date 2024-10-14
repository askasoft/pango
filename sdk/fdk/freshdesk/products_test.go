package freshdesk

import (
	"testing"
)

func TestProducts(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	itcnt := 0
	err := fd.IterProducts(nil, func(p *Product) error {
		itcnt++
		fd.Logger.Debugf("Iterate product #%d: %s", p.ID, p.Name)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Infof("Iterate %d products", itcnt)
}
