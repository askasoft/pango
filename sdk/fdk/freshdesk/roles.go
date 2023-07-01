package freshdesk

import "fmt"

// ---------------------------------------------------
// Role

func (fd *Freshdesk) GetRole(rid int64) (*Role, error) {
	url := fmt.Sprintf("%s/api/v2/roles/%d", fd.Domain, rid)
	role := &Role{}
	err := fd.doGet(url, role)
	return role, err
}

func (fd *Freshdesk) ListRoles() ([]*Role, error) {
	url := fmt.Sprintf("%s/api/v2/roles", fd.Domain)
	roles := []*Role{}
	_, err := fd.doList(url, nil, &roles)
	return roles, err
}
