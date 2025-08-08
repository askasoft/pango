package freshdesk

import "context"

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

func (fd *Freshdesk) CreateTimeEntry(ctx context.Context, tid int64, te *TimeEntryCreate) (*TimeEntry, error) {
	url := fd.Endpoint("/tickets/%d/time_entries", tid)
	result := &TimeEntry{}
	if err := fd.DoPost(ctx, url, te, result); err != nil {
		return nil, err
	}
	return result, nil
}

// List All Time Entries
func (fd *Freshdesk) ListTimeEntries(ctx context.Context, lteo *ListTimeEntriesOption) ([]*TimeEntry, bool, error) {
	url := fd.Endpoint("/time_entries")
	tes := []*TimeEntry{}
	next, err := fd.DoList(ctx, url, lteo, &tes)
	return tes, next, err
}

func (fd *Freshdesk) IterTimeEntries(ctx context.Context, lteo *ListTimeEntriesOption, itef func(*TimeEntry) error) error {
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
		tes, next, err := fd.ListTimeEntries(ctx, lteo)
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
func (fd *Freshdesk) UpdateTimeEntry(ctx context.Context, teid int64, te *TimeEntryUpdate) (*TimeEntry, error) {
	url := fd.Endpoint("/time_entries/%d", teid)
	result := &TimeEntry{}
	if err := fd.DoPut(ctx, url, te, result); err != nil {
		return nil, err
	}
	return result, nil
}

// Start/Stop Timer
// PUT  /api/v2/time_entries/[time_entry_id]/toggle_timer
func (fd *Freshdesk) ToggleTimer(ctx context.Context, teid int64) (*TimeEntry, error) {
	url := fd.Endpoint("/time_entries/%d/toggle_timer", teid)
	result := &TimeEntry{}
	if err := fd.DoPut(ctx, url, nil, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) DeleteTimeEntry(ctx context.Context, teid int64) error {
	url := fd.Endpoint("/time_entries/%d", teid)
	return fd.DoDelete(ctx, url)
}
