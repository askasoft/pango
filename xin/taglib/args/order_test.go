package args

import "testing"

func TestOrderNormalize(t *testing.T) {
	cs := []struct {
		o string
		n []string
		w string
	}{
		{"id", []string{"id", "name"}, "id"},
		{"-id", []string{"id", "name"}, "-id"},
		{"id,tag", []string{"id", "name"}, "id"},
		{"-id,-tag", []string{"id", "name"}, "-id"},
		{"id,tag", []string{"tag", "name"}, "tag"},
		{"-id,-tag", []string{"tag", "name"}, "-tag"},
		{"id,name", []string{"id", "name"}, "id,name"},
		{"-id,-name", []string{"id", "name"}, "-id,-name"},
		{"id,name", nil, "id,name"},
		{"-id,-name", nil, "-id,-name"},
		{"id,name", []string{"tag", "namex"}, ""},
		{"-id,-name", []string{"tag", "namex"}, ""},
	}

	for i, c := range cs {
		ods := Orders{c.o}
		ods.Normalize(c.n...)
		if c.w != ods.Order {
			t.Errorf("#%d Normalize(%q, %v) = %q, want %q", i, c.o, c.n, ods, c.w)
		}
	}
}
