package freshservice

type RequesterGroupType string

const (
	RequesterGroupTypeManual    RequesterGroupType = "manual"
	RequesterGroupTypeRuleBased RequesterGroupType = "rule_based"
)

type RequesterGroup struct {
	ID int64 `json:"id,omitempty"`

	// Name of the requester group.
	Name string `json:"name,omitempty"`

	// Description of the requester group.
	Description string `json:"description,omitempty"`

	// Method of requester addition. “manual” if individual requesters can be chosen manually, and “rule_based” if members are automatically added based on rules.
	Type RequesterGroupType `json:"type,omitempty"`
}

func (rg *RequesterGroup) String() string {
	return toString(rg)
}

type requesterGroupResult struct {
	RequesterGroup  *RequesterGroup   `json:"requester_group,omitempty"`
	RequesterGroups []*RequesterGroup `json:"requester_groups,omitempty"`
}
