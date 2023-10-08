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

	ot := &Ticket{
		Name:        "test",
		Phone:       "09012345678",
		Subject:     "test " + time.Now().String(),
		Description: "description " + time.Now().String(),
		Status:      TicketStatusOpen,
		Priority:    TicketPriorityMedium,
	}

	ct, err := fs.CreateTicket(ot)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	defer fs.DeleteTicket(ct.ID)

	tm := &TimeEntry{
		AgentID:   17001136838,
		TimeSpent: 70,
		Note:      "test time entry",
	}

	ctm, err := fs.CreateTimeEntry(ct.ID, tm)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Debug(ctm)

	tms, _, err := fs.ListTicketTimeEntries(ct.ID, nil)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
	fs.Logger.Debug(tms)

	err = fs.IterTicketTimeEntries(ct.ID, nil, func(te *TimeEntry) error {
		fs.Logger.Debug(te)
		return nil
	})
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}

	err = fs.DeleteTimeEntry(ct.ID, ctm.ID)
	if err != nil {
		t.Fatalf("ERROR: %v", err)
	}
}
