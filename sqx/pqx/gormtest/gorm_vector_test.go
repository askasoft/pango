package gormtest

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/askasoft/pango/sqx/pqx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CREATE USER pango PASSWORD 'pango';
// CREATE DATABASE pango WITH OWNER=pango ENCODING='UTF-8';
// GRANT ALL ON DATABASE pango TO pango;

// statically assert that Vector implements sql.Scanner.
var _ sql.Scanner = (*pqx.Vector)(nil)

// statically assert that Vector implements driver.Valuer.
var _ driver.Valuer = (*pqx.Vector)(nil)

type GormVectorTest struct {
	gorm.Model
	Embedding pqx.Vector `gorm:"type:vector(3)"`
}

func createTestItems(t *testing.T, db *gorm.DB) {
	items := []GormVectorTest{
		{Embedding: pqx.Vector([]float64{1, 1, 1})},
		{Embedding: pqx.Vector([]float64{2, 2, 2})},
		{Embedding: pqx.Vector([]float64{1, 1, 2})},
	}

	result := db.Create(items)

	if result.Error != nil {
		t.Fatal(result.Error)
	}
}

func TestGormVector(t *testing.T) {
	dsn := "host=127.0.0.1 user=pango password=pango dbname=pango port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip(err)
		return
	}

	//	db.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	db.Exec("DROP TABLE IF EXISTS gorm_vector_tests")

	db.AutoMigrate(&GormVectorTest{})

	db.Exec("CREATE INDEX ON gorm_items USING hnsw (embedding vector_l2_ops)")

	createTestItems(t, db)

	var items []GormVectorTest
	db.Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "embedding <-> ?", Vars: []any{pqx.Vector([]float64{1, 1, 1})}},
	}).Limit(5).Find(&items)
	if items[0].ID != 1 || items[1].ID != 3 || items[2].ID != 2 {
		t.Errorf("Bad ids")
	}
	if !reflect.DeepEqual(items[1].Embedding.Slice(), []float64{1, 1, 2}) {
		t.Errorf("Bad embedding")
	}

	var distances []float64
	db.Model(&GormVectorTest{}).Select("embedding <-> ?", pqx.Vector([]float64{1, 1, 1})).Order("id").Find(&distances)
	fmt.Println(distances)
	if distances[0] != 0 || distances[1] != math.Sqrt(3) || distances[2] != 1 {
		t.Errorf("Bad distances")
	}

	var similars []float64
	db.Model(&GormVectorTest{}).Select("1 - (embedding <=> ?)", pqx.Vector([]float64{1, 1, 1})).Order("id").Find(&similars)
	fmt.Println(similars)
	if similars[0] != 1 || similars[1] != 1 || similars[2] < 0.9 {
		t.Errorf("Bad similarity")
	}
}
