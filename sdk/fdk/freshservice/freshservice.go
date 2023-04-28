package freshservice

import (
	"fmt"
	"strings"

	"github.com/askasoft/pango/sdk/fdk"
	"github.com/askasoft/pango/str"
)

type Freshservice fdk.FDK

func (fs *Freshservice) doGet(url string, result any) error {
	return (*fdk.FDK)(fs).DoGet(url, result)
}

func (fs *Freshservice) doList(url string, lo ListOption, ap any) (bool, error) {
	return (*fdk.FDK)(fs).DoList(url, lo, ap)
}

func (fs *Freshservice) doPost(url string, source, result any) error {
	return (*fdk.FDK)(fs).DoPost(url, source, result)
}

func (fs *Freshservice) doPut(url string, source, result any) error {
	return (*fdk.FDK)(fs).DoPut(url, source, result)
}

func (fs *Freshservice) doDelete(url string) error {
	return (*fdk.FDK)(fs).DoDelete(url)
}

func (fs *Freshservice) Download(url string) ([]byte, error) {
	return (*fdk.FDK)(fs).DoDownload(url)
}

func (fs *Freshservice) SaveFile(url string, filename string) error {
	return (*fdk.FDK)(fs).DoSave(url, filename)
}

// GetHelpdeskAttachmentURL return a permlink for helpdesk attachment/avator URL
func (fs *Freshservice) GetHelpdeskAttachmentURL(aid int64) string {
	return fmt.Sprintf("%s/helpdesk/attachments/%d", fs.Domain, aid)
}

// ---------------------------------------------------
// Agent

func (fs *Freshservice) GetAgentRole(id int64) (*AgentRole, error) {
	url := fmt.Sprintf("%s/api/v2/roles/%d", fs.Domain, id)
	result := &agentRoleResult{}
	err := fs.doGet(url, result)
	return result.Role, err
}

func (fs *Freshservice) ListAgentRoles() ([]*AgentRole, error) {
	url := fmt.Sprintf("%s/api/v2/roles", fs.Domain)
	result := &agentRoleResult{}
	err := fs.doGet(url, result)
	return result.Roles, err
}

func (fs *Freshservice) CreateAgentGroup(ag *AgentGroup) (*AgentGroup, error) {
	url := fmt.Sprintf("%s/api/v2/groups", fs.Domain)
	result := &agentGroupResult{}
	err := fs.doPost(url, ag, result)
	return result.Group, err
}

func (fs *Freshservice) GetAgentGroup(id int64) (*AgentGroup, error) {
	url := fmt.Sprintf("%s/api/v2/groups/%d", fs.Domain, id)
	result := &agentGroupResult{}
	err := fs.doGet(url, result)
	return result.Group, err
}

func (fs *Freshservice) ListAgentGroups() ([]*AgentGroup, error) {
	url := fmt.Sprintf("%s/api/v2/groups", fs.Domain)
	result := &agentGroupResult{}
	err := fs.doGet(url, result)
	return result.Groups, err
}

func (fs *Freshservice) UpdateAgentGroup(id int64, ag *AgentGroup) (*AgentGroup, error) {
	url := fmt.Sprintf("%s/api/v2/groups/%d", fs.Domain, id)
	result := &agentGroupResult{}
	err := fs.doPut(url, ag, result)
	return result.Group, err
}

func (fs *Freshservice) DeleteAgentGroup(id int64) error {
	url := fmt.Sprintf("%s/api/v2/groups/%d", fs.Domain, id)
	return fs.doDelete(url)
}

func (fs *Freshservice) CreateAgent(agent *Agent) (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents", fs.Domain)
	result := &agentResult{}
	err := fs.doPost(url, agent, result)
	return result.Agent, err
}

func (fs *Freshservice) GetAgent(id int64) (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents/%d", fs.Domain, id)
	result := &agentResult{}
	err := fs.doGet(url, result)
	return result.Agent, err
}

