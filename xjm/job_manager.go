package xjm

import (
	"time"
)

type JobManager interface {
	// GetJobLogs get job logs
	// set levels to ("I", "W", "E", "F") to filter DEBUG/TRACE logs
	GetJobLogs(jid int64, start, limit int, levels ...string) ([]*JobLog, error)

	AddJobLogs([]*JobLog) error

	// GetJob get job detail
	GetJob(jid int64) (*Job, error)

	// FindJob find job by name, default select all columns.
	// cols: columns to select.
	FindJob(name string, cols ...string) (*Job, error)

	// FindJobs find jobs by name, default select all columns.
	// cols: columns to select.
	FindJobs(name string, start, limit int, cols ...string) ([]*Job, error)

	AppendJob(name, file, param string) (int64, error)

	FindAndAbortJob(name, reason string) error

	AbortJob(jid int64, reason string) error

	CompleteJob(jid int64, result string) error

	CheckoutJob(jid, rid int64) error

	PingJob(jid, rid int64) error

	RunningJob(jid, rid int64, result string) error

	ReappendJobs(before time.Time) (int64, error)

	StartJobs(limit int, run func(*Job)) error

	CleanOutdatedJobs(before time.Time) (int64, int64, error)
}
