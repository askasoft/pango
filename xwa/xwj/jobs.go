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
	"github.com/askasoft/pango/xwa/xwm"
	"gorm.io/gorm"
)

var (
	runnings int32
	mutex    sync.Mutex
)

func RependingJobs(db *gorm.DB, delay time.Duration) {
	ut := time.Now().Add(-delay)
	tx := db.Where("status = ? AND updated_at < ?", xwm.JobStatusRunning, ut)

	job := &xwm.Job{RID: 0, Status: xwm.JobStatusPending, Error: ""}
	r := tx.Select("rid", "status", "error").Updates(job)
	if r.Error != nil {
		log.Errorf("Failed to RependingJobs(): %v", r.Error)
		return
	}
	if r.RowsAffected > 0 {
		log.Infof("RependingJobs: %d", r.RowsAffected)
	}
}

func PendingJob(db *gorm.DB, name string, param string) (int64, error) {
	job := &xwm.Job{Name: name, Param: param, Status: xwm.JobStatusPending}
	r := db.Create(job)
	if r.Error != nil {
		log.Errorf("Failed to pending job %s: %v", job.String(), r.Error)
	}
	return job.ID, r.Error
}

func GetJob(db *gorm.DB, jid int64, details ...bool) (*xwm.Job, error) {
	job := &xwm.Job{}

	tx := db.Model(job)
	if len(details) <= 0 || !details[0] {
		tx = tx.Select("id", "rid", "name", "status")
	}

	r := tx.Take(job, jid)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if r.Error != nil {
		log.Errorf("Failed to get job #%d: %v", jid, r.Error)
		return nil, r.Error
	}
	return job, nil
}

// GetJobLogs get job logs
// set levels to ("I", "W", "E", "F") to filter DEBUG/TRACE logs
func GetJobLogs(db *gorm.DB, jid int64, start, limit int, levels ...string) ([]*xwm.JobLog, error) {
	var jls []*xwm.JobLog

	tx := db.Where("jid = ?", jid)
	if len(levels) > 0 {
		tx.Where("level IN ?", levels)
	}

	r := tx.Order("id asc").Offset(start).Limit(limit).Find(&jls)
	if r.Error != nil {
		log.Errorf("Failed to get job logs #%d: %v", jid, r.Error)
	}
	return jls, r.Error
}

func FindJob(db *gorm.DB, name string, details ...bool) (*xwm.Job, error) {
	job := &xwm.Job{}

	tx := db.Model(job).Where("name = ?", name).Order("id desc")
	if len(details) <= 0 || !details[0] {
		tx = tx.Select("id", "rid", "name", "status")
	}

	r := tx.First(job)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if r.Error != nil {
		log.Errorf("Failed to find job %q: %v", name, r.Error)
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
	job := &xwm.Job{Status: xwm.JobStatusAborted, Error: reason}
	tx := db.Where("id = ? AND status IN ?", jid, []string{xwm.JobStatusPending, xwm.JobStatusRunning})
	r := tx.Select("status", "error").Updates(job)
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected != 1 {
		return ErrJobMissing
	}
	return nil
}

func StartJobs(db *gorm.DB, limit int, run func(*xwm.Job)) {
	mutex.Lock()
	defer mutex.Unlock()

	current := atomic.LoadInt32(&runnings)

	log.Debugf("Current running jobs: %d / %d", current, limit)

	remain := limit - int(current)
	if remain <= 0 {
		return
	}

	var jobs []*xwm.Job

	r := db.Where("status = ?", xwm.JobStatusPending).Order("id asc").Limit(remain).Find(&jobs)
	if r.Error != nil {
		log.Errorf("Failed to find pending jobs: %v", r.Error)
		return
	}

	for _, job := range jobs {
		go SafeRunJob(job, run)
	}
}

func SafeRunJob(job *xwm.Job, run func(*xwm.Job)) {
	n := atomic.AddInt32(&runnings, 1)

	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Panic: %v", err)
		}
		atomic.AddInt32(&runnings, -1)
	}()

	log.Debugf("Start job #%d: #%d %s", n, job.ID, job.Name)

	run(job)
}

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