// List Agents
// Use Agent attributes to filter your list.
// Note:
// 1. Filtered results cannot be sorted. By default it is sorted by created_at in descending order.
// 2. The query must be URL encoded (see example).
// 3. Query can be framed using the name of the agent fields, which can be obtained from the Supported Agent Fields Section.
// 4. Query string must be enclosed between a pair of double quotes and can have up to 512 characters.
// 5. Logical operators AND, OR along with parenthesis( ) can be used to group conditions.
// 6. Relational operators greater than or equal to :> and less than or equal to :< can be used along with date fields and numeric fields.
// 7. Input for date field should be in UTC Format.
// 8. The number of objects returned per page is 30.
// 9. To scroll through the pages add the page parameter to the url. The page number starts with 1 and should not exceed 40.
// 10. To filter for fields with no values assigned, use the null keyword.
// 11. The "~" query operator can be used for "starts with" text searches. "Starts with" search is supported for one or more of the following attributes: first_name, last_name, name, email, mobile_phone_number, work_phone_number. The query format is https://domain.freshservice.com/api/v2/agents?query="~[attribute_1|attribute_2]:'somestring'". The query needs to be URL encoded. This would return a list of users for whom attribute_1 OR attribute_2 starts with "somestring". Refer to examples 11, 12, and 13.
// 12. Please note that any update made to an agent either in Freshservice application or through API may take a few minutes to get indexed, after which the updated results will be available through API.
func (fs *Freshservice) ListAgents(lro *ListAgentsOption) ([]*Agent, bool, error) {
	url := fmt.Sprintf("%s/api/v2/agents", fs.Domain)
	result := &agentResult{}
	next, err := fs.doList(url, lro, result)
	return result.Agents, next, err
}

// Update an Agent
// This operation allows you to modify the profile of a particular agent.
// Note:
// can_see_all_tickets_from_associated_departments will automatically be set to false unless it is explicitly set to true in the payload, irrespective of the previous value of the field.
func (fs *Freshservice) UpdateAgent(id int64, agent *Agent) (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents/%d", fs.Domain, id)
	result := &agentResult{}
	err := fs.doPut(url, agent, result)
	return result.Agent, err
}

// Deactivate a Agent
// This operation allows you to deactivate a agent.
func (fs *Freshservice) DeactivateAgent(id int64) error {
	url := fmt.Sprintf("%s/api/v2/agents/%d", fs.Domain, id)
	return fs.doDelete(url)
}

// Forget a Agent
// This operation allows you to permanently delete a agent and the tickets that they requested.
func (fs *Freshservice) ForgetAgent(id int64) error {
	url := fmt.Sprintf("%s/api/v2/agents/%d/forget", fs.Domain, id)
	return fs.doDelete(url)
}

// Reactivate a Agent
// This operation allows you to reactivate a particular deactivated agent.
func (fs *Freshservice) ReactivateAgent(id int64) (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents/%d/reactivate", fs.Domain, id)
	result := &agentResult{}
	err := fs.doPut(url, nil, result)
	return result.Agent, err
}

// Convert a particular agent into a requester.
func (fs *Freshservice) ConvertAgentToRequester(id int64) (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/agents/%d/convert_to_requester", fs.Domain, id)
	result := &agentResult{}
	err := fs.doPut(url, nil, result)
	return result.Agent, err
}

func (fs *Freshservice) GetAgentFields() ([]*AgentField, error) {
	url := fmt.Sprintf("%s/api/v2/agent_fields", fs.Domain)
	result := &agentFieldResult{}
	err := fs.doGet(url, result)
	return result.AgentFields, err
}

// ---------------------------------------------------
// Requester

func (fs *Freshservice) CreateRequesterGroup(rg *RequesterGroup) (*RequesterGroup, error) {
	url := fmt.Sprintf("%s/api/v2/requester_groups", fs.Domain)
	result := &requesterGroupResult{}
	err := fs.doPost(url, rg, result)
	return result.RequesterGroup, err
}

