package treemap

import (
	"cmp"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/askasoft/pango/cog"
)

func TestTreeMapInterface(t *testing.T) {
	var _ cog.Map[int, int] = NewTreeMap[int, int](cmp.Compare[int])
	var _ cog.IterableMap[int, int] = NewTreeMap[int, int](cmp.Compare[int])
}

func TestTreeMapSet(t *testing.T) {
	tree := NewTreeMap[int, string](cmp.Compare[int])
	tree.Set(5, "e")
	tree.Set(6, "f")
	tree.Set(7, "g")
	tree.Set(3, "c")
	tree.Set(4, "d")
	tree.Set(1, "x")
	tree.Set(2, "b")
	tree.Set(1, "a") //overwrite

	if av := tree.Len(); av != 7 {
		t.Errorf("Got %v expected %v", av, 7)
	}
	if av, ev := fmt.Sprintf("%v", tree.Keys()), "[1 2 3 4 5 6 7]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%v", tree.Values()), "[a b c d e f g]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%v", tree.Entries()), "[{1 a} {2 b} {3 c} {4 d} {5 e} {6 f} {7 g}]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}

	tests1 := [][]any{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, "", false},
	}

	for _, test := range tests1 {
		// retrievals
		av, actualFound := tree.Get(test[0].(int))
		if av != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", av, test[1])
		}
	}
}

func TestTreeMapDebug(t *testing.T) {
	tree := NewTreeMap[int, int](cmp.Compare[int])
	ev := "(empty)"
	av := tree.debug()
	if av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestTreeMapDelete(t *testing.T) {
	tree := NewTreeMap[int, string](cmp.Compare[int])
	tree.Set(5, "e")
	tree.Set(6, "f")
	tree.Set(7, "g")
	tree.Set(3, "c")
	tree.Set(4, "d")
	tree.Set(1, "x")
	tree.Set(2, "b")
	tree.Set(1, "a") //overwrite

	tree.Remove(5)
	tree.RemoveAll(6, 7)
	tree.Remove(8)
	tree.Remove(5)

	if av, ev := fmt.Sprintf("%v", tree.Keys()), "[1 2 3 4]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%v", tree.Values()), "[a b c d]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%v", tree.Values()), "[a b c d]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av := tree.Len(); av != 4 {
		t.Errorf("Got %v expected %v", av, 7)
	}

	tests2 := [][]any{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "", false},
		{6, "", false},
		{7, "", false},
		{8, "", false},
	}

	for _, test := range tests2 {
		av, actualFound := tree.Get(test[0].(int))
		if av != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", av, test[1])
		}
	}

	tree.Remove(1)
	tree.RemoveAll(4, 2, 3, 2)
	tree.Remove(2)

	if av, ev := fmt.Sprintf("%v", tree.Keys()), "[]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%v", tree.Values()), "[]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%v", tree.Entries()), "[]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if empty, size := tree.IsEmpty(), tree.Len(); empty != true || size != -0 {
		t.Errorf("Got %v expected %v", empty, true)
	}
}

func TestTreeMapDelete2(t *testing.T) {
	// │           ┌── 20
	// │       ┌── 19
	// │       │   └── 18
	// │   ┌── 17
	// │   │   │   ┌── 16
	// │   │   └── 15
	// │   │       └── 14
	// └── 13
	//     │       ┌── 12
	//     │   ┌── 11
	//     │   │   └── 10
	//     └── 9
	//         │       ┌── 8
	//         │   ┌── 7
	//         │   │   └── 6
	//         └── 5
	//             │   ┌── 4
	//             └── 3
	//                 └── 2
	//                     └── 1

	for i := 1; i <= 7; i++ {
		tm := NewTreeMap[int, int](cmp.Compare[int])
		for n := 20; n >= 1; n-- {
			tm.Set(n, n)
		}
		//fmt.Println(tm.Graph(false))

		bvs := make([]string, 0)
		for j, k := range tm.Keys() {
			if j != i-1 {
				bvs = append(bvs, fmt.Sprintf("%v", tm.lookup(k)))
			}
		}

		tm.Remove(i)

		avs := make([]string, 0)
		for _, k := range tm.Keys() {
			avs = append(avs, fmt.Sprintf("%v", tm.lookup(k)))
		}

		if !reflect.DeepEqual(bvs, avs) {
			t.Errorf("tm.Remove(%v) failed\nBefore: %v\nAfter : %v", i, bvs, avs)
		}
	}
}

