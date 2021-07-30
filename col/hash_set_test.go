package col

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"github.com/pandafw/pango/str"
)

func TestHashSetInterface(t *testing.T) {
	hs := NewHashSet()

	var c Collection = hs
	if c == nil {
		t.Error("HashSet is not a Collection")
	}

	var s Set = hs
	if s == nil {
		t.Error("HashSet is not a Set")
	}
}

func TestHashSetLazyInit(t *testing.T) {
	{
		hs := &HashSet{}
		if hs.Len() != 0 {
			t.Error("hs.Len() != 0")
		}
		if !hs.IsEmpty() {
			t.Error("hs.IsEmpty() = true")
		}
		if len(hs.Values()) != 0 {
			t.Error("len(hs.Values()) != 0")
		}
		if hs.Contains(1) {
			t.Error("hs.Contains(1) = true, want false")
		}
		if hs.Len() != 0 {
			t.Error("hs.Len() != 0")
		}
		if !hs.ContainsSet(&HashSet{}) {
			t.Error("&HashSet{}.ContainsSet(&HashSet{}) = true, want false")
		}
	}
	{
		hs := &HashSet{}
		hs.Add(1)
		if hs.Len() != 1 {
			t.Errorf("hs.Len() = %v, want 1", hs.Len())
		}
	}
	{
		hs := &HashSet{}
		hs.AddSet(NewHashSet(1))
		if hs.Len() != 1 {
			t.Errorf("hs.Len() = %v, want 1", hs.Len())
		}
	}
	{
		hs := &HashSet{}
		hs.Delete(1)
		if hs.Len() != 0 {
			t.Error("hs.Len() != 0")
		}
	}
	{
		hs := &HashSet{}
		hs.Each(func(v interface{}) {})
	}
	{
		hs := &HashSet{}
		as := hs.Difference(&HashSet{})
		if as.Len() != 0 {
			t.Errorf("hs.Difference(&HashSet{}) == %v, want %v", as, hs)
		}
	}
	{
		hs := &HashSet{}
		as := hs.Intersection(&HashSet{})
		if as.Len() != 0 {
			t.Errorf("hs.Intersection(&HashSet{}) == %v, want %v", as, hs)
		}
	}
}

func TestHashSetSimple(t *testing.T) {
	s := NewHashSet()

	s.Add(5)
	if s.Len() != 1 {
		t.Errorf("Length should be 1")
	}

	if !s.Contains(5) {
		t.Errorf("Membership test failed")
	}

	s.Delete(5)

	if s.Len() != 0 {
		t.Errorf("Length should be 0")
	}

	if s.Contains(5) {
		t.Errorf("The set should be empty")
	}
}

func TestHashSetNewHashSet(t *testing.T) {
	set := NewHashSet(2, 1)

	if av := set.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	if av := set.Contains(1); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := set.Contains(2); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := set.Contains(3); av != false {
		t.Errorf("Got %v expected %v", av, true)
	}
}

func TestHashSetAdd(t *testing.T) {
	set := NewHashSet()
	set.Add()
	set.Add(1)
	set.Add(2)
	set.Add(2, 3)
	set.Add()
	if av := set.IsEmpty(); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := set.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
}

