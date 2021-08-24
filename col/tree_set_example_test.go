package col

import "github.com/pandafw/pango/cmp"

func ExampleNewTreeSet() {
	set := NewTreeSet(cmp.CompareInt) // empty (keys are of type int)
	set.Add(1)                        // 1
	set.Add(2, 2, 3, 4, 5)            // 1, 2, 3, 4, 5 (in order, duplicates ignored)
	set.Delete(4)                     // 1, 2, 3, 5 (in order)
	set.Delete(2, 3)                  // 1, 5 (in order)
	set.Contains(1)                   // true
	set.Contains(1, 5)                // true
	set.Contains(1, 6)                // false
	_ = set.Values()                  // []int{1,5} (in order)
	set.Clear()                       // empty
	set.IsEmpty()                     // true
	set.Len()                         // 0
}