func TestTreeMapFloor(t *testing.T) {
	m := NewTreeMap[int, string](cmp.Compare[int])
	m.Set(7, "g")
	m.Set(3, "c")
	m.Set(1, "a")

	// key,expectedKey,expectedValue,expectedFound
	tests1 := [][]any{
		{-1, nil, nil, false},
		{0, nil, nil, false},
		{1, 1, "a", true},
		{2, 1, "a", true},
		{3, 3, "c", true},
		{4, 3, "c", true},
		{7, 7, "g", true},
		{8, 7, "g", true},
	}

	for _, test := range tests1 {
		// retrievals
		n := m.Floor(test[0].(int))
		if test[3] == true {
			if n == nil {
				t.Fatalf("[%v] Got %v, expected %v", test[0], n, test[3])
			}
			if n.key != test[1] || n.value != test[2] {
				t.Fatalf("[%v] Got %v, expected %v, %v", test[0], n, test[1], test[2])
			}
		} else {
			if n != nil {
				t.Fatalf("[%v] Got %v, expected %v", test[0], n, test[3])
			}
		}
	}
}

func TestTreeMapCeiling(t *testing.T) {
	m := NewTreeMap[int, string](cmp.Compare[int])
	m.Set(7, "g")
	m.Set(3, "c")
	m.Set(1, "a")

	// key,expectedKey,expectedValue,expectedFound
	tests1 := [][]any{
		{-1, 1, "a", true},
		{0, 1, "a", true},
		{1, 1, "a", true},
		{2, 3, "c", true},
		{3, 3, "c", true},
		{4, 7, "g", true},
		{7, 7, "g", true},
		{8, nil, nil, false},
	}

	for _, test := range tests1 {
		// retrievals
		n := m.Ceiling(test[0].(int))
		if test[3] == true {
			if n == nil {
				t.Fatalf("[%v] Got %v, expected %v", test[0], n, test[3])
			}
			if n.key != test[1] || n.value != test[2] {
				t.Fatalf("[%v] Got %v, expected %v, %v", test[0], n, test[1], test[2])
			}
		} else {
			if n != nil {
				t.Fatalf("[%v] Got %v, expected %v", test[0], n, test[3])
			}
		}
	}
}

