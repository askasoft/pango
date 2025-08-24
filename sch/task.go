package sch

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/ref"
)

type Task struct {
	Name           string
	Logger         log.Logger
	Trigger        Trigger
	Callback       func()
	ScheduledTime  time.Time
	ExecutionTime  time.Time
	CompletionTime time.Time
	Error          any

	mutex sync.Mutex
	timer *time.Timer
	count int32
}

func NewTask(name string, trigger Trigger, callback func(), logger ...log.Logger) *Task {
	return &Task{
		Name:     name,
		Trigger:  trigger,
		Callback: callback,
		Logger:   asg.First(logger),
	}
}

func (t *Task) callback() {
	t.Error = nil

	defer func() {
		if err := recover(); err != nil {
			t.Error = err
			if log := t.Logger; log != nil {
				log.Errorf("Task %q %s() run error: %v", t.Name, ref.NameOfFunc(t.Callback), err)
			}
		}
	}()

	if log := t.Logger; log != nil {
		log.Debugf("Task %q %s() start at %s", t.Name, ref.NameOfFunc(t.Callback), t.ExecutionTime.Format(time.RFC3339))
	}

	t.Callback()
}

func (t *Task) run() {
	cnt := atomic.AddInt32(&t.count, 1)
	if cnt > 1 {
		atomic.AddInt32(&t.count, -1)
		if log := t.Logger; log != nil {
			log.Warnf("Task %q %s() is running at %s, SKIP!", t.Name, ref.NameOfFunc(t.Callback), t.ExecutionTime.Format(time.RFC3339))
		}
		return
	}

	t.ExecutionTime = time.Now()
	t.callback()
	t.CompletionTime = time.Now()

	t.start(false)

	atomic.AddInt32(&t.count, -1)
}

func (t *Task) Start() {
	t.start(true)
}

func (t *Task) start(force bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if !force && t.timer == nil {
		return
	}

	t.ScheduledTime = t.Trigger.NextExecutionTime(t)
	if t.ScheduledTime.IsZero() {
		t.timer = nil
		return
	}

	if log := t.Logger; log != nil {
		log.Infof("Schedule task %q %s() at %s", t.Name, ref.NameOfFunc(t.Callback), t.ScheduledTime.Format(time.RFC3339))
	}

	if t.timer == nil {
		t.timer = time.AfterFunc(time.Until(t.ScheduledTime), t.run)
	} else {
		t.timer.Reset(time.Until(t.ScheduledTime))
	}
}

func (t *Task) Stop() bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if timer := t.timer; timer != nil {
		if log := t.Logger; log != nil {
			log.Infof("Stop task %q %s() at %s", t.Name, ref.NameOfFunc(t.Callback), t.ScheduledTime.Format(time.RFC3339))
		}

		t.timer = nil
		t.ScheduledTime = time.Time{}
		return timer.Stop()
	}
	return false
}
