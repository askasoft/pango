package col

import "github.com/pandafw/pango/cmp"

func ExampleNewSortedList() {
	list := NewSortedList(cmp.LessInt) // empty (keys are of type int)
	list.Add(1)                        // 1
	list.Add(2, 2, 3, 4, 5)            // 1, 2, 2, 3, 4, 5 (in order)
	_ = list.Get(0)                    // "a"	// _ = list.Get(100)                     // panic
	list.Delete(4)                     // 1, 2, 2, 3, 5 (in order)
	list.Delete(2, 3)                  // 1, 5 (in order)
	list.Contains(1)                   // true
	list.Contains(1, 5)                // true
	list.Contains(1, 6)                // false
	_ = list.Values()                  // []int{1,5} (in order)
	list.Clear()                       // empty
	list.IsEmpty()                     // true
	list.Len()                         // 0
}
