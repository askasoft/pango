package freshservice

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
func (fs *Freshservice) CreateTimeEntry(tid int64, tm *TimeEntry) (*TimeEntry, error) {
	url := fs.endpoint("/tickets/%d/time_entries", tid)
	result := &timeEntryResult{}
	err := fs.doPost(url, tm, result)
	return result.TimeEntry, err
}

// View a Time Entry
// This API call helps to list a particular Time Entry.
// GET  /api/v2/tickets/[ticket_id]/time_entries/[id]
func (fs *Freshservice) GetTimeEntry(tid, teid int64) (*TimeEntry, error) {
	url := fs.endpoint("/tickets/%d/time_entries/%d", tid, teid)
	result := &timeEntryResult{}
	err := fs.doGet(url, result)
	return result.TimeEntry, err
}

// List All Time Entries of a Ticket
// This API helps to view all time entries of a particular ticket.
// GET  /api/v2/tickets/[ticket_id]/time_entries
func (fs *Freshservice) ListTicketTimeEntries(tid int64, lteo *ListTimeEntriesOption) ([]*TimeEntry, bool, error) {
	url := fs.endpoint("/tickets/%d/time_entries", tid)
	result := &timeEntriesResult{}
	next, err := fs.doList(url, lteo, result)
	return result.TimeEntries, next, err
}

func (fs *Freshservice) IterTicketTimeEntries(tid int64, lteo *ListTimeEntriesOption, itef func(*TimeEntry) error) error {
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
		tms, next, err := fs.ListTicketTimeEntries(tid, lteo)
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
func (fs *Freshservice) UpdateTimeEntry(tid, teid int64, tm *TimeEntry) (*TimeEntry, error) {
	url := fs.endpoint("/tickets/%d/time_entries/%d", tid, teid)
	result := &timeEntryResult{}
	err := fs.doPut(url, tm, result)
	return result.TimeEntry, err
}

func (fs *Freshservice) DeleteTimeEntry(tid, teid int64) error {
	url := fs.endpoint("/tickets/%d/time_entries/%d", tid, teid)
	return fs.doDelete(url)
}
