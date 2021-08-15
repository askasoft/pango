package col

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/pandafw/pango/ars"
	"github.com/pandafw/pango/cmp"
	"github.com/pandafw/pango/str"
)

func TestLinkedHashSetInterface(t *testing.T) {
	var s Set = NewLinkedHashSet()
	if s == nil {
		t.Error("LinkedHashSet is not a Set")
	}
}

func TestLinkedHashSetNew(t *testing.T) {
	lset1 := NewLinkedHashSet()

	if av := lset1.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}

	lset2 := NewLinkedHashSet(1, "b")

	if av := lset2.Len(); av != 2 {
		t.Errorf("Got %v expected %v", av, 2)
	}

	if av, ok := lset2.Get(0); av != 1 || !ok {
		t.Errorf("Got %v expected %v", av, 1)
	}

	if av, ok := lset2.Get(1); av != "b" || !ok {
		t.Errorf("Got %v expected %v", av, "b")
	}

	if av, ok := lset2.Get(2); av != nil || ok {
		t.Errorf("Got %v expected %v", av, nil)
	}
}

func TestLinkedHashSetAdd(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a")
	lset.Add("b", "c")
	if av := lset.IsEmpty(); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := lset.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av, ok := lset.Get(2); av != "c" || !ok {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func TestLinkedHashSetRemove(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a")
	lset.Add("b", "c")
	lset.Remove(2)
	if av, ok := lset.Get(2); av != nil || ok {
		t.Errorf("Got %v expected %v", av, nil)
	}
	lset.Remove(1)
	lset.Remove(0)
	lset.Remove(0) // no effect
	if av := lset.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := lset.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestLinkedHashSetGet(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a")
	lset.Add("b", "c")
	if av, ok := lset.Get(0); av != "a" || !ok {
		t.Errorf("Got %v expected %v", av, "a")
	}
	if av, ok := lset.Get(1); av != "b" || !ok {
		t.Errorf("Got %v expected %v", av, "b")
	}
	if av, ok := lset.Get(2); av != "c" || !ok {
		t.Errorf("Got %v expected %v", av, "c")
	}
	if av, ok := lset.Get(3); av != nil || ok {
		t.Errorf("Got %v expected %v", av, nil)
	}
	lset.Remove(0)
	if av, ok := lset.Get(0); av != "b" || !ok {
		t.Errorf("Got %v expected %v", av, "b")
	}
}

func TestLinkedHashSetSwap(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a")
	lset.Add("b", "c")
	lset.Swap(0, 1)
	if av, ok := lset.Get(0); av != "b" || !ok {
		t.Errorf("Got %v expected %v", av, "c")
	}
}

func TestLinkedHashSetClear(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("e", "f", "g", "a", "b", "c", "d")
	lset.Clear()
	if av := lset.IsEmpty(); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := lset.Len(); av != 0 {
		t.Errorf("Got %v expected %v", av, 0)
	}
}

func TestLinkedHashSetContains(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a")
	lset.Add("b", "c")
	if av := lset.Contains("a"); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := lset.Contains("a", "b", "c"); av != true {
		t.Errorf("Got %v expected %v", av, true)
	}
	if av := lset.Contains("a", "b", "c", "d"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	lset.Clear()
	if av := lset.Contains("a"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
	if av := lset.Contains("a", "b", "c"); av != false {
		t.Errorf("Got %v expected %v", av, false)
	}
}

func TestLinkedHashSetValues(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a")
	lset.Add("b", "c")
	if av, ev := fmt.Sprintf("%s%s%s", lset.Values()...), "abc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestLinkedHashSetInsert(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Insert(0, "b", "c")
	lset.Insert(0, "a")
	lset.Insert(10, "x") // ignore
	if av := lset.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	lset.Insert(3, "d") // append
	if av := lset.Len(); av != 4 {
		t.Errorf("Got %v expected %v", av, 4)
	}
	if av, ev := fmt.Sprintf("%s%s%s%s", lset.Values()...), "abcd"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
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
	lset.Set(4, "d")  // ignore
	lset.Set(1, "bb") // update
	if av := lset.Len(); av != 3 {
		t.Errorf("Got %v expected %v", av, 3)
	}
	if av, ev := fmt.Sprintf("%s%s%s", lset.Values()...), "abbc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	lset.Set(2, "cc") // last to first traversal
	lset.Set(0, "aa") // first to last traversal
	if av, ev := fmt.Sprintf("%s%s%s", lset.Values()...), "aabbcc"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestLinkedHashSetEach(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a", "b", "c")
	index := 0
	lset.Each(func(value interface{}) {
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

func TestLinkedHashSetIteratorNextOnEmpty(t *testing.T) {
	lset := NewLinkedHashSet()
	it := lset.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty lset")
	}
}

func TestLinkedHashSetIteratorNext(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a", "b", "c")
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
	lset := NewLinkedHashSet()
	it := lset.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty lset")
	}
}

func TestLinkedHashSetIteratorPrev(t *testing.T) {
	lset := NewLinkedHashSet()
	lset.Add("a", "b", "c")
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
	lset := NewLinkedHashSet()
	it := lset.Iterator()
	lset.Add("a", "b", "c")
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

func TestLinkedHashSetIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		lset := NewLinkedHashSet()

		for n := 0; n < i; n++ {
			lset.Add(n)
		}

		it := lset.Iterator()

		it.Remove()
		if lset.Len() != i {
			t.Errorf("[%d] lset.Len() == %v, want %v", i, lset.Len(), i)
		}

		// remove middle
		x := rand.Intn(i-4) + 1
		for j := 0; j <= x; j++ {
			it.Next()
		}

		v := it.Value().(int)
		it.Remove()
		if lset.Len() != i-1 {
			t.Errorf("[%d] lset.Len() == %v, want %v", i, lset.Len(), i-1)
		}
		if lset.Contains(v) {
			t.Errorf("[%d] lset.Contains(%v) = true", i, v)
		}

		it.Next()
		if v+1 != it.Value() {
			t.Errorf("[%d] it.Value() = %v, want %v", i, it.Value(), v+1)
		}
		it.Remove()
		if lset.Contains(v + 1) {
			t.Errorf("[%d] lset.Contains(%v) = true", i, v+1)
		}

		it.Prev()
		if v-1 != it.Value() {
			t.Errorf("[%d] it.Value() = %v, want %v", i, it.Value(), v-1)
		}
		it.Remove()
		if lset.Contains(v - 1) {
			t.Errorf("[%d] lset.Contains(%v) = true", i, v-1)
		}

		// remove first
		for it.Prev() {
		}
		it.Remove()
		if lset.Contains(0) {
			t.Errorf("[%d] lset.Contains(%v) = true", i, 0)
		}
		if it.Prev() {
			t.Errorf("[%d] lset.Prev() = true", i)
		}

		// remove last
		for it.Next() {
		}
		it.Remove()
		if lset.Contains(i - 1) {
			t.Errorf("[%d] lset.Contains(%v) = true", i, i-1)
		}
		if it.Next() {
			t.Errorf("[%d] lset.Next() = true", i)
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
		if !lset.IsEmpty() {
			t.Errorf("[%d] lset.IsEmpty() = true", i)
		}
	}
}

func TestLinkedHashSetLazyInit(t *testing.T) {
	{
		lset := &LinkedHashSet{}
		if lset.Len() != 0 {
			t.Error("lset.Len() != 0")
		}
		if !lset.IsEmpty() {
			t.Error("lset.IsEmpty() = true")
		}
		if len(lset.Values()) != 0 {
			t.Error("len(lset.Values()) != 0")
		}
		if lset.Contains(1) {
			t.Error("lset.Contains(1) = true, want false")
		}
		if lset.Len() != 0 {
			t.Error("lset.Len() != 0")
		}
		if lset.Front() != nil {
			t.Error("lset.Front() != nil")
		}
		if lset.Back() != nil {
			t.Error("lset.Back() != nil")
		}
		if i := lset.Search(1); i != nil {
			t.Errorf("lset.Search(1) == %v, want %v", i, nil)
		}
	}
	{
		lset := &LinkedHashSet{}
		lset.PushFront(1)
		if lset.Len() != 1 {
			t.Errorf("lset.Len() = %v, want 1", lset.Len())
		}
	}
	{
		lset := &LinkedHashSet{}
		lset.PushBack(1)
		if lset.Len() != 1 {
			t.Errorf("lset.Len() = %v, want 1", lset.Len())
		}
	}
	{
		lset := &LinkedHashSet{}
		lset.Delete(1)
		if lset.Len() != 0 {
			t.Error("lset.Len() != 0")
		}
	}
	{
		lset := &LinkedHashSet{}
		lset.Each(func(v interface{}) {})
		lset.ReverseEach(func(v interface{}) {})
	}
}

func TestLinkedHashSetSort(t *testing.T) {
	for i := 1; i < 20; i++ {
		lset := NewLinkedHashSet()

		a := make([]interface{}, 0, 20)
		for n := i; n < 20; n++ {
			v := rand.Intn(1000)
			if !ars.Contains(a, v) {
				a = append(a)
			}
		}

		for j := len(a) - 1; j >= 0; j-- {
			lset.Add(a[j])
		}

		lset.Sort(cmp.LessInt)
		sort.Sort(inta(a))

		if !reflect.DeepEqual(a, lset.Values()) {
			t.Fatalf("%v != %v", a, lset.Values())
		}
	}
}

func checkLinkedHashSetLen(t *testing.T, lset *LinkedHashSet, len int) bool {
	if n := lset.Len(); n != len {
		t.Errorf("lset.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkLinkedHashSetPointers(t *testing.T, lset *LinkedHashSet, lis []*LinkedSetItem) {
	root := &lset.root

	if !checkLinkedHashSetLen(t, lset, len(lis)) {
		return
	}

	// zero length lists must be the zero value or properly initialized (sentinel circle)
	if len(lis) == 0 {
		if lset.root.next != nil && lset.root.next != root || lset.root.prev != nil && lset.root.prev != root {
			t.Errorf("lset.root.next = %p, lset.root.prev = %p; both should both be nil or %p", lset.root.next, lset.root.prev, root)
		}
		return
	}
	// len(lis) > 0

	// check internal and external prev/next connections
	for i, e := range lis {
		prev := root
		var Prev *LinkedSetItem
		if i > 0 {
			prev = lis[i-1]
			Prev = prev
		}
		if p := e.prev; p != prev {
			t.Errorf("elt[%d](%p).prev = %p, want %p", i, e, p, prev)
		}
		if p := e.Prev(); p != Prev {
			t.Errorf("elt[%d](%p).Prev() = %p, want %p", i, e, p, Prev)
		}

		next := root
		var Next *LinkedSetItem
		if i < len(lis)-1 {
			next = lis[i+1]
			Next = next
		}
		if n := e.next; n != next {
			t.Errorf("elt[%d](%p).next = %p, want %p", i, e, n, next)
		}
		if n := e.Next(); n != Next {
			t.Errorf("elt[%d](%p).Next() = %p, want %p", i, e, n, Next)
		}
	}
}

func TestLinkedHashSetBasic(t *testing.T) {
	lset := NewLinkedHashSet()
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{})

	// Single item lset
	e := lset.PushFront("a")
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e})
	lset.MoveToFront(e)
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e})
	lset.MoveToBack(e)
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e})
	e.Remove()
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{})

	// Bigger lset
	e2 := lset.PushFront(2)
	e1 := lset.PushFront(1)
	e3 := lset.PushBack(3)
	e4 := lset.PushBack("banana")
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e2, e3, e4})

	e2.Remove()
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e3, e4})

	lset.MoveToFront(e3) // move from middle
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e3, e1, e4})

	lset.MoveToFront(e1)
	lset.MoveToBack(e3) // move from middle
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e4, e3})

	lset.MoveToFront(e3) // move from back
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e3, e1, e4})
	lset.MoveToFront(e3) // should be no-op
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e3, e1, e4})

	lset.MoveToBack(e3) // move from front
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e4, e3})
	lset.MoveToBack(e3) // should be no-op
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e4, e3})

	e2 = lset.InsertBefore(e1, 2) // insert before front
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e2, e1, e4, e3})
	e2.Remove()

	e2 = lset.InsertBefore(e4, 2) // insert before middle
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e2, e4, e3})
	e2.Remove()

	e2 = lset.InsertBefore(e3, 2) // insert before back
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e4, e2, e3})
	e2.Remove()

	e2 = lset.InsertAfter(e1, 2) // insert after front
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e2, e4, e3})
	e2.Remove()

	e2 = lset.InsertAfter(e4, 2) // insert after middle
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e4, e2, e3})
	e2.Remove()

	e2 = lset.InsertAfter(e3, 2) // insert after back
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e4, e3, e2})
	e2.Remove()

	// Check standard iteration.
	sum := 0
	for e := lset.Front(); e != nil; e = e.Next() {
		if i, ok := e.Value().(int); ok {
			sum += i
		}
	}
	if sum != 4 {
		t.Errorf("sum over lset = %d, want 4", sum)
	}

	// Clear all items by iterating
	var next *LinkedSetItem
	for e := lset.Front(); e != nil; e = next {
		next = e.Next()
		e.Remove()
	}
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{})
}

