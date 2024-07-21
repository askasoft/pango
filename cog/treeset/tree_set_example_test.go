//go:build go1.18
// +build go1.18

package treeset

import "github.com/askasoft/pango/cmp"

func ExampleNewTreeSet() {
	set := NewTreeSet(cmp.CompareInt)
	set.Add(1)              // 1
	set.Adds(2, 2, 3, 4, 5) // 1, 2, 3, 4, 5 (in order, duplicates ignored)
	set.Remove(4)           // 1, 2, 3, 5 (in order)
	set.Removes(2, 3)       // 1, 5 (in order)
	set.Contain(1)          // true
	set.Contains(1, 5)      // true
	set.Contains(1, 6)      // false
	_ = set.Values()        // []int{1,5} (in order)
	set.Clear()             // empty
	set.IsEmpty()           // true
	set.Len()               // 0
}
