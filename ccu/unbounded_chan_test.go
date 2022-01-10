package ccu

import (
	"sync"
	"testing"
	"time"
)

func TestUnboundedChanPollEmpty(t *testing.T) {
	rb := NewUnboundedChan(10)

	if _, ok := rb.Poll(); ok {
		t.Error("should return false when polling empty queue")
	}

	rb.Push(1)
	time.Sleep(time.Millisecond * 100)
	if _, ok := rb.Poll(); !ok {
		t.Error("should return true when polling non-empty queue")
	}

	if _, ok := rb.Poll(); ok {
		t.Error("should return false when polling empty queue")
	}
}

func TestUnboundedChanSerial(t *testing.T) {
	ubc := NewUnboundedChan(5)

	for i := 0; i < 1000; i++ {
		ubc.In <- i
	}

	for i := 0; i < 1000; i++ {
		v, ok := <-ubc.Out
		if !ok || v.(int) != i {
			t.Errorf("%d = %v, want %v", i, v, i)
		}
	}

	ubc.Close()
}

func TestUnboundedChanParallel(t *testing.T) {
	ubc := NewUnboundedChan(5)

	cnt := 10_000

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for i := 0; i < cnt; i++ {
			ubc.In <- i
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for i := 0; i < cnt; i++ {
			v, ok := <-ubc.Out
			if !ok || v.(int) != i {
				t.Errorf("%d = %v, want %v", i, v, i)
			}
		}
		wg.Done()
	}()

	wg.Wait()
	ubc.Close()
}
