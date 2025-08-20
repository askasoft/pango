package taglib

import (
	"reflect"
	"testing"
)

func TestAttrsAdd(t *testing.T) {
	cs := []struct {
		k, v string
		a    Attrs
		w    Attrs
	}{
		{"a", "", Attrs{}, Attrs{"a": ""}},
		{"a", "v", Attrs{}, Attrs{"a": "v"}},
		{"a", "", Attrs{"a": ""}, Attrs{"a": ""}},
		{"a", "", Attrs{"a": "v"}, Attrs{"a": "v"}},
		{"a", "v", Attrs{"a": ""}, Attrs{"a": "v"}},
		{"a", "v", Attrs{"a": "u"}, Attrs{"a": "u v"}},
	}

	for i, c := range cs {
		c.a.Add(c.k, c.v)
		if !reflect.DeepEqual(c.a, c.w) {
			t.Errorf("[%d] Add(%q, %q)\n GOT: %v\nWANT: %v\n", i, c.k, c.v, c.a, c.w)
		}
	}
}
