//go:build go1.18
// +build go1.18

package cog

func ExampleArrayList() {
	list := NewArrayList[string]()
	list.Add("a")                         // ["a"]
	list.Adds("c", "b")                   // ["a","c","b"]
	list.Sort(LessString)                 // ["a","b","c"]
	_ = list.Get(0)                       // "a"  //_ = list.Get(100)  --> panic
	_ = list.Contains("a", "b", "c")      // true
	_ = list.Contains("a", "b", "c", "d") // false
	list.Swap(0, 1)                       // ["b","a",c"]
	list.RemoveAt(2)                      // ["b","a"]
	list.RemoveAt(1)                      // ["b"]
	list.RemoveAt(0)                      // []
	_ = list.IsEmpty()                    // true
	_ = list.Len()                        // 0
	list.Add("a")                         // ["a"]
	list.Clear()                          // []
	list.Insert(0, "b")                   // ["b"]
	list.Insert(0, "a")                   // ["a","b"]
}
