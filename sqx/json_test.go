package sqx

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// ### init postgresql database
// ```sql
// CREATE USER pango PASSWORD 'pango';
// CREATE DATABASE pango WITH OWNER=pango ENCODING='UTF-8';
// GRANT ALL ON DATABASE pango TO pango;
// ```
func testOpen() (*sql.DB, error) {
	pgdsn := os.Getenv("SQLX_POSTGRES_DSN")

	if pgdsn == "" {
		return nil, errors.New("SQLX_POSTGRES_DSN is not set")
	}

	return sql.Open("postgres", pgdsn)
}

const (
	sqlCreate = `
CREATE TABLE persons (
	name text,
	props jsonb,
	anys jsonb,
	strs jsonb,
	ints jsonb,
	i64s jsonb
)
`
	sqlDrop = `DROP TABLE persons`
)

type Person struct {
	Name  string
	Props JSONObject
	Anys  JSONArray
	Strs  JSONStringArray
	Ints  JSONIntArray
	I64s  JSONInt64Array
}

func TestJSON(t *testing.T) {
	pgdb, err := testOpen()
	if pgdb == nil {
		t.Skip(err)
	}
	defer pgdb.Close()

	if _, err := pgdb.Exec(sqlCreate); err != nil {
		t.Fatal(err)
	}

	defer pgdb.Exec(sqlDrop)

	p1 := &Person{"1", nil, nil, nil, nil, nil}
	if _, err := pgdb.Exec(`insert into persons values ($1, $2, $3, $4, $5, $6)`, p1.Name, p1.Props, p1.Anys, p1.Strs, p1.Ints, p1.I64s); err != nil {
		t.Fatal(err)
	}

	p2 := &Person{"x", JSONObject{"p": 1}, JSONArray{4}, JSONStringArray{"5"}, JSONIntArray{6}, JSONInt64Array{7}}
	row := pgdb.QueryRow(`select name, props, anys, strs, ints, i64s from persons limit 1`)
	if row.Err() != nil {
		t.Fatal(row.Err())
	}
	if err := row.Scan(&p2.Name, &p2.Props, &p2.Anys, &p2.Strs, &p2.Ints, &p2.I64s); err != nil {
		t.Fatal(err)
	}

	b1, _ := json.Marshal(p1)
	b2, _ := json.Marshal(p2)
	s1, s2 := string(b1), string(b2)
	if s1 != s2 {
		t.Errorf("\n GOT: %s\nWANT: %s", s2, s1)
	}

	p1.Props = JSONObject{"p": "a"}
	p1.Anys = JSONArray{1, "1"}
	p1.Strs = JSONStringArray{"a", "b"}
	p1.Ints = JSONIntArray{11, 22}
	p1.I64s = JSONInt64Array{88, 99}
	if _, err := pgdb.Exec(`update persons set props = $1, anys = $2, strs = $3, ints = $4, i64s = $5`, p1.Props, p1.Anys, p1.Strs, p1.Ints, p1.I64s); err != nil {
		t.Fatal(err)
	}

	p2 = &Person{}
	row = pgdb.QueryRow(`select name, props, anys, strs, ints, i64s from persons limit 1`)
	if row.Err() != nil {
		t.Fatal(row.Err())
	}
	if err := row.Scan(&p2.Name, &p2.Props, &p2.Anys, &p2.Strs, &p2.Ints, &p2.I64s); err != nil {
		t.Fatal(err)
	}

	b1, _ = json.Marshal(p1)
	b2, _ = json.Marshal(p2)
	s1, s2 = string(b1), string(b2)
	if s1 != s2 {
		t.Errorf("\n GOT: %q\nWANT: %q", s2, s1)
	}
}
