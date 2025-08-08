package freshservice

import "context"

// ---------------------------------------------------
// Time Entries

// PerPage: 1 ~ 100, default: 30
type ListTimeEntriesOption = PageOption

// Create a Time Entry
// This API helps to create a Time Entry.
// Note:
// 1. If timer_running is not specified in the request, it is considered as false and time_spent is mandatory in this scenario.
// 2. time_spent can be set only if timer_running is false or not set.
// POST  /api/v2/tickets/[ticket_id]/time_entries
func (fs *Freshservice) CreateTimeEntry(ctx context.Context, tid int64, tm *TimeEntryCreate) (*TimeEntry, error) {
	url := fs.Endpoint("/tickets/%d/time_entries", tid)
	result := &timeEntryResult{}
	if err := fs.DoPost(ctx, url, tm, result); err != nil {
		return nil, err
	}
	return result.TimeEntry, nil
}

// View a Time Entry
// This API call helps to list a particular Time Entry.
// GET  /api/v2/tickets/[ticket_id]/time_entries/[id]
func (fs *Freshservice) GetTimeEntry(ctx context.Context, tid, teid int64) (*TimeEntry, error) {
	url := fs.Endpoint("/tickets/%d/time_entries/%d", tid, teid)
	result := &timeEntryResult{}
	err := fs.DoGet(ctx, url, result)
	return result.TimeEntry, err
}

// List All Time Entries of a Ticket
// This API helps to view all time entries of a particular ticket.
// GET  /api/v2/tickets/[ticket_id]/time_entries
func (fs *Freshservice) ListTicketTimeEntries(ctx context.Context, tid int64, lteo *ListTimeEntriesOption) ([]*TimeEntry, bool, error) {
	url := fs.Endpoint("/tickets/%d/time_entries", tid)
	result := &timeEntriesResult{}
	next, err := fs.DoList(ctx, url, lteo, result)
	return result.TimeEntries, next, err
}

func (fs *Freshservice) IterTicketTimeEntries(ctx context.Context, tid int64, lteo *ListTimeEntriesOption, itef func(*TimeEntry) error) error {
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
		tms, next, err := fs.ListTicketTimeEntries(ctx, tid, lteo)
		if err != nil {
			return err
		}
		for _, t := range tms {
			if err = itef(t); err != nil {
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
// This API can be used to update/modify existing time entry.
// Note:
// 1. For a running timer, time_spent cannot be updated without stopping it.
func (fs *Freshservice) UpdateTimeEntry(ctx context.Context, tid, teid int64, tm *TimeEntryUpdate) (*TimeEntry, error) {
	url := fs.Endpoint("/tickets/%d/time_entries/%d", tid, teid)
	result := &timeEntryResult{}
	if err := fs.DoPut(ctx, url, tm, result); err != nil {
		return nil, err
	}
	return result.TimeEntry, nil
}

func (fs *Freshservice) DeleteTimeEntry(ctx context.Context, tid, teid int64) error {
	url := fs.Endpoint("/tickets/%d/time_entries/%d", tid, teid)
	return fs.DoDelete(ctx, url)
}
