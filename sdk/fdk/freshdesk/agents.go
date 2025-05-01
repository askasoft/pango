package freshdesk

import (
	"context"
	"net/url"
)

// ---------------------------------------------------
// Agent

type AgentState string

const (
	AgentStateFulltime   AgentState = "fulltime"
	AgentStateOccasional AgentState = "occasional"
)

type ListAgentsOption struct {
	Email   string
	Mobile  string
	Phone   string
	State   AgentState // [fulltime/occasional]
	Page    int
	PerPage int
}

func (lao *ListAgentsOption) IsNil() bool {
	return lao == nil
}

func (lao *ListAgentsOption) Values() Values {
	q := Values{}
	q.SetString("email", lao.Email)
	q.SetString("mobile", lao.Mobile)
	q.SetString("phone", lao.Phone)
	q.SetString("state", (string)(lao.State))
	q.SetInt("page", lao.Page)
	q.SetInt("per_page", lao.PerPage)
	return q
}

func (fd *Freshdesk) GetAgent(ctx context.Context, aid int64) (*Agent, error) {
	url := fd.endpoint("/agents/%d", aid)
	agent := &Agent{}
	err := fd.doGet(ctx, url, agent)
	return agent, err
}

func (fd *Freshdesk) ListAgents(ctx context.Context, lao *ListAgentsOption) ([]*Agent, bool, error) {
	url := fd.endpoint("/agents")
	agents := []*Agent{}
	next, err := fd.doList(ctx, url, lao, &agents)
	return agents, next, err
}

func (fd *Freshdesk) IterAgents(ctx context.Context, lao *ListAgentsOption, iaf func(*Agent) error) error {
	if lao == nil {
		lao = &ListAgentsOption{}
	}
	if lao.Page < 1 {
		lao.Page = 1
	}
	if lao.PerPage < 1 {
		lao.PerPage = 100
	}

	for {
		agents, next, err := fd.ListAgents(ctx, lao)
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

func (fd *Freshdesk) CreateAgent(ctx context.Context, agent *AgentCreate) (*Agent, error) {
	url := fd.endpoint("/agents")
	result := &Agent{}
	if err := fd.doPost(ctx, url, agent, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) UpdateAgent(ctx context.Context, aid int64, agent *AgentUpdate) (*Agent, error) {
	url := fd.endpoint("/agents/%d", aid)
	result := &Agent{}
	if err := fd.doPut(ctx, url, agent, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fd *Freshdesk) DeleteAgent(ctx context.Context, aid int64) error {
	url := fd.endpoint("/agents/%d", aid)
	return fd.doDelete(ctx, url)
}

func (fd *Freshdesk) GetCurrentAgent(ctx context.Context) (*Agent, error) {
	url := fd.endpoint("/agents/me")
	agent := &Agent{}
	err := fd.doGet(ctx, url, agent)
	return agent, err
}

func (fd *Freshdesk) SearchAgents(ctx context.Context, keyword string) ([]*Agent, error) {
	url := fd.endpoint("/agents/autocomplete?term=%s", url.QueryEscape(keyword))
	agents := []*Agent{}
	err := fd.doGet(ctx, url, &agents)
	return agents, err
}
