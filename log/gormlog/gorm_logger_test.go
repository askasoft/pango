package gormlog

import (
	"fmt"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
	"gorm.io/driver/postgres" // postgres dialect
	"gorm.io/gorm"
)

type Pango struct {
	gorm.Model
}

func TestGormLogger(t *testing.T) {
	log := log.NewLog()
	logger := log.GetLogger("SQL")

	dsn := "host=127.0.0.1 user=panda password=panda dbname=ptest port=5432 sslmode=disable"
	orm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &GormLogger{
			Logger:        logger,
			SlowThreshold: time.Second,
		},
		SkipDefaultTransaction: true,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := orm.DB()
	if err != nil {
		fmt.Println(err)
		return
	}
	db.SetConnMaxLifetime(time.Minute)

	// migration
	err = orm.AutoMigrate(&Pango{})
	if err != nil {
		fmt.Println(err)
		return
	}
}
