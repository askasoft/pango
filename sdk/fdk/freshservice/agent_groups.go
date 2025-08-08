package freshservice

import "context"

// ---------------------------------------------------
// Agent Group

type ListAgentGroupsOption = PageOption

func (fs *Freshservice) CreateAgentGroup(ctx context.Context, ag *AgentGroupCreate) (*AgentGroup, error) {
	url := fs.Endpoint("/groups")
	result := &agentGroupResult{}
	if err := fs.DoPost(ctx, url, ag, result); err != nil {
		return nil, err
	}
	return result.Group, nil
}

func (fs *Freshservice) GetAgentGroup(ctx context.Context, id int64) (*AgentGroup, error) {
	url := fs.Endpoint("/groups/%d", id)
	result := &agentGroupResult{}
	err := fs.DoGet(ctx, url, result)
	return result.Group, err
}

func (fs *Freshservice) ListAgentGroups(ctx context.Context, lago *ListAgentGroupsOption) ([]*AgentGroup, bool, error) {
	url := fs.Endpoint("/groups")
	result := &agentGroupsResult{}
	next, err := fs.DoList(ctx, url, lago, result)
	return result.Groups, next, err
}

func (fs *Freshservice) IterAgentGroups(ctx context.Context, lago *ListAgentGroupsOption, iagf func(*AgentGroup) error) error {
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
		ags, next, err := fs.ListAgentGroups(ctx, lago)
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

func (fs *Freshservice) UpdateAgentGroup(ctx context.Context, id int64, ag *AgentGroupUpdate) (*AgentGroup, error) {
	url := fs.Endpoint("/groups/%d", id)
	result := &agentGroupResult{}
	if err := fs.DoPut(ctx, url, ag, result); err != nil {
		return nil, err
	}
	return result.Group, nil
}

func (fs *Freshservice) DeleteAgentGroup(ctx context.Context, id int64) error {
	url := fs.Endpoint("/groups/%d", id)
	return fs.DoDelete(ctx, url)
}
