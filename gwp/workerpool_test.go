package gwp

import (
	"sync"
	"testing"
	"time"
)

const max = 20

func newTestWorkerPool(maxWorks, maxWaits int) *WorkerPool {
	return NewWorkerPool(maxWorks, maxWaits)
}

func TestExample(t *testing.T) {
	wp := newTestWorkerPool(2, 0)
	requests := []string{"alpha", "beta", "gamma", "delta", "epsilon"}

	rspChan := make(chan string, len(requests))
	for _, r := range requests {
		r := r
		wp.Submit(func() {
			rspChan <- r
		})
	}

	wp.StopWait()

	close(rspChan)
	rspSet := map[string]struct{}{}
	for rsp := range rspChan {
		rspSet[rsp] = struct{}{}
	}
	if len(rspSet) < len(requests) {
		t.Fatal("Did not handle all requests")
	}
	for _, req := range requests {
		if _, ok := rspSet[req]; !ok {
			t.Fatal("Missing expected values:", req)
		}
	}
}

func TestMaxWorkers(t *testing.T) {
	wp := newTestWorkerPool(0, 0)
	wp.Stop()
	if wp.maxWorks != 1 {
		t.Fatal("should have created one worker")
	}

	wp = newTestWorkerPool(max, 0)
	defer wp.Stop()

	if wp.maxWorks != max {
		t.Fatal("wrong size returned")
	}

	started := make(chan struct{}, max)
	release := make(chan struct{})

	// Start workers, and have them all wait on a channel before completing.
	for i := 0; i < max; i++ {
		wp.Submit(func() {
			started <- struct{}{}
			<-release
		})
	}

	timeout := time.After(5 * time.Second)
	for startCount := 0; startCount < max; {
		select {
		case <-started:
			startCount++
		case <-timeout:
			t.Fatal("timed out waiting for workers to start")
		}
	}

	// Release workers.
	close(release)
}

func TestReuseWorkers(t *testing.T) {
	wp := newTestWorkerPool(5, 0)
	defer wp.Stop()

	release := make(chan struct{})

	// Cause worker to be created, and available for reuse before next task.
	for i := 0; i < 10; i++ {
		wp.Submit(func() { <-release })
		release <- struct{}{}
		time.Sleep(time.Millisecond)
	}
	close(release)

	// If the same worker was always reused, then only one worker would have
	// been created and there should only be one ready.
	if countReady(wp) > 1 {
		t.Fatal("Worker not reused")
	}
}

func TestStop(t *testing.T) {
	// Start workers, and have them all wait on a channel before completing.
	wp := newTestWorkerPool(5, max)

	release := make(chan struct{})
	finished := make(chan struct{}, max)
	for i := 0; i < max; i++ {
		wp.Submit(func() {
			<-release
			finished <- struct{}{}
			time.Sleep(10 * time.Millisecond)
		})
	}

	// Call Stop() and see that only the already running tasks were completed.
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(release)
	}()

	time.Sleep(10 * time.Millisecond)
	wp.Stop()

	var count int
Count:
	for count < max {
		select {
		case <-finished:
			count++
		default:
			break Count
		}
	}
	if count > max-5 {
		t.Fatal("Should not have completed any queued tasks, did", count)
	}

	// Check that calling Stop() again is OK.
	wp.Stop()
}

func TestStopWait(t *testing.T) {
	// Start workers, and have them all wait on a channel before completing.
	wp := newTestWorkerPool(5, max)
	release := make(chan struct{})
	finished := make(chan struct{}, max)
	for i := 0; i < max; i++ {
		wp.Submit(func() {
			<-release
			finished <- struct{}{}
		})
	}

	// Call StopWait() and see that all tasks were completed.
	go func() {
		time.Sleep(10 * time.Millisecond)
		close(release)
	}()

	wp.StopWait()

	for count := 0; count < max; count++ {
		select {
		case <-finished:
		default:
			t.Fatal("Should have completed all queued tasks")
		}
	}

	if anyReady(wp) {
		t.Fatal("should have zero workers after stopwait")
	}

	if wp.Running() {
		t.Fatal("pool should be stopped")
	}
}

func TestStopWait2(t *testing.T) {
	// Make sure that calling StopWait() with no queued tasks is OK.
	wp := newTestWorkerPool(5, 0)
	wp.StopWait()

	if anyReady(wp) {
		t.Fatal("should have zero workers after stopwait")
	}

	// Check that calling StopWait() again is OK.
	wp.StopWait()
}

