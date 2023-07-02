package freshdesk

import (
	"net/url"
)

// ---------------------------------------------------
// Agent

type ListAgentsOption struct {
	Email   string
	Mobile  string
	Phone   string
	State   string // [fulltime/occasional]
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
	q.SetString("state", lao.State)
	q.SetInt("page", lao.Page)
	q.SetInt("per_page", lao.PerPage)
	return q
}

func (fd *Freshdesk) GetAgent(aid int64) (*Agent, error) {
	url := fd.endpoint("/agents/%d", aid)
	agent := &Agent{}
	err := fd.doGet(url, agent)
	return agent, err
}

func (fd *Freshdesk) ListAgents(lao *ListAgentsOption) ([]*Agent, bool, error) {
	url := fd.endpoint("/agents")
	agents := []*Agent{}
	next, err := fd.doList(url, lao, &agents)
	return agents, next, err
}

func (fd *Freshdesk) IterAgents(lao *ListAgentsOption, iaf func(*Agent) error) error {
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
		agents, next, err := fd.ListAgents(lao)
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

func (fd *Freshdesk) CreateAgent(agent *AgentRequest) (*Agent, error) {
	url := fd.endpoint("/agents")
	result := &Agent{}
	err := fd.doPost(url, agent, result)
	return result, err
}

func (fd *Freshdesk) UpdateAgent(aid int64, agent *AgentRequest) (*Agent, error) {
	url := fd.endpoint("/agents/%d", aid)
	result := &Agent{}
	err := fd.doPut(url, agent, result)
	return result, err
}

func (fd *Freshdesk) DeleteAgent(aid int64) error {
	url := fd.endpoint("/agents/%d", aid)
	return fd.doDelete(url)
}

func (fd *Freshdesk) GetCurrentAgent() (*Agent, error) {
	url := fd.endpoint("/agents/me")
	agent := &Agent{}
	err := fd.doGet(url, agent)
	return agent, err
}

func (fd *Freshdesk) SearchAgents(keyword string) ([]*Agent, error) {
	url := fd.endpoint("/agents/autocomplete?term=%s", url.QueryEscape(keyword))
	agents := []*Agent{}
	err := fd.doGet(url, &agents)
	return agents, err
}
