package freshservice

import "context"

// ---------------------------------------------------
// Agent Role

type ListAgentRolesOption = PageOption

func (fs *Freshservice) GetAgentRole(ctx context.Context, id int64) (*AgentRole, error) {
	url := fs.Endpoint("/roles/%d", id)
	result := &agentRoleResult{}
	err := fs.DoGet(ctx, url, result)
	return result.Role, err
}

func (fs *Freshservice) ListAgentRoles(ctx context.Context, laro *ListAgentRolesOption) ([]*AgentRole, bool, error) {
	url := fs.Endpoint("/roles")
	result := &agentRolesResult{}
	next, err := fs.DoList(ctx, url, laro, result)
	return result.Roles, next, err
}

func (fs *Freshservice) IterAgentRoles(ctx context.Context, laro *ListAgentRolesOption, iarf func(*AgentRole) error) error {
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
		ars, next, err := fs.ListAgentRoles(ctx, laro)
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
