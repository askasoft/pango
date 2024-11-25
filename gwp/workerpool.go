package gwp

import (
	"runtime"
	"sync"
	"time"
)

// WorkerPool is a collection of goroutines, where the number of concurrent
// goroutines processing requests does not exceed the specified maximum.
type WorkerPool struct {
	*workerpool
}

type workerpool struct {
	curWorks    int
	maxWorks    int
	idleTimeout time.Duration
	taskChan    chan func()
	workChan    chan func()
	stopChan    chan bool
	slock       sync.Mutex
	running     bool
	waitg       sync.WaitGroup
}

// NewWorkerPool creates and starts a pool of worker goroutines.
//
// The maxWorks parameter specifies the maximum number of workers that can
// execute tasks concurrently. When there are no incoming tasks, workers are
// gradually stopped until there are no remaining workers.
func NewWorkerPool(maxWorks, maxWaits int) *WorkerPool {
	// There must be at least one worker.
	if maxWorks < 1 {
		maxWorks = 1
	}
	if maxWaits < 0 {
		maxWaits = 0
	}

	wp := &workerpool{
		maxWorks:    maxWorks,
		idleTimeout: 2 * time.Second,
		taskChan:    make(chan func(), maxWaits),
		workChan:    make(chan func()),
		stopChan:    make(chan bool, 2),
	}

	WP := &WorkerPool{wp}
	WP.Start()
	runtime.SetFinalizer(WP, finalStop)

	return WP
}

func finalStop(wp *WorkerPool) {
	wp.StopWait()
}

// MaxWaits returns the maximum number of concurrent workers.
func (wp *workerpool) MaxWorks() int {
	return wp.maxWorks
}

// SetMaxWaits set the maximum number of concurrent workers, panic if maxWorks < 1.
func (wp *workerpool) SetMaxWorks(maxWorks int) {
	if maxWorks < 1 {
		panic("WorkerPool: maxWorks must greater than 0")
	}
	wp.maxWorks = maxWorks
}

// SetIdleTimeout set the timeout to stop a idle worker, panic if timeout < 1ms.
func (wp *workerpool) SetIdleTimeout(timeout time.Duration) {
	if timeout < time.Millisecond {
		panic("WorkerPool: timeout must greater than 1ms")
	}
	wp.idleTimeout = timeout
}

// Start start the pool go-routine
func (wp *workerpool) Start() {
	wp.slock.Lock()
	defer wp.slock.Unlock()

	if !wp.running {
		wp.running = true
		wp.waitg.Add(1)
		go wp.run()
	}
}

// Running returns true if this worker pool is running.
func (wp *workerpool) Running() bool {
	wp.slock.Lock()
	defer wp.slock.Unlock()
	return wp.running
}

// Stop stops the worker pool and waits for only currently running tasks to
// complete. Pending tasks that are not currently running are abandoned. Tasks
// must not be submitted to the worker pool after calling stop.
//
// Since creating the worker pool starts at least one goroutine for the
// dispatcher, Stop() should be called when the worker pool is no longer needed.
func (wp *workerpool) Stop() {
	wp.stop(false)
}

// StopWait stops the worker pool and waits for all queued tasks tasks to
// complete. No additional tasks may be submitted, but all pending tasks are
// executed by workers before this function returns.
func (wp *workerpool) StopWait() {
	wp.stop(true)
}

func (wp *workerpool) stop(wait bool) {
	wp.slock.Lock()
	defer wp.slock.Unlock()

	if wp.running {
		wp.stopChan <- wait
		wp.waitg.Wait()
		wp.running = false
	}
}

// Submit enqueues a function for a worker to execute.
//
// Any external values needed by the task function must be captured in a
// closure. Any return values should be returned over a channel that is
// captured in the task function closure.
//
// Submit will block if the task wait channel is full.
//
// As long as no new tasks arrive, one available worker is shutdown each time
// period until there are no more idle workers. Since the time to start new
// go-routines is not significant, there is no need to retain idle workers
// indefinitely.
func (wp *workerpool) Submit(task func()) {
	if task != nil {
		wp.taskChan <- task
	}
}

// SubmitWait enqueues the given function and waits for it to be executed.
func (wp *workerpool) SubmitWait(task func()) {
	if task == nil {
		return
	}

	doneChan := make(chan struct{})
	wp.taskChan <- func() {
		defer close(doneChan)
		task()
	}
	<-doneChan
}

func (wp *workerpool) run() {
	timeout := time.NewTimer(wp.idleTimeout)

	idle, stop := false, false

	defer func() {
		// Stop all workers
		for wp.curWorks > 0 {
			wp.workChan <- nil
			wp.curWorks--
		}

		timeout.Stop()
		wp.waitg.Done()
	}()

	for {
		select {
		case task := <-wp.taskChan:
			// Got a task to do.
			wp.dispatch(task)
			idle = false
		case <-timeout.C:
			// Timed out waiting for work to arrive. Kill a ready worker if
			// pool has been idle for a whole timeout.
			if idle && wp.curWorks > 0 {
				if wp.killIdleWorker() {
					wp.curWorks--
				}
			}
			idle = true
			timeout.Reset(wp.idleTimeout)
		case wait := <-wp.stopChan:
			if !wait {
				return
			}
			stop = true
		}

		if stop && len(wp.taskChan) == 0 {
			return
		}
	}
}

// dispatch sends the next queued task to an available worker.
func (wp *workerpool) dispatch(task func()) {
	if wp.curWorks >= wp.maxWorks {
		// Dispatch task to work queue, if max workers have been created.
		wp.workChan <- task
		return
	}

	select {
	case wp.workChan <- task:
		// Attempt to dispatch the task to a idle workder.
		// If failed then goto the default case to create a new worker.
	default:
		// Create a new worker.
		wp.waitg.Add(1)
		go wp.worker(task)
		wp.curWorks++
	}
}

func (wp *workerpool) killIdleWorker() bool {
	select {
	case wp.workChan <- nil:
		// Sent kill signal to worker.
		return true
	default:
		// No ready workers. All, if any, workers are busy.
		return false
	}
}

// worker executes tasks and stops when it receives a nil task.
func (wp *workerpool) worker(task func()) {
	for task != nil {
		task()
		task = <-wp.workChan
	}
	wp.waitg.Done()
}
