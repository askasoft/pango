package col

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestHashMapInterface(t *testing.T) {
	var m Map = NewHashMap()
	if m == nil {
		t.Error("HashMap is not a Map")
	}
}

func TestHashMapAsHashMap(t *testing.T) {
	m := map[any]any{
		"a": 1,
		"b": 2,
	}
	hm := AsHashMap(m)

	if hm.Len() != len(m) {
		t.Fatal("hm.Len() != len(m)")
	}
	hm.Clear()

	if hm.Len() != 0 {
		t.Fatal("hm.Len() != 0")
	}
	if hm.Len() != len(m) {
		t.Fatal("hm.Len() != len(m)")
	}
}

func TestHashMapSet(t *testing.T) {
	m := NewHashMap()
	m.Set(5, "e")
	m.Set(6, "f")
	m.Set(7, "g")
	m.Set(3, "c")
	m.Set(4, "d")
	m.Set(1, "x")
	m.Set(2, "b")
	m.Set(1, "a") //overwrite

	if av := m.Len(); av != 7 {
		t.Errorf("Got %v expected %v", av, 7)
	}
	if av, ev := m.Keys(), []any{1, 2, 3, 4, 5, 6, 7}; !testHashMapSame(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := m.Values(), []any{"a", "b", "c", "d", "e", "f", "g"}; !testHashMapSame(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}

	// key,ev,expectedFound
	tests1 := [][]any{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, nil, false},
	}

	for _, test := range tests1 {
		// retrievals
		av, actualFound := m.Get(test[0])
		if av != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", av, test[1])
		}
	}
}

func TestHashMapRemove(t *testing.T) {
	m := NewHashMap()
	m.Set(5, "e")
	m.Set(6, "f")
	m.Set(7, "g")
	m.Set(3, "c")
	m.Set(4, "d")
	m.Set(1, "x")
	m.Set(2, "b")
	m.Set(1, "a") //overwrite

	m.Remove(5)
	m.Removes(6, 7, 8)
	m.Remove(5)

	if av, ev := m.Keys(), []any{1, 2, 3, 4}; !testHashMapSame(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}

	if av, ev := m.Values(), []any{"a", "b", "c", "d"}; !testHashMapSame(av, ev) {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av := m.Len(); av != 4 {
		t.Errorf("Got %v expected %v", av, 4)
	}

	tests2 := [][]any{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, nil, false},
		{6, nil, false},
		{7, nil, false},
		{8, nil, false},
	}

	for _, test := range tests2 {
		av, actualFound := m.Get(test[0])
		if av != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", av, test[1])
		}
	}

	m.Remove(1)
	m.Removes(4, 2, 3, 2)
	m.Remove(2)

	if av, ev := fmt.Sprintf("%s", m.Keys()), "[]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%s", m.Values()), "[]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av := m.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
	if av := m.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
}

func TestHashMapJSON(t *testing.T) {
	cs := []struct {
		s string
		w *HashMap
	}{
		{`{}`, NewHashMap()},
		{`{"a":1,"b":2,"c":3}`, NewHashMap([]P{{"a", 1.0}, {"b", 2.0}, {"c", 3.0}}...)},
	}

	for i, c := range cs {
		a := NewHashMap()
		err := json.Unmarshal(([]byte)(c.s), &a)
		if err != nil {
			t.Errorf("[%d] json.Unmarshal(%q) = %v", i, c.s, err)
		}
		if !testHashMapSame(a.Values(), c.w.Values()) {
			t.Errorf("[%d] json.Unmarshal(%q) = %v, want %v", i, c.s, a.Values(), c.w.Values())
		}
		if !testHashMapSame(a.Keys(), c.w.Keys()) {
			t.Errorf("[%d] json.Unmarshal(%q) = %v, want %v", i, c.s, a.Keys(), c.w.Keys())
		}

		bs, err := json.Marshal(a)
		if err != nil {
			t.Errorf("[%d] json.Marshal(%v) = %v", i, a, err)
		}

		m := make(map[string]any)
		err = json.Unmarshal(bs, &m)
		if err != nil {
			t.Errorf("[%d] json.Unmarshal(%q) = %v", i, bs, err)
		}

		ks := make([]any, 0)
		vs := make([]any, 0)
		for k, v := range m {
			ks = append(ks, k)
			vs = append(vs, v)
		}

		if !testHashMapSame(ks, a.Keys()) || !testHashMapSame(vs, a.Values()) {
			t.Errorf("[%d] json.Marshal(%q) = %q, want %q", i, a.String(), string(bs), c.s)
		}
	}
}

func testHashMapSame(a []any, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for _, av := range a {
		found := false
		for _, bv := range b {
			if av == bv {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
