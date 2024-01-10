package sch

import (
	"reflect"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/askasoft/pango/log"
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
			if t.Logger != nil {
				t.Logger.Errorf("Task '%s' %s() run error: %v", t.Name, nameOfCallback(t.Callback), err)
			}
		}
	}()

	if t.Logger != nil {
		t.Logger.Debugf("Task '%s' %s() start at %v", t.Name, nameOfCallback(t.Callback), t.ExecutionTime)
	}

	t.Callback()
}

func (t *Task) run() {
	cnt := atomic.AddInt32(&t.count, 1)
	if cnt > 1 {
		atomic.AddInt32(&t.count, -1)
		if t.Logger != nil {
			t.Logger.Warnf("Task '%s' %s() is count at %v, SKIP!", t.Name, nameOfCallback(t.Callback), t.ExecutionTime)
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

func nameOfCallback(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func (t *Task) schedule() {
	timer := t.timer
	if timer != nil && !t.ScheduledTime.IsZero() {
		if t.Logger != nil {
			t.Logger.Infof("Schedule task '%s' %s() at %v", t.Name, nameOfCallback(t.Callback), t.ScheduledTime)
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
		if t.Logger != nil {
			t.Logger.Infof("Stop task '%s' %s() at %v", t.Name, nameOfCallback(t.Callback), t.ScheduledTime)
		}
		t.timer = nil
		return timer.Stop()
	}
	return false
}
