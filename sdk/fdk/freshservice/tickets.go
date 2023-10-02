package freshservice

import (
	"strings"
)

// ---------------------------------------------------
// Ticket

const (
	TicketIncludeConversations  = "conversations"
	TicketIncludeRequester      = "requester"
	TicketIncludeRequestedFor   = "requested_for"
	TicketIncludeStats          = "stats"
	TicketIncludeProblem        = "problem"
	TicketIncludeAssets         = "assets"
	TicketIncludeChange         = "change"
	TicketIncludeRelatedTickets = "related_tickets"
)

type ListTicketsOption struct {
	Filter       string // The various filters available are new_and_my_open, watching, spam, deleted.
	WorkspaceID  int64
	RequestID    int64
	Type         string
	UpdatedSince Time
	Include      string // stats, requester, requester_for
	OrderType    string // asc, desc (default)
	Page         int
	PerPage      int
}

func (lto *ListTicketsOption) IsNil() bool {
	return lto == nil
}

func (lto *ListTicketsOption) Values() Values {
	q := Values{}
	q.SetString("filter", lto.Filter)
	q.SetInt64("workspace_id", lto.WorkspaceID)
	q.SetInt64("request_id", lto.RequestID)
	q.SetString("type", lto.Type)
	q.SetTime("updated_since", lto.UpdatedSince)
	q.SetString("include", lto.Include)
	q.SetString("order_type", lto.OrderType)
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

func (fs *Freshservice) CreateTicket(ticket *Ticket) (*Ticket, error) {
	url := fs.endpoint("/tickets")
	result := &ticketResult{}
	err := fs.doPost(url, ticket, result)
	return result.Ticket, err
}

// GetTicket Get a Ticket
// include: conversations, requester, requested_for, stats, problem, assets, change, related_tickets
func (fs *Freshservice) GetTicket(tid int64, include ...string) (*Ticket, error) {
	url := fs.endpoint("/tickets/%d", tid)
	if len(include) > 0 {
		s := strings.Join(include, ",")
		url += "?include=" + s
	}
	result := &ticketResult{}
	err := fs.doGet(url, result)
	return result.Ticket, err
}

// Filter Tickets
// Custom ticket fields that you have created in your account can also be used to filter through the tickets and get a list of tickets matching the specified ticket fields.
// Query Format(query) - "(ticket_field:integer OR ticket_field:'string') AND ticket_field:boolean"
// Note:
// 1. The query must be URL encoded
// 2. Query can be framed using the ticket field name in snake case, which can be obtained from Ticket Fields endpoint. Ticket Fields are case sensitive
// 3. Query string must be enclosed between a pair of double quotes and can have up to 512 characters
// 4. Logical operators AND, OR along with parentheses () can be used to group conditions
// 5. Relational operators greater than or equal to :> and less than or equal to :< can be used along with date fields and numeric fields
// 6. Input for date fields should be in UTC Format
// 7. The number of objects returned per page is 30 also the total count of the results will be returned along with the result
// 8. To scroll through the pages add page parameter to the url
// 9. To filter for agent and group with no values assigned, use the null keyword
// 10. By default, only tickets from the primary workspace will be returned for accounts with the 'Workspaces' feature enabled. For tickets from other workspaces, use the workspace_id filter.
// Supported Ticket Fields
// Field	Type	Description
// workspace_id	number	Workspace ID of the ticket. The attribute is applicable only for accounts with the 'Workspaces' feature enabled. The value 1 for workspace_id will return tickets from all workspaces, with only global level fields.
// requester_id	number	User ID of the requester
// email	string	Email address of the requester
// agent_id	number	ID of the agent to whom the ticket has been assigned
// group_id	number	ID of the group to which the ticket has been assigned
// priority	number	Priority of the ticket
// status	number	Status of the ticket
// impact	number	Ticket impact
// urgency	number	Ticket urgency
// tag	string	Tag that has been associated to the tickets
// due_by	date	Date (YYYY-MM-DD) when the ticket is due to be resolved
// fr_due_by	date	Date (YYYY-MM-DD) when the first response is due
// created_at	date	Ticket creation date (YYYY-MM-DD)
// Custom Fields
// Field Type	Type
// Number	number
// Checkbox	boolean
// Dropdown	string
// Date	date(YYYY-MM-DD)
// Note:
// Queries can be combined using AND or OR.
// https://domain.freshservice.com/api/v2/tickets/filter?query="priority: 1 AND status: 2 OR urgency: 3"
// Supported operators
// 1. priority: 1 (priority equal to 1)
// 2. priority:> 1 (priority greater than or equal to 1)
// 3. priority :< 1 (priority less than or equal to 1)
// Formatting
// 1. String fields to be enclosed in single quotes (”)
// 2. Number fields to be given as number without quotes.
// 3. Date and date_time fields to be enclosed in single quotes('yyyy-mm-dd')
// 4. only :> and :< are supported for date and date_time fields. Both fields expect input in the same format as 'yyyy-mm-dd'
func (fs *Freshservice) FilterTickets(fto *FilterTicketsOption) ([]*Ticket, bool, error) {
	url := fs.endpoint("/tickets/filter")
	result := &ticketResult{}
	next, err := fs.doList(url, fto, result)
	return result.Tickets, next, err
}

func (fs *Freshservice) IterFilterTickets(fto *FilterTicketsOption, itf func(*Ticket) error) error {
	if fto == nil {
		fto = &FilterTicketsOption{}
	}
	if fto.Page < 1 {
		fto.Page = 1
	}
	if fto.PerPage < 1 {
		fto.PerPage = 100
	}

	for {
		tickets, next, err := fs.FilterTickets(fto)
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
		fto.Page++
	}
	return nil
}

// List of Tickets
// Use filters to view only specific tickets (those which match the criteria that you choose). By default only tickets that have not been deleted or marked as spam will be returned, unless you use the 'deleted' filter.
// Note:
// 1. By default only tickets that have been created within the past 30 days will be returned. For older tickets, use the updated_since filter.
// 2. Use 'include' to embed additional details in the response. Each include will consume an additional 2 credits. For example if you embed the stats information you will be charged a total of 3 API credits (1 credit for the API call, and 2 credits for the additional stats embedding).
// 3. By default, only tickets from the primary workspace will be returned for accounts with the 'Workspaces' feature enabled. For tickets from other workspaces, use the workspace_id filter.
func (fs *Freshservice) ListTickets(lto *ListTicketsOption) ([]*Ticket, bool, error) {
	url := fs.endpoint("/tickets")
	result := &ticketResult{}
	next, err := fs.doList(url, lto, result)
	return result.Tickets, next, err
}

func (fs *Freshservice) IterTickets(lto *ListTicketsOption, itf func(*Ticket) error) error {
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
		tickets, next, err := fs.ListTickets(lto)
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

// Update a Ticket
// This API lets you make changes to the parameters of a ticket from updating statuses to changing ticket type.
// Note:
// 1. While updating tags, all the tags that need to stay associated with the Ticket should be provided in the PUT request payload.
// 2. The requested_for_id field can be updated only for Service Request tickets.
// Query Parameters	Handle
// bypass_mandatory: To bypass mandatory fields check while updating the ticket except for requester_id, source. Any business rules trying to mandate certain fields will also be bypassed. All fields configured as mandatory upon closing or resolving the ticket will be skipped while updating the ticket. This can only be passed by an admin.
func (fs *Freshservice) UpdateTicket(tid int64, ticket *Ticket) (*Ticket, error) {
	url := fs.endpoint("/tickets/%d", tid)
	result := &ticketResult{}
	err := fs.doPut(url, ticket, result)
	return result.Ticket, err
}

func (fs *Freshservice) DeleteTicket(tid int64) error {
	url := fs.endpoint("/tickets/%d", tid)
	return fs.doDelete(url)
}

func (fs *Freshservice) DeleteTicketAttachment(tid, aid int64) error {
	url := fs.endpoint("/tickets/%d/attachments/%d", tid, aid)
	return fs.doDelete(url)
}

// Restore a Ticket
// The API mentioned previously. If you deleted some tickets and regret doing so now, this API will help you restore them.
func (fs *Freshservice) Restore(tid int64) error {
	url := fs.endpoint("/tickets/%d/restore", tid)
	return fs.doPut(url, nil, nil)
}

// Create a Child Ticket
// This API lets you to create a new child ticket on an existing ticket
// Note:
// 1. Association of child tickets is not available in the Sprout plan.
// 2. Association of child tickets to a service request is not possible.
// 3. Association of child tickets to a deleted or a spammed ticket is not allowed.
// 4. Nesting of a child ticket under another child ticket is not supported.
func (fs *Freshservice) CreateChildTicket(tid int64, ticket *Ticket) (*Ticket, error) {
	url := fs.endpoint("/tickets/%d/create_child_ticket", tid)
	result := &ticketResult{}
	err := fs.doPost(url, ticket, result)
	return result.Ticket, err
}

func (fs *Freshservice) GetTicketFields() ([]*TicketField, error) {
	url := fs.endpoint("/ticket_form_fields")
	result := &ticketFieldResult{}
	err := fs.doGet(url, result)
	return result.TicketFields, err
}

func (fs *Freshservice) GetTicketActivities(tid int64) ([]*TicketActivity, error) {
	url := fs.endpoint("/tickets/%d/activities", tid)
	result := &ticketActivityResult{}
	err := fs.doGet(url, result)
	return result.TicketActivitys, err
}

// ---------------------------------------------------
// Conversation

func (fs *Freshservice) CreateReply(tid int64, reply *Reply) (*Conversation, error) {
	url := fs.endpoint("/tickets/%d/reply", tid)
	result := &conversationResult{}
	err := fs.doPost(url, reply, result)
	return result.Conversation, err
}

func (fs *Freshservice) CreateNote(tid int64, note *Note) (*Conversation, error) {
	url := fs.endpoint("/tickets/%d/notes", tid)
	result := &conversationResult{}
	err := fs.doPost(url, note, result)
	return result.Conversation, err
}

// Update a Conversation
// Only public & private notes can be edited.
func (fs *Freshservice) UpdateConversation(cid int64, conversation *Conversation) (*Conversation, error) {
	url := fs.endpoint("/conversations/%d", cid)
	result := &conversationResult{}
	err := fs.doPut(url, conversation, result)
	return result.Conversation, err
}

func (fs *Freshservice) DeleteConversation(cid int64) error {
	url := fs.endpoint("/conversations/%d", cid)
	return fs.doDelete(url)
}

func (fs *Freshservice) DeleteConversationAttachment(cid, aid int64) error {
	url := fs.endpoint("/conversations/%d/attachments/%d", cid, aid)
	return fs.doDelete(url)
}

func (fs *Freshservice) ListTicketConversations(tid int64, lco *ListConversationsOption) ([]*Conversation, bool, error) {
	url := fs.endpoint("/tickets/%d/conversations", tid)
	conversations := []*Conversation{}
	next, err := fs.doList(url, lco, &conversations)
	return conversations, next, err
}

func (fs *Freshservice) IterTicketConversations(tid int64, lco *ListConversationsOption, icf func(*Conversation) error) error {
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
		conversations, next, err := fs.ListTicketConversations(tid, lco)
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
