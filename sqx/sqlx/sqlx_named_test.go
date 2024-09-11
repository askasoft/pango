package sqlx

import (
	"context"
	"fmt"
	"testing"
)

type Test struct {
	t *testing.T
}

func (t Test) Error(err error, msg ...any) {
	t.t.Helper()
	if err != nil {
		if len(msg) == 0 {
			t.t.Error(err)
		} else {
			t.t.Error(msg...)
		}
	}
}

func (t Test) Errorf(err error, format string, args ...any) {
	t.t.Helper()
	if err != nil {
		t.t.Errorf(format, args...)
	}
}

func TestNamedQueries(t *testing.T) {
	RunWithSchema(defaultSchema, t, func(db *DB, t *testing.T, now string) {
		loadDefaultFixture(db, t)
		test := Test{t}
		var ns *NamedStmt
		var err error

		// Check that invalid preparations fail
		ns, err = db.PrepareNamed("SELECT * FROM person WHERE first_name=:first:name")
		if err == nil {
			t.Error("Expected an error with invalid prepared statement.")
		}

		ns, err = db.PrepareNamed("invalid sql")
		if err == nil {
			t.Error("Expected an error with invalid prepared statement.")
		}

		// Check closing works as anticipated
		ns, err = db.PrepareNamed("SELECT * FROM person WHERE first_name=:first_name")
		test.Error(err)
		err = ns.Close()
		test.Error(err)

		ns, err = db.PrepareNamed(`
			SELECT first_name, last_name, email 
			FROM person WHERE first_name=:first_name AND email=:email`)
		test.Error(err)

		// test Queryx w/ uses Query
		p := Person{FirstName: "Jason", LastName: "Moiron", Email: "jmoiron@jmoiron.net"}

		rows, err := ns.Queryx(p)
		test.Error(err)
		for rows.Next() {
			var p2 Person
			rows.StructScan(&p2)
			if p.FirstName != p2.FirstName {
				t.Errorf("got %s, expected %s", p.FirstName, p2.FirstName)
			}
			if p.LastName != p2.LastName {
				t.Errorf("got %s, expected %s", p.LastName, p2.LastName)
			}
			if p.Email != p2.Email {
				t.Errorf("got %s, expected %s", p.Email, p2.Email)
			}
		}

		// test Select
		people := make([]Person, 0, 5)
		err = ns.Select(&people, p)
		test.Error(err)

		if len(people) != 1 {
			t.Errorf("got %d results, expected %d", len(people), 1)
		}
		if p.FirstName != people[0].FirstName {
			t.Errorf("got %s, expected %s", p.FirstName, people[0].FirstName)
		}
		if p.LastName != people[0].LastName {
			t.Errorf("got %s, expected %s", p.LastName, people[0].LastName)
		}
		if p.Email != people[0].Email {
			t.Errorf("got %s, expected %s", p.Email, people[0].Email)
		}

		// test struct batch inserts
		sls := []Person{
			{FirstName: "Ardie", LastName: "Savea", Email: "asavea@ab.co.nz"},
			{FirstName: "Sonny Bill", LastName: "Williams", Email: "sbw@ab.co.nz"},
			{FirstName: "Ngani", LastName: "Laumape", Email: "nlaumape@ab.co.nz"},
		}

		insert := fmt.Sprintf(
			"INSERT INTO person (first_name, last_name, email, added_at) VALUES (:first_name, :last_name, :email, %v)\n",
			now,
		)
		_, err = db.NamedExec(insert, sls)
		test.Error(err)

		// test map batch inserts
		slsMap := []map[string]any{
			{"first_name": "Ardie", "last_name": "Savea", "email": "asavea@ab.co.nz"},
			{"first_name": "Sonny Bill", "last_name": "Williams", "email": "sbw@ab.co.nz"},
			{"first_name": "Ngani", "last_name": "Laumape", "email": "nlaumape@ab.co.nz"},
		}

		_, err = db.NamedExec(`INSERT INTO person (first_name, last_name, email)
			VALUES (:first_name, :last_name, :email) ;--`, slsMap)
		test.Error(err)

		type A map[string]any

		typedMap := []A{
			{"first_name": "Ardie", "last_name": "Savea", "email": "asavea@ab.co.nz"},
			{"first_name": "Sonny Bill", "last_name": "Williams", "email": "sbw@ab.co.nz"},
			{"first_name": "Ngani", "last_name": "Laumape", "email": "nlaumape@ab.co.nz"},
		}

		_, err = db.NamedExec(`INSERT INTO person (first_name, last_name, email)
			VALUES (:first_name, :last_name, :email) ;--`, typedMap)
		test.Error(err)

		for _, p := range sls {
			dest := Person{}
			err = db.Get(&dest, db.Rebind("SELECT * FROM person WHERE email=?"), p.Email)
			test.Error(err)
			if dest.Email != p.Email {
				t.Errorf("expected %s, got %s", p.Email, dest.Email)
			}
		}

		// test Exec
		ns, err = db.PrepareNamed(`
			INSERT INTO person (first_name, last_name, email)
			VALUES (:first_name, :last_name, :email)`)
		test.Error(err)

		js := Person{
			FirstName: "Julien",
			LastName:  "Savea",
			Email:     "jsavea@ab.co.nz",
		}
		_, err = ns.Exec(js)
		test.Error(err)

		// Make sure we can pull him out again
		p2 := Person{}
		db.Get(&p2, db.Rebind("SELECT * FROM person WHERE email=?"), js.Email)
		if p2.Email != js.Email {
			t.Errorf("expected %s, got %s", js.Email, p2.Email)
		}

		// test Txn NamedStmts
		tx := db.MustBeginx()
		txns := tx.NamedStmt(ns)

		// We're going to add Steven in this txn
		sl := Person{
			FirstName: "Steven",
			LastName:  "Luatua",
			Email:     "sluatua@ab.co.nz",
		}

		_, err = txns.Exec(sl)
		test.Error(err)
		// then rollback...
		tx.Rollback()
		// looking for Steven after a rollback should fail
		err = db.Get(&p2, db.Rebind("SELECT * FROM person WHERE email=?"), sl.Email)
		if err != ErrNoRows {
			t.Errorf("expected no rows error, got %v", err)
		}

		// now do the same, but commit
		tx = db.MustBeginx()
		txns = tx.NamedStmt(ns)
		_, err = txns.Exec(sl)
		test.Error(err)
		tx.Commit()

		// looking for Steven after a Commit should succeed
		err = db.Get(&p2, db.Rebind("SELECT * FROM person WHERE email=?"), sl.Email)
		test.Error(err)
		if p2.Email != sl.Email {
			t.Errorf("expected %s, got %s", sl.Email, p2.Email)
		}

	})
}