func TestTreeMapEach(t *testing.T) {
	m := NewTreeMap[string, int](cmp.Compare[string])
	m.Set("c", 3)
	m.Set("a", 1)
	m.Set("b", 2)

	count := 0
	m.Each(func(key string, value int) bool {
		count++
		if actualValue, expectedValue := count, value; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		switch value {
		case 1:
			if actualValue, expectedValue := key, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := key, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 3:
			if actualValue, expectedValue := key, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
		return true
	})
}

func TestHashMapSeq(t *testing.T) {
	m := NewTreeMap[string, int](cmp.Compare[string])
	m.Set("c", 3)
	m.Set("a", 1)
	m.Set("b", 2)

	count := 0
	for key, value := range m.Seq() {
		count++
		if actualValue, expectedValue := count, value; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		switch value {
		case 1:
			if actualValue, expectedValue := key, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := key, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 3:
			if actualValue, expectedValue := key, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
	}
}

func TestTreeMapHeadAndTail(t *testing.T) {
	tree := NewTreeMap[int, string](cmp.Compare[int])

	if av := tree.Head(); av != nil {
		t.Errorf("Got %v expected %v", av, nil)
	}
	if av := tree.Tail(); av != nil {
		t.Errorf("Got %v expected %v", av, nil)
	}

	tree.Set(1, "a")
	tree.Set(5, "e")
	tree.Set(6, "f")
	tree.Set(7, "g")
	tree.Set(3, "c")
	tree.Set(4, "d")
	tree.Set(1, "x") // overwrite
	tree.Set(2, "b")

	if av, ev := fmt.Sprintf("%d", tree.Head().key), "1"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := tree.Head().value, "x"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}

	if av, ev := fmt.Sprintf("%d", tree.Tail().key), "7"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := tree.Tail().value, "g"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestTreeMapCeilingAndFloor(t *testing.T) {
	tree := NewTreeMap[int, string](cmp.Compare[int])

	if node := tree.Floor(0); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
	if node := tree.Ceiling(0); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	tree.Set(5, "e")
	tree.Set(6, "f")
	tree.Set(7, "g")
	tree.Set(3, "c")
	tree.Set(4, "d")
	tree.Set(1, "x")
	tree.Set(2, "b")

	if node := tree.Floor(4); node.key != 4 {
		t.Errorf("Got %v expected %v", node.key, 4)
	}
	if node := tree.Floor(0); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	if node := tree.Ceiling(4); node.key != 4 {
		t.Errorf("Got %v expected %v", node.key, 4)
	}
	if node := tree.Ceiling(8); node != nil {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
}

func TestTreeMapIteratorNextOnEmpty(t *testing.T) {
	m := NewTreeMap[string, string](cmp.Compare[string])
	it := m.Iterator()
	if it.Next() {
		t.Errorf("Shouldn't iterate on empty map")
	}
}

func TestTreeMapIteratorPrevOnEmpty(t *testing.T) {
	m := NewTreeMap[string, string](cmp.Compare[string])
	it := m.Iterator()
	if it.Prev() {
		t.Errorf("Shouldn't iterate on empty map")
	}
}

func testTreeMapIterNextKey(t *testing.T, cmp cog.Compare[int], kvs []cog.P[int, string], vcnt bool) {
	tree := NewTreeMap(cmp, kvs...)

	count := 0
	for it := tree.Iterator(); it.Next(); {
		count++
		key := it.Key()
		if vcnt {
			key, _ = strconv.Atoi(it.Value())
		}

		switch key {
		case count:
			if av, ev := key, count; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		default:
			if av, ev := key, count; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		}
	}
	if av, ev := count, tree.Len(); av != ev {
		t.Errorf("Size different. Got %v expected %v", av, ev)
	}
}

func testTreeMapIterPrevKey(t *testing.T, cmp cog.Compare[int], kvs []cog.P[int, string], vcnt bool) {
	tree := NewTreeMap(cmp, kvs...)

	countDown := tree.Len()
	for it := tree.Iterator(); it.Prev(); {
		key := it.Key()
		if vcnt {
			key, _ = strconv.Atoi(it.Value())
		}
		switch key {
		case countDown:
			if av, ev := key, countDown; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		default:
			if av, ev := key, countDown; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		}
		countDown--
	}
	if av, ev := countDown, 0; av != ev {
		t.Errorf("Size different. Got %v expected %v", av, ev)
	}
}

func testTreeMapIterRemoveHead2Tail(t *testing.T, cmp cog.Compare[int], kvs []cog.P[int, string], vcnt bool) {
	tree := NewTreeMap(cmp, kvs...)

	sz := tree.Len()
	count := 0
	for it := tree.Iterator(); it.Next(); {
		it.Remove()
		sz--
		count++
		key := it.Key()
		if vcnt {
			key, _ = strconv.Atoi(it.Value())
		}
		switch key {
		case count:
			if av, ev := key, count; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		default:
			if av, ev := key, count; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		}
		if sz != tree.Len() {
			t.Errorf("tree.len = %v, want %v", tree.Len(), sz)
		}
	}
}

func TestTreeMapIteratorReset(t *testing.T) {
	m := NewTreeMap[int, string](cmp.Compare[int])
	it := m.Iterator()
	m.Set(3, "c")
	m.Set(1, "a")
	m.Set(2, "b")
	for it.Next() {
	}
	it.Reset()
	it.Next()
	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
	}

	it.Reset()
	it.Prev()
	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
	}
}

func assertTreeMapIteratorRemove(t *testing.T, i int, it cog.Iterator2[int, int], w *TreeMap[int, int]) int {
	m := it.(*treeMapIterator[int, int]).tree
	k := it.Key()

	//fmt.Printf("[%v] Remove(%v): %p\n%s\n", i, it.Key(), it.(*treeMapIterator).node, w.debug())
	//fmt.Printf("Remove %v\n%s\n", it.Key(), w.debug())

	w.Remove(k)

	it.Remove()

	it.SetValue(9999)

	if m.Contains(k) {
		t.Fatalf("[%d] w.Contains(%v) = true", i, k)
	}

	if m.String() != w.String() {
		t.Fatalf("[%d] (%v) %v != %v", i, k, m.String(), w.String())
	}

	return k
}

func TestTreeMapIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		m := NewTreeMap[int, int](cmp.Compare[int])
		w := NewTreeMap[int, int](cmp.Compare[int])

		for n := 0; n < i; n++ {
			//fmt.Printf("Add %v\n%v\n\n", n, w.debug())
			m.Set(n, -n)
			w.Set(n, -n)
		}

		it := m.Iterator()

		// remove nothing
		it.Remove()
		it.SetValue(9999)
		if m.Len() != i {
			t.Fatalf("[%d] m.Len() == %v, want %v", i, m.Len(), i)
		}

		// remove middle
		for j := 0; j <= m.Len()/2; j++ {
			it.Next()
		}

		v := assertTreeMapIteratorRemove(t, i, it, w)

		it.Next()
		if v+1 != it.Key() {
			t.Fatalf("[%d] it.Key() = %v, want %v", i, it.Key(), v+1)
		}
		assertTreeMapIteratorRemove(t, i, it, w)

		it.Prev()
		if v-1 != it.Key() {
			t.Fatalf("[%d] it.Key() = %v, want %v", i, it.Key(), v-1)
		}
		assertTreeMapIteratorRemove(t, i, it, w)

		// remove first
		for it.Prev() {
		}
		assertTreeMapIteratorRemove(t, i, it, w)

		// remove last
		for it.Next() {
		}
		assertTreeMapIteratorRemove(t, i, it, w)

		// remove all
		it.Reset()
		if i%2 == 0 {
			for it.Prev() {
				assertTreeMapIteratorRemove(t, i, it, w)
			}
		} else {
			for it.Next() {
				assertTreeMapIteratorRemove(t, i, it, w)
			}
		}
		if !m.IsEmpty() {
			t.Fatalf("[%d] m.IsEmpty() = true", i)
		}
	}
}

