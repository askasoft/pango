package xjm

import (
	"time"
)

type JobChainer interface {
	// GetJobChain get job chain detail
	GetJobChain(cid int64) (*JobChain, error)

	// FindJobChain find the latest job chain by name.
	// status: status to filter.
	FindJobChain(name string, asc bool, status ...string) (*JobChain, error)

	// FindJobChains find job chains by name.
	// status: status to filter.
	FindJobChains(name string, start, limit int, asc bool, status ...string) ([]*JobChain, error)

	// IterJobChains find job chains by name and iterate.
	// status: status to filter.
	IterJobChains(it func(*JobChain) error, name string, start, limit int, asc bool, status ...string) error

	// CreateJobChain create a job chain
	CreateJobChain(name, states string) (int64, error)

	// UpcateJobChain update the job chain, ignore empty status, states
	UpdateJobChain(cid int64, status string, states ...string) error

	// CleanOutdatedJobChains delete outdated job chains
	CleanOutdatedJobChains(before time.Time) (int64, error)
}
