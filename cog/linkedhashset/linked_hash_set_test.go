package linkedhashset

import (
	"cmp"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/cog/arraylist"
)

func TestLinkedHashSetInterface(t *testing.T) {
	var _ cog.Set[int] = NewLinkedHashSet[int]()
	var _ cog.SortIF[int] = NewLinkedHashSet[int]()
	var _ cog.Sortable[int] = NewLinkedHashSet[int]()
}

func TestLinkedHashSetNew(t *testing.T) {
	lset1 := NewLinkedHashSet[int]()
	if av := lset1.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}

	lset2 := NewLinkedHashSet(1, 2)
	if av := lset2.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	if av := lset2.Get(0); av != 1 {
		t.Errorf("Got %v expected %v", av, 1)
	}
	if av := lset2.Get(1); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
}

func TestLinkedHashSetAdd(t *testing.T) {
	lset := NewLinkedHashSet[string]()
	lset.Add("a")
	lset.AddAll("b", "c")
	if av := lset.IsEmpty(); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := lset.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av := lset.Get(2); av != "c" {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func TestLinkedHashSetRemove(t *testing.T) {
	lset := NewLinkedHashSet[string]()
	lset.Add("a")
	lset.AddAll("b", "c")
	lset.DeleteAt(2)
	lset.DeleteAt(1)
	lset.DeleteAt(0)
	if av := lset.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := lset.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestLinkedHashSetRemovePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedHashSet("a")
	list.DeleteAt(1)
}

func TestLinkedHashSetGet(t *testing.T) {
	lset := NewLinkedHashSet[string]()
	lset.AddAll("a", "a")
	lset.AddAll("b", "c", "b", "c")
	if av := lset.Get(0); av != "a" {
		t.Errorf("Got %v expected %v", av, "a")
	}
	if av := lset.Get(1); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
	if av := lset.Get(2); av != "c" {
		t.Errorf("Got %v expected %v", av, "c")
	}
	lset.DeleteAt(0)
	if av := lset.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestLinkedHashSetGetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedHashSet("a")
	list.Get(1)
}

func TestLinkedHashSetSwap(t *testing.T) {
	lset := NewLinkedHashSet[string]()
	lset.AddAll("a", "a")
	lset.AddAll("b", "c", "b", "c")
	lset.Swap(0, 1)
	if av := lset.Get(0); av != "b" {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func TestLinkedHashSetClear(t *testing.T) {
	lset := NewLinkedHashSet[string]()
	lset.AddAll("e", "f", "g", "a", "b", "c", "d")
	lset.AddAll("e", "f", "g", "a", "b", "c", "d")
	lset.Clear()
	if av := lset.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := lset.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestLinkedHashSetContains(t *testing.T) {
	list := NewLinkedHashSet[int]()

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

func TestLinkedHashSetRetain(t *testing.T) {
	for n := 0; n < 100; n++ {
		a := []int{}
		list := NewLinkedHashSet[int]()
		for i := 0; i < n; i++ {
			if i&1 == 0 {
				a = append(a, i)
			}
			list.Add(i)

			list.RetainAll(a...)
			vs := list.Values()
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

func TestLinkedHashSetValues(t *testing.T) {
	lset := NewLinkedHashSet[string]()
	lset.AddAll("a", "a")
	lset.AddAll("b", "c", "b", "c")
	if av, ev := fmt.Sprintf("%v", lset.Values()), "[a b c]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestLinkedHashSetInsert(t *testing.T) {
	lset := NewLinkedHashSet[string]()
	lset.Inserts(0, "b", "c", "b", "c")
	lset.Inserts(0, "a", "a")
	if av := lset.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	lset.Insert(3, "d") // append
	if av := lset.Len(); av != 4 {
		t.Errorf("Got %v expected %v", av, 4)
	}
	if av, ev := fmt.Sprintf("%v", lset.Values()), "[a b c d]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestLinkedHashSetInsertPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedHashSet("a")
	list.Insert(2, "b")
}

func TestLinkedHashSetSet(t *testing.T) {
	lset := NewLinkedHashSet("0", "1")
	lset.Set(0, "a")
	lset.Set(1, "b")
	if av := lset.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}
	lset.Add("")
	lset.Set(2, "c") // last
	if av := lset.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	lset.Set(1, "bb") // update
	if av := lset.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av, ev := fmt.Sprintf("%v", lset.Values()), "[a bb c]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	lset.Set(2, "cc") // last to first traversal
	lset.Set(0, "aa") // first to last traversal
	if av, ev := fmt.Sprintf("%v", lset.Values()), "[aa bb cc]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestLinkedHashSetSetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("want out of bounds panic")
		}
	}()

	list := NewLinkedHashSet("a")
	list.Set(1, "b")
}

func TestLinkedHashSetEach(t *testing.T) {
	lset := NewLinkedHashSet[string]()
	lset.AddAll("a", "b", "c")
	lset.AddAll("a", "b", "c")
	lset.Each(func(index int, value string) bool {
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
		return true
	})
}

func TestLinkedHashSetSeq(t *testing.T) {
	ehs := NewLinkedHashSet("a", "b", "c", "a", "b", "c")
	ahs := NewLinkedHashSet[string]()
	for s := range ehs.Seq() {
		ahs.Add(s)
	}

	w := fmt.Sprint(ehs.Values())
	a := fmt.Sprint(ahs.Values())
	if a != w {
		t.Errorf("Each():\nWANT: %s\n GOT: %s", w, a)
	}
}

func TestLinkedHashSetIteratorNextOnEmpty(t *testing.T) {
	lset := NewLinkedHashSet[int]()
	it := lset.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty lset")
	}
}

func TestLinkedHashSetIteratorNext(t *testing.T) {
	lset := NewLinkedHashSet[string]()
	lset.AddAll("a", "b", "c")
	lset.AddAll("a", "b", "c")
	it := lset.Iterator()
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

func TestLinkedHashSetIteratorPrevOnEmpty(t *testing.T) {
	lset := NewLinkedHashSet[int]()
	it := lset.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty lset")
	}
}

func TestLinkedHashSetIteratorPrev(t *testing.T) {
	lset := NewLinkedHashSet[string]()
	lset.AddAll("a", "b", "c")
	lset.AddAll("a", "b", "c")
	it := lset.Iterator()
	count := 0
	index := lset.Len()
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

func TestLinkedHashSetIteratorReset(t *testing.T) {
	lset := NewLinkedHashSet[string]()
	it := lset.Iterator()
	lset.AddAll("a", "b", "c")
	lset.AddAll("a", "b", "c")
	for it.Next() {
	}
	it.Reset()
	it.Next()
	if value := it.Value(); value != "a" {
		t.Errorf("Got %v expected %v", value, "a")
	}

	it.Reset()
	it.Prev()
	if value := it.Value(); value != "c" {
		t.Errorf("Got %v expected %v", value, "c")
	}
}

func assertLinkedHashSetIteratorRemove(t *testing.T, i int, it cog.Iterator[int], w *LinkedHashSet[int]) int {
	it.Remove()

	v := it.Value()
	w.Remove(v)

	it.SetValue(9999)

	lset := it.(*linkedHashSetIterator[int]).lset
	if lset.Contains(v) {
		t.Fatalf("[%d] lset.Contains(%v) = true", i, v)
	}

	if lset.String() != w.String() {
		t.Fatalf("[%d] (%v) %v != %v", i, v, lset.String(), w.String())
	}

	return v
}

func TestLinkedHashSetIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		lset := NewLinkedHashSet[int]()
		wset := NewLinkedHashSet[int]()

		for n := 0; n < i; n++ {
			lset.Add(n)
			wset.Add(n)
		}

		it := lset.Iterator()

		it.Remove()
		if lset.Len() != i {
			t.Fatalf("[%d] lset.Len() == %v, want %v", i, lset.Len(), i)
		}

		// remove middle
		for j := 0; j <= lset.Len()/2; j++ {
			it.Next()
		}

		v := assertLinkedHashSetIteratorRemove(t, i, it, wset)

		it.Next()
		if v+1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v+1)
		}
		assertLinkedHashSetIteratorRemove(t, i, it, wset)

		it.Prev()
		if v-1 != it.Value() {
			t.Fatalf("[%d] it.Value() = %v, want %v", i, it.Value(), v-1)
		}
		assertLinkedHashSetIteratorRemove(t, i, it, wset)

		// remove first
		for it.Prev() {
		}
		assertLinkedHashSetIteratorRemove(t, i, it, wset)

		// remove last
		for it.Next() {
		}
		assertLinkedHashSetIteratorRemove(t, i, it, wset)

		// remove all
		it.Reset()
		if i%2 == 0 {
			for it.Prev() {
				assertLinkedHashSetIteratorRemove(t, i, it, wset)
			}
		} else {
			for it.Next() {
				assertLinkedHashSetIteratorRemove(t, i, it, wset)
			}
		}
		if !lset.IsEmpty() {
			t.Fatalf("[%d] lset.IsEmpty() = true", i)
		}
	}
}

func TestLinkedHashSetIteratorSetValue(t *testing.T) {
	ls := NewLinkedHashSet[int]()
	for i := 1; i <= 100; i++ {
		ls.Add(i)
	}

	// forward (1->2, 3->4, ... )
	for it := ls.Iterator(); it.Next(); {
		v := it.Value() + 1
		it.SetValue(v)
	}
	for i := 1; i <= ls.Len(); i++ {
		v := ls.Get(i - 1)
		w := i * 2
		if v != w {
			t.Fatalf("Set[%d] = %v, want %v", i-1, v, w)
		}
	}

	// backward (100 -> 98, 96 -> 94)
	for it := ls.Iterator(); it.Prev(); {
		v := it.Value() - 2
		it.SetValue(v)
	}
	for i := 1; i <= ls.Len(); i++ {
		v := ls.Get(i - 1)
		w := i*4 - 2
		if v != w {
			t.Fatalf("Set[%d] = %v, want %v", i-1, v, w)
		}
	}
}

func TestLinkedHashSetSort(t *testing.T) {
	for i := 1; i < 20; i++ {
		lset := NewLinkedHashSet[int]()

		a := make([]int, 0, 20)
		for n := i; n < 20; n++ {
			v := rand.Intn(1000)
			if !asg.Contains(a, v) {
				a = append(a, v)
			}
		}

		for j := len(a) - 1; j >= 0; j-- {
			lset.Add(a[j])
		}

		lset.Sort(cmp.Less[int])
		sort.Sort(inta(a))

		if !reflect.DeepEqual(a, lset.Values()) {
			t.Fatalf("%v != %v", a, lset.Values())
		}
	}
}

func checkLinkedHashSetLen[T comparable](t *testing.T, lset *LinkedHashSet[T], len int) bool {
	if n := lset.Len(); n != len {
		t.Errorf("lset.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkLinkedHashSet(t *testing.T, lset *LinkedHashSet[int], evs []int) {
	if !checkLinkedHashSetLen(t, lset, len(evs)) {
		return
	}

	for i, it := 0, lset.Iterator(); it.Next(); i++ {
		v := it.Value()
		if v != evs[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, v, evs[i])
		}
	}

	avs := lset.Values()
	for i, v := range avs {
		if v != evs[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, v, evs[i])
		}
	}
}

func TestLinkedHashSetExtending(t *testing.T) {
	l1 := NewLinkedHashSet(1, 2, 3)
	l2 := NewLinkedHashSet[int]()
	l2.PushTail(4)
	l2.PushTail(5)

	l3 := NewLinkedHashSet[int]()
	l3.PushTailCol(l1)
	checkLinkedHashSet(t, l3, []int{1, 2, 3})
	l3.PushTailCol(l2)
	checkLinkedHashSet(t, l3, []int{1, 2, 3, 4, 5})

	l3 = NewLinkedHashSet[int]()
	l3.PushHeadCol(l2)
	checkLinkedHashSet(t, l3, []int{4, 5})
	l3.PushHeadCol(l1)
	checkLinkedHashSet(t, l3, []int{1, 2, 3, 4, 5})

	checkLinkedHashSet(t, l1, []int{1, 2, 3})
	checkLinkedHashSet(t, l2, []int{4, 5})

	l3 = NewLinkedHashSet[int]()
	l3.PushTailCol(l1)
	checkLinkedHashSet(t, l3, []int{1, 2, 3})
	l3.PushTailCol(l3)
	checkLinkedHashSet(t, l3, []int{1, 2, 3})

	l3 = NewLinkedHashSet[int]()
	l3.PushHeadCol(l1)
	checkLinkedHashSet(t, l3, []int{1, 2, 3})
	l3.PushHeadCol(l3)
	checkLinkedHashSet(t, l3, []int{1, 2, 3})

	l3 = NewLinkedHashSet[int]()
	l1.PushTailCol(l3)
	checkLinkedHashSet(t, l1, []int{1, 2, 3})
	l1.PushHeadCol(l3)
	checkLinkedHashSet(t, l1, []int{1, 2, 3})

	l1.Clear()
	l2.Clear()
	l3.Clear()
	l1.PushTails(1, 2, 3)
	checkLinkedHashSet(t, l1, []int{1, 2, 3})
	l2.PushTails(4, 5)
	checkLinkedHashSet(t, l2, []int{4, 5})
	l3.PushTailCol(l1)
	checkLinkedHashSet(t, l3, []int{1, 2, 3})
	l3.PushTails(4, 5)
	checkLinkedHashSet(t, l3, []int{1, 2, 3, 4, 5})
	l3.PushHeads(4, 5)
	checkLinkedHashSet(t, l3, []int{1, 2, 3, 4, 5})
}

func TestLinkedHashSetDelete(t *testing.T) {
	lset := NewLinkedHashSet[int]()

	for i := 1; i <= 100; i++ {
		lset.PushTail(i)
	}

	lset.RemoveFunc(func(d int) bool {
		return d == 101
	})
	if lset.Len() != 100 {
		t.Error("LinkedHashSet.Remove(101) should do nothing")
	}
	for i := 1; i <= 100; i++ {
		lset.Remove(i)
		lset.RemoveAll(i, i)
		if lset.Len() != 100-i {
			t.Errorf("LinkedHashSet.Remove(%v) failed, lset.Len() = %v, want %v", i, lset.Len(), 100-i)
		}
	}

	if !lset.IsEmpty() {
		t.Error("LinkedHashSet.IsEmpty() should return true")
	}
}

func TestLinkedHashSetDelete2(t *testing.T) {
	lset := NewLinkedHashSet[int]()

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			lset.PushTail(i)
		}
	}

	n := lset.Len()
	lset.Remove(100)
	if lset.Len() != n {
		t.Errorf("LinkedHashSet.Remove(100).Len() = %v, want %v", lset.Len(), n)
	}

	for i := 0; i < 100; i++ {
		n = lset.Len()
		lset.Remove(i)
		z := i % 10
		if z > 1 {
			z = 1
		}
		a := n - lset.Len()
		if a != z {
			t.Errorf("LinkedHashSet.Remove(%v) = %v, want %v", i, a, z)
		}
	}

	if !lset.IsEmpty() {
		t.Error("LinkedHashSet.IsEmpty() should return true")
	}
}

