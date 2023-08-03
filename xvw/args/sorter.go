package args

import (
	"github.com/askasoft/pango/ars"
	"github.com/askasoft/pango/str"
)

type Sorter struct {
	Col string `form:"c"`
	Dir string `form:"d"`
}

func (s *Sorter) String() string {
	return s.Col + " " + s.Dir
}

func (s *Sorter) Normalize(columns ...string) {
	if len(columns) > 0 {
		if !ars.ContainsString(columns, s.Col) {
			s.Col = columns[0]
		}
	}

	s.Dir = str.ToLower(s.Dir)
	if s.Dir != "asc" && s.Dir != "desc" {
		s.Dir = "asc"
	}
}
