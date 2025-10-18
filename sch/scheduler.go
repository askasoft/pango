package sch

import (
	"sync"

	"github.com/askasoft/pango/cog/treemap"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/str"
)

// Scheduler task scheduler
type Scheduler struct {
	Logger log.Logger
	tasks  *treemap.TreeMap[string, *Task]
	mutex  sync.Mutex
}

// GetTask get task by task name
func (s *Scheduler) GetTask(name string) (*Task, bool) {
	return s.tasks.Get(name)
}

// AddTask add or replace a task with same name
func (s *Scheduler) AddTask(task *Task) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.tasks == nil {
		s.tasks = treemap.NewTreeMap[string, *Task](str.Compare)
	}
	s.tasks.Set(task.Name, task)
}

// RemoveTask remove a task by name
// Returns the removed task and whether it was found.
func (s *Scheduler) RemoveTask(name string) (*Task, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.tasks != nil {
		return s.tasks.Remove(name)
	}
	return nil, false
}

// Schedule schedule a task and start it
func (s *Scheduler) Schedule(name string, trigger Trigger, callback func()) {
	task := NewTask(name, trigger, callback, s.Logger)

	task.Start()

	s.AddTask(task)
}

// Start start all tasks
func (s *Scheduler) Start() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for it := s.tasks.Iterator(); it.Next(); {
		it.Value().Start()
	}
}

// Stop stop all tasks
func (s *Scheduler) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for it := s.tasks.Iterator(); it.Next(); {
		it.Value().Stop()
	}
}
