package xjm

import (
	"errors"
	"time"
)

var (
	ErrJobChainMissing = errors.New("jobchain missing")
)

type JobChainer interface {
	// GetJobChain get a job chain
	GetJobChain(cid int64) (*JobChain, error)

	// FindJobChain find a job chain
	// name: name to filter (optional)
	// status: status to filter (optional)
	FindJobChain(name string, asc bool, status ...string) (*JobChain, error)

	// FindJobChains find job chains
	// name: name to filter (optional)
	// status: status to filter (optional)
	FindJobChains(name string, start, limit int, asc bool, status ...string) ([]*JobChain, error)

	// IterJobChains find job chains and iterate
	// name: name to filter (optional)
	// status: status to filter (optional)
	IterJobChains(it func(*JobChain) error, name string, start, limit int, asc bool, status ...string) error

	// CreateJobChain create a job chain
	CreateJobChain(name, states string) (int64, error)

	// UpcateJobChain update the job chain, ignore empty status, states
	UpdateJobChain(cid int64, status string, states ...string) error

	// DeleteJobChains delete job chains
	DeleteJobChains(cids ...int64) (int64, error)

	// CleanOutdatedJobChains delete outdated job chains
	CleanOutdatedJobChains(before time.Time) (int64, error)
}