func (fs *Freshservice) GetRequesterGroup(id int64) (*RequesterGroup, error) {
	url := fmt.Sprintf("%s/api/v2/requester_groups/%d", fs.Domain, id)
	result := &requesterGroupResult{}
	err := fs.doGet(url, result)
	return result.RequesterGroup, err
}

func (fs *Freshservice) ListRequesterGroups() ([]*RequesterGroup, error) {
	url := fmt.Sprintf("%s/api/v2/requester_groups", fs.Domain)
	result := &requesterGroupResult{}
	err := fs.doGet(url, result)
	return result.RequesterGroups, err
}

// Note:
// Only groups of type “manual” can be updated through this API.
func (fs *Freshservice) UpdateRequesterGroup(id int64, rg *RequesterGroup) (*RequesterGroup, error) {
	url := fmt.Sprintf("%s/api/v2/requester_groups/%d", fs.Domain, id)
	result := &requesterGroupResult{}
	err := fs.doPut(url, rg, result)
	return result.RequesterGroup, err
}

// Delete Requester Group
// Note:
// 1. Deleting a Requester Group will only disband the requester group and will not delete its members.
// 2. Deleted requester groups cannot be restored.
func (fs *Freshservice) DeleteRequesterGroup(id int64) error {
	url := fmt.Sprintf("%s/api/v2/requester_groups/%d", fs.Domain, id)
	return fs.doDelete(url)
}

// Add Requester to Requester Group
// Note:
// 1.Requesters can be added only to manual requester groups.
// 2.Requester can be added one at a time.
func (fs *Freshservice) AddRequesterToRequesterGroup(rgid, rid int64) error {
	url := fmt.Sprintf("%s/api/v2/requester_groups/%d/members/%d", fs.Domain, rgid, rid)
	return fs.doPost(url, nil, nil)
}

// Delete Requester from Requester Group
// Note:
// 1.Requesters can be removed only from manual requester groups.
// 2.Requester can be removed one at a time.
func (fs *Freshservice) DeleteRequesterFromRequesterGroup(rgid, rid int64) error {
	url := fmt.Sprintf("%s/api/v2/requester_groups/%d/members/%d", fs.Domain, rgid, rid)
	return fs.doDelete(url)
}

func (fs *Freshservice) ListRequesterGroupMembers(rgid int64, lmo *ListMembersOption) ([]*Requester, bool, error) {
	url := fmt.Sprintf("%s/api/v2/requester_groups/%d/members", fs.Domain, rgid)
	result := &requesterResult{}
	next, err := fs.doList(url, lmo, result)
	return result.Requesters, next, err
}

func (fs *Freshservice) CreateRequester(requester *Requester) (*Requester, error) {
	url := fmt.Sprintf("%s/api/v2/requesters", fs.Domain)
	result := &requesterResult{}
	err := fs.doPost(url, requester, result)
	return result.Requester, err
}

func (fs *Freshservice) GetRequester(id int64) (*Requester, error) {
	url := fmt.Sprintf("%s/api/v2/requesters/%d", fs.Domain, id)
	result := &requesterResult{}
	err := fs.doGet(url, result)
	return result.Requester, err
}

