//go:build go1.18
// +build go1.18

package cog

func ExampleHashMap() {
	m := NewHashMap[int, string]() // empty
	m.Set(1, "x")                  // 1->x
	m.Set(2, "b")                  // 2->b, 1->x (random order)
	m.Set(1, "a")                  // 2->b, 1->a (random order)
	_, _ = m.Get(2)                // b, true
	_, _ = m.Get(3)                // nil, false
	_ = m.Values()                 // []interface {}{"b", "a"} (random order)
	_ = m.Keys()                   // []interface {}{1, 2} (random order)
	m.Delete(1)                    // 2->b
	m.Clear()                      // empty
	m.IsEmpty()                    // true
	m.Len()                        // 0
}
