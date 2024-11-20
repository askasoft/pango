package freshdesk

import (
	"testing"
	"time"
)

func TestTimeEntriesAPIs(t *testing.T) {
	fd := testNewFreshdesk(t)
	if fd == nil {
		return
	}

	ot := &Ticket{
		Name:        "test",
		Phone:       "09012345678",
		Subject:     "test " + time.Now().String(),
		Description: "description " + time.Now().String(),
		Status:      TicketStatusClosed,
		Priority:    TicketPriorityMedium,
	}

	ct, err := fd.CreateTicket(ctxbg, ot)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	defer fd.DeleteTicket(ctxbg, ct.ID)

	tm := &TimeEntry{
		AgentID:   2043035275047,
		TimeSpent: 70,
		Note:      "test time entry",
	}

	ctm, err := fd.CreateTimeEntry(ctxbg, ct.ID, tm)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(ctm)

	tms, _, err := fd.ListTimeEntries(ctxbg, &ListTimeEntriesOption{
		AgentID: 2043035275047,
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fd.Logger.Debug(tms)

	err = fd.IterTimeEntries(ctxbg, &ListTimeEntriesOption{
		AgentID: 2043035275047,
	}, func(te *TimeEntry) error {
		fd.Logger.Debug(te)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	err = fd.DeleteTimeEntry(ctxbg, ctm.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}
