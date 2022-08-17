package sch

// default scheduler instance
var _sch = &Scheduler{}

//----------------------------------------------------
// package functions

// Default returns the default Scheduler instance used by the package-level functions.
func Default() *Scheduler {
	return _sch
}

// Schedule schedule a task
func Schedule(trigger Trigger, callback func()) {
	_sch.Schedule(trigger, callback)
}

// Shutdown stop all task
func Shutdown() {
	_sch.Shutdown()
}
