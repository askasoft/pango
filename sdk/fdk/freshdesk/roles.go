package freshdesk

import "context"

// ---------------------------------------------------
// Role

type ListRolesOption = PageOption

func (fd *Freshdesk) GetRole(ctx context.Context, rid int64) (*Role, error) {
	url := fd.Endpoint("/roles/%d", rid)
	role := &Role{}
	err := fd.DoGet(ctx, url, role)
	return role, err
}

func (fd *Freshdesk) ListRoles(ctx context.Context, lro *ListRolesOption) ([]*Role, bool, error) {
	url := fd.Endpoint("/roles")
	roles := []*Role{}
	next, err := fd.DoList(ctx, url, lro, &roles)
	return roles, next, err
}

func (fd *Freshdesk) IterRoles(ctx context.Context, lro *ListRolesOption, irf func(*Role) error) error {
	if lro == nil {
		lro = &ListRolesOption{}
	}
	if lro.Page < 1 {
		lro.Page = 1
	}
	if lro.PerPage < 1 {
		lro.PerPage = 100
	}

	for {
		roles, next, err := fd.ListRoles(ctx, lro)
		if err != nil {
			return err
		}
		for _, g := range roles {
			if err = irf(g); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lro.Page++
	}
	return nil
}
