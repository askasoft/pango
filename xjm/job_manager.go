package xjm

import (
	"time"
)

type JobManager interface {
	// CountJobLogs count job logs
	CountJobLogs(jid int64, levels ...string) (int64, error)

	// GetJobLogs get job logs
	// set levels to ("I", "W", "E", "F") to filter DEBUG/TRACE logs
	GetJobLogs(jid int64, min, max int64, asc bool, limit int, levels ...string) ([]*JobLog, error)

	// AddJobLogs append job logs
	AddJobLogs([]*JobLog) error

	// GetJob get job detail
	GetJob(jid int64) (*Job, error)

	// FindJob find the latest job by name, default select all columns.
	// cols: columns to select.
	FindJob(name string, cols ...string) (*Job, error)

	// FindJobs find jobs by name, default select all columns.
	// cols: columns to select.
	FindJobs(name string, start, limit int, cols ...string) ([]*Job, error)

	// AppendJob append a pendding job
	AppendJob(name, file, param string) (int64, error)

	// FindAndAbortJob find latest running or pendding job to abort
	FindAndAbortJob(name, reason string) error

	// AbortJob abort the job
	AbortJob(jid int64, reason string) error

	// CompleteJob complete the job
	CompleteJob(jid int64, result string) error

	// CheckoutJob checkout the job to the running status
	CheckoutJob(jid, rid int64) error

	// PingJob update the job updated_at to now
	PingJob(jid, rid int64) error

	// RunningJob update the running job state
	RunningJob(jid, rid int64, state string) error

	// ReappendJobs reappend the interrupted runnings job to the pennding status
	ReappendJobs(before time.Time) (int64, error)

	// StartJobs start to run jobs
	StartJobs(limit int, run func(*Job)) error

	// CleanOutdatedJobs delete outdated jobs
	CleanOutdatedJobs(before time.Time) (int64, int64, error)
}
