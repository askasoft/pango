package hashset

func ExampleHashSet() {
	set := NewHashSet[int]()
	set.Add(1)                // 1
	set.AddAll(2, 2, 3, 4, 5) // 3, 1, 2, 4, 5 (random order, duplicates ignored)
	set.Remove(4)             // 5, 3, 2, 1 (random order)
	set.RemoveAll(2, 3)       // 1, 5 (random order)
	set.Contains(1)           // true
	set.ContainsAll(1, 5)     // true
	set.ContainsAll(1, 6)     // false
	_ = set.Values()          // []int{5,1} (random order)
	set.Clear()               // empty
	set.IsEmpty()             // true
	set.Len()                 // 0
}
