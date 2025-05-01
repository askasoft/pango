package freshservice

import (
	"testing"
	"time"
)

func TestTimeEntriesAPIs(t *testing.T) {
	fs := testNewFreshservice(t)
	if fs == nil {
		return
	}

	tc := &TicketCreate{
		Name:        "test",
		Phone:       "09012345678",
		Subject:     "test " + time.Now().String(),
		Description: "description " + time.Now().String(),
		Status:      TicketStatusOpen,
		Priority:    TicketPriorityMedium,
	}

	ct, err := fs.CreateTicket(ctxbg, tc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer fs.DeleteTicket(ctxbg, ct.ID)

	tmc := &TimeEntryCreate{
		AgentID:   17001136838,
		TimeSpent: 70,
		Note:      "test time entry",
	}

	ctm, err := fs.CreateTimeEntry(ctxbg, ct.ID, tmc)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Debug(ctm)

	tms, _, err := fs.ListTicketTimeEntries(ctxbg, ct.ID, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Debug(tms)

	err = fs.IterTicketTimeEntries(ctxbg, ct.ID, nil, func(te *TimeEntry) error {
		fs.Logger.Debug(te)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	err = fs.DeleteTimeEntry(ctxbg, ct.ID, ctm.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}
