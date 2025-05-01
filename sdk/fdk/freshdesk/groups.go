package freshdesk

import "context"

// ---------------------------------------------------
// Group

type ListGroupsOption = PageOption

func (fd *Freshdesk) GetGroup(ctx context.Context, gid int64) (*Group, error) {
	url := fd.endpoint("/groups/%d", gid)
	group := &Group{}
	err := fd.doGet(ctx, url, group)
	return group, err
}

func (fd *Freshdesk) CreateGroup(ctx context.Context, group *GroupCreate) (*Group, error) {
	url := fd.endpoint("/groups")
	result := &Group{}
	if err := fd.doPost(ctx, url, group, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) ListGroups(ctx context.Context, lgo *ListGroupsOption) ([]*Group, bool, error) {
	url := fd.endpoint("/groups")
	groups := []*Group{}
	next, err := fd.doList(ctx, url, lgo, &groups)
	return groups, next, err
}

func (fd *Freshdesk) IterGroups(ctx context.Context, lgo *ListGroupsOption, igf func(*Group) error) error {
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
		groups, next, err := fd.ListGroups(ctx, lgo)
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

func (fd *Freshdesk) UpdateGroup(ctx context.Context, gid int64, group *GroupUpdate) (*Group, error) {
	url := fd.endpoint("/groups/%d", gid)
	result := &Group{}
	if err := fd.doPut(ctx, url, group, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) DeleteGroup(ctx context.Context, gid int64) error {
	url := fd.endpoint("/groups/%d", gid)
	return fd.doDelete(ctx, url)
}
