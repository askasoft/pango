package col

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/pandafw/pango/cmp"
)

func TestArrayListInterface(t *testing.T) {
	var l List = NewArrayList()
	if l == nil {
		t.Error("ArrayList is not a List")
	}
}

func TestArrayListRoundup(t *testing.T) {
	cs := []struct {
		s int
		w int
	}{
		{0, 0},
		{10, 32},
		{20, 32},
		{31, 32},
		{32, 32},
		{33, 64},
	}

	al := &ArrayList{}
	for i, c := range cs {
		a := al.roundup(c.s)
		if a != c.w {
			t.Errorf("[%d] roundup(%d) = %d, want %d", i, c.s, a, c.w)
		}
	}
}

func TestArrayListNew(t *testing.T) {
	list1 := NewArrayList()

	if av := list1.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}

	list2 := NewArrayList(1, "b")

	if av := list2.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}

	if av, ok := list2.Get(0); av != 1 || !ok {
		t.Errorf("Got %v expected %v", av, 1)
	}

	if av, ok := list2.Get(1); av != "b" || !ok {
		t.Errorf("Got %v expected %v", av, "b")
	}

	if av, ok := list2.Get(2); av != nil || ok {
		t.Errorf("Got %v expected %v", av, nil)
	}
}

func TestArrayListAdd(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Add("b", "c")
	if av := list.IsEmpty(); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := list.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av, ok := list.Get(2); av != "c" || !ok {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func TestArrayListGrow(t *testing.T) {
	list := NewArrayList()

	for i := 0; i < 1000; i++ {
		list.Add(i)
		if l := list.Len(); l != i+1 {
			t.Errorf("list.Len() = %v, want %v", l, i+1)
		}

		wc := list.roundup(i + 1)
		ac := cap(list.data)
		if wc != ac {
			t.Errorf("list.Cap(%d) = %v, want %v", i+1, ac, wc)
		}

		for n := 0; n <= i; n++ {
			if v, ok := list.Get(n); v != n || !ok {
				t.Errorf("list.Get(%d) = %v, want %v", n, v, n)
			}
		}
	}
}

func TestArrayListIndex(t *testing.T) {
	list := NewArrayList()

	expectedIndex := -1
	if index := list.Index("a"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	list.Add("a")
	list.Add("b", "c")

	expectedIndex = 0
	if index := list.Index("a"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 1
	if index := list.Index("b"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 2
	if index := list.Index("c"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}
}

func TestArrayListRemove(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Add("b", "c")
	list.Remove(2)
	if av, ok := list.Get(2); av != nil || ok {
		t.Errorf("Got %v expected %v", av, nil)
	}
	list.Remove(1)
	list.Remove(0)
	list.Remove(0) // no effect
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestArrayListDelete(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Add("b", "c")
	list.Delete("c")
	if av, ok := list.Get(2); av != nil || ok {
		t.Errorf("Got %v expected %v", av, nil)
	}
	list.Delete("b", "a")
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
	list.Delete("a") // no effect
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestArrayListGet(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Add("b", "c")
	if av, ok := list.Get(0); av != "a" || !ok {
		t.Errorf("Got %v expected %v", av, "a")
	}
	if av, ok := list.Get(1); av != "b" || !ok {
		t.Errorf("Got %v expected %v", av, "b")
	}
	if av, ok := list.Get(2); av != "c" || !ok {
		t.Errorf("Got %v expected %v", av, "c")
	}
	if av, ok := list.Get(3); av != nil || ok {
		t.Errorf("Got %v expected %v", av, nil)
	}
	list.Remove(0)
	if av, ok := list.Get(0); av != "b" || !ok {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestArrayListSwap(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Add("b", "c")
	list.Swap(0, 1)
	if av, ok := list.Get(0); av != "b" || !ok {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestArrayListSort(t *testing.T) {
	for i := 1; i < 100; i++ {
		l := NewArrayList()

		a := make([]interface{}, 0, 100)
		for n := i; n < 100; n++ {
			a = append(a, rand.Intn(20))
		}

		for j := len(a) - 1; j >= 0; j-- {
			l.Add(a[j])
		}

		l.Sort(cmp.LessInt)
		sort.Sort(inta(a))

		if !reflect.DeepEqual(a, l.Values()) {
			t.Errorf("%v != %v", a, l.Values())
		}
	}
}

func TestArrayListClear(t *testing.T) {
	list := NewArrayList()
	list.Add("e", "f", "g", "a", "b", "c", "d")
	list.Clear()
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestArrayListContains(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Add("b", "c")
	if av := list.Contains("a"); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Contains("a", "b", "c"); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Contains("a", "b", "c", "d"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	list.Clear()
	if av := list.Contains("a"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := list.Contains("a", "b", "c"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
}

func TestArrayListValues(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Add("b", "c")
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "abc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestArrayListInsert(t *testing.T) {
	list := NewArrayList()
	list.Insert(0, "b", "c")
	list.Insert(0, "a")
	list.Insert(10, "x") // ignore
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

func TestArrayListSet(t *testing.T) {
	list := NewArrayList("", "")
	list.Set(0, "a")
	list.Set(1, "b")
	if av := list.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	list.Add("")
	list.Set(2, "c")
	if av := list.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	list.Set(4, "d")  // ignore
	list.Set(1, "bb") // update
	if av := list.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "abbc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestArrayListEach(t *testing.T) {
	list := NewArrayList()
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

func TestArrayListIteratorNextOnEmpty(t *testing.T) {
	list := NewArrayList()
	it := list.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestArrayListIteratorNext(t *testing.T) {
	list := NewArrayList()
	list.Add("a", "b", "c")
	it := list.Iterator()
	count := 0
	for index := 0; it.Next(); index++ {
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
	}
	if av, ev := count, 3; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestArrayListIteratorPrevOnEmpty(t *testing.T) {
	list := NewArrayList()
	it := list.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestArrayListIteratorPrev(t *testing.T) {
	list := NewArrayList()
	list.Add("a", "b", "c")
	it := list.Iterator()
	count := 0
	index := list.Len()
	for it.Prev() {
		index--
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
	}
	if av, ev := count, 3; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestArrayListIteratorReset(t *testing.T) {
	list := NewArrayList()

	list.Add("a", "b", "c")
	it := list.Iterator()

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

func TestArrayListIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		l := NewArrayList()

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

func TestArrayListJSON(t *testing.T) {
	cs := []struct {
		s string
		a *ArrayList
	}{
		{`[]`, NewArrayList()},
		{`["a","b","c"]`, NewArrayList("a", "b", "c")},
	}

	for i, c := range cs {
		a := NewArrayList()
		json.Unmarshal(([]byte)(c.s), &a)
		if !reflect.DeepEqual(a.Values(), c.a.Values()) {
			t.Errorf("[%d] json.Unmarshal(%q) = %v, want %v", i, c.s, a.Values(), c.a.Values())
		}

		bs, _ := json.Marshal(a)
		if string(bs) != c.s {
			t.Errorf("[%d] json.Marshal(%v) = %q, want %q", i, a.Values(), string(bs), c.s)
		}
	}
}
