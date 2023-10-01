package freshdesk

// ---------------------------------------------------
// Group

type ListGroupsOption = PageOption

func (fd *Freshdesk) GetGroup(gid int64) (*Group, error) {
	url := fd.endpoint("/groups/%d", gid)
	group := &Group{}
	err := fd.doGet(url, group)
	return group, err
}

func (fd *Freshdesk) CreateGroup(group *Group) (*Group, error) {
	url := fd.endpoint("/groups")
	result := &Group{}
	err := fd.doPost(url, group, result)
	return result, err
}

func (fd *Freshdesk) ListGroups(lgo *ListGroupsOption) ([]*Group, bool, error) {
	url := fd.endpoint("/groups")
	groups := []*Group{}
	next, err := fd.doList(url, lgo, &groups)
	return groups, next, err
}

func (fd *Freshdesk) IterGroups(lgo *ListGroupsOption, igf func(*Group) error) error {
	if lgo == nil {
		lgo = &ListGroupsOption{}
	}
	if lgo.Page < 1 {
		lgo.Page = 1
	}
	if lgo.PerPage < 1 {
		lgo.PerPage = 100
	}

	for {
		groups, next, err := fd.ListGroups(lgo)
		if err != nil {
			return err
		}
		for _, g := range groups {
			if err = igf(g); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lgo.Page++
	}
	return nil
}

func (fd *Freshdesk) UpdateGroup(gid int64, group *Group) (*Group, error) {
	url := fd.endpoint("/groups/%d", gid)
	result := &Group{}
	err := fd.doPut(url, group, result)
	return result, err
}

func (fd *Freshdesk) DeleteGroup(gid int64) error {
	url := fd.endpoint("/groups/%d", gid)
	return fd.doDelete(url)
}
