package freshservice

type Approval struct {
	ApprovalType int `json:"approval_type,omitempty"`

	ApproverIDs []int64 `json:"approver_ids,omitempty"`
}

func (a *Approval) String() string {
	return toString(a)
}