func checkLinkedHashSet(t *testing.T, lset *LinkedHashSet, evs []interface{}) {
	if !checkLinkedHashSetLen(t, lset, len(evs)) {
		return
	}

	for i, e := 0, lset.Front(); e != nil; i, e = i+1, e.Next() {
		v := e.Value().(int)
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
	l2 := NewLinkedHashSet()
	l2.PushBack(4)
	l2.PushBack(5)

	l3 := NewLinkedHashSet()
	l3.PushBackAll(l1)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3})
	l3.PushBackAll(l2)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3, 4, 5})

	l3 = NewLinkedHashSet()
	l3.PushFrontAll(l2)
	checkLinkedHashSet(t, l3, []interface{}{4, 5})
	l3.PushFrontAll(l1)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3, 4, 5})

	checkLinkedHashSet(t, l1, []interface{}{1, 2, 3})
	checkLinkedHashSet(t, l2, []interface{}{4, 5})

	l3 = NewLinkedHashSet()
	l3.PushBackAll(l1)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3})
	l3.PushBackAll(l3)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3})

	l3 = NewLinkedHashSet()
	l3.PushFrontAll(l1)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3})
	l3.PushFrontAll(l3)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3})

	l3 = NewLinkedHashSet()
	l1.PushBackAll(l3)
	checkLinkedHashSet(t, l1, []interface{}{1, 2, 3})
	l1.PushFrontAll(l3)
	checkLinkedHashSet(t, l1, []interface{}{1, 2, 3})

	l1.Clear()
	l2.Clear()
	l3.Clear()
	l1.PushBack(1, 2, 3)
	checkLinkedHashSet(t, l1, []interface{}{1, 2, 3})
	l2.PushBack(4, 5)
	checkLinkedHashSet(t, l2, []interface{}{4, 5})
	l3.PushBackAll(l1)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3})
	l3.PushBack(4, 5)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3, 4, 5})
	l3.PushFront(4, 5)
	checkLinkedHashSet(t, l3, []interface{}{1, 2, 3, 4, 5})
}

