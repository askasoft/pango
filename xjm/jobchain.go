package xjm

import (
	"time"
)

const (
	JobChainAborted   = "A"
	JobChainCompleted = "C"
	JobChainPending   = "P"
	JobChainRunning   = "R"
)

var (
	JobChainPendingRunning   = []string{JobChainPending, JobChainRunning}
	JobChainAbortedCompleted = []string{JobChainAborted, JobChainCompleted}
)

type JobChain struct {
	ID        int64     `gorm:"not null;primaryKey;autoIncrement" form:"id" json:"id"`
	Name      string    `gorm:"size:250;not null" form:"title,strip" json:"title"`
	Status    string    `gorm:"size:1;not null" form:"status,strip" json:"status"`
	States    string    `gorm:"not null" form:"states" json:"states,omitempty"`
	CreatedAt time.Time `gorm:"not null;<-:create" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

func (jc *JobChain) String() string {
	return toString(jc)
}
