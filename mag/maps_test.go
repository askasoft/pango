//go:build go1.18
// +build go1.18

package mag

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	cs := []struct {
		s []map[int]int
		w map[int]int
	}{
		{[]map[int]int{{1: 10, 2: 20, 3: 30}, {1: 11, 2: 22}}, map[int]int{1: 11, 2: 22, 3: 30}},
	}

	for i, c := range cs {
		a := Merge(c.s...)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] Merge(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestReverse(t *testing.T) {
	cs := []struct {
		s map[int]int
		w map[int]int
	}{
		{map[int]int{1: 10, 2: 20, 3: 30}, map[int]int{10: 1, 20: 2, 30: 3}},
	}

	for i, c := range cs {
		a := Reverse(c.s)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] Reverse(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}
