package col

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/pandafw/pango/cmp"
	"github.com/pandafw/pango/str"
)

type inta []interface{}

// Len is the number of elements in the collection.
func (a inta) Len() int {
	return len(a)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (a inta) Less(i, j int) bool {
	return a[i].(int) < a[j].(int)
}

// Swap swaps the elements with indexes i and j.
func (a inta) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func TestSortedListInterface(t *testing.T) {
	var l List = NewSortedList(cmp.LessInt)
	if l == nil {
		t.Error("SortedList is not a List")
	}
}

func TestSortedListAsc(t *testing.T) {
	for i := 1; i < 100; i++ {
		sl := NewSortedList(cmp.LessInt)
		a := make([]interface{}, 0, 100)
		for n := i; n < 100; n++ {
			a = append(a, n)
		}

		for j := 0; j < len(a); j++ {
			sl.Add(a[j])
			if !sort.IsSorted(inta(sl.Values())) {
				t.Errorf("[%d] sort.IsSorted(%v) = %v, want %v", j, sl.Values(), false, true)
			}
		}

		if !reflect.DeepEqual(a, sl.Values()) {
			t.Errorf("%v != %v", a, sl.Values())
		}
	}
}

func TestSortedListDesc(t *testing.T) {
	for i := 1; i < 100; i++ {
		sl := NewSortedList(cmp.LessInt)
		a := make([]interface{}, 0, 100)
		for n := i; n < 100; n++ {
			a = append(a, n)
		}
		for j := len(a) - 1; j >= 0; j-- {
			sl.Add(a[j])

			if !sort.IsSorted(inta(sl.Values())) {
				t.Errorf("[%d] sort.IsSorted(%v) = %v, want %v", j, sl.Values(), false, true)
			}
		}

		if !reflect.DeepEqual(a, sl.Values()) {
			t.Errorf("%v != %v", a, sl.Values())
		}
	}
}

func TestSortedListRandom(t *testing.T) {
	for i := 1; i < 100; i++ {
		sl := NewSortedList(cmp.LessInt)
		a := make([]interface{}, 0, 100)
		for n := i; n < 100; n++ {
			a = append(a, rand.Intn(20))
		}
		for j := len(a) - 1; j >= 0; j-- {
			sl.Add(a[j])

			if !sort.IsSorted(inta(sl.Values())) {
				t.Errorf("[%d] sort.IsSorted(%v) = %v, want %v", j, sl.Values(), false, true)
			}
		}

		sort.Sort(inta(a))

		if !reflect.DeepEqual(a, sl.Values()) {
			t.Errorf("%v != %v", a, sl.Values())
		}
	}
}

func TestSortedListContains(t *testing.T) {
	sl := NewSortedList(cmp.LessString, "1", "11", "111", "1", "11", "111")
	n := str.Repeat("1", 3)

	if !sl.Contains(n) {
		t.Errorf("SortedList [%v] should contains %v", sl, n)
	}

	n += "1"
	if sl.Contains(n) {
		t.Errorf("SortedList [%v] should not contains %v", sl, n)
	}
}

func TestSetSearch(t *testing.T) {
	sl := NewSortedList(cmp.LessInt, 1, 11)

	n := 10 + 1
	sn, se := sl.Search(n)
	if se == nil || sn != 1 {
		t.Errorf("SortedList [%v] should contains %v", sl, n)
	}

	n = 1 + 10 + 100
	for i := 0; i < 100; i++ {
		sl.Add(111)
		sn, se = sl.Search(n)
		if se == nil || sn != 2 {
			t.Errorf("SortedList [%v] should contains %v", sl, n)
		}
	}

	n++
	sn, se = sl.Search(n)
	if se != nil || sn != -1 {
		t.Errorf("SortedList [%v] should not contains %v", sl, n)
	}
}

func TestSortedListDelete(t *testing.T) {
	sl := NewSortedList(cmp.LessInt)

	for i := 1; i <= 100; i++ {
		sl.Add(i)
	}

	sl.Delete(101)
	if sl.Len() != 100 {
		t.Error("SortedList.Delete(101) should do nothing")
	}
	for i := 1; i <= 100; i++ {
		sl.Delete(i)
		if sl.Len() != 100-i {
			t.Errorf("SortedList.Delete(%v) failed, l.Len() = %v, want %v", i, sl.Len(), 100-i)
		}
	}
	if !sl.IsEmpty() {
		t.Error("sl.IsEmpty()=false, want true")
	}
}

func TestSortedListDeleteAll(t *testing.T) {
	sl := NewSortedList(cmp.LessInt)

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			sl.Add(i)
		}
	}

	n := sl.Len()
	sl.Delete(100)
	if sl.Len() != n {
		t.Errorf("SortedList.Delete(100).Len() = %v, want %v", sl.Len(), n)
	}
	for i := 0; i < 100; i++ {
		n = sl.Len()
		z := i % 10
		sl.Delete(i)
		a := n - sl.Len()
		if a != z {
			t.Errorf("SortdList.Delete(%v) = %v, want %v", i, a, z)
		}
	}
	if !sl.IsEmpty() {
		t.Error("sl.IsEmpty()=false, want true")
	}
}

func TestSortedListIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		l := NewSortedList(cmp.LessInt)

		for n := 0; n < i; n++ {
			l.Add(n)
		}

		it := l.Iterator()

		it.Remove()
		if l.Len() != i {
			t.Errorf("[%d] l.Len() == %v, want %v", i, l.Len(), i)
		}

		// remove middle
		x := rand.Intn(i-4) + 1
		for j := 0; j <= x; j++ {
			it.Next()
		}

		v := it.Value().(int)
		it.Remove()
		if l.Len() != i-1 {
			t.Errorf("[%d] l.Len() == %v, want %v", i, l.Len(), i-1)
		}
		if l.Contains(v) {
			t.Errorf("[%d] l.Contains(%v) = true", i, v)
		}

		it.Next()
		if v+1 != it.Value() {
			t.Errorf("[%d] it.Value() = %v, want %v", i, it.Value(), v+1)
		}
		it.Remove()
		if l.Contains(v + 1) {
			t.Errorf("[%d] l.Contains(%v) = true", i, v+1)
		}

		it.Prev()
		if v-1 != it.Value() {
			t.Errorf("[%d] it.Value() = %v, want %v", i, it.Value(), v-1)
		}
		it.Remove()
		if l.Contains(v - 1) {
			t.Errorf("[%d] l.Contains(%v) = true", i, v-1)
		}

		// remove first
		for it.Prev() {
		}
		it.Remove()
		if l.Contains(0) {
			t.Errorf("[%d] l.Contains(%v) = true", i, 0)
		}
		if it.Prev() {
			t.Errorf("[%d] l.Prev() = true", i)
		}

		// remove last
		for it.Next() {
		}
		it.Remove()
		if l.Contains(i - 1) {
			t.Errorf("[%d] l.Contains(%v) = true", i, i-1)
		}
		if it.Next() {
			t.Errorf("[%d] l.Next() = true", i)
		}

		// remove all
		it.Reset()
		if i%2 == 0 {
			for it.Prev() {
				it.Remove()
			}
		} else {
			for it.Next() {
				it.Remove()
			}
		}
		if !l.IsEmpty() {
			t.Errorf("[%d] l.IsEmpty() = true", i)
		}
	}
}

func TestSortedListString(t *testing.T) {
	w := "[1,2,3]"
	a := fmt.Sprintf("%s", NewSortedList(cmp.LessInt, 1, 3, 2))
	if w != a {
		t.Errorf("SortedList.String() = %v, want %v", a, w)
	}
}

func TestSortedListMarshalJSON(t *testing.T) {
	type Case struct {
		list *SortedList
		json string
	}

	cs := []Case{
		{NewSortedList(cmp.LessString, "2", "1", "0"), `["0","1","2"]`},
	}

	for i, c := range cs {
		bs, err := json.Marshal(c.list)

		if err != nil {
			t.Fatalf("[%d] json.Marshal(%v) error: %v", i, c.list, err)
		}

		a := string(bs)
		if !reflect.DeepEqual(c.json, a) {
			t.Fatalf("[%d] json.Marshal(%v) = %v, want %v", i, c.list, a, c.json)
		}
	}
}

func TestSortedListUnmarshalJSON(t *testing.T) {
	cs := []struct {
		json string
		list List
	}{
		{`["2","0","1"]`, NewLinkedList("0", "1", "2")},
	}

	for i, c := range cs {
		a := NewSortedList(cmp.LessString)
		err := json.Unmarshal([]byte(c.json), a)

		if err != nil {
			t.Fatalf("[%d] json.Unmarshal(%v) error: %v", i, c.json, err)
		}

		if !reflect.DeepEqual(a.Values(), c.list.Values()) {
			t.Fatalf("[%d] json.Marshal(%v) = %v, want %v", i, c.json, a.Values(), c.list.Values())
		}
	}
}
