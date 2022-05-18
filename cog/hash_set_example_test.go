//go:build go1.18
// +build go1.18

package cog

func ExampleHashSet() {
	set := NewHashSet[int]() // empty
	set.Add(1)               // 1
	set.Add(2, 2, 3, 4, 5)   // 3, 1, 2, 4, 5 (random order, duplicates ignored)
	set.Delete(4)            // 5, 3, 2, 1 (random order)
	set.Delete(2, 3)         // 1, 5 (random order)
	set.Contains(1)          // true
	set.Contains(1, 5)       // true
	set.Contains(1, 6)       // false
	_ = set.Values()         // []int{5,1} (random order)
	set.Clear()              // empty
	set.IsEmpty()            // true
	set.Len()                // 0
}
