package sch

import (
	"sync"
)

// default scheduler instance
var _sch = Scheduler{}

// Scheduler task scheduler
type Scheduler struct {
	tasks []*Task
	mutex sync.Mutex
}

func (s *Scheduler) addTask(task *Task) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.tasks = append(s.tasks, task)
}

// Schedule schedule a task
func (s *Scheduler) Schedule(callback func(), trigger Trigger) {
	task := &Task{
		Callback: callback,
		Trigger:  trigger,
	}

	task.Start()

	s.addTask(task)
}

// Shutdown stop all task
func (s *Scheduler) Shutdown() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, t := range s.tasks {
		t.Stop()
	}
}

//----------------------------------------------------
// package functions

// Schedule schedule a task
func Schedule(callback func(), trigger Trigger) {
	_sch.Schedule(callback, trigger)
}

// Shutdown stop all task
func Shutdown() {
	_sch.Shutdown()
}
