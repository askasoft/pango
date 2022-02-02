package col

import "github.com/pandafw/pango/cmp"

func ExampleArrayList() {
	list := NewArrayList()
	list.Add("a")                         // ["a"]
	list.Add("c", "b")                    // ["a","c","b"]
	list.Sort(cmp.LessString)             // ["a","b","c"]
	_ = list.Get(0)                       // "a"  //_ = list.Get(100)  --> panic
	_ = list.Contains("a", "b", "c")      // true
	_ = list.Contains("a", "b", "c", "d") // false
	list.Swap(0, 1)                       // ["b","a",c"]
	list.Remove(2)                        // ["b","a"]
	list.Remove(1)                        // ["b"]
	list.Remove(0)                        // []
	list.Remove(0)                        // [] (ignored)
	_ = list.IsEmpty()                    // true
	_ = list.Len()                        // 0
	list.Add("a")                         // ["a"]
	list.Clear()                          // []
	list.Insert(0, "b")                   // ["b"]
	list.Insert(0, "a")                   // ["a","b"]
}
