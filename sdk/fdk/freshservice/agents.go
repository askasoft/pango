package freshservice

import "context"

// ---------------------------------------------------
// Agent

type AgentActive string
type AgentState string

const (
	AgentStateFulltime   AgentState = "fulltime"
	AgentStateOccasional AgentState = "occasional"

	AgentActiveTrue  AgentActive = "true"
	AgentActiveFalse AgentActive = "false"
	AgentActiveALL   AgentActive = ""
)

type ListAgentsOption struct {
	Email             string
	MobilePhoneNumber string
	WorkPhoneNumber   string
	Active            AgentActive
	State             AgentState // fulltime or occasional
	Page              int
	PerPage           int
}

func (lao *ListAgentsOption) IsNil() bool {
	return lao == nil
}

func (lao *ListAgentsOption) Values() Values {
	q := Values{}
	q.SetString("email", lao.Email)
	q.SetString("mobile_phone_number", lao.MobilePhoneNumber)
	q.SetString("work_phone_number", lao.WorkPhoneNumber)
	q.SetString("active", string(lao.Active))
	q.SetString("state", string(lao.State))
	q.SetInt("page", lao.Page)
	q.SetInt("per_page", lao.PerPage)
	return q
}

type FilterAgentsOption = FilterOption

func (fs *Freshservice) CreateAgent(ctx context.Context, agent *AgentCreate) (*Agent, error) {
	url := fs.Endpoint("/agents")
	result := &agentResult{}
	if err := fs.DoPost(ctx, url, agent, result); err != nil {
		return nil, err
	}
	return result.Agent, nil
}

func (fs *Freshservice) GetAgent(ctx context.Context, id int64) (*Agent, error) {
	url := fs.Endpoint("/agents/%d", id)
	result := &agentResult{}
	err := fs.DoGet(ctx, url, result)
	return result.Agent, err
}

func (fs *Freshservice) ListAgents(ctx context.Context, lao *ListAgentsOption) ([]*Agent, bool, error) {
	url := fs.Endpoint("/agents")
	result := &agentsResult{}
	next, err := fs.DoList(ctx, url, lao, result)
	return result.Agents, next, err
}

func (fs *Freshservice) IterAgents(ctx context.Context, lao *ListAgentsOption, iaf func(*Agent) error) error {
	if lao == nil {
		lao = &ListAgentsOption{}
	}
	if lao.Page < 1 {
		lao.Page = 1
	}
	if lao.PerPage < 1 {
		lao.PerPage = 100
	}

	for {
		agents, next, err := fs.ListAgents(ctx, lao)
		if err != nil {
			return err
		}
		for _, c := range agents {
			if err = iaf(c); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lao.Page++
	}
	return nil
}

// FilterAgents Use Agent attributes to filter your list.
// See: https://api.freshservice.com/#filter_agents
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
// Supported Agent Fields:
// Field	Type	Description
// first_name	string	First name of the agent.
// last_name	string	Last name of the agent.
// name	string	Concatenation of first_name and last_name with single space in-between fields.
// job_title	string	Title of the agent.
// email	string	Email address of the agent.
// work_phone_number	string	Work phone of the agent.
// mobile_phone_number	string	Mobile phone of the agent.
// department_id	integer	ID of the department(s) assigned to the agent.
// reporting_manager_id	integer	ID of the reporting manager.
// time_zone	string	ID of the department.
// language	string	Language code(Eg. en, ja-JP).
// location_id	integer	ID of the location.
// created_at	date	Date (YYYY-MM-DD) when the agent is created.
// updated_at	date	Date (YYYY-MM-DD) when the agent is updated.
func (fs *Freshservice) FilterAgents(ctx context.Context, fao *FilterAgentsOption) ([]*Agent, bool, error) {
	url := fs.Endpoint("/agents")
	result := &agentsResult{}
	next, err := fs.DoList(ctx, url, fao, result)
	return result.Agents, next, err
}

func (fs *Freshservice) IterFilterAgents(ctx context.Context, fao *FilterAgentsOption, iaf func(*Agent) error) error {
	if fao == nil {
		fao = &FilterAgentsOption{}
	}
	if fao.Page < 1 {
		fao.Page = 1
	}
	if fao.PerPage < 1 {
		fao.PerPage = 100
	}

	for {
		agents, next, err := fs.FilterAgents(ctx, fao)
		if err != nil {
			return err
		}
		for _, c := range agents {
			if err = iaf(c); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		fao.Page++
	}
	return nil
}

// Update an Agent
// This operation allows you to modify the profile of a particular agent.
// Note:
// can_see_all_tickets_from_associated_departments will automatically be set to false unless it is explicitly set to true in the payload, irrespective of the previous value of the field.
func (fs *Freshservice) UpdateAgent(ctx context.Context, id int64, agent *AgentUpdate) (*Agent, error) {
	url := fs.Endpoint("/agents/%d", id)
	result := &agentResult{}
	if err := fs.DoPut(ctx, url, agent, result); err != nil {
		return nil, err
	}
	return result.Agent, nil
}

// Deactivate a Agent
// This operation allows you to deactivate a agent.
func (fs *Freshservice) DeactivateAgent(ctx context.Context, id int64) error {
	url := fs.Endpoint("/agents/%d", id)
	return fs.DoDelete(ctx, url)
}

// Forget a Agent
// This operation allows you to permanently delete a agent and the tickets that they requested.
func (fs *Freshservice) ForgetAgent(ctx context.Context, id int64) error {
	url := fs.Endpoint("/agents/%d/forget", id)
	return fs.DoDelete(ctx, url)
}

// Reactivate a Agent
// This operation allows you to reactivate a particular deactivated agent.
func (fs *Freshservice) ReactivateAgent(ctx context.Context, id int64) (*Agent, error) {
	url := fs.Endpoint("/agents/%d/reactivate", id)
	result := &agentResult{}
	if err := fs.DoPut(ctx, url, nil, result); err != nil {
		return nil, err
	}
	return result.Agent, nil
}

// Convert a particular agent into a requester.
func (fs *Freshservice) ConvertAgentToRequester(ctx context.Context, id int64) (*Agent, error) {
	url := fs.Endpoint("/agents/%d/convert_to_requester", id)
	result := &agentResult{}
	if err := fs.DoPut(ctx, url, nil, result); err != nil {
		return nil, err
	}
	return result.Agent, nil
}

func (fs *Freshservice) GetAgentFields(ctx context.Context) ([]*AgentField, error) {
	url := fs.Endpoint("/agent_fields")
	result := &agentFieldsResult{}
	err := fs.DoGet(ctx, url, result)
	return result.AgentFields, err
}
