package ars

import (
	"math"
	"testing"
)

var equalIntTests = []struct {
	s1, s2 []int
	want   bool
}{
	{
		[]int{1},
		nil,
		false,
	},
	{
		[]int{},
		nil,
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3},
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		false,
	},
}

var equalFloatTests = []struct {
	s1, s2       []float64
	wantEqual    bool
	wantEqualNaN bool
}{
	{
		[]float64{1, 2},
		[]float64{1, 2},
		true,
		true,
	},
	{
		[]float64{1, 2, math.NaN()},
		[]float64{1, 2, math.NaN()},
		false,
		true,
	},
}

func TestEqualInts(t *testing.T) {
	for _, test := range equalIntTests {
		if got := EqualInts(test.s1, test.s2); got != test.want {
			t.Errorf("EqualInts(%v, %v) = %t, want %t", test.s1, test.s2, got, test.want)
		}
	}
}

func TestEqualFloat64s(t *testing.T) {
	for _, test := range equalFloatTests {
		if got := EqualFloat64s(test.s1, test.s2); got != test.wantEqual {
			t.Errorf("EqualFloat64s(%v, %v) = %t, want %t", test.s1, test.s2, got, test.wantEqual)
		}
	}
}

var indexTests = []struct {
	s    []int
	v    int
	want int
}{
	{
		nil,
		0,
		-1,
	},
	{
		[]int{},
		0,
		-1,
	},
	{
		[]int{1, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 3, 2},
		2,
		1,
	},
}

func TestIndexInt(t *testing.T) {
	for _, test := range indexTests {
		if got := IndexInt(test.s, test.v); got != test.want {
			t.Errorf("IndexInt(%v, %v) = %d, want %d", test.s, test.v, got, test.want)
		}
	}
}

func TestContainsInt(t *testing.T) {
	for _, test := range indexTests {
		if got := ContainsInt(test.s, test.v); got != (test.want != -1) {
			t.Errorf("ContainsInt(%v, %v) = %t, want %t", test.s, test.v, got, test.want != -1)
		}
	}
}

var deleteTests = []struct {
	s    []int
	i, j int
	want []int
}{
	{
		[]int{1, 2, 3},
		0,
		0,
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		0,
		1,
		[]int{2, 3},
	},
	{
		[]int{1, 2, 3},
		3,
		3,
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		0,
		2,
		[]int{3},
	},
	{
		[]int{1, 2, 3},
		0,
		3,
		[]int{},
	},
}

func TestDeleteInts(t *testing.T) {
	for _, test := range deleteTests {
		copy := append([]int{}, test.s...)
		if got := DeleteInts(copy, test.i, test.j); !EqualInts(got, test.want) {
			t.Errorf("DeleteInts(%v, %d, %d) = %v, want %v", test.s, test.i, test.j, got, test.want)
		}
	}
}
