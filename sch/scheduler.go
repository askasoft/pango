package sch

import (
	"sync"
)

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
func (s *Scheduler) Schedule(trigger Trigger, callback func()) {
	task := &Task{
		Trigger:  trigger,
		Callback: callback,
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
