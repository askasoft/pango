package xwj

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/askasoft/pango/bye"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/str"
	"gorm.io/gorm"
)

var (
	runnings int32
	mutex    sync.Mutex
)

func LocalRunningJobs() int {
	return int(atomic.LoadInt32(&runnings))
}

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

func RependingJobs(db *gorm.DB, delay time.Duration) (int, error) {
	ut := time.Now().Add(-delay)
	tx := db.Where("status = ? AND updated_at < ?", JobStatusRunning, ut)

	job := &Job{RID: 0, Status: JobStatusPending, Error: ""}
	r := tx.Select("rid", "status", "error").Updates(job)
	if r.Error != nil {
		return 0, r.Error
	}
	return int(r.RowsAffected), nil
}

func PendingJob(db *gorm.DB, name string, param string) (int64, error) {
	job := &Job{Name: name, Param: param, Status: JobStatusPending}
	r := db.Create(job)
	return job.ID, r.Error
}

func GetJob(db *gorm.DB, jid int64, details ...bool) (*Job, error) {
	job := &Job{}

	tx := db.Model(job)
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
func GetJobLogs(db *gorm.DB, jid int64, start, limit int, levels ...string) ([]*JobLog, error) {
	var jls []*JobLog

	tx := db.Where("jid = ?", jid)
	if len(levels) > 0 {
		tx.Where("level IN ?", levels)
	}

	r := tx.Order("id asc").Offset(start).Limit(limit).Find(&jls)
	return jls, r.Error
}

func FindJob(db *gorm.DB, name string, details ...bool) (*Job, error) {
	job := &Job{}

	tx := db.Model(job).Where("name = ?", name).Order("id desc")
	if len(details) <= 0 || !details[0] {
		tx = tx.Select("id", "rid", "name", "status", "created_at", "updated_at")
	}

	r := tx.First(job)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return job, r.Error
}

func FindAndAbortJob(db *gorm.DB, name, reason string) error {
	job, err := FindJob(db, name)
	if err != nil {
		return err
	}

	return AbortJob(db, job.ID, reason)
}

func AbortJob(db *gorm.DB, jid int64, reason string) error {
	job := &Job{Status: JobStatusAborted, Error: reason}
	tx := db.Where("id = ? AND status IN ?", jid, []string{JobStatusPending, JobStatusRunning})
	r := tx.Select("status", "error").Updates(job)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected != 1 {
		return ErrJobMissing
	}
	return nil
}

func StartJobs(db *gorm.DB, limit int, run func(*Job), loggers ...log.Logger) error {
	mutex.Lock()
	defer mutex.Unlock()

	current := atomic.LoadInt32(&runnings)

	logger := getLogger(loggers...)
	logger.Debugf("Current running jobs: %d / %d", current, limit)

	remain := limit - int(current)
	if remain <= 0 {
		return nil
	}

	var jobs []*Job

	r := db.Where("status = ?", JobStatusPending).Order("id asc").Limit(remain).Find(&jobs)
	if r.Error != nil {
		return r.Error
	}

	for _, job := range jobs {
		go SafeRunJob(job, run, logger)
	}

	return nil
}

func SafeRunJob(job *Job, run func(*Job), logger log.Logger) {
	n := atomic.AddInt32(&runnings, 1)

	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Panic: %v", err)
		}
		atomic.AddInt32(&runnings, -1)
	}()

	logger.Debugf("Start job #%d: #%d %s", n, job.ID, job.Name)

	run(job)
}

func CleanOutdatedJobs(db *gorm.DB, due time.Time, loggers ...log.Logger) error {
	logger := getLogger(loggers...)

	ss := []string{JobStatusAborted, JobStatusCompleted}
	r := db.Where("jid IN (SELECT id FROM jobs WHERE status IN ? AND updated_at < ?)", ss, due).Delete(&JobLog{})
	if r.Error != nil {
		logger.Errorf("Failed to delete outdated job logs: %v", r.Error)
		return r.Error
	}
	if r.RowsAffected > 0 {
		logger.Infof("Delete outdated job logs: %d", r.RowsAffected)
	}

	r = db.Where("status IN ? AND updated_at < ?", ss, due).Delete(&Job{})
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
