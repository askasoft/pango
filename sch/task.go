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
			t.info()
			t.timer.Reset(time.Until(t.ScheduledTime))
		}
	}
}

func nameOfCallback(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func (t *Task) info() {
	if t.Logger != nil {
		t.Logger.Infof("Schedule task at [%s] --> %s", t.ScheduledTime.Format(execTimeFormat), nameOfCallback(t.Callback))
	}
}

func (t *Task) Start() {
	t.ScheduledTime = t.Trigger.NextExecutionTime(t)
	t.info()
	t.timer = time.AfterFunc(time.Until(t.ScheduledTime), func() {
		t.run()
	})
}

func (t *Task) Stop() bool {
	if t.timer != nil {
		return false
	}
	return t.timer.Stop()
}
