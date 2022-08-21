package sch

import (
	"reflect"
	"runtime"
	"time"

	"github.com/pandafw/pango/log"
)

const execTimeFormat = "2006-01-02T15:04:05"

type Task struct {
	Logger         log.Logger // Error logger
	Trigger        Trigger
	Callback       func()
	ScheduledTime  time.Time
	ExecutionTime  time.Time
	CompletionTime time.Time
	Error          any
	timer          *time.Timer
}

func (t *Task) run() {
	defer func() {
		if err := recover(); err != nil {
			t.Error = err
			if t.Logger != nil {
				t.Logger.Errorf("Task error %s: %v", nameOfCallback(t.Callback), err)
			}
		}
	}()

	t.ExecutionTime = time.Now()
	t.Callback()
	t.CompletionTime = time.Now()

	if t.Error == nil {
		st := t.Trigger.NextExecutionTime(t)
		if !st.IsZero() {
			t.ScheduledTime = st
			t.schedule()
		}
	}
}

func nameOfCallback(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func (t *Task) schedule() {
	if t.Logger != nil {
		t.Logger.Infof("Schedule task at [%s] --> %s", t.ScheduledTime.Format(execTimeFormat), nameOfCallback(t.Callback))
	}
	t.timer.Reset(time.Until(t.ScheduledTime))
}

func (t *Task) Start() {
	t.ScheduledTime = t.Trigger.NextExecutionTime(t)

	// create a fake timer to get timer instance
	t.timer = time.AfterFunc(time.Hour, func() {
		t.run()
	})

	// reset timer
	t.schedule()
}

func (t *Task) Stop() bool {
	if t.timer != nil {
		return false
	}
	return t.timer.Stop()
}
