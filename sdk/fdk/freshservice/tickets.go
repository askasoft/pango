package freshservice

import (
	"context"
	"strings"
)

// ---------------------------------------------------
// Ticket

const (
	TicketIncludeTags               = "tags"
	TicketIncludeConversations      = "conversations"
	TicketIncludeRequester          = "requester"
	TicketIncludeRequestedFor       = "requested_for"
	TicketIncludeStats              = "stats"
	TicketIncludeProblem            = "problem"
	TicketIncludeAssets             = "assets"
	TicketIncludeChanges            = "changes"
	TicketIncludeRelatedTickets     = "related_tickets"
	TicketIncludeOnboardingContext  = "onboarding_context"
	TicketIncludeOffboardingContext = "offboarding_context"

	TicketTypeIncident       = "Incident"
	TicketTypeServiceRequest = "Service Request"
)

type ListTicketsOption struct {
	Filter       string // The various filters available are new_and_my_open, watching, spam, deleted.
	WorkspaceID  int64
	RequestID    int64
	Email        string
	Type         string // Incident, Service Request
	UpdatedSince Time
	Include      string    // stats, requester, requester_for
	OrderType    OrderType // asc, desc (default)
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
	q.SetString("email", lto.Email)
	q.SetString("type", lto.Type)
	q.SetTime("updated_since", lto.UpdatedSince)
	q.SetString("include", lto.Include)
	q.SetString("order_type", string(lto.OrderType))
	q.SetInt("page", lto.Page)
	q.SetInt("per_page", lto.PerPage)
	return q
}

type FilterTicketsOption = FilterOption

// PerPage: 1 ~ 100, default: 30
type ListConversationsOption = PageOption

func (fs *Freshservice) CreateTicket(ctx context.Context, ticket *TicketCreate) (*Ticket, error) {
	url := fs.Endpoint("/tickets")
	result := &ticketResult{}
	if err := fs.DoPost(ctx, url, ticket, result); err != nil {
		return nil, err
	}
	return result.Ticket, nil
}

