package sch

import (
	"fmt"
	"testing"
	"time"

	"github.com/pandafw/pango/log"
)

func testPeriodicTriggerTask(t *testing.T, logger log.Logger) {
	cnt := 0
	task := Task{
		Logger:  logger,
		Trigger: &PeriodicTrigger{Period: time.Second},
		Callback: func() {
			cnt++
			fmt.Printf("[%d] call %s\n", cnt, time.Now().Format(testTimeFormat))
			if cnt > 1 {
				panic("panic for stop task")
			}
		},
	}
	task.Start()

	time.Sleep(time.Second * 3)
	if cnt != 2 {
		t.Errorf("task execute count = %d, want 2", cnt)
	}
	if task.Stop() {
		t.Error("task timer still active")
	}
}

func TestPeriodicTriggerTask(t *testing.T) {
	testPeriodicTriggerTask(t, nil)
}

func TestPeriodicTriggerTaskWithLog(t *testing.T) {
	testPeriodicTriggerTask(t, log.NewLog().GetLogger("TASK"))
}