func TestNamedContextQueries(t *testing.T) {
	RunWithSchema(defaultSchema, t, func(db *DB, t *testing.T, now string) {
		loadDefaultFixture(db, t)
		test := Test{t}
		var ns *NamedStmt
		var err error

		ctx := context.Background()

		// Check that invalid preparations fail
		ns, err = db.PrepareNamedContext(ctx, "SELECT * FROM person WHERE first_name=:first:name")
		if err == nil {
			t.Error("Expected an error with invalid prepared statement.")
		}

		ns, err = db.PrepareNamedContext(ctx, "invalid sql")
		if err == nil {
			t.Error("Expected an error with invalid prepared statement.")
		}

		// Check closing works as anticipated
		ns, err = db.PrepareNamedContext(ctx, "SELECT * FROM person WHERE first_name=:first_name")
		test.Error(err)
		err = ns.Close()
		test.Error(err)

		ns, err = db.PrepareNamedContext(ctx, `
			SELECT first_name, last_name, email
			FROM person WHERE first_name=:first_name AND email=:email`)
		test.Error(err)

		// test Queryx w/ uses Query
		p := Person{FirstName: "Jason", LastName: "Moiron", Email: "jmoiron@jmoiron.net"}

		rows, err := ns.QueryxContext(ctx, p)
		test.Error(err)
		for rows.Next() {
			var p2 Person
			rows.StructScan(&p2)
			if p.FirstName != p2.FirstName {
				t.Errorf("got %s, expected %s", p.FirstName, p2.FirstName)
			}
			if p.LastName != p2.LastName {
				t.Errorf("got %s, expected %s", p.LastName, p2.LastName)
			}
			if p.Email != p2.Email {
				t.Errorf("got %s, expected %s", p.Email, p2.Email)
			}
		}

		// test Select
		people := make([]Person, 0, 5)
		err = ns.SelectContext(ctx, &people, p)
		test.Error(err)

		if len(people) != 1 {
			t.Errorf("got %d results, expected %d", len(people), 1)
		}
		if p.FirstName != people[0].FirstName {
			t.Errorf("got %s, expected %s", p.FirstName, people[0].FirstName)
		}
		if p.LastName != people[0].LastName {
			t.Errorf("got %s, expected %s", p.LastName, people[0].LastName)
		}
		if p.Email != people[0].Email {
			t.Errorf("got %s, expected %s", p.Email, people[0].Email)
		}

		// test Exec
		ns, err = db.PrepareNamedContext(ctx, `
			INSERT INTO person (first_name, last_name, email)
			VALUES (:first_name, :last_name, :email)`)
		test.Error(err)

		js := Person{
			FirstName: "Julien",
			LastName:  "Savea",
			Email:     "jsavea@ab.co.nz",
		}
		_, err = ns.ExecContext(ctx, js)
		test.Error(err)

		// Make sure we can pull him out again
		p2 := Person{}
		db.GetContext(ctx, &p2, db.Rebind("SELECT * FROM person WHERE email=?"), js.Email)
		if p2.Email != js.Email {
			t.Errorf("expected %s, got %s", js.Email, p2.Email)
		}

		// test Txn NamedStmts
		tx := db.MustBeginTx(ctx, nil)
		txns := tx.NamedStmtContext(ctx, ns)

		// We're going to add Steven in this txn
		sl := Person{
			FirstName: "Steven",
			LastName:  "Luatua",
			Email:     "sluatua@ab.co.nz",
		}

		_, err = txns.ExecContext(ctx, sl)
		test.Error(err)
		// then rollback...
		tx.Rollback()
		// looking for Steven after a rollback should fail
		err = db.GetContext(ctx, &p2, db.Rebind("SELECT * FROM person WHERE email=?"), sl.Email)
		if err != ErrNoRows {
			t.Errorf("expected no rows error, got %v", err)
		}

		// now do the same, but commit
		tx = db.MustBeginTx(ctx, nil)
		txns = tx.NamedStmtContext(ctx, ns)
		_, err = txns.ExecContext(ctx, sl)
		test.Error(err)
		tx.Commit()

		// looking for Steven after a Commit should succeed
		err = db.GetContext(ctx, &p2, db.Rebind("SELECT * FROM person WHERE email=?"), sl.Email)
		test.Error(err)
		if p2.Email != sl.Email {
			t.Errorf("expected %s, got %s", sl.Email, p2.Email)
		}

	})
}
