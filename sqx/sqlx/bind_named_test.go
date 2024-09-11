package sqlx

import (
	"testing"
)

func TestCompileQuery(t *testing.T) {
	table := []struct {
		Q, R, D, T, N string
		V             []string
	}{
		// basic test for named parameters, invalid char ',' terminating
		{
			Q: `INSERT INTO foo (a,b,c,d) VALUES (:name, :age, :first, :last)`,
			R: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?)`,
			D: `INSERT INTO foo (a,b,c,d) VALUES ($1, $2, $3, $4)`,
			T: `INSERT INTO foo (a,b,c,d) VALUES (@p1, @p2, @p3, @p4)`,
			N: `INSERT INTO foo (a,b,c,d) VALUES (:name, :age, :first, :last)`,
			V: []string{"name", "age", "first", "last"},
		},
		// This query tests a named parameter ending the string as well as numbers
		{
			Q: `SELECT * FROM a WHERE first_name=:name1 AND last_name=:name2`,
			R: `SELECT * FROM a WHERE first_name=? AND last_name=?`,
			D: `SELECT * FROM a WHERE first_name=$1 AND last_name=$2`,
			T: `SELECT * FROM a WHERE first_name=@p1 AND last_name=@p2`,
			N: `SELECT * FROM a WHERE first_name=:name1 AND last_name=:name2`,
			V: []string{"name1", "name2"},
		},
		{
			Q: `SELECT "::foo" FROM a WHERE first_name=:name1 AND last_name=:name2`,
			R: `SELECT ":foo" FROM a WHERE first_name=? AND last_name=?`,
			D: `SELECT ":foo" FROM a WHERE first_name=$1 AND last_name=$2`,
			T: `SELECT ":foo" FROM a WHERE first_name=@p1 AND last_name=@p2`,
			N: `SELECT ":foo" FROM a WHERE first_name=:name1 AND last_name=:name2`,
			V: []string{"name1", "name2"},
		},
		{
			Q: `SELECT 'a::b::c' || first_name, '::::ABC::_::' FROM person WHERE first_name=:first_name AND last_name=:last_name`,
			R: `SELECT 'a:b:c' || first_name, '::ABC:_:' FROM person WHERE first_name=? AND last_name=?`,
			D: `SELECT 'a:b:c' || first_name, '::ABC:_:' FROM person WHERE first_name=$1 AND last_name=$2`,
			T: `SELECT 'a:b:c' || first_name, '::ABC:_:' FROM person WHERE first_name=@p1 AND last_name=@p2`,
			N: `SELECT 'a:b:c' || first_name, '::ABC:_:' FROM person WHERE first_name=:first_name AND last_name=:last_name`,
			V: []string{"first_name", "last_name"},
		},
		{
			Q: `SELECT @name := "name", :age, :first, :last`,
			R: `SELECT @name := "name", ?, ?, ?`,
			D: `SELECT @name := "name", $1, $2, $3`,
			N: `SELECT @name := "name", :age, :first, :last`,
			T: `SELECT @name := "name", @p1, @p2, @p3`,
			V: []string{"age", "first", "last"},
		},
		{
			Q: `INSERT INTO foo (a,b,c,d) VALUES (:あ, :b, :キコ, :名前)`,
			R: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?)`,
			D: `INSERT INTO foo (a,b,c,d) VALUES ($1, $2, $3, $4)`,
			N: `INSERT INTO foo (a,b,c,d) VALUES (:あ, :b, :キコ, :名前)`,
			T: `INSERT INTO foo (a,b,c,d) VALUES (@p1, @p2, @p3, @p4)`,
			V: []string{"あ", "b", "キコ", "名前"},
		},
	}

	for _, test := range table {
		qr, names, err := compileNamedQuery(BindQuestion, test.Q)
		if err != nil {
			t.Error(err)
		}
		if qr != test.R {
			t.Errorf("expected %s, got %s", test.R, qr)
		}
		if len(names) != len(test.V) {
			t.Errorf("expected %#v, got %#v", test.V, names)
		} else {
			for i, name := range names {
				if name != test.V[i] {
					t.Errorf("expected %dth name to be %s, got %s", i+1, test.V[i], name)
				}
			}
		}
		qd, _, _ := compileNamedQuery(BindDollar, test.Q)
		if qd != test.D {
			t.Errorf("\nexpected: `%s`\ngot:      `%s`", test.D, qd)
		}

		qt, _, _ := compileNamedQuery(BindAt, test.Q)
		if qt != test.T {
			t.Errorf("\nexpected: `%s`\ngot:      `%s`", test.T, qt)
		}

		qq, _, _ := compileNamedQuery(BindColon, test.Q)
		if qq != test.N {
			t.Errorf("\nexpected: `%s`\ngot:      `%s`\n(len: %d vs %d)", test.N, qq, len(test.N), len(qq))
		}
	}
}

