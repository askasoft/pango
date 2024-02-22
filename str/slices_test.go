package str

import (
	"fmt"
	"reflect"
	"testing"
)

func TestToLowers(t *testing.T) {
	cs := []struct {
		s []string
		w []string
	}{
		{[]string{" A "}, []string{" a "}},
	}

	for i, c := range cs {
		a := ToLowers(c.s)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] ToLowers(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestToUppers(t *testing.T) {
	cs := []struct {
		s []string
		w []string
	}{
		{[]string{" a "}, []string{" A "}},
	}

	for i, c := range cs {
		a := ToUppers(c.s)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] ToUppers(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestTrimSpaces(t *testing.T) {
	cs := []struct {
		s []string
		w []string
	}{
		{[]string{" a "}, []string{"a"}},
	}

	for i, c := range cs {
		a := TrimSpaces(c.s)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] TrimSpaces(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestRemoveEmpties(t *testing.T) {
	cs := []struct {
		s []string
		w []string
	}{
		{[]string{"a", ""}, []string{"a"}},
		{[]string{"a", "", "b", "", "c"}, []string{"a", "b", "c"}},
	}

	for i, c := range cs {
		a := RemoveEmpties(c.s)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] RemoveEmpties(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestJoinInts(t *testing.T) {
	cs := []struct {
		s []int
		f []func(int) string
		w string
	}{
		{[]int{}, nil, ""},
		{[]int{1}, nil, "1"},
		{[]int{1, 2}, nil, "1 2"},
		{[]int{1, 10}, []func(int) string{func(i int) string { return fmt.Sprintf("0x%x", i) }}, "0x1 0xa"},
	}

	for i, c := range cs {
		a := JoinInts(c.s, " ", c.f...)
		if a != c.w {
			t.Errorf("[%d] JoinInts(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestJoinInt64s(t *testing.T) {
	cs := []struct {
		s []int64
		f []func(int64) string
		w string
	}{
		{[]int64{}, nil, ""},
		{[]int64{1}, nil, "1"},
		{[]int64{1, 2}, nil, "1 2"},
		{[]int64{1, 10}, []func(int64) string{func(i int64) string { return fmt.Sprintf("0x%x", i) }}, "0x1 0xa"},
	}

	for i, c := range cs {
		a := JoinInt64s(c.s, " ", c.f...)
		if a != c.w {
			t.Errorf("[%d] JoinInt64s(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}
