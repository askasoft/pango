//go:build go1.18
// +build go1.18

package cog

func ExampleNewLinkedHashSet() {
	set := NewLinkedHashSet[int]()
	set.Add(5)              // 5
	set.Adds(4, 4, 3, 2, 1) // 5, 4, 3, 2, 1 (in insertion-order, duplicates ignored)
	set.Add(4)              // 5, 4, 3, 2, 1 (duplicates ignored, insertion-order unchanged)
	set.Remove(4)           // 5, 3, 2, 1 (in insertion-order)
	set.Removes(2, 3)       // 5, 1 (in insertion-order)
	set.Contains(1)         // true
	set.Contains(1, 5)      // true
	set.Contains(1, 6)      // false
	_ = set.Values()        // []int{5, 1} (in insertion-order)
	set.Clear()             // empty
	set.IsEmpty()           // true
	set.Len()               // 0
}
