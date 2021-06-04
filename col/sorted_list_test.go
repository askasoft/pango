package col

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/pandafw/pango/str"
	"github.com/stretchr/testify/assert"
)

type inta []interface{}

// Len is the number of elements in the collection.
func (a inta) Len() int {
	return len(a)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (a inta) Less(i, j int) bool {
	return a[i].(int) < a[j].(int)
}

// Swap swaps the elements with indexes i and j.
func (a inta) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func TestSortedListAsc(t *testing.T) {
	for i := 1; i < 100; i++ {
		sl := NewSortedList(func(a, b interface{}) bool {
			return a.(int) < b.(int)
		})
		a := make([]interface{}, 0, 100)
		for n := i; n < 100; n++ {
			a = append(a, n)
		}

		for j := 0; j < len(a); j++ {
			sl.Add(a[j])
			assert.True(t, sort.IsSorted(inta(sl.Values())), j)
		}

		assert.Equal(t, a, sl.Values())
	}
}

func TestSortedListDesc(t *testing.T) {
	for i := 1; i < 100; i++ {
		sl := NewSortedList(LessInt)
		a := make([]interface{}, 0, 100)
		for n := i; n < 100; n++ {
			a = append(a, n)
		}
		for j := len(a) - 1; j >= 0; j-- {
			sl.Add(a[j])
			assert.True(t, sort.IsSorted(inta(sl.Values())), j)
		}

		assert.Equal(t, a, sl.Values())
	}
}

func TestSortedListRandom(t *testing.T) {
	for i := 1; i < 100; i++ {
		sl := NewSortedList(func(a, b interface{}) bool {
			return a.(int) < b.(int)
		})
		a := make([]interface{}, 0, 100)
		for n := i; n < 100; n++ {
			a = append(a, rand.Intn(20))
		}
		for j := len(a) - 1; j >= 0; j-- {
			sl.Add(a[j])
			assert.True(t, sort.IsSorted(inta(sl.Values())), j)
		}

		sort.Sort(inta(a))
		assert.Equal(t, a, sl.Values())
	}
}

func TestSortedListContains(t *testing.T) {
	sl := NewSortedList(LessString, "1", "11", "111", "1", "11", "111")

	n := str.Repeat("1", 3)

	if !sl.Contains(n) {
		t.Errorf("SortedList [%v] should contains %v", sl, n)
	}

	n += "1"
	if sl.Contains(n) {
		t.Errorf("SortedList [%v] should not contains %v", sl, n)
	}
}

func TestSetSearch(t *testing.T) {
	sl := NewSortedList(LessInt, 1, 11)

	n := 10 + 1
	sn, se := sl.Search(n)
	if se == nil || sn != 1 {
		t.Errorf("SortedList [%v] should contains %v", sl, n)
	}

	n = 1 + 10 + 100
	for i := 0; i < 100; i++ {
		n111 := sl.Add(111)
		sn, se = sl.Search(n)
		if se != n111 || sn != 2 {
			t.Errorf("SortedList [%v] should contains %v", sl, n)
		}
	}

	n++
	sn, se = sl.Search(n)
	if se != nil || sn != -1 {
		t.Errorf("SortedList [%v] should not contains %v", sl, n)
	}
}

func TestSortedListDelete(t *testing.T) {
	sl := NewSortedList(LessInt)

	for i := 0; i < 100; i++ {
		sl.Add(i)
	}

	assert.False(t, sl.Delete(100))
	for i := 0; i < 100; i++ {
		assert.True(t, sl.Delete(i))
	}
	assert.True(t, sl.IsEmpty())
}

func TestSortedListDeleteAll(t *testing.T) {
	sl := NewSortedList(LessInt)

	for i := 0; i < 100; i++ {
		z := i % 10
		for j := 0; j < z; j++ {
			sl.Add(i)
		}
	}

	assert.Equal(t, 0, sl.DeleteAll(100))
	for i := 0; i < 100; i++ {
		z := i % 10
		assert.Equal(t, z, sl.DeleteAll(i), i)
	}
	assert.True(t, sl.IsEmpty())
}

func TestSortedListString(t *testing.T) {
	assert.Equal(t, "[1,2,3]", fmt.Sprintf("%s", NewSortedList(LessInt, 1, 3, 2)))
}

func TestSortedListMarshalJSON(t *testing.T) {
	type Case struct {
		list *SortedList
		json string
	}

	cs := []Case{
		{NewSortedList(LessString, "2", "1", "0"), `["0","1","2"]`},
	}

	for i, c := range cs {
		bs, err := json.Marshal(c.list)
		assert.Nil(t, err)
		assert.Equal(t, c.json, string(bs), fmt.Sprintf("Marshal [%d]", i))
	}
}

func TestSortedListUnmarshalJSON(t *testing.T) {
	type Case struct {
		json string
		list *List
	}

	cs := []Case{
		{`["2","0","1"]`, NewList("0", "1", "2")},
	}

	for i, c := range cs {
		a := NewSortedList(LessString)
		err := json.Unmarshal([]byte(c.json), a)
		assert.Nil(t, err)
		if !reflect.DeepEqual(a.list, c.list) {
			t.Fatalf("Unmarshal List [%d] not deeply equal: %#v expected %#v", i, a, c.list)
		}
	}
}
