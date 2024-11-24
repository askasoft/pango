package xjm

import (
	"time"

	"github.com/askasoft/pango/asg"
)

const (
	JobStatusAborted  = "A"
	JobStatusCanceled = "C"
	JobStatusFinished = "F"
	JobStatusPending  = "P"
	JobStatusRunning  = "R"
)

var (
	JobDoneStatus   = []string{JobStatusAborted, JobStatusCanceled, JobStatusFinished}
	JobUndoneStatus = []string{JobStatusPending, JobStatusRunning}
)

type Job struct {
	ID        int64     `gorm:"not null;primaryKey;autoIncrement" json:"id,omitempty"`
	CID       int64     `gorm:"column:cid;not null" json:"cid,omitempty"`
	RID       int64     `gorm:"column:rid;not null" json:"rid,omitempty"`
	Name      string    `gorm:"size:250;not null;index:idx_jobs_name" json:"name,omitempty"`
	Status    string    `gorm:"size:1;not null" json:"status,omitempty"`
	Locale    string    `gorm:"size:20;not null" json:"locale,omitempty"`
	File      string    `gorm:"not null" json:"file,omitempty"`
	Param     string    `gorm:"not null" json:"param,omitempty"`
	State     string    `gorm:"not null" form:"state" json:"state,omitempty"`
	Result    string    `gorm:"not null" json:"result,omitempty"`
	Error     string    `gorm:"not null" json:"error,omitempty"`
	CreatedAt time.Time `gorm:"not null;<-:create" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

func (j *Job) IsAborted() bool {
	return j.Status == JobStatusAborted
}

func (j *Job) IsCanceled() bool {
	return j.Status == JobStatusCanceled
}

func (j *Job) IsFinished() bool {
	return j.Status == JobStatusFinished
}

func (j *Job) IsPending() bool {
	return j.Status == JobStatusPending
}

func (j *Job) IsRunning() bool {
	return j.Status == JobStatusRunning
}

func (j *Job) IsDone() bool {
	return asg.Contains(JobDoneStatus, j.Status)
}

func (j *Job) IsUndone() bool {
	return asg.Contains(JobUndoneStatus, j.Status)
}

func (j *Job) String() string {
	return toString(j)
}
