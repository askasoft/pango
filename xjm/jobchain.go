package xjm

import (
	"time"

	"github.com/askasoft/pango/asg"
)

type JobChain struct {
	ID        int64     `gorm:"not null;primaryKey;autoIncrement" json:"id,omitempty"`
	Name      string    `gorm:"size:250;not null;index:idx_job_chains_name" json:"name,omitempty"`
	Status    string    `gorm:"size:1;not null" json:"status,omitempty"`
	States    string    `gorm:"not null" json:"states,omitempty"`
	CreatedAt time.Time `gorm:"not null;<-:create" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

func (jc *JobChain) IsAborted() bool {
	return jc.Status == JobStatusAborted
}

func (jc *JobChain) IsCanceled() bool {
	return jc.Status == JobStatusCanceled
}

func (jc *JobChain) IsFinished() bool {
	return jc.Status == JobStatusFinished
}

func (jc *JobChain) IsPending() bool {
	return jc.Status == JobStatusPending
}

func (jc *JobChain) IsRunning() bool {
	return jc.Status == JobStatusRunning
}

func (jc *JobChain) IsDone() bool {
	return asg.Contains(JobDoneStatus, jc.Status)
}

func (jc *JobChain) IsUndone() bool {
	return asg.Contains(JobUndoneStatus, jc.Status)
}

func (jc *JobChain) String() string {
	return toString(jc)
}
