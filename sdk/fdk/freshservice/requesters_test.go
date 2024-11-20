package freshservice

import (
	"testing"
)

func TestIterRequesters(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	itcnt := 0
	err := fs.IterRequesters(ctxbg, nil, func(r *Requester) error {
		fs.Logger.Debugf("Iterate requester [%d] %s %s <%s> (%s, %s)",
			r.ID, r.FirstName, r.LastName, r.PrimaryEmail, r.MobilePhoneNumber, r.WorkPhoneNumber)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Infof("Iterate %d requesters", itcnt)
}
