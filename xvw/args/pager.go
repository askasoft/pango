package args

import (
	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/num"
)

type Pager struct {
	Page  int `json:"p,omitempty" form:"p"` // page no
	Limit int `json:"l,omitempty" form:"l"` // page limit
	Count int `json:"-"`                    // page item count
	Total int `json:"-"`                    // total count
}

func (p *Pager) Start() int {
	return (p.Page - 1) * p.Limit
}

func (p *Pager) Begin() int {
	return p.Start() + 1
}

func (p *Pager) End() int {
	return p.Start() + p.Count
}

// Pages return the pages according to the total, limit.
func (p *Pager) Pages() int {
	if p.Total <= 0 {
		return 0
	}

	if p.Limit <= 0 {
		return 1
	}

	return ((p.Total - 1) / p.Limit) + 1
}

func (p *Pager) Normalize(limits ...int) {
	p.Limit = num.AbsInt(p.Limit)
	if len(limits) > 0 {
		if !asg.Contains(limits, p.Limit) {
			p.Limit = limits[0]
		}
	}

	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Page > p.Pages() {
		p.Page = p.Pages()
	}
}
