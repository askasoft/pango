package sch

import (
	"cmp"
	"sync"

	"github.com/askasoft/pango/cog/treemap"
	"github.com/askasoft/pango/log"
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

// AddTask add a task
func (s *Scheduler) AddTask(task *Task) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.tasks == nil {
		s.tasks = treemap.NewTreeMap[string, *Task](cmp.Compare[string])
	}
	s.tasks.Set(task.Name, task)
}

// Schedule schedule a task and start it
func (s *Scheduler) Schedule(name string, trigger Trigger, callback func()) {
	task := &Task{
		Name:     name,
		Logger:   s.Logger,
		Trigger:  trigger,
		Callback: callback,
	}

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
