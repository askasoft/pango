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
		err := ArraySet(c.a, c.i, c.v)
		if err != nil || !reflect.DeepEqual(c.a, c.w) {
			t.Errorf("[%d] ArraySet(%v, %d, %v) = _, %v, want: %v", i, c.a, c.i, c.v, err, c.w)
		}
	}
}

func TestToSlice(t *testing.T) {
	cs := []struct {
		a any
		w any
	}{
		{[1]string{"1"}, []string{"1"}},
		{[2]string{"1"}, []string{"1", ""}},
	}

	for i, c := range cs {
		a, err := ToSlice(c.a)
		if err != nil || !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] ToSlice(%v) = (%v, %v), want: %v", i, c.a, a, err, c.w)
		}
	}
}

func TestSliceAdd(t *testing.T) {
	cs := []struct {
		s any
		a []any
		w any
	}{
		{[]string{"a"}, []any{1, "2"}, []string{"a", "1", "2"}},
		{[1]string{"a"}, []any{1, "2"}, []string{"a", "1", "2"}},
		{[2]string{"a"}, []any{1, "2"}, []string{"a", "", "1", "2"}},
	}

	for i, c := range cs {
		a, err := SliceAdd(c.s, c.a...)
		if err != nil || !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] SliceAdd(%v, %v) = (%v, %v), want: %v", i, c.s, c.a, a, err, c.w)
		}
	}
}
