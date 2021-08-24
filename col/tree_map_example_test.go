package col

import "github.com/pandafw/pango/cmp"

func ExampleTreeMap() {
	m := NewTreeMap(cmp.CompareInt) // empty (keys are of type int)
	m.Set(1, "x")                   // 1->x
	m.Set(2, "b")                   // 1->x, 2->b (in order)
	m.Set(1, "a")                   // 1->a, 2->b (in order)
	_, _ = m.Get(2)                 // b, true
	_, _ = m.Get(3)                 // nil, false
	_ = m.Values()                  // []interface {}{"a", "b"} (in order)
	_ = m.Keys()                    // []interface {}{1, 2} (in order)
	m.Delete(1)                     // 2->b
	m.Clear()                       // empty
	m.IsEmpty()                     // true
	m.Len()                         // 0

	// Other:
	m.Front() // Returns the minimum key and its value from map.
	m.Back()  // Returns the maximum key and its value from map.
}
