package col

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestLinkedListInterface(t *testing.T) {
	var l List = NewLinkedList()
	if l == nil {
		t.Error("LinkedList is not a List")
	}

	var q Queue = NewLinkedList()
	if q == nil {
		t.Error("LinkedList is not a Queue")
	}

	var dq Deque = NewLinkedList()
	if dq == nil {
		t.Error("LinkedList is not a Deque")
	}

	var st Sortable = NewLinkedList()
	if st == nil {
		t.Error("LinkedList is not a Sortable")
	}
}

func TestLinkedListNew(t *testing.T) {
	list1 := NewLinkedList()

	if av := list1.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}

	list2 := NewLinkedList(1, "b")

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

func TestLinkedListAdd(t *testing.T) {
	list := NewLinkedList()
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

func TestLinkedListIndex(t *testing.T) {
	list := NewLinkedList()

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

func TestLinkedListDelete(t *testing.T) {
	l := NewLinkedList()

	for i := 1; i <= 100; i++ {
		l.Add(i)
	}

	l.RemoveIf(func(d any) bool {
		return d == 101
	})
	if l.Len() != 100 {
		t.Error("LinkedList.Delete(101) should do nothing")
	}
	for i := 1; i <= 100; i++ {
		l.Remove(i)
		l.Removes(i, i)
		if l.Len() != 100-i {
			t.Errorf("LinkedList.Remove(%v) failed, l.Len() = %v, want %v", i, l.Len(), 100-i)
		}
	}

	if !l.IsEmpty() {
		t.Error("LinkedList.IsEmpty() should return true")
	}
}

func TestLinkedListDelete2(t *testing.T) {
	l := NewLinkedList()

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			l.Add(i)
		}
	}

	n := l.Len()
	l.Remove(100)
	if l.Len() != n {
		t.Errorf("LinkedList.Remove(100).Len() = %v, want %v", l.Len(), n)
	}
	for i := 0; i < 100; i++ {
		n = l.Len()
		z := i % 10
		l.Remove(i)
		a := n - l.Len()
		if a != z {
			t.Errorf("LinkedList.Remove(%v) = %v, want %v", i, a, z)
		}
	}

	if !l.IsEmpty() {
		t.Error("LinkedList.IsEmpty() should return true")
	}
}

func TestLinkedListDeleteAll(t *testing.T) {
	l := NewLinkedList()

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			l.Add(i)
		}
	}

	n := l.Len()
	l.Remove(100)
	if l.Len() != n {
		t.Errorf("LinkedList.Remove(100).Len() = %v, want %v", l.Len(), n)
	}
	for i := 0; i < 100; i++ {
		n = l.Len()
		z := i % 10
		l.RemoveCol(NewArrayList(i, i))
		a := n - l.Len()
		if a != z {
			t.Errorf("LinkedList.Remove(%v) = %v, want %v", i, a, z)
		}
	}

	if !l.IsEmpty() {
		t.Error("LinkedList.IsEmpty() should return true")
	}
}

func TestLinkedListDeleteAt(t *testing.T) {
	list := NewLinkedList()
	list.Add("a")
	list.Adds("b", "c")
	list.DeleteAt(2)
	list.DeleteAt(1)
	list.DeleteAt(0)
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestLinkedListDeleteAtPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedList("a")
	list.DeleteAt(1)
}

func TestLinkedListGet(t *testing.T) {
	list := NewLinkedList()
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
	list.DeleteAt(0)
	if av := list.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestLinkedListGetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedList("a")
	list.Get(1)
}

