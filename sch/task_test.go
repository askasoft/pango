package sch

import (
	"fmt"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
)

const testTimeFormat = "2006-01-02T15:04:05.000"

func TestDelayedTriggerTask(t *testing.T) {
	cnt := 0
	task := Task{
		Logger:  log.NewLog().GetLogger("TASK"),
		Trigger: &DelayedTrigger{Delay: time.Millisecond * 200},
		Callback: func() {
			cnt++
			fmt.Printf("[%d] call %s\n", cnt, time.Now().Format(testTimeFormat))
			panic("panic test")
		},
	}
	task.Start()

	time.Sleep(time.Millisecond * 600)
	if cnt != 1 {
		t.Errorf("task execute count = %d, want 1", cnt)
	}
	task.Stop()
}

func TestRepeatTriggerTask(t *testing.T) {
	cnt := 0
	task := Task{
		Logger:  log.NewLog().GetLogger("TASK"),
		Trigger: &RepeatTrigger{Duration: time.Millisecond * 500},
		Callback: func() {
			cnt++
			fmt.Printf("[%d] call %s\n", cnt, time.Now().Format(testTimeFormat))
		},
	}
	task.Start()

	time.Sleep(time.Millisecond * 600)
	if cnt != 2 {
		t.Errorf("task execute count = %d, want 2", cnt)
	}
	task.Stop()
}

func TestRepeatTriggerTaskStop(t *testing.T) {
	cnt := 0
	task := Task{
		Logger:  log.NewLog().GetLogger("TASK"),
		Trigger: &RepeatTrigger{Duration: time.Millisecond * 500},
		Callback: func() {
			cnt++
			fmt.Printf("[%d] call %s\n", cnt, time.Now().Format(testTimeFormat))
			time.Sleep(time.Second)
		},
	}
	task.Start()

	time.Sleep(time.Millisecond * 1100)
	fmt.Println("stop timer")
	task.Stop()

	time.Sleep(time.Millisecond * 1500)
	if cnt != 1 {
		t.Errorf("task execute count = %d, want 1", cnt)
	}
}

func TestRepeatTriggerTaskStopStartInRunning(t *testing.T) {
	cnt := 0
	task := Task{
		Logger:  log.NewLog().GetLogger("TASK"),
		Trigger: &RepeatTrigger{Duration: time.Millisecond * 300},
		Callback: func() {
			cnt++
			fmt.Printf("[%d] begin %s\n", cnt, time.Now().Format(testTimeFormat))
			time.Sleep(time.Millisecond * 900)
			fmt.Printf("[%d]  end  %s\n", cnt, time.Now().Format(testTimeFormat))
			time.Sleep(time.Millisecond * 100)
		},
	}
	task.Start()

	time.Sleep(time.Millisecond * 100)
	fmt.Println("stop timer in running")
	task.Stop()

	fmt.Println("start timer in running")
	task.Start()

	time.Sleep(time.Millisecond * 1000)
	if cnt != 1 {
		t.Errorf("task execute count = %d, want 1", cnt)
	}

	task.Stop()
	time.Sleep(time.Millisecond * 100)
}

func TestRepeatTriggerTaskStopStartRunned(t *testing.T) {
	cnt := 0
	task := Task{
		Logger:  log.NewLog().GetLogger("TASK"),
		Trigger: &RepeatTrigger{Duration: time.Hour},
		Callback: func() {
			cnt++
			fmt.Printf("[%d] call %s\n", cnt, time.Now().Format(testTimeFormat))
		},
	}
	task.Start()

	time.Sleep(time.Millisecond * 100)

	fmt.Println("stop timer after callback")
	task.Stop()
	if cnt != 1 {
		t.Errorf("task execute count = %d, want 1", cnt)
	}

	fmt.Println("start timer again")
	task.Start()

	time.Sleep(time.Millisecond * 100)
	if cnt != 2 {
		t.Errorf("task execute count = %d, want 2", cnt)
	}

	task.Stop()
	time.Sleep(time.Millisecond * 100)
}

func TestFixedRateRepeatTriggerTask(t *testing.T) {
	cnt := 0
	task := Task{
		Logger:  log.NewLog().GetLogger("TASK"),
		Trigger: &RepeatTrigger{Duration: time.Millisecond * 400, FixedRate: true},
		Callback: func() {
			cnt++
			fmt.Printf("[%d] call start %s\n", cnt, time.Now().Format(testTimeFormat))
			time.Sleep(time.Second)
			fmt.Printf("[%d] call end %s\n", cnt, time.Now().Format(testTimeFormat))
		},
	}
	task.Start()

	time.Sleep(time.Millisecond * 1100)
	fmt.Println("stop timer")
	task.Stop()

	time.Sleep(time.Millisecond * 1500)
	if cnt != 1 {
		t.Errorf("task execute count = %d, want 1", cnt)
	}
}
