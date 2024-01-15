package col

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/askasoft/pango/str"
)

func TestArrayListInterface(t *testing.T) {
	var l List = NewArrayList()
	if l == nil {
		t.Error("ArrayList is not a List")
	}

	var q Queue = NewArrayList()
	if q == nil {
		t.Error("ArrayList is not a Queue")
	}

	var dq Queue = NewArrayList()
	if dq == nil {
		t.Error("ArrayList is not a Deque")
	}

	var s Sortable = NewArrayList()
	if s == nil {
		t.Error("ArrayList is not a Sortable")
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

	if av := list2.Get(0); av != 1 {
		t.Errorf("Got %v expected %v", av, 1)
	}

	if av := list2.Get(1); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestArrayListAdd(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Adds("b", "c")
	if av := list.IsEmpty(); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := list.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av := list.Get(2); av != "c" {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func calcArrayCap(n int) int {
	c := minArrayCap
	for c < n {
		c <<= 1
	}
	return c
}

func TestArrayListGrow(t *testing.T) {
	list := &ArrayList{}

	for i := 0; i < 1000; i++ {
		list.Add(i)
		if l := list.Len(); l != i+1 {
			t.Errorf("list.Len() = %v, want %v", l, i+1)
		}

		wc := calcArrayCap(list.Len())
		ac := list.Cap()
		if wc != ac {
			t.Errorf("list.Cap(%d) = %v, want %v", i+1, ac, wc)
		}

		for n := 0; n <= i; n++ {
			if v := list.Get(n); v != n {
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
	list.Adds("b", "c")

	expectedIndex = 0
	if index := list.Index("a"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 1
	if index := list.Index("b"); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}

	expectedIndex = 2
	if index := list.IndexFunc(func(v any) bool {
		return v == "c"
	}); index != expectedIndex {
		t.Errorf("Got %v expected %v", index, expectedIndex)
	}
}

func TestArrayListDelete(t *testing.T) {
	l := NewArrayList()

	for i := 1; i <= 100; i++ {
		l.Add(i)
	}

	l.RemoveIf(func(d any) bool {
		return d == 101
	})
	if l.Len() != 100 {
		t.Error("ArrayList.Remove(101) should do nothing")
	}
	for i := 1; i <= 100; i++ {
		l.Remove(i)
		l.Removes(i, i)
		if l.Len() != 100-i {
			t.Errorf("ArrayList.Remove(%v) failed, l.Len() = %v, want %v", i, l.Len(), 100-i)
		}
	}

	if !l.IsEmpty() {
		t.Error("ArrayList.IsEmpty() should return true")
	}
}

func TestArrayListDelete2(t *testing.T) {
	l := NewArrayList()

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			l.Add(i)
		}
	}

	n := l.Len()
	l.Remove(100)
	if l.Len() != n {
		t.Errorf("ArrayList.Remove(100).Len() = %v, want %v", l.Len(), n)
	}
	for i := 0; i < 100; i++ {
		n = l.Len()
		z := i % 10
		l.Remove(i)
		a := n - l.Len()
		if a != z {
			t.Errorf("ArrayList.Remove(%v) = %v, want %v", i, a, z)
		}
	}

	if !l.IsEmpty() {
		t.Error("ArrayList.IsEmpty() should return true")
	}
}

func TestArrayListDeleteAll(t *testing.T) {
	l := NewArrayList()

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			l.Add(i)
		}
	}

	n := l.Len()
	l.Remove(100)
	if l.Len() != n {
		t.Errorf("ArrayList.Remove(100).Len() = %v, want %v", l.Len(), n)
	}
	for i := 0; i < 100; i++ {
		n = l.Len()
		z := i % 10
		l.RemoveCol(NewArrayList(i, i))
		a := n - l.Len()
		if a != z {
			t.Errorf("ArrayList.Remove(%v) = %v, want %v", i, a, z)
		}
	}

	if !l.IsEmpty() {
		t.Error("ArrayList.IsEmpty() should return true")
	}
}

func TestArrayListRemoveAt(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Adds("b", "c")
	list.RemoveAt(2)
	list.RemoveAt(1)
	list.RemoveAt(0)
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestArrayListRemoveAtPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewArrayList("a")
	list.RemoveAt(1)
}

func TestArrayListGet(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Adds("b", "c")
	if av := list.Get(0); av != "a" {
		t.Errorf("Got %v expected %v", av, "a")
	}
	if av := list.Get(1); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
	if av := list.Get(2); av != "c" {
		t.Errorf("Got %v expected %v", av, "c")
	}
	list.RemoveAt(0)
	if av := list.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestArrayListGetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewArrayList("a")
	list.Get(1)
}

func TestArrayListSwap(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Adds("b", "c")
	list.Swap(0, 1)
	if av := list.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestArrayListClear(t *testing.T) {
	list := NewArrayList()
	list.Adds("e", "f", "g", "a", "b", "c", "d")
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

	a := []T{}
	for i := 0; i < 100; i++ {
		a = append(a, i)
		list.Add(i)
	}
	a = append(a, 1000)

	for i := 0; i < 100; i++ {
		if !list.Contains(i) {
			t.Errorf("%d Contains() should return true", i)
		}
		if !list.Contains(a[0 : i+1]...) {
			t.Errorf("%d Contains(...) should return true", i)
		}
		if list.Contains(a...) {
			t.Errorf("%d Contains(...) should return false", i)
		}
		if !list.ContainCol(AsArrayList(a[0 : i+1])) {
			t.Errorf("%d ContainCol(...) should return true", i)
		}
		if list.ContainCol(AsArrayList(a)) {
			t.Errorf("%d ContainCol(...) should return false", i)
		}
	}

	list.Clear()
	if av := list.Contain("a"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := list.Contains("a", "b", "c"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
}

func TestArrayListRetain(t *testing.T) {
	for n := 0; n < 100; n++ {
		a := []T{}
		list := NewArrayList()
		for i := 0; i < n; i++ {
			if i&1 == 0 {
				a = append(a, i)
			}
			list.Add(i)

			list.Retains(a...)
			vs := list.Values()
			if !reflect.DeepEqual(vs, a) {
				t.Fatalf("%d Retains() = %v, want %v", i, vs, a)
			}
		}

		{
			a = []T{}
			list.Retains()
			vs := list.Values()
			if len(vs) > 0 {
				t.Fatalf("%d Retains() = %v, want %v", n, vs, a)
			}
		}

		a = []T{}
		list.Clear()
		for i := 0; i < n; i++ {
			if i&1 == 0 {
				a = append(a, i)
			}
			list.Add(i)

			list.RetainCol(AsArrayList(a))
			vs := list.Values()
			if !reflect.DeepEqual(vs, a) {
				t.Fatalf("%d RetainCol() = %v, want %v", i, vs, a)
			}
		}

		{
			a = []T{}
			list.RetainCol(AsArrayList(a))
			vs := list.Values()
			if len(vs) > 0 {
				t.Fatalf("%d Retains() = %v, want %v", n, vs, a)
			}
		}
	}
}

func TestArrayListValues(t *testing.T) {
	list := NewArrayList()
	list.Add("a")
	list.Adds("b", "c")
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "abc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestArrayListInsert(t *testing.T) {
	cs := []struct {
		s *ArrayList
		i int
		v []any
		w []any
	}{
		{NewArrayList(), 0, []any{"a", "b"}, []any{"a", "b"}},
		{NewArrayList("a", "b", "c"), 0, []any{"x", "y"}, []any{"x", "y", "a", "b", "c"}},
	}

	for i, c := range cs {
		c.s.Inserts(c.i, c.v...)
		a := c.s.Values()
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] Inserts(%v, %v) = %v, want %v", i, c.i, c.v, a, c.w)
		}
	}
}

func TestArrayListInsertPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewArrayList("a")
	list.Insert(2, "b")
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
	list.Set(1, "bb") // update
	if av := list.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "abbc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestArrayListSetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewArrayList("a")
	list.Set(1, "b")
}

func TestArrayListEach(t *testing.T) {
	list := NewArrayList()
	list.Adds("a", "b", "c")
	index := 0
	list.Each(func(value T) {
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

func TestArrayListIteratorPrevOnEmpty(t *testing.T) {
	list := NewArrayList()
	it := list.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestArrayListIteratorNextOnEmpty(t *testing.T) {
	list := NewArrayList()
	it := list.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestArrayListIteratorPrev(t *testing.T) {
	list := NewArrayList()
	list.Adds("a", "b", "c")
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

func TestArrayListIteratorNext(t *testing.T) {
	list := NewArrayList()
	list.Adds("a", "b", "c")
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

func TestArrayListIteratorReset(t *testing.T) {
	list := NewArrayList()

	it := list.Iterator()
	list.Adds("a", "b", "c")

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

func assertArrayListIteratorRemove(t *testing.T, i int, it Iterator, w *ArrayList) int {
	v := it.Value()

	it.Remove()

	w.Remove(v)

	it.SetValue(9999)

	l := it.(*arrayListIterator).list
	if l.Contains(v) {
		t.Fatalf("[%d] l.Contains(%v) = true", i, v)
	}

	if l.String() != w.String() {
		t.Fatalf("[%d] (%v) %v != %v", i, v, l.String(), w.String())
	}

	return v.(int)
}

func TestArrayListIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		l := NewArrayList()
		w := NewArrayList()

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

		v := assertArrayListIteratorRemove(t, i, it, w)

		it.Next()
		if v+1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v+1)
		}
		assertArrayListIteratorRemove(t, i, it, w)

		it.Prev()
		if v-1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v-1)
		}
		assertArrayListIteratorRemove(t, i, it, w)

		// remove first
		for it.Prev() {
		}
		assertArrayListIteratorRemove(t, i, it, w)

		// remove last
		for it.Next() {
		}
		assertArrayListIteratorRemove(t, i, it, w)

		// remove all
		it.Reset()
		if i%2 == 0 {
			for it.Prev() {
				assertArrayListIteratorRemove(t, i, it, w)
			}
		} else {
			for it.Next() {
				assertArrayListIteratorRemove(t, i, it, w)
			}
		}
		if !l.IsEmpty() {
			t.Fatalf("[%d] l.IsEmpty() = true", i)
		}
	}
}

func TestArrayListIteratorSetValue(t *testing.T) {
	l := NewArrayList()
	for i := 1; i <= 100; i++ {
		l.Add(i)
	}

	// forward
	for it := l.Iterator(); it.Next(); {
		it.SetValue(it.Value().(int) + 100)
	}
	for i := 1; i <= l.Len(); i++ {
		v := l.Get(i - 1).(int)
		w := i + 100
		if v != w {
			t.Fatalf("List[%d] = %v, want %v", i-1, v, w)
		}
	}

	// backward
	for it := l.Iterator(); it.Prev(); {
		it.SetValue(it.Value().(int) + 100)
	}
	for i := 1; i <= l.Len(); i++ {
		v := l.Get(i - 1).(int)
		w := i + 200
		if v != w {
			t.Fatalf("List[%d] = %v, want %v", i-1, v, w)
		}
	}
}

func TestArrayListSort(t *testing.T) {
	for i := 1; i < 100; i++ {
		l := NewArrayList()

		a := make([]T, 0, 100)
		for n := i; n < 100; n++ {
			a = append(a, rand.Intn(20))
		}

		for j := len(a) - 1; j >= 0; j-- {
			l.Add(a[j])
		}

		l.Sort(LessInt)
		sort.Sort(inta(a))

		if !reflect.DeepEqual(a, l.Values()) {
			t.Errorf("%v != %v", a, l.Values())
		}
	}
}

func checkArrayListLen(t *testing.T, l *ArrayList, len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkArrayList(t *testing.T, l *ArrayList, evs []T) {
	if !checkArrayListLen(t, l, len(evs)) {
		return
	}

	for i, it := 0, l.Iterator(); it.Next(); i++ {
		v := it.Value().(int)
		if v != evs[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, v, evs[i])
		}
	}

	avs := l.Values()
	for i, v := range avs {
		if v != evs[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, v, evs[i])
		}
	}
}

func TestArrayListExtending(t *testing.T) {
	l1 := NewArrayList(1, 2, 3)
	l2 := NewArrayList()
	l2.Add(4)
	l2.Add(5)

	l3 := NewArrayList()
	l3.AddCol(l1)
	checkArrayList(t, l3, []T{1, 2, 3})
	l3.AddCol(l2)
	checkArrayList(t, l3, []T{1, 2, 3, 4, 5})

	l3 = NewArrayList()
	l3.PushHeadCol(l2)
	checkArrayList(t, l3, []T{4, 5})
	l3.PushHeadCol(l1)
	checkArrayList(t, l3, []T{1, 2, 3, 4, 5})

	checkArrayList(t, l1, []T{1, 2, 3})
	checkArrayList(t, l2, []T{4, 5})

	l3 = NewArrayList()
	l3.PushTailCol(l1)
	checkArrayList(t, l3, []T{1, 2, 3})
	l3.PushTailCol(l3)
	checkArrayList(t, l3, []T{1, 2, 3, 1, 2, 3})

	l3 = NewArrayList()
	l3.PushHeadCol(l1)
	checkArrayList(t, l3, []T{1, 2, 3})
	l3.PushHeadCol(l3)
	checkArrayList(t, l3, []T{1, 2, 3, 1, 2, 3})

	l3 = NewArrayList()
	l1.PushTailCol(l3)
	checkArrayList(t, l1, []T{1, 2, 3})
	l1.PushHeadCol(l3)
	checkArrayList(t, l1, []T{1, 2, 3})

	l1.Clear()
	l2.Clear()
	l3.Clear()
	l1.PushTails(1, 2, 3)
	checkArrayList(t, l1, []T{1, 2, 3})
	l2.PushTails(4, 5)
	checkArrayList(t, l2, []T{4, 5})
	l3.PushTailCol(l1)
	checkArrayList(t, l3, []T{1, 2, 3})
	l3.PushTails(4, 5)
	checkArrayList(t, l3, []T{1, 2, 3, 4, 5})
	l3.PushHeads(4, 5)
	checkArrayList(t, l3, []T{4, 5, 1, 2, 3, 4, 5})
}

func TestArrayListContains2(t *testing.T) {
	l := NewArrayList(1, 11, 111, "1", "11", "111")

	n := (100+1)/101 + 110

	if !l.Contains(n) {
		t.Errorf("ArrayList [%v] should contains %v", l, n)
	}

	n++
	if l.Contains(n) {
		t.Errorf("ArrayList [%v] should not contains %v", l, n)
	}

	s := str.Repeat("1", 3)

	if !l.Contain(s) {
		t.Errorf("ArrayList [%v] should contains %v", l, s)
	}

	s += "0"
	if l.Contains(s) {
		t.Errorf("ArrayList [%v] should not contains %v", l, s)
	}
}

func TestArrayListQueue(t *testing.T) {
	q := NewArrayList()

	if _, ok := q.Peek(); ok {
		t.Error("should return false when peeking empty queue")
	}

	for i := 0; i < 100; i++ {
		q.Push(i)
	}

	for i := 0; i < 100; i++ {
		v, _ := q.Peek()
		if v.(int) != i {
			t.Errorf("Peek(%d) = %v, want %v", i, v, i)
		}

		x, _ := q.Poll()
		if x != i {
			t.Errorf("Poll(%d) = %v, want %v", i, x, i)
		}
	}

	if _, ok := q.Poll(); ok {
		t.Error("should return false when removing empty queue")
	}
}

func TestArrayListDeque(t *testing.T) {
	q := NewArrayList()

	if _, ok := q.PeekHead(); ok {
		t.Error("should return false when peeking empty queue")
	}

	if _, ok := q.PeekTail(); ok {
		t.Error("should return false when peeking empty queue")
	}

	for i := 0; i < 100; i++ {
		if i&1 == 0 {
			q.PushHead(i)
		} else {
			q.PushTail(i)
		}
	}

	for i := 0; i < 100; i++ {
		if i&1 == 0 {
			w := 100 - i - 2
			v, _ := q.PeekHead()
			if v != w {
				t.Errorf("PeekHead(%d) = %v, want %v", i, v, w)
			}

			x, _ := q.PollHead()
			if x != w {
				t.Errorf("PeekHead(%d) = %v, want %v", i, x, w)
			}
		} else {
			w := 100 - i
			v, _ := q.PeekTail()
			if v != w {
				t.Errorf("PeekTail(%d) = %v, want %v", i, v, w)
			}

			x, _ := q.PollTail()
			if x != w {
				t.Errorf("PoolTail(%d) = %v, want %v", i, x, w)
			}
		}
	}

	if _, ok := q.PollHead(); ok {
		t.Error("should return false when removing empty queue")
	}

	if _, ok := q.PollTail(); ok {
		t.Error("should return false when removing empty queue")
	}
}

func TestArrayListJSON(t *testing.T) {
	cs := []struct {
		s string
		a *ArrayList
	}{
		{`[]`, NewArrayList()},
		{`["a","b","c",[1,2,3]]`, NewArrayList("a", "b", "c", JSONArray{float64(1), float64(2), float64(3)})},
	}

	for i, c := range cs {
		a := NewArrayList()
		err := json.Unmarshal(([]byte)(c.s), &a)
		if err != nil {
			t.Errorf("[%d] json.Unmarshal(%q) = %v", i, c.s, err)
			continue
		}
		if !reflect.DeepEqual(a.Values(), c.a.Values()) {
			t.Errorf("[%d] json.Unmarshal(%q) = %v, want %v", i, c.s, a.Values(), c.a.Values())
		}

		bs, _ := json.Marshal(a)
		if string(bs) != c.s {
			t.Errorf("[%d] json.Marshal(%v) = %q, want %q", i, a.Values(), string(bs), c.s)
		}
	}
}
