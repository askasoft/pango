package treeset

import "cmp"

func ExampleNewTreeSet() {
	set := NewTreeSet(cmp.Compare[int])
	set.Add(1)                // 1
	set.AddAll(2, 2, 3, 4, 5) // 1, 2, 3, 4, 5 (in order, duplicates ignored)
	set.Remove(4)             // 1, 2, 3, 5 (in order)
	set.RemoveAll(2, 3)       // 1, 5 (in order)
	set.Contains(1)           // true
	set.ContainsAll(1, 5)     // true
	set.ContainsAll(1, 6)     // false
	_ = set.Values()          // []int{1,5} (in order)
	set.Clear()               // empty
	set.IsEmpty()             // true
	set.Len()                 // 0
}
