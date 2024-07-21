//go:build go1.18
// +build go1.18

package linkedlist

import "github.com/askasoft/pango/cmp"

func ExampleNewLinkedList() {
	list := NewLinkedList[string]()
	list.Add("a")                         // ["a"]
	list.Adds("c", "b")                   // ["a","c","b"]
	list.Sort(cmp.LessString)             // ["a","b","c"]
	_ = list.Get(0)                       // "a"  //_ = list.Get(100)  --> panic
	_ = list.Contains("a", "b", "c")      // true
	_ = list.Contains("a", "b", "c", "d") // false
	list.Swap(0, 1)                       // ["b","a",c"]
	list.DeleteAt(2)                      // ["b","a"]
	list.DeleteAt(1)                      // ["b"]
	list.DeleteAt(0)                      // []
	_ = list.IsEmpty()                    // true
	_ = list.Len()                        // 0
	list.Add("a")                         // ["a"]
	list.Clear()                          // []
	list.Insert(0, "b")                   // ["b"]
	list.Insert(0, "a")                   // ["a","b"]
}
