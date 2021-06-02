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

	n111 := sl.Add(111)

	n := 1 + 10 + 100
	sn, se := sl.Search(n)
	if n111 != se || sn != 2 {
		t.Errorf("SortedList [%v] should contains %v", sl, n)
	}

	n++
	sn, se = sl.Search(n)
	if se != nil || sn != -1 {
		t.Errorf("SortedList [%v] should not contains %v", sl, n)
	}
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