func TestLinkedHashSetRemove2(t *testing.T) {
	lset := NewLinkedHashSet()
	e1 := lset.PushBack(1)
	e2 := lset.PushBack(2)
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e2})
	e := lset.Front()
	e.Remove()
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e2})
	e.Remove()
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e2})
}

func TestLinkedHashSetInsertBefore(t *testing.T) {
	l1 := NewLinkedHashSet()
	l1.PushBack(1)
	l1.PushBack(2)

	l2 := NewLinkedHashSet()
	l2.PushBack(3)
	l2.PushBack(4)

	e := l1.Front()
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(e, 8)
	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestLinkedHashSetRemove1(t *testing.T) {
	lset := NewLinkedHashSet(1, 2, 3)

	e := lset.Front().Next()
	e.Remove()
	if e.Value() != 2 {
		t.Errorf("e.value = %d, want 2", e.Value())
	}
	if e.Next() != lset.Back() {
		t.Errorf("e.Next() != lset.Back()")
	}
	if e.Prev() != lset.Front() {
		t.Errorf("e.Prev() != lset.Front()")
	}
}

func TestLinkedHashSetMove(t *testing.T) {
	lset := NewLinkedHashSet()
	e1 := lset.PushBack(1)
	e2 := lset.PushBack(2)
	e3 := lset.PushBack(3)
	e4 := lset.PushBack(4)

	lset.MoveAfter(e3, e3)
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e2, e3, e4})
	lset.MoveBefore(e2, e2)
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e2, e3, e4})

	lset.MoveAfter(e2, e3)
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e2, e3, e4})
	lset.MoveBefore(e2, e2)
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e2, e3, e4})

	lset.MoveBefore(e4, e2)
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e3, e2, e4})
	e2, e3 = e3, e2

	lset.MoveBefore(e1, e4)
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e4, e1, e2, e3})
	e1, e2, e3, e4 = e4, e1, e2, e3

	lset.MoveAfter(e1, e4)
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e4, e2, e3})
	e2, e3, e4 = e4, e2, e3

	lset.MoveAfter(e3, e2)
	checkLinkedHashSetPointers(t, lset, []*LinkedSetItem{e1, e3, e2, e4})
	e2, e3 = e3, e2
}

