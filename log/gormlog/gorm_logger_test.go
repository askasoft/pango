package gormlog

import (
	"fmt"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormLoggerTest struct {
	gorm.Model
}

// CREATE USER pango PASSWORD 'pango';
// CREATE DATABASE pango WITH OWNER=pango ENCODING='UTF-8';
// GRANT ALL ON DATABASE pango TO pango;

func TestGormLogger(t *testing.T) {
	log := log.NewLog()
	logger := log.GetLogger("SQL")

	dsn := "host=127.0.0.1 user=pango password=pango dbname=pango port=5432 sslmode=disable"
	orm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &GormLogger{
			Logger:        logger,
			SlowThreshold: time.Second,
		},
		SkipDefaultTransaction: true,
	})
	if err != nil {
		fmt.Println(err)
		t.Skip(err)
	}

	db, err := orm.DB()
	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
	}
	db.SetConnMaxLifetime(time.Minute)

	// migration
	err = orm.AutoMigrate(&GormLoggerTest{})
	if err != nil {
		fmt.Println(err)
	}
}
