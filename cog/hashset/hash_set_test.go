package hashset

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/cog/arraylist"
	"github.com/askasoft/pango/str"
)

func TestHashSetInterface(t *testing.T) {
	hs := NewHashSet(0)

	var _ cog.Collection[int] = hs
	var _ cog.Set[int] = hs
}

func TestHashSetLazyInit(t *testing.T) {
	{
		hs := &HashSet[int]{}
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
		if !hs.ContainsCol(&HashSet[int]{}) {
			t.Error("&HashSet{}.ContainsCol(&HashSet{}) = true, want false")
		}
	}
	{
		hs := &HashSet[int]{}
		hs.Add(1)
		if hs.Len() != 1 {
			t.Errorf("hs.Len() = %v, want 1", hs.Len())
		}
	}
	{
		hs := &HashSet[int]{}
		hs.AddCol(NewHashSet(1))
		if hs.Len() != 1 {
			t.Errorf("hs.Len() = %v, want 1", hs.Len())
		}
	}
	{
		hs := &HashSet[int]{}
		hs.Remove(1)
		if hs.Len() != 0 {
			t.Error("hs.Len() != 0")
		}
	}
	{
		hs := &HashSet[int]{}
		hs.Each(func(i, v int) bool { return true })
	}
	{
		hs := &HashSet[int]{}
		as := hs.Difference(&HashSet[int]{})
		if as.Len() != 0 {
			t.Errorf("hs.Difference(&HashSet{}) == %v, want %v", as, hs)
		}
	}
	{
		hs := &HashSet[int]{}
		as := hs.Intersection(&HashSet[int]{})
		if as.Len() != 0 {
			t.Errorf("hs.Intersection(&HashSet{}) == %v, want %v", as, hs)
		}
	}
}

func TestHashSetSimple(t *testing.T) {
	s := NewHashSet[int]()

	s.Add(5)
	if s.Len() != 1 {
		t.Errorf("Length should be 1")
	}
	if !s.Contains(5) {
		t.Errorf("Membership test failed")
	}

	s.Remove(5)
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
	set := NewHashSet[int]()
	set.AddAll()
	set.Add(1)
	set.Add(2)
	set.AddAll(2, 3)
	set.AddAll()
	if av := set.IsEmpty(); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := set.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
}

func TestHashSetContains(t *testing.T) {
	list := NewHashSet[int]()

	a := []int{}
	for i := 0; i < 100; i++ {
		a = append(a, i)
		list.Add(i)
	}
	a = append(a, 1000)

	for i := 0; i < 100; i++ {
		if !list.Contains(i) {
			t.Errorf("%d Contains() should return true", i)
		}
		if !list.ContainsAll(a[0 : i+1]...) {
			t.Errorf("%d ContainsAll(...) should return true", i)
		}
		if list.ContainsAll(a...) {
			t.Errorf("%d ContainsAll(...) should return false", i)
		}
		if !list.ContainsCol(arraylist.AsArrayList(a[0 : i+1])) {
			t.Errorf("%d ContainsCol(...) should return true", i)
		}
		if list.ContainsCol(arraylist.AsArrayList(a)) {
			t.Errorf("%d ContainsCol(...) should return false", i)
		}
	}

	list.Clear()
	if av := list.Contains(0); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := list.ContainsAll(0, 1, 2); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
}