// Test that a lset lset is not modified when calling InsertBefore with a mark that is not an item of lset.
func TestLinkedHashSetInsertBeforeUnknownMark(t *testing.T) {
	lset := NewLinkedHashSet(1, 2, 3)
	lset.InsertBefore(new(LinkedSetItem), 1)
	checkLinkedHashSet(t, lset, []interface{}{1, 2, 3})
}

// Test that a lset lset is not modified when calling InsertAfter with a mark that is not an item of lset.
func TestLinkedHashSetInsertAfterUnknownMark(t *testing.T) {
	lset := NewLinkedHashSet(1, 2, 3)
	lset.InsertAfter(new(LinkedSetItem), 1)
	checkLinkedHashSet(t, lset, []interface{}{1, 2, 3})
}

// Test that a lset lset is not modified when calling MoveAfter or MoveBefore with a mark that is not an item of lset.
func TestLinkedHashSetMoveUnknownMark(t *testing.T) {
	l1 := NewLinkedHashSet()
	e1 := l1.PushBack(1)

	l2 := NewLinkedHashSet()
	e2 := l2.PushBack(2)

	l1.MoveAfter(e2, e1)
	checkLinkedHashSet(t, l1, []interface{}{1})
	checkLinkedHashSet(t, l2, []interface{}{2})

	l1.MoveBefore(e2, e1)
	checkLinkedHashSet(t, l1, []interface{}{1})
	checkLinkedHashSet(t, l2, []interface{}{2})
}

