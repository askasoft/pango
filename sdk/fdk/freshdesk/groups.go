package freshdesk

// ---------------------------------------------------
// Group

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

func (fd *Freshdesk) ListGroups() ([]*Group, error) {
	url := fd.endpoint("/groups")
	groups := []*Group{}
	_, err := fd.doList(url, nil, &groups)
	return groups, err
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
