package freshdesk

import "fmt"

// ---------------------------------------------------
// Group

func (fd *Freshdesk) GetGroup(gid int64) (*Group, error) {
	url := fmt.Sprintf("%s/api/v2/groups/%d", fd.Domain, gid)
	group := &Group{}
	err := fd.doGet(url, group)
	return group, err
}

func (fd *Freshdesk) CreateGroup(group *Group) (*Group, error) {
	url := fmt.Sprintf("%s/api/v2/groups", fd.Domain)
	result := &Group{}
	err := fd.doPost(url, group, result)
	return result, err
}

func (fd *Freshdesk) ListGroups() ([]*Group, error) {
	url := fmt.Sprintf("%s/api/v2/groups", fd.Domain)
	groups := []*Group{}
	_, err := fd.doList(url, nil, &groups)
	return groups, err
}

func (fd *Freshdesk) UpdateGroup(gid int64, group *Group) (*Group, error) {
	url := fmt.Sprintf("%s/api/v2/groups/%d", fd.Domain, gid)
	result := &Group{}
	err := fd.doPut(url, group, result)
	return result, err
}

func (fd *Freshdesk) DeleteGroup(gid int64) error {
	url := fmt.Sprintf("%s/api/v2/groups/%d", fd.Domain, gid)
	return fd.doDelete(url)
}
