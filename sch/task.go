package sch

import (
	"sync/atomic"
	"time"

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

	timer *time.Timer
	count int32
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

	st := t.Trigger.NextExecutionTime(t)
	if !st.IsZero() {
		t.ScheduledTime = st
		t.schedule()
	}

	atomic.AddInt32(&t.count, -1)
}

func (t *Task) schedule() {
	timer := t.timer
	if timer != nil && !t.ScheduledTime.IsZero() {
		if log := t.Logger; log != nil {
			log.Infof("Schedule task %q %s() at %s", t.Name, ref.NameOfFunc(t.Callback), t.ScheduledTime.Format(time.RFC3339))
		}

		d := time.Until(t.ScheduledTime)
		timer.Reset(d)
	}
}

func (t *Task) Start() {
	t.ScheduledTime = t.Trigger.NextExecutionTime(t)

	if t.timer == nil {
		// create a fake timer to get timer instance
		t.timer = time.AfterFunc(time.Hour, t.run)
	}

	// reset timer
	t.schedule()
}

func (t *Task) Stop() bool {
	timer := t.timer
	if timer != nil {
		if log := t.Logger; log != nil {
			log.Infof("Stop task %q %s() at %s", t.Name, ref.NameOfFunc(t.Callback), t.ScheduledTime.Format(time.RFC3339))
		}
		t.timer = nil
		return timer.Stop()
	}
	return false
}
