package col

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/pandafw/pango/cmp"
)

func TestTreeMapSet(t *testing.T) {
	tree := NewTreeMap(cmp.CompareInt)
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
	if av, ev := fmt.Sprintf("%d%d%d%d%d%d%d", tree.Keys()...), "1234567"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%s%s%s%s%s%s%s", tree.Values()...), "abcdefg"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}

	tests1 := [][]interface{}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, nil, false},
	}

	for _, test := range tests1 {
		// retrievals
		av, actualFound := tree.Get(test[0])
		if av != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", av, test[1])
		}
	}
}

func TestTreeMapDelete(t *testing.T) {
	tree := NewTreeMap(cmp.CompareInt)
	tree.Set(5, "e")
	tree.Set(6, "f")
	tree.Set(7, "g")
	tree.Set(3, "c")
	tree.Set(4, "d")
	tree.Set(1, "x")
	tree.Set(2, "b")
	tree.Set(1, "a") //overwrite

	tree.Delete(5)
	tree.Delete(6)
	tree.Delete(7)
	tree.Delete(8)
	tree.Delete(5)

	if av, ev := fmt.Sprintf("%d%d%d%d", tree.Keys()...), "1234"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%s%s%s%s", tree.Values()...), "abcd"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%s%s%s%s", tree.Values()...), "abcd"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av := tree.Len(); av != 4 {
		t.Errorf("Got %v expected %v", av, 7)
	}

	tests2 := [][]interface{}{
		{1, "a", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, nil, false},
		{6, nil, false},
		{7, nil, false},
		{8, nil, false},
	}

	for _, test := range tests2 {
		av, actualFound := tree.Get(test[0])
		if av != test[1] || actualFound != test[2] {
			t.Errorf("Got %v expected %v", av, test[1])
		}
	}

	tree.Delete(1)
	tree.Delete(4)
	tree.Delete(2)
	tree.Delete(3)
	tree.Delete(2)
	tree.Delete(2)

	if av, ev := fmt.Sprintf("%s", tree.Keys()), "[]"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%s", tree.Values()), "[]"; av != ev {
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
		tm := NewTreeMap(cmp.CompareInt)
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

		tm.Delete(i)

		avs := make([]string, 0)
		for _, k := range tm.Keys() {
			avs = append(avs, fmt.Sprintf("%v", tm.lookup(k)))
		}

		if !reflect.DeepEqual(bvs, avs) {
			t.Errorf("tm.Delete(%v) failed\nBefore: %v\nAfter : %v", i, bvs, avs)
		}
	}
}

