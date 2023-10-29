package freshservice

// ---------------------------------------------------
// Agent Group

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
	result := &agentGroupsResult{}
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
