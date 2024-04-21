package freshservice

import "github.com/askasoft/pango/num"

type ApprovalType int
type ApprovalStatus int

const (
	ApprovalTypeEveryone  ApprovalType = 1
	ApprovalTypeAnyone    ApprovalType = 2
	ApprovalTypeMajority  ApprovalType = 3
	ApprovalTypeResponder ApprovalType = 4

	ApprovalStatusRequested ApprovalStatus = 0
	ApprovalStatusApproved  ApprovalStatus = 1
	ApprovalStatusRejected  ApprovalStatus = 2
	ApprovalStatusCanceled  ApprovalStatus = 3
)

func (as ApprovalStatus) String() string {
	switch as {
	case ApprovalStatusRequested:
		return "requested"
	case ApprovalStatusApproved:
		return "approved"
	case ApprovalStatusRejected:
		return "rejected"
	case ApprovalStatusCanceled:
		return "canceled"
	default:
		return num.Itoa(int(as))
	}
}

type ApprovalSetting struct {
	ApprovalType int     `json:"approval_type,omitempty"`
	ApproverIDs  []int64 `json:"approver_ids,omitempty"`
}

func (a *ApprovalSetting) String() string {
	return toString(a)
}

type ApprovalInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (a *ApprovalInfo) String() string {
	return toString(a)
}

type Delegatee struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (d *Delegatee) String() string {
	return toString(d)
}

type Approval struct {
	ID             int64         `json:"id,omitempty"`
	Parent         string        `json:"parent,omitempty"`
	ParentID       int64         `json:"parent_id,omitempty"`
	ApproverID     int64         `json:"approver_id,omitempty"`
	ApproverName   string        `json:"approver_name,omitempty"`
	ApprovalType   ApprovalType  `json:"approval_type,omitempty"`
	Level          int           `json:"level,omitempty"`
	UserID         int64         `json:"user_id,omitempty"`
	UserName       string        `json:"user_name,omitempty"`
	MemberID       int64         `json:"member_id,omitempty"`
	MemberName     string        `json:"member_name,omitempty"`
	ApprovalStatus *ApprovalInfo `json:"approval_status,omitempty"`
	Delegatee      *Delegatee    `json:"delegatee,omitempty"`
	LatestRemark   string        `json:"latest_remark,omitempty"`
	CreatedAt      *Time         `json:"created_at,omitempty"`
	UpdatedAt      *Time         `json:"updated_at,omitempty"`
}

func (a *Approval) String() string {
	return toString(a)
}

type approvalsResult struct {
	Approvals []*Approval `json:"approvals,omitempty"`
}
