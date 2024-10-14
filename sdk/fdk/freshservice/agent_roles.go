package freshservice

// ---------------------------------------------------
// Agent Role

type ListAgentRolesOption = PageOption

func (fs *Freshservice) GetAgentRole(id int64) (*AgentRole, error) {
	url := fs.endpoint("/roles/%d", id)
	result := &agentRoleResult{}
	err := fs.doGet(url, result)
	return result.Role, err
}

func (fs *Freshservice) ListAgentRoles(laro *ListAgentRolesOption) ([]*AgentRole, bool, error) {
	url := fs.endpoint("/roles")
	result := &agentRolesResult{}
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
