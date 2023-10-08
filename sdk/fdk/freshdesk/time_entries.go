package freshdesk

// ---------------------------------------------------
// Time Entries

type TimeEntryBillable string

const (
	TimeEntryBillableTrue  TimeEntryBillable = "true"
	TimeEntryBillableFalse TimeEntryBillable = "false"
	TimeEntryBillableALL   TimeEntryBillable = ""
)

type ListTimeEntriesOption struct {
	CompanyID      int64
	AgentID        int64
	ExecutedAfter  *Time
	ExecutedBefore *Time
	Billable       TimeEntryBillable
	Page           int
	PerPage        int
}

func (lteo *ListTimeEntriesOption) IsNil() bool {
	return lteo == nil
}

func (lteo *ListTimeEntriesOption) Values() Values {
	q := Values{}
	q.SetInt64("company_id", lteo.CompanyID)
	q.SetInt64("agent_id", lteo.AgentID)
	q.SetTimePtr("executed_after", lteo.ExecutedAfter)
	q.SetTimePtr("executed_before", lteo.ExecutedBefore)
	q.SetString("billable", string(lteo.Billable))
	q.SetInt("page", lteo.Page)
	q.SetInt("per_page", lteo.PerPage)
	return q
}

func (fd *Freshdesk) CreateTimeEntry(tid int64, te *TimeEntry) (*TimeEntry, error) {
	url := fd.endpoint("/tickets/%d/time_entries", tid)
	result := &TimeEntry{}
	err := fd.doPost(url, te, result)
	return result, err
}

// List All Time Entries
func (fd *Freshdesk) ListTimeEntries(lteo *ListTimeEntriesOption) ([]*TimeEntry, bool, error) {
	url := fd.endpoint("/time_entries")
	tes := []*TimeEntry{}
	next, err := fd.doList(url, lteo, &tes)
	return tes, next, err
}

func (fd *Freshdesk) IterTimeEntries(lteo *ListTimeEntriesOption, itef func(*TimeEntry) error) error {
	if lteo == nil {
		lteo = &ListTimeEntriesOption{}
	}
	if lteo.Page < 1 {
		lteo.Page = 1
	}
	if lteo.PerPage < 1 {
		lteo.PerPage = 100
	}

	for {
		tes, next, err := fd.ListTimeEntries(lteo)
		if err != nil {
			return err
		}
		for _, te := range tes {
			if err = itef(te); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lteo.Page++
	}
	return nil
}

// Update a Time Entry
// PUT  /api/v2/time_entries/[id]
// Note:
// 1. The start_time cannot be updated if the timer is already running
// 2. The start_time cannot be be updated unless the timer_running attribute is updated from 'true' to 'false'
// 3. The start_time cannot be greater than the current time
// 4. The timer_running attribute cannot be set to the same value as before
// 5. The agent_id cannot be updated if the timer is already running
func (fd *Freshdesk) UpdateTimeEntry(teid int64, te *TimeEntry) (*TimeEntry, error) {
	url := fd.endpoint("/time_entries/%d", teid)
	result := &TimeEntry{}
	err := fd.doPut(url, te, result)
	return result, err
}

// Start/Stop Timer
// PUT  /api/v2/time_entries/[time_entry_id]/toggle_timer
func (fd *Freshdesk) ToggleTimer(teid int64) (*TimeEntry, error) {
	url := fd.endpoint("/time_entries/%d/toggle_timer", teid)
	result := &TimeEntry{}
	err := fd.doPut(url, nil, result)
	return result, err
}

func (fd *Freshdesk) DeleteTimeEntry(teid int64) error {
	url := fd.endpoint("/time_entries/%d", teid)
	return fd.doDelete(url)
}