func TestHashSetRetain(t *testing.T) {
	for n := 0; n < 100; n++ {
		a := []int{}
		list := NewHashSet[int]()
		for i := 0; i < n; i++ {
			if i&1 == 0 {
				a = append(a, i)
			}
			list.Add(i)

			list.RetainAll(a...)
			vs := list.Values()
			sort.Sort(inta(vs))
			if !reflect.DeepEqual(vs, a) {
				t.Fatalf("%d RetainAll() = %v, want %v", i, vs, a)
			}
		}

		{
			a = []int{}
			list.RetainAll()
			vs := list.Values()
			if len(vs) > 0 {
				t.Fatalf("%d RetainAll() = %v, want %v", n, vs, a)
			}
		}

		a = []int{}
		list.Clear()
		for i := 0; i < n; i++ {
			if i&1 == 0 {
				a = append(a, i)
			}
			list.Add(i)

			list.RetainCol(arraylist.AsArrayList(a))
			vs := list.Values()
			sort.Sort(inta(vs))
			if !reflect.DeepEqual(vs, a) {
				t.Fatalf("%d RetainCol() = %v, want %v", i, vs, a)
			}
		}

		{
			a = []int{}
			list.RetainCol(arraylist.AsArrayList(a))
			vs := list.Values()
			if len(vs) > 0 {
				t.Fatalf("%d RetainAll() = %v, want %v", n, vs, a)
			}
		}
	}
}

func TestHashSetDelete(t *testing.T) {
	set := NewHashSet[int]()
	set.AddAll(3, 1, 2)
	set.RemoveAll()
	if av := set.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	set.RemoveFunc(func(d int) bool {
		return d == 1
	})
	if av := set.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	set.RemoveAll(3, 3)
	set.RemoveAll()
	set.RemoveCol(arraylist.NewArrayList(2, 2))
	if av := set.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestHashSetContainsAll(t *testing.T) {
	s1 := NewHashSet[int]()
	s2 := NewHashSet[int]()

	if !s1.ContainsCol(s1) {
		t.Errorf("set should be a subset of itself")
	}

	if !s1.ContainsCol(s2) {
		t.Errorf("empty set should contains another empty set")
	}

	s1.Add(1)
	if !s1.ContainsCol(s2) {
		t.Errorf("set should contains another empty set")
	}

	s2.Add(1)
	if !s1.ContainsCol(s2) {
		t.Errorf("set should contains another same set")
	}

	s1.Add(2)
	if !s1.ContainsCol(s2) {
		t.Errorf("set should contains another small set")
	}

	s2.Add(3)
	if s1.ContainsCol(s2) {
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

	if !s3.ContainsAll(1, 2, 3) {
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

	if !s3.ContainsAll(4, 5, 6) {
		t.Errorf("Set should contain 4, 5, 6")
	}
}

func TestHashSetAddSet(t *testing.T) {
	// AddSet
	s1 := NewHashSet(4, 5, 6)
	s2 := NewHashSet(7, 8, 9)
	s1.AddCol(s2)

	if s1.Len() != 6 {
		t.Errorf("Length should be 6 after union")
	}

	for i := 4; i <= 9; i++ {
		if !(s1.Contains(i)) {
			t.Errorf("Set should contains %d", i)
		}
	}
}

func TestHashSetEach(t *testing.T) {
	ehs := NewHashSet("a", "b", "c")
	ahs := NewHashSet[string]()
	ehs.Each(func(i int, s string) bool {
		ahs.Add(s)
		return true
	})

	evs := ehs.Values()
	sort.Strings(evs)
	w := fmt.Sprint(evs)

	avs := ahs.Values()
	sort.Strings(avs)
	a := fmt.Sprint(avs)
	if a != w {
		t.Errorf("Each():\nWANT: %s\n GOT: %s", w, a)
	}
}

func TestHashSetSeq(t *testing.T) {
	ehs := NewHashSet("a", "b", "c")
	ahs := NewHashSet[string]()
	for s := range ehs.Seq() {
		ahs.Add(s)
	}

	evs := ehs.Values()
	sort.Strings(evs)
	w := fmt.Sprint(evs)

	avs := ahs.Values()
	sort.Strings(avs)
	a := fmt.Sprint(avs)
	if a != w {
		t.Errorf("Each():\nWANT: %s\n GOT: %s", w, a)
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
		hset *HashSet[int]
		json string
	}{
		{NewHashSet(0, 1, -1, -2, 4, 3), `[-2,-1,0,1,3,4]`},
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
		hset *HashSet[int]
	}{
		{`[-2,-1,0,1,3,4]`, NewHashSet(0, 1, -1, -2, 4, 3)},
	}

	for i, c := range cs {
		a := NewHashSet[int]()
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