// GetTicket Get a Ticket
// include: conversations, requester, requested_for, stats, problem, assets, change, related_tickets
func (fs *Freshservice) GetTicket(ctx context.Context, tid int64, include ...string) (*Ticket, error) {
	url := fs.Endpoint("/tickets/%d", tid)
	if len(include) > 0 {
		s := strings.Join(include, ",")
		url += "?include=" + s
	}
	result := &ticketResult{}
	err := fs.DoGet(ctx, url, result)
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
func (fs *Freshservice) FilterTickets(ctx context.Context, fto *FilterTicketsOption) ([]*Ticket, bool, error) {
	url := fs.Endpoint("/tickets/filter")
	result := &ticketResult{}
	next, err := fs.DoList(ctx, url, fto, result)
	return result.Tickets, next, err
}

func (fs *Freshservice) IterFilterTickets(ctx context.Context, fto *FilterTicketsOption, itf func(*Ticket) error) error {
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
		tickets, next, err := fs.FilterTickets(ctx, fto)
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
func (fs *Freshservice) ListTickets(ctx context.Context, lto *ListTicketsOption) ([]*Ticket, bool, error) {
	url := fs.Endpoint("/tickets")
	result := &ticketResult{}
	next, err := fs.DoList(ctx, url, lto, result)
	return result.Tickets, next, err
}

func (fs *Freshservice) IterTickets(ctx context.Context, lto *ListTicketsOption, itf func(*Ticket) error) error {
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
		tickets, next, err := fs.ListTickets(ctx, lto)
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
func (fs *Freshservice) UpdateTicket(ctx context.Context, tid int64, ticket *TicketUpdate) (*Ticket, error) {
	url := fs.Endpoint("/tickets/%d", tid)
	result := &ticketResult{}
	if err := fs.DoPut(ctx, url, ticket, result); err != nil {
		return nil, err
	}
	return result.Ticket, nil
}

func (fs *Freshservice) DeleteTicket(ctx context.Context, tid int64) error {
	url := fs.Endpoint("/tickets/%d", tid)
	return fs.DoDelete(ctx, url)
}

func (fs *Freshservice) DeleteTicketAttachment(ctx context.Context, tid, aid int64) error {
	url := fs.Endpoint("/tickets/%d/attachments/%d", tid, aid)
	return fs.DoDelete(ctx, url)
}

// Restore a Ticket
// The API mentioned previously. If you deleted some tickets and regret doing so now, this API will help you restore them.
func (fs *Freshservice) Restore(ctx context.Context, tid int64) error {
	url := fs.Endpoint("/tickets/%d/restore", tid)
	return fs.DoPut(ctx, url, nil, nil)
}

// Create a Child Ticket
// This API lets you to create a new child ticket on an existing ticket
// Note:
// 1. Association of child tickets is not available in the Sprout plan.
// 2. Association of child tickets to a service request is not possible.
// 3. Association of child tickets to a deleted or a spammed ticket is not allowed.
// 4. Nesting of a child ticket under another child ticket is not supported.
func (fs *Freshservice) CreateChildTicket(ctx context.Context, tid int64, ticket *Ticket) (*Ticket, error) {
	url := fs.Endpoint("/tickets/%d/create_child_ticket", tid)
	result := &ticketResult{}
	if err := fs.DoPost(ctx, url, ticket, result); err != nil {
		return nil, err
	}
	return result.Ticket, nil
}

func (fs *Freshservice) ListTicketFields(ctx context.Context) ([]*TicketField, error) {
	url := fs.Endpoint("/ticket_form_fields")
	result := &ticketFieldResult{}
	err := fs.DoGet(ctx, url, result)
	return result.TicketFields, err
}

func (fs *Freshservice) GetTicketActivities(ctx context.Context, tid int64) ([]*TicketActivity, error) {
	url := fs.Endpoint("/tickets/%d/activities", tid)
	result := &ticketActivitiesResult{}
	err := fs.DoGet(ctx, url, result)
	return result.TicketActivities, err
}

// ---------------------------------------------------
// Conversation

func (fs *Freshservice) CreateReply(ctx context.Context, tid int64, reply *Reply) (*Conversation, error) {
	url := fs.Endpoint("/tickets/%d/reply", tid)
	result := &conversationResult{}
	if err := fs.DoPost(ctx, url, reply, result); err != nil {
		return nil, err
	}
	return result.Conversation, nil
}

func (fs *Freshservice) CreateNote(ctx context.Context, tid int64, note *Note) (*Conversation, error) {
	url := fs.Endpoint("/tickets/%d/notes", tid)
	result := &conversationResult{}
	if err := fs.DoPost(ctx, url, note, result); err != nil {
		return nil, err
	}
	return result.Conversation, nil
}

// Update a Conversation
// Only public & private notes can be edited.
func (fs *Freshservice) UpdateConversation(ctx context.Context, cid int64, note *Note) (*Conversation, error) {
	url := fs.Endpoint("/conversations/%d", cid)
	result := &conversationResult{}
	if err := fs.DoPut(ctx, url, note, result); err != nil {
		return nil, err
	}
	return result.Conversation, nil
}

func (fs *Freshservice) DeleteConversation(ctx context.Context, cid int64) error {
	url := fs.Endpoint("/conversations/%d", cid)
	return fs.DoDelete(ctx, url)
}

func (fs *Freshservice) DeleteConversationAttachment(ctx context.Context, cid, aid int64) error {
	url := fs.Endpoint("/conversations/%d/attachments/%d", cid, aid)
	return fs.DoDelete(ctx, url)
}

func (fs *Freshservice) ListTicketConversations(ctx context.Context, tid int64, lco *ListConversationsOption) ([]*Conversation, bool, error) {
	url := fs.Endpoint("/tickets/%d/conversations", tid)
	result := &conversationsResult{}
	next, err := fs.DoList(ctx, url, lco, result)
	return result.Conversations, next, err
}

func (fs *Freshservice) IterTicketConversations(ctx context.Context, tid int64, lco *ListConversationsOption, icf func(*Conversation) error) error {
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
		conversations, next, err := fs.ListTicketConversations(ctx, tid, lco)
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
