package freshservice

// ---------------------------------------------------
// Agent

type AgentState string

const (
	AgentStateFulltime   AgentState = "fulltime"
	AgentStateOccasional AgentState = "occasional"
)

type ListAgentRolesOption = PageOption
type ListAgentGroupsOption = PageOption

type ListAgentsOption struct {
	FirstName          string
	LastName           string
	Name               string // Concatenation of first_name and last_name with single space in-between fields.
	JobTitle           string
	Email              string
	MobilePhoneNumber  string
	WorkPhoneNumber    string
	Active             bool
	State              AgentState // fulltime or occasional
	DepartmentID       int64
	ReportingManagerID int64
	TimeZone           string
	Language           string
	LocationID         int64
	CreatedAt          Date // Date (YYYY-MM-DD) when the requester is created.
	UpdatedAt          Date // Date (YYYY-MM-DD) when the requester is updated.
	Page               int
	PerPage            int
}

func (lao *ListAgentsOption) IsNil() bool {
	return lao == nil
}

func (lao *ListAgentsOption) Values() Values {
	q := Values{}
	q.SetString("first_name", lao.FirstName)
	q.SetString("last_name", lao.LastName)
	q.SetString("name", lao.Name)
	q.SetString("job_title", lao.JobTitle)
	q.SetString("email", lao.Email)
	q.SetString("mobile_phone_number", lao.MobilePhoneNumber)
	q.SetString("work_phone_number", lao.WorkPhoneNumber)
	q.SetBool("active", lao.Active)
	q.SetString("state", (string)(lao.State))
	q.SetInt64("department_id", lao.DepartmentID)
	q.SetInt64("reporting_manager_id", lao.ReportingManagerID)
	q.SetString("time_zone", lao.TimeZone)
	q.SetString("language", lao.Language)
	q.SetInt64("location_id", lao.LocationID)
	q.SetDate("created_at", lao.CreatedAt)
	q.SetDate("updated_at", lao.UpdatedAt)
	q.SetInt("page", lao.Page)
	q.SetInt("per_page", lao.PerPage)
	return q
}

func (fs *Freshservice) GetAgentRole(id int64) (*AgentRole, error) {
	url := fs.endpoint("/roles/%d", id)
	result := &agentRoleResult{}
	err := fs.doGet(url, result)
	return result.Role, err
}

func (fs *Freshservice) ListAgentRoles(laro *ListAgentRolesOption) ([]*AgentRole, bool, error) {
	url := fs.endpoint("/roles")
	result := &agentRoleResult{}
	next, err := fs.doList(url, laro, result)
	return result.Roles, next, err
}

func (fs *Freshservice) IterAgentRoles(laro *ListAgentRolesOption, iarf func(*AgentRole) error) error {
	if laro == nil {
		laro = &ListAgentRolesOption{}
	}
	if laro.Page < 1 {
		laro.Page = 1
	}
	if laro.PerPage < 1 {
		laro.PerPage = 100
	}

	for {
		ars, next, err := fs.ListAgentRoles(laro)
		if err != nil {
			return err
		}
		for _, ar := range ars {
			if err = iarf(ar); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		laro.Page++
	}
	return nil
}

func (fs *Freshservice) CreateAgentGroup(ag *AgentGroup) (*AgentGroup, error) {
	url := fs.endpoint("/groups")
	result := &agentGroupResult{}
	err := fs.doPost(url, ag, result)
	return result.Group, err
}

func (fs *Freshservice) GetAgentGroup(id int64) (*AgentGroup, error) {
	url := fs.endpoint("/groups/%d", id)
	result := &agentGroupResult{}
	err := fs.doGet(url, result)
	return result.Group, err
}

func (fs *Freshservice) ListAgentGroups(lago *ListAgentGroupsOption) ([]*AgentGroup, bool, error) {
	url := fs.endpoint("/groups")
	result := &agentGroupResult{}
	next, err := fs.doList(url, lago, result)
	return result.Groups, next, err
}

func (fs *Freshservice) IterAgentGroups(lago *ListAgentGroupsOption, iagf func(*AgentGroup) error) error {
	if lago == nil {
		lago = &ListAgentRolesOption{}
	}
	if lago.Page < 1 {
		lago.Page = 1
	}
	if lago.PerPage < 1 {
		lago.PerPage = 100
	}

	for {
		ags, next, err := fs.ListAgentGroups(lago)
		if err != nil {
			return err
		}
		for _, ag := range ags {
			if err = iagf(ag); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lago.Page++
	}
	return nil
}

func (fs *Freshservice) UpdateAgentGroup(id int64, ag *AgentGroup) (*AgentGroup, error) {
	url := fs.endpoint("/groups/%d", id)
	result := &agentGroupResult{}
	err := fs.doPut(url, ag, result)
	return result.Group, err
}

func (fs *Freshservice) DeleteAgentGroup(id int64) error {
	url := fs.endpoint("/groups/%d", id)
	return fs.doDelete(url)
}

func (fs *Freshservice) CreateAgent(agent *Agent) (*Agent, error) {
	url := fs.endpoint("/agents")
	result := &agentResult{}
	err := fs.doPost(url, agent, result)
	return result.Agent, err
}

func (fs *Freshservice) GetAgent(id int64) (*Agent, error) {
	url := fs.endpoint("/agents/%d", id)
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
func (fs *Freshservice) ListAgents(lao *ListAgentsOption) ([]*Agent, bool, error) {
	url := fs.endpoint("/agents")
	result := &agentResult{}
	next, err := fs.doList(url, lao, result)
	return result.Agents, next, err
}

func (fs *Freshservice) IterAgents(lao *ListAgentsOption, iaf func(*Agent) error) error {
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
		agents, next, err := fs.ListAgents(lao)
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

// Update an Agent
// This operation allows you to modify the profile of a particular agent.
// Note:
// can_see_all_tickets_from_associated_departments will automatically be set to false unless it is explicitly set to true in the payload, irrespective of the previous value of the field.
func (fs *Freshservice) UpdateAgent(id int64, agent *Agent) (*Agent, error) {
	url := fs.endpoint("/agents/%d", id)
	result := &agentResult{}
	err := fs.doPut(url, agent, result)
	return result.Agent, err
}

// Deactivate a Agent
// This operation allows you to deactivate a agent.
func (fs *Freshservice) DeactivateAgent(id int64) error {
	url := fs.endpoint("/agents/%d", id)
	return fs.doDelete(url)
}

// Forget a Agent
// This operation allows you to permanently delete a agent and the tickets that they requested.
func (fs *Freshservice) ForgetAgent(id int64) error {
	url := fs.endpoint("/agents/%d/forget", id)
	return fs.doDelete(url)
}

// Reactivate a Agent
// This operation allows you to reactivate a particular deactivated agent.
func (fs *Freshservice) ReactivateAgent(id int64) (*Agent, error) {
	url := fs.endpoint("/agents/%d/reactivate", id)
	result := &agentResult{}
	err := fs.doPut(url, nil, result)
	return result.Agent, err
}

// Convert a particular agent into a requester.
func (fs *Freshservice) ConvertAgentToRequester(id int64) (*Agent, error) {
	url := fs.endpoint("/agents/%d/convert_to_requester", id)
	result := &agentResult{}
	err := fs.doPut(url, nil, result)
	return result.Agent, err
}

func (fs *Freshservice) GetAgentFields() ([]*AgentField, error) {
	url := fs.endpoint("/agent_fields")
	result := &agentFieldResult{}
	err := fs.doGet(url, result)
	return result.AgentFields, err
}