// List Requesters
// Use Requester attributes to filter your list.
// Note:
// 1. Filtered results cannot be sorted. By default it is sorted by created_at in descending order.
// 2. Adding "include_agents=true" to the query string will include agents in the response. The default response includes only requesters and not agents. Only users who also have the "Manage Agents" permission will be able to use this modifier.
// 3. The query must be URL encoded (see example).
// 4. Query can be framed using the name of the requester fields, which can be obtained from the Supported Requester Fields Section.
// 5. Query string must be enclosed between a pair of double quotes and can have up to 512 characters.
// 6. Logical operators AND, OR along with parenthesis( ) can be used to group conditions.
// 7. Relational operators greater than or equal to :> and less than or equal to :< can be used along with date fields and numeric fields.
// 8. Input for date field should be in UTC Format.
// 9. The number of objects returned per page is 30.
// 10. To scroll through the pages add the page parameter to the url. The page number starts with 1 and should not exceed 40.
// 11. To filter for fields with no values assigned, use the null keyword.
// 12. The "~" query operator can be used for "starts with" text searches. "Starts with" search is supported for one or more of the following attributes: first_name, last_name, name, primary_email, mobile_phone_number, work_phone_number. The query format is https://domain.freshservice.com/api/v2/requesters?query="~[attribute_1|attribute_2]:'somestring'". The query needs to be URL encoded. This would return a list of users for whom attribute_1 OR attribute_2 starts with "somestring". Refer to examples 13, 14, and 15.
// 13. Please note that any update made to requester either in Freshservice application or through API may take a few minutes to get indexed, after which the updated results will be available through API.
// == Custom Fields Supported	Type
// Single line text	string
// Number	integer
// Dropdown	string
// Date	date
// Phone number	string
func (fs *Freshservice) ListRequesters(lro *ListRequestersOption) ([]*Requester, bool, error) {
	url := fmt.Sprintf("%s/api/v2/requesters", fs.Domain)
	result := &requesterResult{}
	next, err := fs.doList(url, lro, result)
	return result.Requesters, next, err
}

func (fs *Freshservice) GetRequesterFields() ([]*RequesterField, error) {
	url := fmt.Sprintf("%s/api/v2/requester_fields", fs.Domain)
	result := &requesterFieldResult{}
	err := fs.doGet(url, result)
	return result.RequesterFields, err
}

// Update a Requester
// This operation allows you to modify the profile of a particular requester.
// Note:
// can_see_all_tickets_from_associated_departments will automatically be set to false unless it is explicitly set to true in the payload, irrespective of the previous value of the field.
func (fs *Freshservice) UpdateRequester(id int64, requester *Requester) (*Requester, error) {
	url := fmt.Sprintf("%s/api/v2/requesters/%d", fs.Domain, id)
	result := &requesterResult{}
	err := fs.doPut(url, requester, result)
	return result.Requester, err
}

// Deactivate a Requester
// This operation allows you to deactivate a requester.
func (fs *Freshservice) DeactivateRequester(id int64) error {
	url := fmt.Sprintf("%s/api/v2/requesters/%d", fs.Domain, id)
	return fs.doDelete(url)
}

// Forget a Requester
// This operation allows you to permanently delete a requester and the tickets that they requested.
func (fs *Freshservice) ForgetRequester(id int64) error {
	url := fmt.Sprintf("%s/api/v2/requesters/%d/forget", fs.Domain, id)
	return fs.doDelete(url)
}

// Convert a requester to an occasional agent with SD Agent role and no group memberships.
func (fs *Freshservice) ConvertRequesterToAgent(id int64) (*Agent, error) {
	url := fmt.Sprintf("%s/api/v2/requesters/%d/convert_to_agent", fs.Domain, id)
	result := &agentResult{}
	err := fs.doPut(url, nil, result)
	return result.Agent, err
}

// Merge secondary requesters into a primary requester.
func (fs *Freshservice) MergeRequesters(id int64, ids ...int64) (*Requester, error) {
	url := fmt.Sprintf("%s/api/v2/requesters/%d/merge?secondary_requesters=%s", fs.Domain, id, str.JoinInt64s(ids, ","))
	result := &requesterResult{}
	err := fs.doPut(url, nil, result)
	return result.Requester, err
}

// Reactivate a Requester
// This operation allows you to reactivate a particular deactivated requester.
func (fs *Freshservice) ReactivateRequester(id int64) (*Requester, error) {
	url := fmt.Sprintf("%s/api/v2/requesters/%d/reactivate", fs.Domain, id)
	result := &requesterResult{}
	err := fs.doPut(url, nil, result)
	return result.Requester, err
}

// ---------------------------------------------------
// Ticket

