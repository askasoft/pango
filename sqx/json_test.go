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

func testEqual(t *testing.T, w, a any) {
	bw, _ := json.Marshal(w)
	ba, _ := json.Marshal(a)

	sw, sa := string(bw), string(ba)

	// fmt.Println(sa)
	if sw != sa {
		t.Errorf("\n GOT: %s\nWANT: %s", sa, sw)
	}
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

	p1i := &Person{"z", nil, nil, nil, nil, nil}
	if _, err := pgdb.Exec(`insert into persons values ($1, $2, $3, $4, $5, $6)`, p1i.Name, p1i.Props, p1i.Anys, p1i.Strs, p1i.Ints, p1i.I64s); err != nil {
		t.Fatal(err)
	}

	p1s := &Person{"", JSONObject{}, JSONArray{}, nil, nil, nil}
	row := pgdb.QueryRow(`select name, props, anys, strs, ints, i64s from persons where name = 'z' and props is null and anys is null`)
	if err := row.Scan(&p1s.Name, &p1s.Props, &p1s.Anys, &p1s.Strs, &p1s.Ints, &p1s.I64s); err != nil {
		t.Fatal(err)
	}
	testEqual(t, p1i, p1s)

	p2i := &Person{"x", JSONObject{"p": 1}, JSONArray{4}, JSONStringArray{"5"}, JSONIntArray{6}, JSONInt64Array{7}}
	if _, err := pgdb.Exec(`insert into persons values ($1, $2, $3, $4, $5, $6)`, p2i.Name, p2i.Props, p2i.Anys, p2i.Strs, p2i.Ints, p2i.I64s); err != nil {
		t.Fatal(err)
	}

	p2s := &Person{}
	row = pgdb.QueryRow(`select name, props, anys, strs, ints, i64s from persons where name = 'x'`)
	if row.Err() != nil {
		t.Fatal(row.Err())
	}
	if err := row.Scan(&p2s.Name, &p2s.Props, &p2s.Anys, &p2s.Strs, &p2s.Ints, &p2s.I64s); err != nil {
		t.Fatal(err)
	}
	testEqual(t, p2i, p2s)

	p3u := &Person{"x", JSONObject{"p": "a"}, JSONArray{1, "1"}, JSONStringArray{"a", "b"}, JSONIntArray{11, 22}, JSONInt64Array{88, 99}}
	if _, err := pgdb.Exec(`update persons set props = $1, anys = $2, strs = $3, ints = $4, i64s = $5 where name = 'x'`, p3u.Props, p3u.Anys, p3u.Strs, p3u.Ints, p3u.I64s); err != nil {
		t.Fatal(err)
	}

	p3s := &Person{}
	row = pgdb.QueryRow(`select name, props, anys, strs, ints, i64s from persons where name = 'x'`)
	if row.Err() != nil {
		t.Fatal(row.Err())
	}
	if err := row.Scan(&p3s.Name, &p3s.Props, &p3s.Anys, &p3s.Strs, &p3s.Ints, &p3s.I64s); err != nil {
		t.Fatal(err)
	}
	testEqual(t, p3u, p3s)
}