func TestLinkedListSwap(t *testing.T) {
	list := NewLinkedList()
	list.Add("a")
	list.Adds("b", "c")
	list.Swap(0, 1)
	if av := list.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func TestLinkedListClear(t *testing.T) {
	list := NewLinkedList()
	list.Adds("e", "f", "g", "a", "b", "c", "d")
	list.Clear()
	if av := list.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := list.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestLinkedListContains(t *testing.T) {
	list := NewLinkedList()

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
	if av := list.Contains("a"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := list.Contains("a", "b", "c"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
}

func TestLinkedListRetain(t *testing.T) {
	for n := 0; n < 100; n++ {
		a := []T{}
		list := NewLinkedList()
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

func TestLinkedListValues(t *testing.T) {
	list := NewLinkedList()
	list.Add("a")
	list.Adds("b", "c")
	if av, ev := fmt.Sprintf("%s%s%s", list.Values()...), "abc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestLinkedListInsert(t *testing.T) {
	list := NewLinkedList()
	list.Inserts(0, "b", "c")
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

func TestLinkedListInsertPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedList("a")
	list.Insert(2, "b")
}

func TestLinkedListSet(t *testing.T) {
	list := NewLinkedList("", "")
	list.Set(0, "a")
	list.Set(1, "b")
	if av := list.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	list.Add("")
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

func TestLinkedListSetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedList("a")
	list.Set(1, "b")
}

func TestLinkedListEach(t *testing.T) {
	list := NewLinkedList()
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

func TestLinkedListIteratorPrevOnEmpty(t *testing.T) {
	list := NewLinkedList()
	it := list.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestLinkedListIteratorNextOnEmpty(t *testing.T) {
	list := NewLinkedList()
	it := list.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty list")
	}
}

func TestLinkedListIteratorPrev(t *testing.T) {
	list := NewLinkedList()
	list.Adds("a", "b", "c")
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

func TestLinkedListIteratorNext(t *testing.T) {
	list := NewLinkedList()
	list.Adds("a", "b", "c")
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

func TestLinkedListIteratorReset(t *testing.T) {
	list := NewLinkedList()

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

func assertLinkedListIteratorRemove(t *testing.T, i int, it Iterator, w *LinkedList) int {
	it.Remove()

	v := it.Value()
	w.Remove(v)

	it.SetValue(9999)

	l := it.(*linkedListIterator).list
	if l.Contains(v) {
		t.Fatalf("[%d] l.Contains(%v) = true", i, v)
	}

	if l.String() != w.String() {
		t.Fatalf("[%d] (%v) %v != %v", i, v, l.String(), w.String())
	}

	return v.(int)
}

func TestLinkedListIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		l := NewLinkedList()
		w := NewLinkedList()

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

		v := assertLinkedListIteratorRemove(t, i, it, w)

		it.Next()
		if v+1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v+1)
		}
		assertLinkedListIteratorRemove(t, i, it, w)

		it.Prev()
		if v-1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v-1)
		}
		assertLinkedListIteratorRemove(t, i, it, w)

		// remove first
		for it.Prev() {
		}
		assertLinkedListIteratorRemove(t, i, it, w)

		// remove last
		for it.Next() {
		}
		assertLinkedListIteratorRemove(t, i, it, w)

		// remove all
		it.Reset()
		if i%2 == 0 {
			for it.Prev() {
				assertLinkedListIteratorRemove(t, i, it, w)
			}
		} else {
			for it.Next() {
				assertLinkedListIteratorRemove(t, i, it, w)
			}
		}
		if !l.IsEmpty() {
			t.Fatalf("[%d] l.IsEmpty() = true", i)
		}
	}
}

func TestLinkedListIteratorSetValue(t *testing.T) {
	l := NewLinkedList()
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

func TestLinkedListSort(t *testing.T) {
	for i := 1; i < 100; i++ {
		l := NewLinkedList()

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

func checkLinkedListLen(t *testing.T, l *LinkedList, len int) bool {
	if n := l.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkLinkedList(t *testing.T, l *LinkedList, evs []T) {
	if !checkLinkedListLen(t, l, len(evs)) {
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

func TestLinkedListExtending(t *testing.T) {
	l1 := NewLinkedList(1, 2, 3)
	l2 := NewLinkedList()
	l2.PushTail(4)
	l2.PushTail(5)

	l3 := NewLinkedList()
	l3.PushTailCol(l1)
	checkLinkedList(t, l3, []T{1, 2, 3})
	l3.PushTailCol(l2)
	checkLinkedList(t, l3, []T{1, 2, 3, 4, 5})

	l3 = NewLinkedList()
	l3.PushHeadCol(l2)
	checkLinkedList(t, l3, []T{4, 5})
	l3.PushHeadCol(l1)
	checkLinkedList(t, l3, []T{1, 2, 3, 4, 5})

	checkLinkedList(t, l1, []T{1, 2, 3})
	checkLinkedList(t, l2, []T{4, 5})

	l3 = NewLinkedList()
	l3.PushTailCol(l1)
	checkLinkedList(t, l3, []T{1, 2, 3})
	l3.PushTailCol(l3)
	checkLinkedList(t, l3, []T{1, 2, 3, 1, 2, 3})

	l3 = NewLinkedList()
	l3.PushHeadCol(l1)
	checkLinkedList(t, l3, []T{1, 2, 3})
	l3.PushHeadCol(l3)
	checkLinkedList(t, l3, []T{1, 2, 3, 1, 2, 3})

	l3 = NewLinkedList()
	l1.PushTailCol(l3)
	checkLinkedList(t, l1, []T{1, 2, 3})
	l1.PushHeadCol(l3)
	checkLinkedList(t, l1, []T{1, 2, 3})

	l1.Clear()
	l2.Clear()
	l3.Clear()
	l1.PushTails(1, 2, 3)
	checkLinkedList(t, l1, []T{1, 2, 3})
	l2.PushTails(4, 5)
	checkLinkedList(t, l2, []T{4, 5})
	l3.PushTailCol(l1)
	checkLinkedList(t, l3, []T{1, 2, 3})
	l3.PushTails(4, 5)
	checkLinkedList(t, l3, []T{1, 2, 3, 4, 5})
	l3.PushHeads(4, 5)
	checkLinkedList(t, l3, []T{4, 5, 1, 2, 3, 4, 5})
}

func TestLinkedListQueue(t *testing.T) {
	q := NewLinkedList()

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

func TestLinkedListDeque(t *testing.T) {
	q := NewLinkedList()

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

func TestLinkedListString(t *testing.T) {
	e := "[1,3,2]"
	a := NewLinkedList(1, 3, 2).String()
	if a != e {
		t.Errorf(`fmt.Sprintf("%%s", NewLinkedList(1, 3, 2)) = %v, want %v`, a, e)
	}
}

func TestLinkedListMarshalJSON(t *testing.T) {
	cs := []struct {
		list *LinkedList
		json string
	}{
		{NewLinkedList(0, 1, "0", "1", 0.0, 1.0, true, false), `[0,1,"0","1",0,1,true,false]`},
		{NewLinkedList(0, "1", 2.0, []int{1, 2}, map[int]int{1: 10, 2: 20}), `[0,"1",2,[1,2],{"1":10,"2":20}]`},
	}

	for i, c := range cs {
		bs, err := json.Marshal(c.list)
		if err != nil {
			t.Errorf("[%d] json.Marshal(%v) error: %v", i, c.list, err)
		}

		a := string(bs)
		if a != c.json {
			t.Errorf("[%d] json.Marshal(%v) = %q, want %q", i, c.list, a, c.list)
		}
	}
}

func TestLinkedListUnmarshalJSON(t *testing.T) {
	type Case struct {
		json string
		list *LinkedList
	}

	cs := []Case{
		{`["0","1",0,1,true,false]`, NewLinkedList("0", "1", 0.0, 1.0, true, false)},
		{`["1",2,[1,2],{"1":10,"2":20}]`, NewLinkedList("1", 2.0, JSONArray{1.0, 2.0}, JSONObject{"1": 10.0, "2": 20.0})},
	}

	for i, c := range cs {
		a := NewLinkedList()
		err := json.Unmarshal([]byte(c.json), a)

		if err != nil {
			t.Errorf("[%d] json.Unmarshal(%v) error: %v", i, c.json, err)
		}

		if !reflect.DeepEqual(a, c.list) {
			t.Errorf("[%d] json.Unmarshal(%q) = %v, want %q", i, c.json, a, c.list)
		}
	}
}
