package freshdesk

// ---------------------------------------------------
// Role

func (fd *Freshdesk) GetRole(rid int64) (*Role, error) {
	url := fd.endpoint("/roles/%d", rid)
	role := &Role{}
	err := fd.doGet(url, role)
	return role, err
}

func (fd *Freshdesk) ListRoles() ([]*Role, error) {
	url := fd.endpoint("/roles")
	roles := []*Role{}
	_, err := fd.doList(url, nil, &roles)
	return roles, err
}
