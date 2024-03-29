package xjm

import (
	"time"
)

const (
	JobStatusAborted   = "A"
	JobStatusCompleted = "C"
	JobStatusPending   = "P"
	JobStatusRunning   = "R"
)

var (
	JobPendingRunning   = []string{JobStatusPending, JobStatusRunning}
	JobAbortedCompleted = []string{JobStatusAborted, JobStatusCompleted}
)

type Job struct {
	ID        int64     `gorm:"not null;primaryKey;autoIncrement" uri:"id" form:"id" json:"id,omitempty"`
	RID       int64     `gorm:"column:rid;not null" form:"rid" json:"rid,omitempty"` // job runner id
	Name      string    `gorm:"size:250;not null;index" json:"name,omitempty"`
	Status    string    `gorm:"size:1;not null" form:"status" json:"status,omitempty"`
	File      string    `gorm:"not null" json:"file,omitempty"`
	Param     string    `gorm:"not null" json:"param,omitempty"`
	State     string    `gorm:"not null" form:"state" json:"state,omitempty"`
	Result    string    `gorm:"not null" json:"result,omitempty"`
	Error     string    `gorm:"not null" json:"error,omitempty"`
	CreatedAt time.Time `gorm:"<-:create;not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

func (j *Job) IsAborted() bool {
	return j.Status == JobStatusAborted
}

func (j *Job) IsCompleted() bool {
	return j.Status == JobStatusCompleted
}

func (j *Job) IsPending() bool {
	return j.Status == JobStatusPending
}

func (j *Job) IsRunning() bool {
	return j.Status == JobStatusRunning
}

func (j *Job) String() string {
	return toString(j)
}

func (j *Job) Params() (m map[string]any) {
	return toMap(j.Param)
}

func (j *Job) Results() (m map[string]any) {
	return toMap(j.Result)
}
