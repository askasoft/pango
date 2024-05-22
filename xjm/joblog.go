package xjm

import (
	"time"

	"github.com/askasoft/pango/log"
)

var (
	JobLogLevelFatal = log.LevelFatal.Prefix()
	JobLogLevelError = log.LevelError.Prefix()
	JobLogLevelWarn  = log.LevelWarn.Prefix()
	JobLogLevelInfo  = log.LevelInfo.Prefix()
	JobLogLevelDebug = log.LevelDebug.Prefix()
	JobLogLevelTrace = log.LevelTrace.Prefix()
)

type JobLog struct {
	ID      int64     `gorm:"not null;primaryKey;autoIncrement" json:"id,omitempty"`
	JID     int64     `gorm:"column:jid;not null;index:idx_job_logs_jid" json:"jid,omitempty"`
	Time    time.Time `gorm:"not null" json:"time,omitempty"`
	Level   string    `gorm:"size:1;not null" json:"level,omitempty"`
	Message string    `gorm:"not null" json:"message,omitempty"`
}

func (jl *JobLog) String() string {
	return toString(jl)
}
