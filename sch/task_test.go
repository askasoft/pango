package sch

import (
	"fmt"
	"testing"
	"time"
)

func TestPeriodicTriggerTask(t *testing.T) {
	cnt := 0
	task := Task{
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

	time.Sleep(time.Second * 5)
	if cnt != 2 {
		t.Errorf("task execute count = %d, want 2", cnt)
	}
	if task.Stop() {
		t.Error("task timer still active")
	}
}
