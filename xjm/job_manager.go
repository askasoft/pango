package xjm

import (
	"time"
)

type JobManager interface {
	// CountJobLogs count job logs
	CountJobLogs(jid int64, levels ...string) (int64, error)

	// GetJobLogs get job logs
	// set levels to ("I", "W", "E", "F") to filter DEBUG/TRACE logs
	GetJobLogs(jid int64, minLid, maxLid int64, asc bool, limit int, levels ...string) ([]*JobLog, error)

	// AddJobLogs append job logs
	AddJobLogs([]*JobLog) error

	// AddJobLog append a job log
	AddJobLog(jid int64, time time.Time, level string, message string) error

	// GetJob get a job
	GetJob(jid int64) (*Job, error)

	// FindJob find a job
	// name: name to filter (optional)
	// status: status to filter (optional)
	FindJob(name string, asc bool, status ...string) (*Job, error)

	// FindJobs find jobs
	// name: name to filter (optional)
	// status: status to filter (optional)
	FindJobs(name string, start, limit int, asc bool, status ...string) ([]*Job, error)

	// IterJobs find jobs and iterate
	// name: name to filter (optional)
	// status: status to filter (optional)
	IterJobs(it func(job *Job) error, name string, start, limit int, asc bool, status ...string) error

	// AppendJob append a pendding job
	AppendJob(name, file, param string) (int64, error)

	// AbortJob abort the job
	AbortJob(jid int64, reason string) error

	// CompleteJob update job status to completed
	CompleteJob(jid int64) error

	// CheckoutJob change job status from pending to running
	CheckoutJob(jid, rid int64) error

	// PingJob update the job updated_at to now
	PingJob(jid, rid int64) error

	// SetJobState update the running job state
	SetJobState(jid, rid int64, state string) error

	// AddJobResult append result to the running job
	AddJobResult(jid, rid int64, result string) error

	// ReappendJobs reappend the interrupted runnings job to the pennding status
	ReappendJobs(before time.Time) (int64, error)

	// StartJobs start to run jobs
	StartJobs(limit int, start func(*Job)) error

	// DeleteJobs delete jobs
	DeleteJobs(jids ...int64) (int64, int64, error)

	// CleanOutdatedJobs delete outdated jobs
	CleanOutdatedJobs(before time.Time) (int64, int64, error)
}
