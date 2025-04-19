package snowflake

import (
	"testing"
	"time"
)

func TestNewNode0(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("error creating NewNode, %s", err)
		}
	}()
	_ = NewNode(0)
}

func TestNewNode5000(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("no error creating NewNode, %s", err)
		}
	}()
	_ = NewNode(5000)
}

func TestNext(t *testing.T) {
	node := NewNode(1)

	x := node.NextID()
	for range 1000000 {
		now := time.Now()
		y := node.NextID()
		if y.Node() != 1 {
			t.Fatal("node != 1")
		}
		if y.UnixMilli() < x.UnixMilli() {
			t.Fatal("unixmilli error")
		}
		if y.UnixMilli() == x.UnixMilli() && y.Step()-x.Step() != 1 {
			t.Fatal("step error")
		}

		if x.Int64() == y.Int64() {
			t.Fatalf("x(%d) & y(%d) are the same", x.Int64(), y.Int64())
		}

		d := y.Time().Sub(now)
		if d > 3*time.Millisecond || d < -3*time.Millisecond {
			t.Errorf("time error: %d", d/time.Millisecond)
		}

		x = y
	}
}

// I feel like there's probably a better way
func TestRace(t *testing.T) {
	node := NewNode(1)

	go func() {
		for range 1000000000 {
			NewNode(1)
		}
	}()

	for range 4000 {
		node.NextID()
	}
}
