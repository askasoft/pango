package goqulog

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/lib/pq"
)

type Schema struct {
	SchemaName string `db:"schema_name"`
}

// CREATE USER pango PASSWORD 'pango';
// CREATE DATABASE pango WITH OWNER=pango ENCODING='UTF-8';
// GRANT ALL ON DATABASE pango TO pango;

func TestGoquLogger(t *testing.T) {
	log := log.NewLog()

	logger := log.GetLogger("SQL")

	dsn := "host=127.0.0.1 user=pango password=pango dbname=pango port=5432 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
		t.Skip(err)
	}

	db.SetConnMaxLifetime(time.Minute)

	gd := goqu.Dialect("postgres").DB(db)
	gd.Logger(NewGoquLogger(logger))

	schemas := []Schema{}
	err = gd.From(goqu.S("information_schema").Table("schemata")).Select("schema_name").ScanStructs(&schemas)
	if err != nil {
		fmt.Println(err)
		t.Skip(err)
	}

	for _, s := range schemas {
		fmt.Println(s.SchemaName)
	}
}