func TestTreeMapFloor(t *testing.T) {
	m := NewTreeMap(cmp.CompareInt)
	m.Set(7, "g")
	m.Set(3, "c")
	m.Set(1, "a")

	// key,expectedKey,expectedValue,expectedFound
	tests1 := [][]interface{}{
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
		n := m.Floor(test[0])
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
	m := NewTreeMap(cmp.CompareInt)
	m.Set(7, "g")
	m.Set(3, "c")
	m.Set(1, "a")

	// key,expectedKey,expectedValue,expectedFound
	tests1 := [][]interface{}{
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
		n := m.Ceiling(test[0])
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
	m := NewTreeMap(cmp.CompareString)
	m.Set("c", 3)
	m.Set("a", 1)
	m.Set("b", 2)
	count := 0
	m.Each(func(key interface{}, value interface{}) {
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
	})
}

func TestTreeMapFrontAndBack(t *testing.T) {
	tree := NewTreeMap(cmp.CompareInt)

	if av := tree.Front(); av != nil {
		t.Errorf("Got %v expected %v", av, nil)
	}
	if av := tree.Back(); av != nil {
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

	if av, ev := fmt.Sprintf("%d", tree.Front().key), "1"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%s", tree.Front().value), "x"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}

	if av, ev := fmt.Sprintf("%d", tree.Back().key), "7"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
	if av, ev := fmt.Sprintf("%s", tree.Back().value), "g"; av != ev {
		t.Errorf("Got %v expected %v", av, ev)
	}
}

func TestTreeMapCeilingAndFloor(t *testing.T) {
	tree := NewTreeMap(cmp.CompareInt)

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
	m := NewTreeMap(cmp.CompareString)
	it := m.Iterator()
	if it.Next() {
		t.Errorf("Shouldn't iterate on empty map")
	}
}

func TestTreeMapIteratorPrevOnEmpty(t *testing.T) {
	m := NewTreeMap(cmp.CompareString)
	it := m.Iterator()
	if it.Prev() {
		t.Errorf("Shouldn't iterate on empty map")
	}
}

func testTreeMapIterNextKey(t *testing.T, cmp cmp.Compare, kvs []P, vcnt bool) {
	tree := NewTreeMap(cmp, kvs...)

	count := 0
	for it := tree.Iterator(); it.Next(); {
		count++
		key := it.Key()
		if vcnt {
			key = it.Value()
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

func testTreeMapIterPrevKey(t *testing.T, cmp cmp.Compare, kvs []P, vcnt bool) {
	tree := NewTreeMap(cmp, kvs...)

	countDown := tree.Len()
	for it := tree.Iterator(); it.Prev(); {
		key := it.Key()
		if vcnt {
			key = it.Value()
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

func testTreeMapIterRemoveFront2Back(t *testing.T, cmp cmp.Compare, kvs []P, vcnt bool) {
	tree := NewTreeMap(cmp, kvs...)

	// remove from front to back
	len := tree.Len()
	count := 0
	for it := tree.Iterator(); it.Next(); {
		it.Remove()
		len--
		count++
		key := it.Key()
		if vcnt {
			key = it.Value()
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
		if len != tree.Len() {
			t.Errorf("tree.len = %v, want %v", tree.Len(), len)
		}
	}
}

func TestTreeMapIteratorReset(t *testing.T) {
	m := NewTreeMap(cmp.CompareInt)
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

func assertTreeMapIteratorRemove(t *testing.T, i int, it Iterator2, w *TreeMap) int {
	m := it.(*treeMapIterator).tree
	k := it.Key()

	//fmt.Printf("[%v] Remove(%v): %p\n%s\n", i, it.Key(), it.(*treeMapIterator).node, w.debug())
	//fmt.Printf("Remove %v\n%s\n", it.Key(), w.debug())

	w.Delete(k)

	it.Remove()

	it.SetValue(9999)

	if m.Contains(k) {
		t.Fatalf("[%d] w.Contains(%v) = true", i, k)
	}

	if m.String() != w.String() {
		t.Fatalf("[%d] (%v) %v != %v", i, k, m.String(), w.String())
	}

	return k.(int)
}

func TestTreeMapIteratorRemove(t *testing.T) {
	for i := 20; i < 50; i++ {
		m := NewTreeMap(cmp.CompareInt)
		w := NewTreeMap(cmp.CompareInt)

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
	m := NewTreeMap(cmp.CompareInt)
	for i := 1; i <= 100; i++ {
		m.Set(i, i)
	}

	// forward
	for it := m.Iterator(); it.Next(); {
		it.SetValue(it.Value().(int) + 100)
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
		it.SetValue(it.Value().(int) + 100)
	}
	for i := 1; i <= m.Len(); i++ {
		v, _ := m.Get(i)
		w := i + 200
		if v != w {
			t.Fatalf("Hash[%d] = %v, want %v", i, v, w)
		}
	}
}

func testTreeMapIterRemoveBack2Front(t *testing.T, cmp cmp.Compare, kvs []P, vcnt bool) {
	tree := NewTreeMap(cmp, kvs...)

	// remove from back to front
	countDown := tree.Len()
	for it := tree.Iterator(); it.Prev(); {
		it.Remove()

		key := it.Key()
		if vcnt {
			key = it.Value()
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

func testTreeMapIterRemoveMiddle(t *testing.T, cmp cmp.Compare, kvs []P) {
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

func testTreeMapIterate(t *testing.T, cmp cmp.Compare, kvs []P, vcnt bool) {
	testTreeMapIterNextKey(t, cmp, kvs, vcnt)

	testTreeMapIterPrevKey(t, cmp, kvs, vcnt)

	testTreeMapIterRemoveFront2Back(t, cmp, kvs, vcnt)

	testTreeMapIterRemoveBack2Front(t, cmp, kvs, vcnt)

	testTreeMapIterRemoveMiddle(t, cmp, kvs)
}

func TestTreeMapIterator1(t *testing.T) {
	kvs := []P{{5, "e"}, {6, "f"}, {7, "g"}, {3, "c"}, {4, "d"}, {1, "x"}, {2, "b"}, {1, "a"}} //overwrite
	// │   ┌── 7
	// └── 6
	//     │   ┌── 5
	//     └── 4
	//         │   ┌── 3
	//         └── 2
	//             └── 1

	testTreeMapIterate(t, cmp.CompareInt, kvs, false)
}

func TestTreeMapIterator2(t *testing.T) {
	kvs := []P{{3, "c"}, {1, "a"}, {2, "b"}}

	testTreeMapIterate(t, cmp.CompareInt, kvs, false)
}

func TestTreeMapIterator3(t *testing.T) {
	kvs := []P{{1, "a"}}

	testTreeMapIterate(t, cmp.CompareInt, kvs, false)
}

func TestTreeMapIterator4(t *testing.T) {
	kvs := []P{{13, 5}, {8, 3}, {17, 7}, {1, 1}, {11, 4}, {15, 6}, {25, 9}, {6, 2}, {22, 8}, {27, 10}}
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

	testTreeMapIterate(t, cmp.CompareInt, kvs, true)
}

func TestTreeMapMarshal(t *testing.T) {
	assert := func(m *TreeMap, txt string, t *testing.T) {
		if actualValue := m.Keys(); false ||
			actualValue[0].(string) != "a" ||
			actualValue[1].(string) != "b" ||
			actualValue[2].(string) != "c" ||
			actualValue[3].(string) != "d" ||
			actualValue[4].(string) != "e" {
			t.Errorf("[%s] Got %v expected %v", txt, actualValue, "[a,b,c,d,e]")
		}
		if actualValue := m.Values(); false ||
			actualValue[0].(string) != "1" ||
			actualValue[1].(string) != "2" ||
			actualValue[2].(string) != "3" ||
			actualValue[3].(string) != "4" ||
			actualValue[4].(string) != "5" {
			t.Errorf("[%s] Got %v expected %v", txt, actualValue, "[1,2,3,4,5]")
		}
		if actualValue, expectedValue := m.Len(), 5; actualValue != expectedValue {
			t.Errorf("[%s] Got %v expected %v", txt, actualValue, expectedValue)
		}
	}

	for i := 0; i < 10; i++ {
		original := NewTreeMap(cmp.CompareString)
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

		deserialized := NewTreeMap(cmp.CompareString)
		err = json.Unmarshal(serialized, &deserialized)
		if err != nil {
			t.Errorf("Got error %v", err)
		}
		assert(deserialized, "C", t)
	}
}
