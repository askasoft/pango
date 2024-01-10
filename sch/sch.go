package sch

// default scheduler instance
var _sch = &Scheduler{}

//----------------------------------------------------
// package functions

// Default returns the default Scheduler instance used by the package-level functions.
func Default() *Scheduler {
	return _sch
}

// Schedule schedule a task and start it
func Schedule(name string, trigger Trigger, callback func()) {
	_sch.Schedule(name, trigger, callback)
}

// Start start all task
func Start() {
	_sch.Start()
}

// Stop stop all task
func Stop() {
	_sch.Stop()
}

// GetTask get task by task name
func GetTask(name string) (*Task, bool) {
	return _sch.GetTask(name)
}

// AddTask add a task
func AddTask(task *Task) {
	_sch.AddTask(task)
}