func TestHashSetContains(t *testing.T) {
	set := NewHashSet()
	set.Add(3, 1, 2)
	set.Add(2, 3)
	set.Add()
	if av := set.Contains(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := set.Contains(1); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := set.Contains(1, 2, 3); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := set.Contains(1, 2, 3, 4); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
}

func TestHashSetDelete(t *testing.T) {
	set := NewHashSet()
	set.Add(3, 1, 2)
	set.Delete()
	if av := set.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	set.Delete(1)
	if av := set.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	set.Delete(3)
	set.Delete(3)
	set.Delete()
	set.Delete(2)
	if av := set.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestHashSetContainsSet(t *testing.T) {
	s1 := NewHashSet()
	s2 := NewHashSet()

	if !s1.ContainsSet(s1) {
		t.Errorf("set should be a subset of itself")
	}

	if !s1.ContainsSet(s2) {
		t.Errorf("empty set should contains another empty set")
	}

	s1.Add(1)
	if !s1.ContainsSet(s2) {
		t.Errorf("set should contains another empty set")
	}

	s2.Add(1)
	if !s1.ContainsSet(s2) {
		t.Errorf("set should contains another same set")
	}

	s1.Add(2)
	if !s1.ContainsSet(s2) {
		t.Errorf("set should contains another small set")
	}

	s2.Add(3)
	if !s1.ContainsSet(s2) {
		t.Errorf("set should not contains another different set")
	}
}

func TestHashSetDifference(t *testing.T) {
	// Difference
	s1 := NewHashSet(1, 2, 3, 4, 5, 6)
	s2 := NewHashSet(4, 5, 6)
	s3 := s1.Difference(s2)

	if s3.Len() != 3 {
		t.Errorf("Length should be 3")
	}

	if !(s3.Contains(1) && s3.Contains(2) && s3.Contains(3)) {
		t.Errorf("Set should only contain 1, 2, 3")
	}
}

func TestHashSetIntersection(t *testing.T) {
	s1 := NewHashSet(1, 2, 3, 4, 5, 6)
	s2 := NewHashSet(4, 5, 6)

	// Intersection
	s3 := s1.Intersection(s2)
	if s3.Len() != 3 {
		t.Errorf("Length should be 3 after intersection")
	}

	if !(s3.Contains(4) && s3.Contains(5) && s3.Contains(6)) {
		t.Errorf("Set should contain 4, 5, 6")
	}
}

func TestHashSetAddSet(t *testing.T) {
	// AddSet
	s1 := NewHashSet(4, 5, 6)
	s2 := NewHashSet(7, 8, 9)
	s1.AddSet(s2)

	if s1.Len() != 6 {
		t.Errorf("Length should be 6 after union")
	}

	for i := 4; i <= 9; i++ {
		if !(s1.Contains(i)) {
			t.Errorf("Set should contains %d", i)
		}
	}
}

func sortSetJSON(s string) string {
	s = str.TrimRight(str.TrimLeft(s, "["), "]")
	a := str.FieldsAny(s, ",")
	sort.Strings(a)
	return str.Join(a, ",")
}

func TestHashSetMarshalJSON(t *testing.T) {
	cs := []struct {
		hset *HashSet
		json string
	}{
		{NewHashSet(0, 1, "0", "1", 0.1, 1.2, true, false, "0", "1"), `[0,1,"0","1",0.1,1.2,true,false]`},
		//TODO		{NewHashSet(0, "1", 2.0, 0, "1", 2.0, []int{1, 2}, map[int]int{1: 10, 2: 20}), `[0,"1",2,[1,2],{"1":10,"2":20}]`},
	}

	for i, c := range cs {
		bs, err := json.Marshal(c.hset)
		if err != nil {
			t.Errorf("[%d] Failed to Mashal HashSet %v : %v", i, c.hset, err)
			continue
		}

		act := sortSetJSON(string(bs))
		exp := sortSetJSON(c.json)

		if !reflect.DeepEqual(exp, act) {
			t.Errorf("[%d] Mashal HashSet (%v) = %v, want %v", i, c.hset, act, exp)
		}
	}
}

func TestHashSetUnmarshalJSON(t *testing.T) {
	cs := []struct {
		json string
		hset *HashSet
	}{
		{`["0","1",0,1,true,false]`, NewHashSet("0", "1", 0.0, 1.0, true, false)},
		//TODO		{`["1",2,[1,2],{"1":10,"2":20}]`, NewList("1", 2.0, NewList(1.0, 2.0), map[string]interface{}{"1": 10.0, "2": 20.0})},
	}

	for i, c := range cs {
		a := NewHashSet()
		err := json.Unmarshal([]byte(c.json), a)
		if err != nil {
			t.Errorf("[%d] Failed to Unmashal HashSet %v : %v", i, c.json, err)
			continue
		}
		if !reflect.DeepEqual(a, c.hset) {
			t.Errorf("[%d] Unmarshal List not deeply equal: %#v expected %#v", i, a, c.hset)
		}
	}
}