func TestEscapedColons(t *testing.T) {
	t.Skip("not sure it is possible to support this in general case without an SQL parser")
	var qs = `SELECT * FROM testtable WHERE timeposted BETWEEN (now() AT TIME ZONE 'utc') AND
	(now() AT TIME ZONE 'utc') - interval '01:30:00') AND name = '\'this is a test\'' and id = :id`
	_, _, err := compileNamedQuery(BindDollar, qs)
	if err != nil {
		t.Error("Didn't handle colons correctly when inside a string")
	}
}

func TestFixBounds(t *testing.T) {
	table := []struct {
		name, query, expect string
		loop                int
	}{
		{
			name:   `named syntax`,
			query:  `INSERT INTO foo (a,b,c,d) VALUES (:name, :age, :first, :last)`,
			expect: `INSERT INTO foo (a,b,c,d) VALUES (:name, :age, :first, :last),(:name, :age, :first, :last)`,
			loop:   2,
		},
		{
			name:   `mysql syntax`,
			query:  `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?)`,
			expect: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?),(?, ?, ?, ?)`,
			loop:   2,
		},
		{
			name:   `named syntax w/ trailer`,
			query:  `INSERT INTO foo (a,b,c,d) VALUES (:name, :age, :first, :last) ;--`,
			expect: `INSERT INTO foo (a,b,c,d) VALUES (:name, :age, :first, :last),(:name, :age, :first, :last) ;--`,
			loop:   2,
		},
		{
			name:   `mysql syntax w/ trailer`,
			query:  `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?) ;--`,
			expect: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?),(?, ?, ?, ?) ;--`,
			loop:   2,
		},
		{
			name:   `not found test`,
			query:  `INSERT INTO foo (a,b,c,d) (:name, :age, :first, :last)`,
			expect: `INSERT INTO foo (a,b,c,d) (:name, :age, :first, :last)`,
			loop:   2,
		},
		{
			name:   `found twice test`,
			query:  `INSERT INTO foo (a,b,c,d) VALUES (:name, :age, :first, :last) VALUES (:name, :age, :first, :last)`,
			expect: `INSERT INTO foo (a,b,c,d) VALUES (:name, :age, :first, :last),(:name, :age, :first, :last) VALUES (:name, :age, :first, :last)`,
			loop:   2,
		},
		{
			name:   `nospace`,
			query:  `INSERT INTO foo (a,b) VALUES(:a, :b)`,
			expect: `INSERT INTO foo (a,b) VALUES(:a, :b),(:a, :b)`,
			loop:   2,
		},
		{
			name:   `lowercase`,
			query:  `INSERT INTO foo (a,b) values(:a, :b)`,
			expect: `INSERT INTO foo (a,b) values(:a, :b),(:a, :b)`,
			loop:   2,
		},
		{
			name:   `on duplicate key using VALUES`,
			query:  `INSERT INTO foo (a,b) VALUES (:a, :b) ON DUPLICATE KEY UPDATE a=VALUES(a)`,
			expect: `INSERT INTO foo (a,b) VALUES (:a, :b),(:a, :b) ON DUPLICATE KEY UPDATE a=VALUES(a)`,
			loop:   2,
		},
		{
			name:   `single column`,
			query:  `INSERT INTO foo (a) VALUES (:a)`,
			expect: `INSERT INTO foo (a) VALUES (:a),(:a)`,
			loop:   2,
		},
		{
			name:   `call now`,
			query:  `INSERT INTO foo (a, b) VALUES (:a, NOW())`,
			expect: `INSERT INTO foo (a, b) VALUES (:a, NOW()),(:a, NOW())`,
			loop:   2,
		},
		{
			name:   `two level depth function call`,
			query:  `INSERT INTO foo (a, b) VALUES (:a, YEAR(NOW()))`,
			expect: `INSERT INTO foo (a, b) VALUES (:a, YEAR(NOW())),(:a, YEAR(NOW()))`,
			loop:   2,
		},
		{
			name:   `missing closing bracket`,
			query:  `INSERT INTO foo (a, b) VALUES (:a, YEAR(NOW())`,
			expect: `INSERT INTO foo (a, b) VALUES (:a, YEAR(NOW())`,
			loop:   2,
		},
		{
			name:   `table with "values" at the end`,
			query:  `INSERT INTO table_values (a, b) VALUES (:a, :b)`,
			expect: `INSERT INTO table_values (a, b) VALUES (:a, :b),(:a, :b)`,
			loop:   2,
		},
		{
			name: `multiline indented query`,
			query: `INSERT INTO foo (
		a,
		b,
		c,
		d
	) VALUES (
		:name,
		:age,
		:first,
		:last
	)`,
			expect: `INSERT INTO foo (
		a,
		b,
		c,
		d
	) VALUES (
		:name,
		:age,
		:first,
		:last
	),(
		:name,
		:age,
		:first,
		:last
	)`,
			loop: 2,
		},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			res := fixBound(tc.query, tc.loop)
			if res != tc.expect {
				t.Errorf("mismatched results")
			}
		})
	}
}
