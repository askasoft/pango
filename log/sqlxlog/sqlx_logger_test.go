package sqlxlog

import (
	"fmt"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/sqx"
	"github.com/askasoft/pango/sqx/sqlx"
	_ "github.com/lib/pq"
)

type Schema struct {
	SchemaName string `db:"schema_name"`
}

// CREATE USER pango PASSWORD 'pango';
// CREATE DATABASE pango WITH OWNER=pango ENCODING='UTF-8';
// GRANT ALL ON DATABASE pango TO pango;

func TestSqlxLogger(t *testing.T) {
	log := log.NewLog()

	slg := &SqlxLogger{
		Logger:        log.GetLogger("SQL"),
		SlowThreshold: time.Second,
	}

	dsn := "host=127.0.0.1 user=pango password=pango dbname=pango port=5432 sslmode=disable"
	sdb, err := sqlx.Connect("postgres", dsn, slg.Trace)
	if err != nil {
		fmt.Println(err)
		t.Skip(err)
	}

	sdb.DB().SetConnMaxLifetime(time.Minute)

	sqb := &sqx.Builder{}
	sqb.Select("schema_name").From("information_schema.schemata").Where("schema_name <> ?", "test")
	sql, args := sqb.Build()
	sql = sdb.Rebind(sql)

	schemas := []Schema{}
	err = sdb.Select(&schemas, sql, args...)
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range schemas {
		fmt.Println(s.SchemaName)
	}
}