func (fs *Freshservice) CreateTicket(ticket *Ticket) (*Ticket, error) {
	url := fmt.Sprintf("%s/api/v2/tickets", fs.Domain)
	result := &ticketResult{}
	err := fs.doPost(url, ticket, result)
	return result.Ticket, err
}

// GetTicket Get a Ticket
// include: conversations, requester, requested_for, stats, problem, assets, change, related_tickets
func (fs *Freshservice) GetTicket(tid int64, include ...string) (*Ticket, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fs.Domain, tid)
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
	url := fmt.Sprintf("%s/api/v2/tickets/filter", fs.Domain)
	result := &ticketResult{}
	next, err := fs.doList(url, fto, result)
	return result.Tickets, next, err
}

// List of Tickets
// Use filters to view only specific tickets (those which match the criteria that you choose). By default only tickets that have not been deleted or marked as spam will be returned, unless you use the 'deleted' filter.
// Note:
// 1. By default only tickets that have been created within the past 30 days will be returned. For older tickets, use the updated_since filter.
// 2. Use 'include' to embed additional details in the response. Each include will consume an additional 2 credits. For example if you embed the stats information you will be charged a total of 3 API credits (1 credit for the API call, and 2 credits for the additional stats embedding).
// 3. By default, only tickets from the primary workspace will be returned for accounts with the 'Workspaces' feature enabled. For tickets from other workspaces, use the workspace_id filter.
func (fs *Freshservice) ListTickets(lto *ListTicketsOption) ([]*Ticket, bool, error) {
	url := fmt.Sprintf("%s/api/v2/tickets", fs.Domain)
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
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fs.Domain, tid)
	result := &ticketResult{}
	err := fs.doPut(url, ticket, result)
	return result.Ticket, err
}

func (fs *Freshservice) DeleteTicket(tid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d", fs.Domain, tid)
	return fs.doDelete(url)
}

func (fs *Freshservice) DeleteTicketAttachment(tid, aid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/attachments/%d", fs.Domain, tid, aid)
	return fs.doDelete(url)
}

// Restore a Ticket
// The API mentioned previously. If you deleted some tickets and regret doing so now, this API will help you restore them.
func (fs *Freshservice) Restore(tid int64) error {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/restore", fs.Domain, tid)
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
	url := fmt.Sprintf("%s/api/v2/tickets/%d/create_child_ticket", fs.Domain, tid)
	result := &ticketResult{}
	err := fs.doPost(url, ticket, result)
	return result.Ticket, err
}

func (fs *Freshservice) GetTicketFields() ([]*TicketField, error) {
	url := fmt.Sprintf("%s/api/v2/ticket_form_fields", fs.Domain)
	result := &ticketFieldResult{}
	err := fs.doGet(url, result)
	return result.TicketFields, err
}

func (fs *Freshservice) GetTicketActivities(tid int64) ([]*TicketActivity, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/activities", fs.Domain, tid)
	result := &ticketActivityResult{}
	err := fs.doGet(url, result)
	return result.TicketActivitys, err
}

// ---------------------------------------------------
// Conversation

func (fs *Freshservice) CreateReply(tid int64, reply *Reply) (*Conversation, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/reply", fs.Domain, tid)
	result := &conversationResult{}
	err := fs.doPost(url, reply, result)
	return result.Conversation, err
}

func (fs *Freshservice) CreateNote(tid int64, note *Note) (*Conversation, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/note", fs.Domain, tid)
	result := &conversationResult{}
	err := fs.doPost(url, note, result)
	return result.Conversation, err
}

// Update a Conversation
// Only public & private notes can be edited.
func (fs *Freshservice) UpdateConversation(cid int64, conversation *Conversation) (*Conversation, error) {
	url := fmt.Sprintf("%s/api/v2/conversations/%d", fs.Domain, cid)
	result := &conversationResult{}
	err := fs.doPut(url, conversation, result)
	return result.Conversation, err
}

func (fs *Freshservice) DeleteConversation(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/conversations/%d", fs.Domain, cid)
	return fs.doDelete(url)
}

