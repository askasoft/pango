package ref

import (
	"reflect"
	"testing"
)

func TestArrayGet(t *testing.T) {
	a := [][]string{
		{"0,0", "0,1"},
		{"1,0", "1,1"},
	}

	cs := []struct {
		w any
		p []int
	}{
		{"0,0", []int{0, 0}},
		{"0,1", []int{0, 1}},
		{"1,0", []int{1, 0}},
		{"1,1", []int{1, 1}},
	}

	for i, c := range cs {
		a, err := ArrayGet(a, c.p...)
		if err != nil || a != c.w {
			t.Errorf("[%d] ArrayGet(%v) = %v, %v", i, c.p, a, err)
		}
	}
}

func TestArraySet(t *testing.T) {
	cs := []struct {
		a any
		w any
		i int
		v any
	}{
		{[]string{"", ""}, []string{"", "1"}, 1, 1},
	}

	for i, c := range cs {
		_, err := ArraySet(c.a, c.i, c.v)
		if err != nil || !reflect.DeepEqual(c.a, c.w) {
			t.Errorf("[%d] ArraySet(%v, %d, %v) = _, %v, want: %v", i, c.a, c.i, c.v, err, c.w)
		}
	}
}
