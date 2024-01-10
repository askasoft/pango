package sch

import "time"

type Trigger interface {
	NextExecutionTime(task *Task) time.Time
}

var ZeroTrigger = &zeroTrigger{}

type zeroTrigger struct {
}

func (zt *zeroTrigger) NextExecutionTime(task *Task) (zero time.Time) {
	return
}

type DelayedTrigger struct {
	Delay time.Duration
}

func (dt *DelayedTrigger) NextExecutionTime(task *Task) time.Time {
	if task.ScheduledTime.IsZero() {
		return time.Now().Add(dt.Delay)
	}

	// disable next run
	return time.Time{}
}

type PeriodicTrigger struct {
	Period       time.Duration
	InitialDelay time.Duration
	FixedRate    bool
}

func (pt *PeriodicTrigger) NextExecutionTime(task *Task) time.Time {
	if task.ScheduledTime.IsZero() {
		return time.Now().Add(pt.InitialDelay)
	}

	if pt.FixedRate {
		for {
			st := task.ScheduledTime.Add(pt.Period)
			if st.After(task.CompletionTime) {
				return st
			}
		}
	}

	date := task.CompletionTime
	if date.IsZero() {
		date = time.Now()
	}
	return date.Add(pt.Period)
}

type CronTrigger struct {
	CronSequencer
}

func (ct *CronTrigger) NextExecutionTime(task *Task) time.Time {
	date := task.CompletionTime
	if date.IsZero() {
		date = time.Now()
	} else {
		if !task.ScheduledTime.IsZero() && date.Before(task.ScheduledTime) {
			// Previous task apparently executed too early...
			// Let's simply use the last calculated execution time then,
			// in order to prevent accidental re-fires in the same second.
			date = task.ScheduledTime
		}
	}

	return ct.Next(date)
}

func NewCronTrigger(cron string) *CronTrigger {
	ct := &CronTrigger{}
	if err := ct.Parse(cron); err != nil {
		panic(err)
	}
	return ct
}