func TestSubmitWait(t *testing.T) {
	wp := newTestWorkerPool(1, 0)
	defer wp.Stop()

	// Check that these are noop.
	wp.Submit(nil)
	wp.SubmitWait(nil)

	done1 := make(chan struct{})
	wp.Submit(func() {
		time.Sleep(100 * time.Millisecond)
		close(done1)
	})
	select {
	case <-done1:
		t.Fatal("Submit did not return immediately")
	default:
	}

	done2 := make(chan struct{})
	wp.SubmitWait(func() {
		time.Sleep(100 * time.Millisecond)
		close(done2)
	})
	select {
	case <-done2:
	default:
		t.Fatal("SubmitWait did not wait for function to execute")
	}
}

func TestStopRace(t *testing.T) {
	wp := newTestWorkerPool(max, max)
	defer wp.Stop()

	workRelChan := make(chan struct{})

	var started sync.WaitGroup
	started.Add(max)

	// Start workers, and have them all wait on a channel before completing.
	for i := 0; i < max; i++ {
		wp.Submit(func() {
			started.Done()
			<-workRelChan
		})
	}

	started.Wait()

	const doneCallers = 5
	stopDone := make(chan struct{}, doneCallers)
	for i := 0; i < doneCallers; i++ {
		go func() {
			wp.Stop()
			stopDone <- struct{}{}
		}()
	}

	select {
	case <-stopDone:
		t.Fatal("Stop should not return in any goroutine")
	default:
	}

	close(workRelChan)

	timeout := time.After(time.Second)
	for i := 0; i < doneCallers; i++ {
		select {
		case <-stopDone:
		case <-timeout:
			wp.Stop()
			t.Fatal("timedout waiting for Stop to return")
		}
	}
}

func anyReady(w *WorkerPool) bool {
	release := make(chan struct{})
	wait := func() {
		<-release
	}
	select {
	case w.workChan <- wait:
		close(release)
		return true
	default:
	}
	return false
}

func countReady(w *WorkerPool) int {
	// Try to stop max workers.
	timeout := time.After(100 * time.Millisecond)
	release := make(chan struct{})
	wait := func() {
		<-release
	}
	var readyCount int
	for i := 0; i < max; i++ {
		select {
		case w.workChan <- wait:
			readyCount++
		case <-timeout:
			i = max
		}
	}

	close(release)
	return readyCount
}

/*

Run benchmarking with: go test -bench '.'

*/

func BenchmarkEnqueue(b *testing.B) {
	wp := newTestWorkerPool(1, 1)
	defer wp.Stop()
	releaseChan := make(chan struct{})

	b.ResetTimer()

	// Start workers, and have them all wait on a channel before completing.
	for i := 0; i < b.N; i++ {
		wp.Submit(func() { <-releaseChan })
	}
	close(releaseChan)
}

func BenchmarkEnqueue2(b *testing.B) {
	wp := newTestWorkerPool(2, 2)
	defer wp.Stop()

	b.ResetTimer()

	// Start workers, and have them all wait on a channel before completing.
	for i := 0; i < b.N; i++ {
		releaseChan := make(chan struct{})
		for i := 0; i < 64; i++ {
			wp.Submit(func() { <-releaseChan })
		}
		close(releaseChan)
	}
}

func BenchmarkExecute1Worker(b *testing.B) {
	benchmarkExecWorkers(1, b)
}

func BenchmarkExecute2Worker(b *testing.B) {
	benchmarkExecWorkers(2, b)
}

func BenchmarkExecute4Workers(b *testing.B) {
	benchmarkExecWorkers(4, b)
}

func BenchmarkExecute16Workers(b *testing.B) {
	benchmarkExecWorkers(16, b)
}

func BenchmarkExecute64Workers(b *testing.B) {
	benchmarkExecWorkers(64, b)
}

func BenchmarkExecute1024Workers(b *testing.B) {
	benchmarkExecWorkers(1024, b)
}

func benchmarkExecWorkers(n int, b *testing.B) {
	wp := newTestWorkerPool(n, n)
	defer wp.Stop()
	var allDone sync.WaitGroup
	allDone.Add(b.N * n)

	b.ResetTimer()

	// Start workers, and have them all wait on a channel before completing.
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			wp.Submit(func() {
				//time.Sleep(100 * time.Microsecond)
				allDone.Done()
			})
		}
	}
	allDone.Wait()
}