func TestTreeMapIteratorSetValue(t *testing.T) {
	m := NewTreeMap[int, int](cmp.Compare[int])
	for i := 1; i <= 100; i++ {
		m.Set(i, i)
	}

	// forward
	for it := m.Iterator(); it.Next(); {
		it.SetValue(it.Value() + 100)
	}
	for i := 1; i <= m.Len(); i++ {
		v, _ := m.Get(i)
		w := i + 100
		if v != w {
			t.Fatalf("Hash[%d] = %v, want %v", i, v, w)
		}
	}

	// backward
	for it := m.Iterator(); it.Prev(); {
		it.SetValue(it.Value() + 100)
	}
	for i := 1; i <= m.Len(); i++ {
		v, _ := m.Get(i)
		w := i + 200
		if v != w {
			t.Fatalf("Hash[%d] = %v, want %v", i, v, w)
		}
	}
}

func testTreeMapIterRemoveTail2Head(t *testing.T, cmp cog.Compare[int], kvs []cog.P[int, string], vcnt bool) {
	tree := NewTreeMap(cmp, kvs...)

	countDown := tree.Len()
	for it := tree.Iterator(); it.Prev(); {
		it.Remove()

		key := it.Key()
		if vcnt {
			key, _ = strconv.Atoi(it.Value())
		}
		switch key {
		case countDown:
			if av, ev := key, countDown; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		default:
			if av, ev := key, countDown; av != ev {
				t.Errorf("Got %v expected %v", av, ev)
			}
		}
		countDown--
		if countDown != tree.Len() {
			t.Errorf("tree.len = %v, want %v", tree.Len(), countDown)
		}
	}
}

func testTreeMapIterRemoveMiddle(t *testing.T, cmp cog.Compare[int], kvs []cog.P[int, string]) {
	tree := NewTreeMap(cmp, kvs...)

	// remove from middle
	it := tree.Iterator()

	it.Next()
	for i := 1; i < tree.Len()/2; i++ {
		it.Next()
	}

	countDown := tree.Len()
	for countDown > 0 {
		//fmt.Printf("delete: %v (%p)\n%s\n", it.Key(), it.(*treeMapIterator).node, tree.debug())

		it.Remove()
		if countDown%2 == 0 {
			if !it.Next() {
				it.Prev()
			}
		} else {
			if !it.Prev() {
				it.Next()
			}
		}

		countDown--
		if countDown != tree.Len() {
			t.Fatalf("tree.len = %v, want %v", tree.Len(), countDown)
		}
	}

	countDown = 0
	if countDown != tree.Len() {
		t.Fatalf("tree.len = %v, want %v", tree.Len(), countDown)
	}
}

