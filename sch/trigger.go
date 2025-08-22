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

type RepeatTrigger struct {
	Duratiuon    time.Duration
	InitialDelay time.Duration
	FixedRate    bool
}

func (rt *RepeatTrigger) NextExecutionTime(task *Task) time.Time {
	if task.ScheduledTime.IsZero() {
		return time.Now().Add(rt.InitialDelay)
	}

	if rt.FixedRate {
		for {
			st := task.ScheduledTime.Add(rt.Duratiuon)
			if st.After(task.CompletionTime) {
				return st
			}
		}
	}

	date := task.CompletionTime
	if date.IsZero() {
		date = time.Now()
	}
	return date.Add(rt.Duratiuon)
}

type CronTrigger struct {
	cron Cron
}

func (ct *CronTrigger) Cron() string {
	return ct.cron.String()
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

	return ct.cron.Next(date)
}

func NewCronTrigger(expr string, location ...*time.Location) (*CronTrigger, error) {
	cron, err := ParseCron(expr, location...)
	if err != nil {
		return nil, err
	}
	return &CronTrigger{cron}, nil
}

type PeriodicTrigger struct {
	periodic Periodic
	crontrig *CronTrigger
}

func (pt *PeriodicTrigger) Periodic() string {
	return pt.periodic.String()
}

func (pt *PeriodicTrigger) Cron() string {
	return pt.crontrig.Cron()
}

func (pt *PeriodicTrigger) NextExecutionTime(task *Task) time.Time {
	return pt.crontrig.NextExecutionTime(task)
}

func NewPeriodicTrigger(periodic string, location ...*time.Location) (*PeriodicTrigger, error) {
	p, err := ParsePeriodic(periodic)
	if err != nil {
		return nil, err
	}

	ct, err := NewCronTrigger(p.Cron(), location...)
	if err != nil {
		return nil, err
	}

	pt := &PeriodicTrigger{
		periodic: p,
		crontrig: ct,
	}
	return pt, nil
}
