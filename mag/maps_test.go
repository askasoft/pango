//go:build go1.18
// +build go1.18

package mag

import (
	"reflect"
	"testing"
)

func TestEqual(t *testing.T) {
	var a, b map[int]int
	if !Equal(a, b) {
		t.Errorf("a = %v\nb = %v", a, b)
	}

	a = make(map[int]int)
	if !Equal(a, b) || !Equal(b, a) {
		t.Errorf("a = %v\nb = %v", a, b)
	}

	b = make(map[int]int)
	for i := 0; i < 100; i++ {
		a[i] = i * 100
		b[99-i] = (99 - i) * 100
	}

	if !Equal(a, b) || !Equal(b, a) {
		t.Errorf("a = %v\nb = %v", a, b)
	}
}

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