func TestLinkedHashSetContains2(t *testing.T) {
	lset := NewLinkedHashSet(1, 11, 111, "1", "11", "111")

	n := (100+1)/101 + 110

	if !lset.Contains(n) {
		t.Errorf("LinkedHashSet [%v] should contains %v", lset, n)
	}

	n++
	if lset.Contains(n) {
		t.Errorf("LinkedHashSet [%v] should not contains %v", lset, n)
	}

	s := str.Repeat("1", 3)

	if !lset.Contains(s) {
		t.Errorf("LinkedHashSet [%v] should contains %v", lset, s)
	}

	s += "0"
	if lset.Contains(s) {
		t.Errorf("LinkedHashSet [%v] should not contains %v", lset, s)
	}
}

func TestLinkedHashSetItem(t *testing.T) {
	lset := NewLinkedHashSet()

	for i := 0; i < 100; i++ {
		lset.PushBack(i)
	}

	if lset.Item(100) != nil {
		t.Error("LinkedHashSet.Item(100) != nil")
	}
	if lset.Item(-101) != nil {
		t.Error("LinkedHashSet.Item(-101) != nil")
	}

	for n := 0; n < 100; n++ {
		var a interface{}

		w := n
		s := n
		i := lset.Item(s)
		if i != nil {
			a = i.Value()
		}
		if i == nil || a != w {
			t.Errorf("LinkedHashSet.Item(%v).Value = %v, want %v", s, a, w)
		}

		s = n - 100
		i = lset.Item(s)
		if i != nil {
			a = i.Value()
		}
		if i == nil || a != w {
			t.Errorf("LinkedHashSet.Item(%v).Value = %v, want %v", s, a, w)
		}
	}
}

func TestLinkedHashSetSwapItem(t *testing.T) {
	lset := NewLinkedHashSet()

	for i := 0; i < 100; i++ {
		lset.PushBack(i)
	}

	if lset.SwapItem(lset.Item(0), lset.Front()) {
		t.Error("lset.SwapItem(lset.Item(0), lset.Front()) = true, want false")
	}

	for i := 0; i < 50; i++ {
		ia := lset.Item(i)
		ib := lset.Item(-i - 1)
		a := lset.SwapItem(ia, ib)
		if !a {
			t.Errorf("LinkedHashSet.SwapItem(%v, %v) = %v, want %v", ia, ib, a, true)
		}
	}

	for i := 0; i < 100; i++ {
		e := 100 - 1 - i
		a := lset.Item(i).Value()
		if a != e {
			t.Errorf("LinkedHashSet.Item(%v).Value = %v, want %v", i, a, e)
		}
	}
}

