package sch

import (
	"time"
)

type Task struct {
	ScheduledTime  time.Time
	ExecutionTime  time.Time
	CompletionTime time.Time
	Callback       func()
	Trigger        Trigger
	Error          any
	timer          *time.Timer
}

func (t *Task) run() {
	defer func() {
		if err := recover(); err != nil {
			t.Error = err
		}
	}()

	t.ExecutionTime = time.Now()
	t.Callback()
	t.CompletionTime = time.Now()

	if t.Error == nil {
		st := t.Trigger.NextExecutionTime(t)
		if !st.IsZero() {
			t.ScheduledTime = st
			t.timer.Reset(time.Until(t.ScheduledTime))
		}
	}
}

func (t *Task) Start() {
	t.ScheduledTime = t.Trigger.NextExecutionTime(t)
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
