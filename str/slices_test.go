package str

import (
	"reflect"
	"testing"
)

func TestToLowers(t *testing.T) {
	cs := []struct {
		s []string
		w []string
	}{
		{[]string{" A ", " B "}, []string{" a ", " b "}},
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
		{[]string{" a ", " b "}, []string{" A ", " B "}},
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
		{[]string{" a ", " b "}, []string{"a", "b"}},
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
		{[]string{"", "a", "b"}, []string{"a", "b"}},
		{[]string{"", "a", "b", "", "c"}, []string{"a", "b", "c"}},
	}

	for i, c := range cs {
		a := RemoveEmpties(c.s)
		if !reflect.DeepEqual(a, c.w) {
			t.Errorf("[%d] RemoveEmpties(%v) = %v, want %v", i, c.s, a, c.w)
		}
	}
}
