package freshdesk

import (
	"context"
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

type FilterTicketsOption = FilterOption

type FilterTicketsResult struct {
	Total   int       `json:"total"`
	Results []*Ticket `json:"results"`
}

// PerPage: 1 ~ 100, default: 30
type ListConversationsOption = PageOption

func (fd *Freshdesk) CreateTicket(ctx context.Context, ticket *TicketCreate) (*Ticket, error) {
	url := fd.Endpoint("/tickets")
	result := &Ticket{}
	if err := fd.DoPost(ctx, url, ticket, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) CreateOutboundEmail(ctx context.Context, ticket *OutboundEmail) (*Ticket, error) {
	url := fd.Endpoint("/tickets/outbound_email")
	result := &Ticket{}
	if err := fd.DoPost(ctx, url, ticket, result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetTicket Get a Ticket
// include: conversations, requester, company, stats
func (fd *Freshdesk) GetTicket(ctx context.Context, tid int64, include ...string) (*Ticket, error) {
	url := fd.Endpoint("/tickets/%d", tid)
	if len(include) > 0 {
		s := strings.Join(include, ",")
		url += "?include=" + s
	}

	ticket := &Ticket{}
	err := fd.DoGet(ctx, url, ticket)
	return ticket, err
}

// List All Tickets
// Use filters to view only specific tickets (those which match the criteria that you choose).
// By default, only tickets that have not been deleted or marked as spam will be returned, unless you use the 'deleted' filter.
// Note:
// 1. By default, only tickets that have been created within the past 30 days will be returned. For older tickets, use the updated_since filter
// 2. A maximum of 300 pages (30000 tickets) will be returned.
// 3. When using filters, the query string must be URL encoded - see example
// 4. Use 'include' to embed additional details in the response. Each include will consume an additional 2 credits. For example if you embed the stats information you will be charged a total of 3 API credits for the call.
// 5. For accounts created after 2018-11-30, you will have to use include to get description.
func (fd *Freshdesk) ListTickets(ctx context.Context, lto *ListTicketsOption) ([]*Ticket, bool, error) {
	url := fd.Endpoint("/tickets")
	tickets := []*Ticket{}
	next, err := fd.DoList(ctx, url, lto, &tickets)
	return tickets, next, err
}

func (fd *Freshdesk) IterTickets(ctx context.Context, lto *ListTicketsOption, itf func(*Ticket) error) error {
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
		tickets, next, err := fd.ListTickets(ctx, lto)
		if err != nil {
			return err
		}
		for _, t := range tickets {
			if err = itf(t); err != nil {
				return err
			}
		}
		if !next || lto.Page >= 300 {
			break
		}
		lto.Page++
	}
	return nil
}

// FilterTickets
// Use custom ticket fields that you have created in your account to filter through the tickets and get a list of tickets matching the specified ticket fields.
// Query Format: "(ticket_field:integer OR ticket_field:'string') AND ticket_field:boolean"
// Note:
// 1. Archived tickets will not be included in the results
// 2. The query must be URL encoded
// 3. Query can be framed using the name of the ticket fields, which can be obtained from Ticket Fields endpoint. Ticket Fields are case sensitive
// 4. Query string must be enclosed between a pair of double quotes and can have up to 512 characters
// 5. Logical operators AND, OR along with parentheses () can be used to group conditions
// 6. Relational operators greater than or equal to :> and less than or equal to :< can be used along with date fields and numeric fields
// 7. Input for date fields should be in UTC Format
// 8. The number of objects returned per page is 30 also the total count of the results will be returned along with the result
// 9. To scroll through the pages add page parameter to the url. The page number starts with 1 and should not exceed 10
// 10. To filter for fields with no values assigned, use the null keyword
// 11. Please note that the updates will take a few minutes to get indexed, after which it will be available through API
func (fd *Freshdesk) FilterTickets(ctx context.Context, fto *FilterTicketsOption) ([]*Ticket, int, error) {
	url := fd.Endpoint("/search/tickets")
	ftr := &FilterTicketsResult{}
	_, err := fd.DoList(ctx, url, fto, ftr)
	return ftr.Results, ftr.Total, err
}

func (fd *Freshdesk) IterFilterTickets(ctx context.Context, fto *FilterTicketsOption, itf func(*Ticket) error) error {
	if fto == nil {
		fto = &FilterTicketsOption{}
	}
	if fto.Page < 1 {
		fto.Page = 1
	}

	for {
		tickets, total, err := fd.FilterTickets(ctx, fto)
		if err != nil {
			return err
		}
		for _, t := range tickets {
			if err = itf(t); err != nil {
				return err
			}
		}
		if len(tickets) < 30 || fto.Page >= 10 || (fto.Page-1)*30+len(tickets) >= total {
			break
		}
		fto.Page++
	}
	return nil
}

func (fd *Freshdesk) UpdateTicket(ctx context.Context, tid int64, ticket *TicketUpdate) (*Ticket, error) {
	url := fd.Endpoint("/tickets/%d", tid)
	result := &Ticket{}
	if err := fd.DoPut(ctx, url, ticket, result); err != nil {
		return nil, err
	}
	return result, nil
}

// BulkUpdateTickets returns job id
func (fd *Freshdesk) BulkUpdateTickets(ctx context.Context, tids []int64, properties *TicketProperties) (string, error) {
	url := fd.Endpoint("/tickets/bulk_update")
	data := map[string]any{
		"bulk_action": map[string]any{
			"ids":        tids,
			"properties": properties,
		},
	}
	result := map[string]string{}
	err := fd.DoPut(ctx, url, data, &result)
	return result["job_id"], err
}

func (fd *Freshdesk) ForwardTicket(ctx context.Context, tid int64, tf *TicketForward) (*ForwardResult, error) {
	url := fd.Endpoint("/tickets/%d/forward", tid)
	result := &ForwardResult{}
	if err := fd.DoPost(ctx, url, tf, result); err != nil {
		return nil, err
	}
	return result, nil
}

// MergeTickets
// Sometimes, a customer might try to get your attention regarding a particular issue by contacting you through separate channels.
// Sometimes, the same issue might be reported by different people in the team or someone might accidentally open a new ticket instead of following up on an existing one.
// To avoid conflicts, you can merge all related tickets together and keep the communication streamlined.
func (fd *Freshdesk) MergeTickets(ctx context.Context, tm *TicketsMerge) error {
	url := fd.Endpoint("/tickets/merge")
	err := fd.DoPut(ctx, url, tm, nil)
	return err
}

func (fd *Freshdesk) ListTicketWatchers(ctx context.Context, tid int64) ([]int64, error) {
	url := fd.Endpoint("/tickets/%d/watchers", tid)
	result := &TicketWatchers{}
	err := fd.DoGet(ctx, url, result)
	return result.WatcherIDs, err
}

func (fd *Freshdesk) AddTicketWatcher(ctx context.Context, tid, uid int64) error {
	url := fd.Endpoint("/tickets/%d/watchers", tid)
	data := map[string]any{
		"user_id": uid,
	}
	return fd.DoPost(ctx, url, data, nil)
}

func (fd *Freshdesk) UnwatchTicket(ctx context.Context, tid int64) error {
	url := fd.Endpoint("/tickets/%d/unwatch", tid)
	return fd.DoPut(ctx, url, nil, nil)
}

func (fd *Freshdesk) BulkWatchTickets(ctx context.Context, tids []int64, uid int64) error {
	url := fd.Endpoint("/tickets/buld_watch")
	data := map[string]any{
		"ids":     tids,
		"user_id": uid,
	}
	return fd.DoPut(ctx, url, data, nil)
}

func (fd *Freshdesk) BulkUnwatchTickets(ctx context.Context, tids []int64, uid int64) error {
	url := fd.Endpoint("/tickets/buld_unwatch")
	data := map[string]any{
		"ids":     tids,
		"user_id": uid,
	}
	return fd.DoPut(ctx, url, data, nil)
}

func (fd *Freshdesk) RestoreTicket(ctx context.Context, tid int64) error {
	url := fd.Endpoint("/tickets/%d/restore", tid)
	return fd.DoPut(ctx, url, nil, nil)
}

func (fd *Freshdesk) DeleteTicket(ctx context.Context, tid int64) error {
	url := fd.Endpoint("/tickets/%d", tid)
	return fd.DoDelete(ctx, url)
}

func (fd *Freshdesk) BulkDeleteTickets(ctx context.Context, tids []int64) (string, error) {
	url := fd.Endpoint("/tickets/bulk_delete")
	data := map[string]any{
		"bulk_action": map[string]any{
			"ids": tids,
		},
	}
	result := map[string]string{}
	err := fd.DoPut(ctx, url, data, &result)
	return result["job_id"], err
}

// ---------------------------------------------------
// Conversation

func (fd *Freshdesk) ListTicketConversations(ctx context.Context, tid int64, lco *ListConversationsOption) ([]*Conversation, bool, error) {
	url := fd.Endpoint("/tickets/%d/conversations", tid)
	conversations := []*Conversation{}
	next, err := fd.DoList(ctx, url, lco, &conversations)
	return conversations, next, err
}

func (fd *Freshdesk) IterTicketConversations(ctx context.Context, tid int64, lco *ListConversationsOption, icf func(*Conversation) error) error {
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
		conversations, next, err := fd.ListTicketConversations(ctx, tid, lco)
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

func (fd *Freshdesk) CreateReply(ctx context.Context, tid int64, reply *ReplyCreate) (*Reply, error) {
	url := fd.Endpoint("/tickets/%d/reply", tid)
	result := &Reply{}
	if err := fd.DoPost(ctx, url, reply, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) CreateNote(ctx context.Context, tid int64, note *NoteCreate) (*Note, error) {
	url := fd.Endpoint("/tickets/%d/notes", tid)
	result := &Note{}
	if err := fd.DoPost(ctx, url, note, result); err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateConversation only public & private notes can be edited.
func (fd *Freshdesk) UpdateConversation(ctx context.Context, cid int64, note *NoteUpdate) (*Conversation, error) {
	url := fd.Endpoint("/conversations/%d", cid)
	result := &Conversation{}
	if err := fd.DoPut(ctx, url, note, result); err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteConversation delete a conversation (Incoming Reply can not be deleted)
func (fd *Freshdesk) DeleteConversation(ctx context.Context, cid int64) error {
	url := fd.Endpoint("/conversations/%d", cid)
	return fd.DoDelete(ctx, url)
}

func (fd *Freshdesk) ReplyToForward(ctx context.Context, tid int64, rf *ReplyForward) (*ForwardResult, error) {
	url := fd.Endpoint("/tickets/%d/reply_to_forward", tid)
	result := &ForwardResult{}
	err := fd.DoPost(ctx, url, rf, result)
	return result, err
}
