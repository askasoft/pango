package args

import (
	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/str"
)

type Sorter struct {
	Col string `json:"c,omitempty" form:"c"`
	Dir string `json:"d,omitempty" form:"d,lower"`
}

// String return Col + " " + Dir
func (s *Sorter) String() string {
	return s.Col + " " + s.Dir
}

func (s *Sorter) IsAsc() bool {
	return s.Dir == "asc"
}

func (s *Sorter) IsDesc() bool {
	return s.Dir == "desc"
}

// Order return the sql order by string (quote + Col + quote + " " + Dir)
// the default quote string is ""
func (s *Sorter) Order(quote ...string) string {
	q := ""
	if len(quote) > 0 {
		q = quote[0]
	}
	return q + s.Col + q + " " + s.Dir
}

func (s *Sorter) Normalize(columns ...string) {
	if len(columns) > 0 {
		if !asg.Contains(columns, s.Col) {
			s.Col = columns[0]
		}
	}

	s.Dir = str.ToLower(s.Dir)
	if s.Dir != "asc" && s.Dir != "desc" {
		s.Dir = "asc"
	}
}
