package freshdesk

import (
	"strings"
)

// ---------------------------------------------------
// Ticket

type TicketOrderBy string

const (
	TicketIncludeDescription   = "description"   // Will return the ticket description and description_text.
	TicketIncludeCompany       = "company"       // Will return the company's id and name.
	TicketIncludeConversations = "conversations" // Will return up to 10 conversations sorted by "created_at" in ascending order
	TicketIncludeRequester     = "requester"     // Will return the requester's email, id, mobile, name, and phone
	TicketIncludeStats         = "stats"         // Will return the ticket's closed_at, resolved_at and first_responded_at time

	TicketOrderByCreatedAt TicketOrderBy = "created_at"
	TicketOrderByDueBy     TicketOrderBy = "due_by"
	TicketOrderByUpdatedAt TicketOrderBy = "updated_at"
	TicketOrderByStatus    TicketOrderBy = "status"
)

type ListTicketsOption struct {
	Filter           string // The various filters available are new_and_my_open, watching, spam, deleted.
	RequestID        int64
	Email            string
	UniqueExternalID string
	CompanyID        int64
	UpdatedSince     Time
	Include          string        // stats, requester, description
	OrderBy          TicketOrderBy // created_at, due_by, updated_at, status
	OrderType        OrderType     // asc, desc (default)
	Page             int
	PerPage          int
}

func (lto *ListTicketsOption) IsNil() bool {
	return lto == nil
}

func (lto *ListTicketsOption) Values() Values {
	q := Values{}
	q.SetString("filter", lto.Filter)
	q.SetInt64("request_id", lto.RequestID)
	q.SetString("email", lto.Email)
	q.SetString("unique_external_id", lto.UniqueExternalID)
	q.SetInt64("company_id", lto.CompanyID)
	q.SetTime("updated_since", lto.UpdatedSince)
	q.SetString("include", lto.Include)
	q.SetString("order_by", (string)(lto.OrderBy))
	q.SetString("order_type", (string)(lto.OrderType))
	q.SetInt("page", lto.Page)
	q.SetInt("per_page", lto.PerPage)
	return q
}

type FilterTicketsOption struct {
	Query   string
	Page    int
	PerPage int
}

func (fto *FilterTicketsOption) IsNil() bool {
	return fto == nil
}

func (fto *FilterTicketsOption) Values() Values {
	q := Values{}
	q.SetString("query", fto.Query)
	q.SetInt("page", fto.Page)
	q.SetInt("per_page", fto.PerPage)
	return q
}

// PerPage: 1 ~ 100, default: 30
type ListConversationsOption = PageOption

func (fd *Freshdesk) CreateTicket(ticket *Ticket) (*Ticket, error) {
	url := fd.endpoint("/tickets")
	result := &Ticket{}
	err := fd.doPost(url, ticket, result)
	return result, err
}

// GetTicket Get a Ticket
// include: conversations, requester, company, stats
func (fd *Freshdesk) GetTicket(tid int64, include ...string) (*Ticket, error) {
	url := fd.endpoint("/tickets/%d", tid)
	if len(include) > 0 {
		s := strings.Join(include, ",")
		url += "?include=" + s
	}

	ticket := &Ticket{}
	err := fd.doGet(url, ticket)
	return ticket, err
}

// List All Tickets
// Use filters to view only specific tickets (those which match the criteria that you choose).
// By default, only tickets that have not been deleted or marked as spam will be returned, unless you use the 'deleted' filter.
// Note:
// 1. By default, only tickets that have been created within the past 30 days will be returned. For older tickets, use the updated_since filter
// 2. A maximum of 300 pages (9000 tickets) will be returned.
// 3. When using filters, the query string must be URL encoded - see example
// 4. Use 'include' to embed additional details in the response. Each include will consume an additional 2 credits. For example if you embed the stats information you will be charged a total of 3 API credits for the call.
// 5. For accounts created after 2018-11-30, you will have to use include to get description.
func (fd *Freshdesk) ListTickets(lto *ListTicketsOption) ([]*Ticket, bool, error) {
	url := fd.endpoint("/tickets")
	tickets := []*Ticket{}
	next, err := fd.doList(url, lto, &tickets)
	return tickets, next, err
}