func (fs *Freshservice) DeleteConversationAttachment(cid, aid int64) error {
	url := fmt.Sprintf("%s/api/v2/conversations/%d/attachments/%d", fs.Domain, cid, aid)
	return fs.doDelete(url)
}

func (fs *Freshservice) ListTicketConversations(tid int64, lco *ListConversationsOption) ([]*Conversation, bool, error) {
	url := fmt.Sprintf("%s/api/v2/tickets/%d/conversations", fs.Domain, tid)
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

// ---------------------------------------------------
// Solution

func (fs *Freshservice) CreateCategory(category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories", fs.Domain)
	result := &categoryResult{}
	err := fs.doPost(url, category, result)
	return result.Category, err
}

func (fs *Freshservice) UpdateCategory(cid int64, category *Category) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d", fs.Domain, cid)
	result := &categoryResult{}
	err := fs.doPut(url, category, result)
	return result.Category, err
}

func (fs *Freshservice) GetCategory(cid int64) (*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d", fs.Domain, cid)
	result := &categoryResult{}
	err := fs.doGet(url, result)
	return result.Category, err
}

func (fs *Freshservice) ListCategories() ([]*Category, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/categories", fs.Domain)
	result := &categoryResult{}
	err := fs.doGet(url, result)
	return result.Categories, err
}

func (fs *Freshservice) DeleteCategory(cid int64) error {
	url := fmt.Sprintf("%s/api/v2/solutions/categories/%d", fs.Domain, cid)
	return fs.doDelete(url)
}

func (fs *Freshservice) CreateFolder(folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders", fs.Domain)
	result := &folderResult{}
	err := fs.doPost(url, folder, result)
	return result.Foler, err
}

func (fs *Freshservice) UpdateFolder(fid int64, folder *Folder) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fs.Domain, fid)
	result := &folderResult{}
	err := fs.doPut(url, folder, result)
	return result.Foler, err
}

func (fs *Freshservice) GetFolder(fid int64) (*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fs.Domain, fid)
	result := &folderResult{}
	err := fs.doGet(url, result)
	return result.Foler, err
}

func (fs *Freshservice) ListCategoryFolders(cid int64) ([]*Folder, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders?category_id=%d", fs.Domain, cid)
	result := &folderResult{}
	err := fs.doGet(url, result)
	return result.Folders, err
}

func (fs *Freshservice) DeleteFolder(fid int64) error {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fs.Domain, fid)
	return fs.doDelete(url)
}

func (fs *Freshservice) CreateArticle(fid int64, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/folders/%d", fs.Domain, fid)
	result := &articleResult{}
	err := fs.doPost(url, article, result)
	return result.Article, err
}

func (fs *Freshservice) SendArticleToApproval(aid int64) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d/send_for_approval", fs.Domain, aid)
	result := &articleResult{}
	err := fs.doPut(url, nil, result)
	return result.Article, err
}

func (fs *Freshservice) UpdateArticle(aid int64, article *Article) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d", fs.Domain, aid)
	result := &articleResult{}
	err := fs.doPut(url, article, result)
	return result.Article, err
}

func (fs *Freshservice) GetArticle(aid int64) (*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d", fs.Domain, aid)
	result := &articleResult{}
	err := fs.doGet(url, result)
	return result.Article, err
}

func (fs *Freshservice) ListFolderArticles(fid int64) ([]*Article, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles?folder_id=%d", fs.Domain, fid)
	result := &articleResult{}
	err := fs.doGet(url, result)
	return result.Articles, err
}

func (fs *Freshservice) DeleteArticle(aid int64) error {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/%d", fs.Domain, aid)
	return fs.doDelete(url)
}

func (fs *Freshservice) SearchArticles(sao *SearchArticlesOption) ([]*Article, bool, error) {
	url := fmt.Sprintf("%s/api/v2/solutions/articles/search", fs.Domain)
	result := &articleResult{}
	next, err := fs.doList(url, sao, result)
	return result.Articles, next, err
}
