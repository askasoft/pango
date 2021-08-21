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
				t.Fatalf("[%d] sort.IsSorted(%v) = %v, want %v", j, sl.Values(), false, true)
			}
		}

		if !reflect.DeepEqual(a, sl.Values()) {
			t.Fatalf("%v != %v", a, sl.Values())
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
				t.Fatalf("[%d] sort.IsSorted(%v) = %v, want %v", j, sl.Values(), false, true)
				return
			}
		}

		if !reflect.DeepEqual(a, sl.Values()) {
			t.Fatalf("%v != %v", a, sl.Values())
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

func TestSortedListRemove(t *testing.T) {
	list := NewSortedList(cmp.LessString)
	list.Add("a")
	list.Add("b", "c")
	list.Remove(2)
	list.Remove(1)
	list.Remove(0)
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestSortedListRemovePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewSortedList(cmp.LessString, "a")
	list.Remove(1)
}

func TestSortedListGet(t *testing.T) {
	list := NewSortedList(cmp.LessString)
	list.Add("a")
	list.Add("b", "c")
	if av := list.Get(0); av != "a" {
		t.Errorf("Got %v expected %v", av, "a")
	}
	if av := list.Get(1); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
	if av := list.Get(2); av != "c" {
		t.Errorf("Got %v expected %v", av, "c")
	}
	list.Remove(0)
	if av := list.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestSortedListGetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewSortedList(cmp.LessString, "a")
	list.Get(1)
}

func TestSortedListClear(t *testing.T) {
	list := NewSortedList(cmp.LessString)
	list.Add("e", "f", "g", "a", "b", "c", "d")
	list.Clear()
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
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

func TestSortedListValues(t *testing.T) {
	list := NewSortedList(cmp.LessString)
	list.Add("a")
	list.Add("b", "c")
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "abc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestSortedListInsert(t *testing.T) {
	list := NewSortedList(cmp.LessString)
	list.Insert(0, "b", "c")
	list.Insert(0, "a")
	if av := list.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	list.Insert(3, "d") // append
	if av := list.Len(); av != 4 {
		t.Errorf("Got %v expected %v", av, 4)
	}
	if av, ev := fmt.Sprintf("%s%s%s%s", list.Values()...), "abcd"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestSortedListInsertPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewSortedList(cmp.LessString, "a")
	list.Insert(2, "b")
}

func TestSortedListSet(t *testing.T) {
	list := NewSortedList(cmp.LessString, "z", "z")
	list.Set(0, "a")
	list.Set(1, "b")
	if av := list.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	list.Add("z")
	list.Set(2, "c") // last
	if av := list.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	list.Set(1, "bb") // update
	if av := list.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "abbc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	list.Set(2, "cc") // last to first traversal
	list.Set(0, "aa") // first to last traversal
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "aabbcc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestSortedListSetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewSortedList(cmp.LessString, "a")
	list.Set(1, "b")
}

func TestSortedListEach(t *testing.T) {
	list := NewSortedList(cmp.LessString)
	list.Add("a", "b", "c")
	index := 0
	list.Each(func(value interface{}) {
		switch index {
		case 0:
			if av, ev := value, "a"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 1:
			if av, ev := value, "b"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 2:
			if av, ev := value, "c"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		default:
			t.Errorf("Too many")
		}
		index++
	})
}

func TestSortedListIteratorPrevOnEmpty(t *testing.T) {
	list := NewSortedList(cmp.LessString)
	it := list.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestSortedListIteratorNextOnEmpty(t *testing.T) {
	list := NewSortedList(cmp.LessString)
	it := list.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestSortedListIteratorPrev(t *testing.T) {
	list := NewSortedList(cmp.LessString)
	list.Add("a", "b", "c")
	it := list.Iterator()
	count := 0
	index := list.Len()
	for it.Prev() {
		count++
		index--
		value := it.Value()
		switch index {
		case 0:
			if av, ev := value, "a"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 1:
			if av, ev := value, "b"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 2:
			if av, ev := value, "c"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		default:
			t.Errorf("Too many")
		}
	}
	if av, ev := count, 3; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestSortedListIteratorNext(t *testing.T) {
	list := NewSortedList(cmp.LessString)
	list.Add("a", "b", "c")
	it := list.Iterator()
	count := 0
	index := 0
	for it.Next() {
		count++
		value := it.Value()
		switch index {
		case 0:
			if av, ev := value, "a"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 1:
			if av, ev := value, "b"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		case 2:
			if av, ev := value, "c"; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		default:
			t.Errorf("Too many")
		}
		index++
	}
	if av, ev := count, 3; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestSortedListIteratorReset(t *testing.T) {
	list := NewSortedList(cmp.LessString)

	it := list.Iterator()
	list.Add("a", "b", "c")

	for it.Next() {
	}
	it.Reset()
	it.Next()
	if value := it.Value(); value != "a" {
		t.Errorf("Got %v expected %v", value, "a")
	}

	for it.Prev() {
	}
	it.Reset()
	it.Prev()
	if value := it.Value(); value != "c" {
		t.Errorf("Got %v expected %v", value, "c")
	}
}

func assertSortedListIteratorRemove(t *testing.T, i int, it Iterator, w *SortedList) int {
	it.Remove()

	v := it.Value()
	w.Delete(v)

	it.SetValue(9999)

	l := it.(*sortedListIterator).list
	if l.Contains(v) {
		t.Fatalf("[%d] l.Contains(%v) = true", i, v)
	}

	if l.String() != w.String() {
		t.Fatalf("[%d] (%v) %v != %v", i, v, l.String(), w.String())
	}

	return v.(int)
}

func TestSortedListIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		l := NewSortedList(cmp.LessInt)
		w := NewSortedList(cmp.LessInt)

		for n := 0; n < i; n++ {
			l.Add(n)
			w.Add(n)
		}

		it := l.Iterator()

		it.Remove()
		it.SetValue(9999)
		if l.Len() != i {
			t.Fatalf("[%d] l.Len() == %v, want %v", i, l.Len(), i)
		}

		// remove middle
		for j := 0; j <= l.Len()/2; j++ {
			it.Next()
		}

		v := assertSortedListIteratorRemove(t, i, it, w)

		it.Next()
		if v+1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v+1)
		}
		assertSortedListIteratorRemove(t, i, it, w)

		it.Prev()
		if v-1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v-1)
		}
		assertSortedListIteratorRemove(t, i, it, w)

		// remove first
		for it.Prev() {
		}
		assertSortedListIteratorRemove(t, i, it, w)

		// remove last
		for it.Next() {
		}
		assertSortedListIteratorRemove(t, i, it, w)

		// remove all
		it.Reset()
		if i%2 == 0 {
			for it.Prev() {
				assertSortedListIteratorRemove(t, i, it, w)
			}
		} else {
			for it.Next() {
				assertSortedListIteratorRemove(t, i, it, w)
			}
		}
		if !l.IsEmpty() {
			t.Fatalf("[%d] l.IsEmpty() = true", i)
		}
	}
}

func TestSortedListIteratorSetValue(t *testing.T) {
	l := NewSortedList(cmp.LessInt)
	for i := 1; i <= 100; i++ {
		l.Add(i)
	}

	// forward
	for it := l.Iterator(); it.Next() && it.Value().(int) <= 100; it.Reset() {
		v := it.Value().(int) + 100
		it.SetValue(v)
		if it.Next() {
			t.Fatalf("SortedList(%v).Next() should false", v)
		}
	}
	//fmt.Println(l)
	for i := 1; i <= l.Len(); i++ {
		v := l.Get(i - 1).(int)
		w := i + 100
		if v != w {
			t.Fatalf("SortedList[%d] = %v, want %v", i-1, v, w)
		}
	}

	// backward
	for it := l.Iterator(); it.Prev() && it.Value().(int) <= 200; {
		it.SetValue(it.Value().(int) + 100)
	}
	for i := 1; i <= l.Len(); i++ {
		v := l.Get(i - 1).(int)
		w := i + 200
		if v != w {
			t.Fatalf("SortedList[%d] = %v, want %v", i-1, v, w)
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
		{`["2","0","1"]`, NewSortedList(cmp.LessString, "0", "1", "2")},
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
