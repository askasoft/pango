package mag

import (
	"reflect"
	"sort"
	"testing"
)

func TestKeys(t *testing.T) {
	cs := []struct {
		m map[int]string
		w []int
	}{
		{map[int]string{1: "11", 2: "22"}, []int{1, 2}},
	}

	for i, c := range cs {
		a := Keys(c.m)
		sort.Ints(a)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] Keys(%v) = %v, want %v", i, c.m, a, c.w)
		}
	}
}

func TestValues(t *testing.T) {
	cs := []struct {
		m map[int]string
		w []string
	}{
		{map[int]string{1: "11", 2: "22"}, []string{"11", "22"}},
	}

	for i, c := range cs {
		a := Values(c.m)
		sort.Strings(a)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] Values(%v) = %v, want %v", i, c.m, a, c.w)
		}
	}
}

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

func TestCopy(t *testing.T) {
	cs := []struct {
		d map[int]int
		s map[int]int
		w map[int]int
	}{
		{map[int]int{1: 11, 2: 22}, map[int]int{2: 20, 3: 30}, map[int]int{1: 11, 2: 20, 3: 30}},
	}

	for i, c := range cs {
		Copy(c.d, c.s)
		a := c.d
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] Copy(%v) = %v, want %v", i, c.s, a, c.w)
		}
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
