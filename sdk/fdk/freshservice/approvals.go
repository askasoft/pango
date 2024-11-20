package freshservice

import "context"

type ListApprovalsOption struct {
	Parent      string
	ParentID    int64
	Status      string
	ApproverID  int64
	Level       int
	DelegateeID int64
	Page        int
	PerPage     int
}

func (lao *ListApprovalsOption) IsNil() bool {
	return lao == nil
}

func (lao *ListApprovalsOption) Values() Values {
	q := Values{}
	q.SetString("parent", lao.Parent)
	q.SetInt64("parent_id", lao.ParentID)
	q.SetString("status", lao.Status)
	q.SetInt64("approver_id", lao.ApproverID)
	q.SetInt("level", lao.Level)
	q.SetInt("delegatee_id", int(lao.DelegateeID))
	q.SetInt("page", lao.Page)
	q.SetInt("per_page", lao.PerPage)
	return q
}

func (fs *Freshservice) ListApprovals(ctx context.Context, lao *ListApprovalsOption) ([]*Approval, bool, error) {
	url := fs.endpoint("/approvals")
	result := &approvalsResult{}
	next, err := fs.doList(ctx, url, lao, result)
	return result.Approvals, next, err
}

func (fs *Freshservice) IterApprovals(ctx context.Context, lao *ListApprovalsOption, iaf func(*Approval) error) error {
	if lao == nil {
		lao = &ListApprovalsOption{}
	}
	if lao.Page < 1 {
		lao.Page = 1
	}
	if lao.PerPage < 1 {
		lao.PerPage = 100
	}

	for {
		agents, next, err := fs.ListApprovals(ctx, lao)
		if err != nil {
			return err
		}
		for _, c := range agents {
			if err = iaf(c); err != nil {
				return err
			}
		}
		if !next {
			break
		}
		lao.Page++
	}
	return nil
}