func testTreeMapIterate(t *testing.T, cmp cog.Compare[int], kvs []cog.P[int, string], vcnt bool) {
	testTreeMapIterNextKey(t, cmp, kvs, vcnt)

	testTreeMapIterPrevKey(t, cmp, kvs, vcnt)

	testTreeMapIterRemoveHead2Tail(t, cmp, kvs, vcnt)

	testTreeMapIterRemoveTail2Head(t, cmp, kvs, vcnt)

	testTreeMapIterRemoveMiddle(t, cmp, kvs)
}

func TestTreeMapIterator1(t *testing.T) {
	kvs := []cog.P[int, string]{
		{Key: 5, Val: "e"},
		{Key: 6, Val: "f"},
		{Key: 7, Val: "g"},
		{Key: 3, Val: "c"},
		{Key: 4, Val: "d"},
		{Key: 1, Val: "x"},
		{Key: 2, Val: "b"},
		{Key: 1, Val: "a"},
	} //overwrite
	// │   ┌── 7
	// └── 6
	//     │   ┌── 5
	//     └── 4
	//         │   ┌── 3
	//         └── 2
	//             └── 1

	testTreeMapIterate(t, cmp.Compare[int], kvs, false)
}

func TestTreeMapIterator2(t *testing.T) {
	kvs := []cog.P[int, string]{
		{Key: 3, Val: "c"},
		{Key: 1, Val: "a"},
		{Key: 2, Val: "b"},
	}

	testTreeMapIterate(t, cmp.Compare[int], kvs, false)
}

func TestTreeMapIterator3(t *testing.T) {
	kvs := []cog.P[int, string]{{Key: 1, Val: "a"}}

	testTreeMapIterate(t, cmp.Compare[int], kvs, false)
}

func TestTreeMapIterator4(t *testing.T) {
	kvs := []cog.P[int, string]{
		{Key: 13, Val: "5"},
		{Key: 8, Val: "3"},
		{Key: 17, Val: "7"},
		{Key: 1, Val: "1"},
		{Key: 11, Val: "4"},
		{Key: 15, Val: "6"},
		{Key: 25, Val: "9"},
		{Key: 6, Val: "2"},
		{Key: 22, Val: "8"},
		{Key: 27, Val: "10"}}
	// │           ┌── 27
	// │       ┌── 25
	// │       │   └── 22
	// │   ┌── 17
	// │   │   └── 15
	// └── 13
	//     │   ┌── 11
	//     └── 8
	//         │   ┌── 6
	//         └── 1

	testTreeMapIterate(t, cmp.Compare[int], kvs, true)
}

func TestTreeMapMarshal(t *testing.T) {
	assert := func(m *TreeMap[string, string], txt string, t *testing.T) {
		if actualValue := m.Keys(); false ||
			actualValue[0] != "a" ||
			actualValue[1] != "b" ||
			actualValue[2] != "c" ||
			actualValue[3] != "d" ||
			actualValue[4] != "e" {
			t.Errorf("[%s] Got %v expected %v", txt, actualValue, "[a,b,c,d,e]")
		}
		if actualValue := m.Values(); false ||
			actualValue[0] != "1" ||
			actualValue[1] != "2" ||
			actualValue[2] != "3" ||
			actualValue[3] != "4" ||
			actualValue[4] != "5" {
			t.Errorf("[%s] Got %v expected %v", txt, actualValue, "[1,2,3,4,5]")
		}
		if actualValue, expectedValue := m.Len(), 5; actualValue != expectedValue {
			t.Errorf("[%s] Got %v expected %v", txt, actualValue, expectedValue)
		}
	}

	for i := 0; i < 10; i++ {
		original := NewTreeMap[string, string](cmp.Compare[string])
		original.Set("d", "4")
		original.Set("e", "5")
		original.Set("c", "3")
		original.Set("b", "2")
		original.Set("a", "1")

		assert(original, "A", t)

		serialized, err := json.Marshal(original)
		if err != nil {
			t.Errorf("Got error %v", err)
		}
		assert(original, "B", t)

		deserialized := NewTreeMap[string, string](cmp.Compare[string])
		err = json.Unmarshal(serialized, &deserialized)
		if err != nil {
			t.Errorf("Got error %v", err)
		}
		assert(deserialized, "C", t)
	}
}
