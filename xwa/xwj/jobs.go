package xwj

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/str"
	"gorm.io/gorm"
)

func Encode(v any) string {
	if v == nil {
		return ""
	}

	if s, ok := v.(string); ok {
		return s
	}

	if bs, ok := v.([]byte); ok {
		return base64.StdEncoding.EncodeToString(bs)
	}

	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bye.UnsafeString(bs)
}

func Decode(p string, v any) error {
	if ps, ok := v.(*string); ok {
		*ps = p
		return nil
	}

	if pb, ok := v.(*[]byte); ok {
		bs, err := base64.StdEncoding.DecodeString(p)
		if err != nil {
			return err
		}

		*pb = bs
		return nil
	}

	return json.Unmarshal(str.UnsafeBytes(p), v)
}

func GetJob(db *gorm.DB, table string, jid int64, details ...bool) (*Job, error) {
	job := &Job{}

	tx := db.Table(table)
	if len(details) <= 0 || !details[0] {
		tx = tx.Select("id", "rid", "name", "status", "created_at", "updated_at")
	}

	r := tx.Take(job, jid)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if r.Error != nil {
		return nil, r.Error
	}
	return job, nil
}

// GetJobLogs get job logs
// set levels to ("I", "W", "E", "F") to filter DEBUG/TRACE logs
func GetJobLogs(db *gorm.DB, table string, jid int64, start, limit int, levels ...string) ([]*JobLog, error) {
	var jls []*JobLog

	tx := db.Table(table).Where("jid = ?", jid)
	if len(levels) > 0 {
		tx.Where("level IN ?", levels)
	}

	r := tx.Order("id asc").Offset(start).Limit(limit).Find(&jls)
	return jls, r.Error
}

func FindJob(db *gorm.DB, table string, name string, details ...bool) (*Job, error) {
	job := &Job{}

	tx := db.Table(table).Where("name = ?", name).Order("id desc")
	if len(details) <= 0 || !details[0] {
		tx = tx.Select("id", "rid", "name", "status", "created_at", "updated_at")
	}

	r := tx.First(job)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return job, r.Error
}

func AppendJob(db *gorm.DB, table string, name string, param string) (int64, error) {
	job := &Job{Name: name, Param: param, Status: JobStatusPending}
	r := db.Table(table).Create(job)
	return job.ID, r.Error
}

func FindAndAbortJob(db *gorm.DB, table string, name, reason string, loggers ...log.Logger) error {
	logger := getLogger(loggers...)

	job, err := FindJob(db, table, name)
	if err != nil {
		logger.Errorf("Failed to find job '%s': %v", name, err)
		return err
	}

	return AbortJob(db, table, job.ID, reason, logger)
}

func AbortJob(db *gorm.DB, table string, jid int64, reason string, loggers ...log.Logger) error {
	logger := getLogger(loggers...)

	logger.Infof("Abort job #%d: %s", jid, reason)

	job := &Job{Status: JobStatusAborted, Error: reason}
	jss := []string{JobStatusPending, JobStatusRunning}

	tx := db.Table(table).Where("id = ? AND status IN ?", jid, jss)
	r := tx.Select("status", "error").Updates(job)
	if r.Error != nil {
		logger.Errorf("Failed to abort job #%d: %v", jid, r.Error)
		return r.Error
	}
	if r.RowsAffected != 1 {
		logger.Warnf("Unable to abort job #%d: %d, %v", jid, r.RowsAffected, ErrJobMissing)
		return ErrJobMissing
	}
	return nil
}

func ReappendJobs(db *gorm.DB, table string, before time.Time, loggers ...log.Logger) error {
	logger := getLogger(loggers...)

	job := &Job{RID: 0, Status: JobStatusPending, Error: ""}

	tx := db.Table(table).Where("status = ? AND updated_at < ?", JobStatusRunning, before)
	r := tx.Select("rid", "status", "error").Updates(job)
	if r.Error != nil {
		logger.Errorf("Failed to ReappendJobs(): %v", r.Error)
		return r.Error
	}
	if r.RowsAffected > 0 {
		logger.Infof("Job reappended: %d", r.RowsAffected)
	}
	return nil
}

func StartJobs(db *gorm.DB, table string, limit int, run func(*Job), loggers ...log.Logger) error {
	logger := getLogger(loggers...)

	var jobs []*Job

	r := db.Table(table).Where("status = ?", JobStatusPending).Order("id asc").Limit(limit).Find(&jobs)
	if r.Error != nil {
		logger.Errorf("Failed to find pendding job: %v", r.Error)
		return r.Error
	}

	for _, job := range jobs {
		go SafeRunJob(job, run, logger)
	}

	return nil
}

func SafeRunJob(job *Job, run func(*Job), loggers ...log.Logger) {
	logger := getLogger(loggers...)

	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Job #%d '%s' panic: %v", job.ID, job.Name, err)
		}
	}()

	logger.Debugf("Start job #%d '%s'", job.ID, job.Name)

	run(job)
}

func CleanOutdatedJobs(db *gorm.DB, jobTable, logTable string, before time.Time, loggers ...log.Logger) error {
	logger := getLogger(loggers...)

	jss := []string{JobStatusAborted, JobStatusCompleted}
	where := "jid IN (SELECT id FROM " + jobTable + " WHERE status IN ? AND updated_at < ?)"

	r := db.Table(logTable).Where(where, jss, before).Delete(&JobLog{})
	if r.Error != nil {
		logger.Errorf("Failed to delete outdated job logs: %v", r.Error)
		return r.Error
	}
	if r.RowsAffected > 0 {
		logger.Infof("Delete outdated job logs: %d", r.RowsAffected)
	}

	r = db.Table(jobTable).Where("status IN ? AND updated_at < ?", jss, before).Delete(&Job{})
	if r.Error != nil {
		logger.Errorf("Failed to delete outdated jobs: %v", r.Error)
		return r.Error
	}
	if r.RowsAffected > 0 {
		logger.Infof("Delete outdated jobs: %d", r.RowsAffected)
	}

	return nil
}

func getLogger(loggers ...log.Logger) log.Logger {
	if len(loggers) > 0 {
		return loggers[0]
	}

	return log.GetLogger("XWJ")
}
