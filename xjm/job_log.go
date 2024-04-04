package xjm

import (
	"time"
)

type JobLog struct {
	ID      int64     `gorm:"not null;primaryKey;autoIncrement" uri:"id" form:"id" json:"id,omitempty"`
	JID     int64     `gorm:"column:jid;not null;index:idx_job_logs_jid" uri:"jid" form:"jid" json:"jid,omitempty"`
	Time    time.Time `gorm:"not null" json:"time,omitempty"`
	Level   string    `gorm:"size:1;not null" json:"level,omitempty"`
	Message string    `gorm:"not null" json:"message,omitempty"`
}

func (jl *JobLog) String() string {
	return toString(jl)
}