func TestLinkedHashSetDeleteAll(t *testing.T) {
	lset := NewLinkedHashSet[int]()

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			lset.PushTail(i)
		}
	}

	n := lset.Len()
	lset.Remove(100)
	if lset.Len() != n {
		t.Errorf("LinkedHashSet.Remove(100).Len() = %v, want %v", lset.Len(), n)
	}

	for i := 0; i < 100; i++ {
		n = lset.Len()
		lset.RemoveCol(arraylist.NewArrayList(i, i))
		z := i % 10
		if z > 1 {
			z = 1
		}
		a := n - lset.Len()
		if a != z {
			t.Errorf("LinkedHashSet.Remove(%v) = %v, want %v", i, a, z)
		}
	}

	if !lset.IsEmpty() {
		t.Error("LinkedHashSet.IsEmpty() should return true")
	}
}

func TestLinkedHashSetQueue(t *testing.T) {
	q := NewLinkedHashSet[int]()

	if _, ok := q.Peek(); ok {
		t.Error("should return false when peeking empty queue")
	}

	for i := 0; i < 100; i++ {
		q.Push(i)
	}

	for i := 0; i < 100; i++ {
		v, _ := q.Peek()
		if v != i {
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

func TestLinkedHashSetDeque(t *testing.T) {
	q := NewLinkedHashSet[int]()

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

func TestLinkedHashSetString(t *testing.T) {
	e := "[1,3,2]"
	a := NewLinkedHashSet(1, 3, 2).String()
	if a != e {
		t.Errorf(`fmt.Sprintf("%%s", NewLinkedHashSet(1, 3, 2)) = %v, want %v`, a, e)
	}
}

func TestLinkedHashSetMarshalJSON(t *testing.T) {
	cs := []struct {
		lset *LinkedHashSet[int]
		json string
	}{
		{NewLinkedHashSet(0, 1, 3, 2, 1, 0), `[0,1,3,2]`},
	}

	for i, c := range cs {
		bs, err := json.Marshal(c.lset)
		if err != nil {
			t.Errorf("[%d] json.Marshal(%v) error: %v", i, c.lset, err)
		}

		a := string(bs)
		if a != c.json {
			t.Errorf("[%d] json.Marshal(%v) = %q, want %q", i, c.lset, a, c.lset)
		}
	}
}

func TestLinkedHashSetUnmarshalJSON(t *testing.T) {
	type Case struct {
		json string
		lset *LinkedHashSet[int]
	}

	cs := []Case{
		{`[0,1,3,2,1,0,1]`, NewLinkedHashSet(0, 1, 3, 2)},
	}

	for i, c := range cs {
		a := NewLinkedHashSet[int]()
		err := json.Unmarshal([]byte(c.json), a)

		if err != nil {
			t.Errorf("[%d] json.Unmarshal(%v) error: %v", i, c.json, err)
		}

		if !reflect.DeepEqual(a, c.lset) {
			t.Errorf("[%d] json.Unmarshal(%q) = %v, want %q", i, c.json, a, c.lset)
		}
	}
}