func TestLinkedHashSetSearch(t *testing.T) {
	lset := NewLinkedHashSet(1, 11)

	n111 := lset.PushBack(111)

	lset.PushBack("1", "11")
	s111 := lset.PushBack("111")

	n := (100+1)/101 + 110

	se := lset.Search(n)
	if n111 != se {
		t.Errorf("LinkedHashSet [%v] should contains %v", lset, n)
	}

	n++
	se = lset.Search(n)
	if se != nil {
		t.Errorf("LinkedHashSet [%v] should not contains %v", lset, n)
	}

	s := str.Repeat("1", 3)

	se = lset.Search(s)
	if s111 != se {
		t.Errorf("LinkedHashSet [%v] should contains %v", lset, s)
	}

	s += "0"
	se = lset.Search(s)
	if se != nil {
		t.Errorf("LinkedHashSet [%v] should not contains %v", lset, s)
	}
}

func TestLinkedHashSetDelete(t *testing.T) {
	lset := NewLinkedHashSet()

	for i := 1; i <= 100; i++ {
		lset.PushBack(i)
	}

	lset.Delete(101)
	if lset.Len() != 100 {
		t.Error("LinkedHashSet.Delete(101) should do nothing")
	}
	for i := 1; i <= 100; i++ {
		lset.Delete(i)
		if lset.Len() != 100-i {
			t.Errorf("LinkedHashSet.Delete(%v) failed, lset.Len() = %v, want %v", i, lset.Len(), 100-i)
		}
	}

	if !lset.IsEmpty() {
		t.Error("LinkedHashSet.IsEmpty() should return true")
	}
}

func TestLinkedHashSetDeleteAll(t *testing.T) {
	lset := NewLinkedHashSet()

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			lset.PushBack(i)
		}
	}

	n := lset.Len()
	lset.Delete(100)
	if lset.Len() != n {
		t.Errorf("LinkedHashSet.Delete(100).Len() = %v, want %v", lset.Len(), n)
	}

	for i := 0; i < 100; i++ {
		n = lset.Len()
		lset.Delete(i)
		z := i % 10
		if z > 1 {
			z = 1
		}
		a := n - lset.Len()
		if a != z {
			t.Errorf("LinkedHashSet.Delete(%v) = %v, want %v", i, a, z)
		}
	}

	if !lset.IsEmpty() {
		t.Error("LinkedHashSet.IsEmpty() should return true")
	}
}

func TestLinkedHashSetString(t *testing.T) {
	e := "[1,3,2]"
	a := fmt.Sprintf("%s", NewLinkedHashSet(1, 3, 2))
	if a != e {
		t.Errorf(`fmt.Sprintf("%%s", NewLinkedHashSet(1, 3, 2)) = %v, want %v`, a, e)
	}
}

func TestLinkedHashSetMarshalJSON(t *testing.T) {
	cs := []struct {
		lset *LinkedHashSet
		json string
	}{
		{NewLinkedHashSet(0, 1, "0", "1", 0.0, 1.0, true, false), `[0,1,"0","1",0,1,true,false]`},
		{NewLinkedHashSet(0, "1", 2.0, [2]int{1, 2}), `[0,"1",2,[1,2]]`},
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
		lset *LinkedHashSet
	}

	cs := []Case{
		{`["0","1",0,1,true,false]`, NewLinkedHashSet("0", "1", 0.0, 1.0, true, false)},
		{`["1",2]`, NewLinkedHashSet("1", 2.0)},
	}

	for i, c := range cs {
		a := NewLinkedHashSet()
		err := json.Unmarshal([]byte(c.json), a)

		if err != nil {
			t.Errorf("[%d] json.Unmarshal(%v) error: %v", i, c.json, err)
		}

		if !reflect.DeepEqual(a, c.lset) {
			t.Errorf("[%d] json.Unmarshal(%q) = %v, want %q", i, c.json, a, c.lset)
		}
	}
}

func TestLinkedHashSetItemString(t *testing.T) {
	cs := []struct {
		e string
		s interface{}
	}{
		{"a", "a"},
		{"1", 1},
	}

	for _, c := range cs {
		i := &LinkedSetItem{value: c.s}
		a := i.String()
		if a != c.e {
			t.Errorf("LinkedSetItem(%v).String() = %q, want %q", c.s, a, c.e)
		}
	}
}