func (fd *Freshdesk) IterTickets(lto *ListTicketsOption, itf func(*Ticket) error) error {
	if lto == nil {
		lto = &ListTicketsOption{}
	}
	if lto.Page < 1 {
		lto.Page = 1
	}
	if lto.PerPage < 1 {
		lto.PerPage = 100
	}

	for {
		tickets, next, err := fd.ListTickets(lto)
		if err != nil {
			return err
		}
		for _, t := range tickets {
			if err = itf(t); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lto.Page++
	}
	return nil
}

// FilterTickets
// Use custom ticket fields that you have created in your account to filter through the tickets and get a list of tickets matching the specified ticket fields.
// Query Format: "(ticket_field:integer OR ticket_field:'string') AND ticket_field:boolean"
func (fd *Freshdesk) FilterTickets(fto *FilterTicketsOption) ([]*Ticket, bool, error) {
	url := fd.endpoint("/search/tickets")
	tickets := []*Ticket{}
	next, err := fd.doList(url, fto, &tickets)
	return tickets, next, err
}

func (fd *Freshdesk) UpdateTicket(tid int64, ticket *Ticket) (*Ticket, error) {
	url := fd.endpoint("/tickets/%d", tid)
	result := &Ticket{}
	err := fd.doPut(url, ticket, result)
	return result, err
}

// BulkUpdateTickets returns job id
func (fd *Freshdesk) BulkUpdateTickets(tids []int64, properties *TicketProperties) (string, error) {
	url := fd.endpoint("/tickets/bulk_update")
	data := map[string]any{
		"bulk_action": map[string]any{
			"ids":        tids,
			"properties": properties,
		},
	}
	result := map[string]string{}
	err := fd.doPut(url, data, &result)
	return result["job_id"], err
}

func (fd *Freshdesk) ForwardTicket(tid int64, tf *TicketForward) (*ForwardResult, error) {
	url := fd.endpoint("/tickets/%d/forward", tid)
	result := &ForwardResult{}
	err := fd.doPost(url, tf, result)
	return result, err
}

// MergeTickets
// Sometimes, a customer might try to get your attention regarding a particular issue by contacting you through separate channels.
// Sometimes, the same issue might be reported by different people in the team or someone might accidentally open a new ticket instead of following up on an existing one.
// To avoid conflicts, you can merge all related tickets together and keep the communication streamlined.
func (fd *Freshdesk) MergeTickets(tm *TicketsMerge) error {
	url := fd.endpoint("/tickets/merge")
	err := fd.doPut(url, tm, nil)
	return err
}

func (fd *Freshdesk) ListTicketWatchers(tid int64) ([]int64, error) {
	url := fd.endpoint("/tickets/%d/watchers", tid)
	result := &TicketWatchers{}
	err := fd.doGet(url, result)
	return result.WatcherIDs, err
}

func (fd *Freshdesk) AddTicketWatcher(tid, uid int64) error {
	url := fd.endpoint("/tickets/%d/watchers", tid)
	data := map[string]any{
		"user_id": uid,
	}
	return fd.doPost(url, data, nil)
}

func (fd *Freshdesk) UnwatchTicket(tid int64) error {
	url := fd.endpoint("/tickets/%d/unwatch", tid)
	return fd.doPut(url, nil, nil)
}

func (fd *Freshdesk) BulkWatchTickets(tids []int64, uid int64) error {
	url := fd.endpoint("/tickets/buld_watch")
	data := map[string]any{
		"ids":     tids,
		"user_id": uid,
	}
	return fd.doPut(url, data, nil)
}

func (fd *Freshdesk) BulkUnwatchTickets(tids []int64, uid int64) error {
	url := fd.endpoint("/tickets/buld_unwatch")
	data := map[string]any{
		"ids":     tids,
		"user_id": uid,
	}
	return fd.doPut(url, data, nil)
}

func (fd *Freshdesk) RestoreTicket(tid int64) error {
	url := fd.endpoint("/tickets/%d/restore", tid)
	return fd.doPut(url, nil, nil)
}

func (fd *Freshdesk) DeleteTicket(tid int64) error {
	url := fd.endpoint("/tickets/%d", tid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) BulkDeleteTickets(tids []int64) (string, error) {
	url := fd.endpoint("/tickets/bulk_delete")
	data := map[string]any{
		"bulk_action": map[string]any{
			"ids": tids,
		},
	}
	result := map[string]string{}
	err := fd.doPut(url, data, &result)
	return result["job_id"], err
}

// ---------------------------------------------------
// Conversation

func (fd *Freshdesk) ListTicketConversations(tid int64, lco *ListConversationsOption) ([]*Conversation, bool, error) {
	url := fd.endpoint("/tickets/%d/conversations", tid)
	conversations := []*Conversation{}
	next, err := fd.doList(url, lco, &conversations)
	return conversations, next, err
}

func (fd *Freshdesk) IterTicketConversations(tid int64, lco *ListConversationsOption, icf func(*Conversation) error) error {
	if lco == nil {
		lco = &ListConversationsOption{}
	}
	if lco.Page < 1 {
		lco.Page = 1
	}
	if lco.PerPage < 1 {
		lco.PerPage = 100
	}

	for {
		conversations, next, err := fd.ListTicketConversations(tid, lco)
		if err != nil {
			return err
		}
		for _, c := range conversations {
			if err = icf(c); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lco.Page++
	}
	return nil
}

func (fd *Freshdesk) CreateReply(tid int64, reply *Reply) (*Reply, error) {
	url := fd.endpoint("/tickets/%d/reply", tid)
	result := &Reply{}
	err := fd.doPost(url, reply, result)
	return result, err
}

func (fd *Freshdesk) CreateNote(tid int64, note *Note) (*Note, error) {
	url := fd.endpoint("/tickets/%d/notes", tid)
	result := &Note{}
	err := fd.doPost(url, note, result)
	return result, err
}

// UpdateConversation only public & private notes can be edited.
func (fd *Freshdesk) UpdateConversation(cid int64, conversation *Conversation) (*Conversation, error) {
	url := fd.endpoint("/conversations/%d", cid)
	result := &Conversation{}
	err := fd.doPut(url, conversation, result)
	return result, err
}

// DeleteConversation delete a conversation (Incoming Reply can not be deleted)
func (fd *Freshdesk) DeleteConversation(cid int64) error {
	url := fd.endpoint("/conversations/%d", cid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) ReplyToForward(tid int64, rf *ReplyForward) (*ForwardResult, error) {
	url := fd.endpoint("/tickets/%d/reply_to_forward", tid)
	result := &ForwardResult{}
	err := fd.doPost(url, rf, result)
	return result, err
}
