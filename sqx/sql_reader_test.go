package sqx

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestReadSql(t *testing.T) {
	raws := `
/*
 comment; comment
 */
-- comment; comment

SELECT 'a''b' FROM 1;

SELECT 'a''b' 
FROM 2;
	`

	exps := []string{
		`SELECT 'a''b' FROM 1`,
		`SELECT 'a''b' 
FROM 2`,
	}
	acts := []string{}
	sr := NewSqlReader(strings.NewReader(raws))
	for {
		sql, err := sr.ReadSql()
		if errors.Is(err, io.EOF) {
			break
		}
		acts = append(acts, sql)
	}

	if !reflect.DeepEqual(acts, exps) {
		t.Errorf("\nReadSql(): %q\n     WANT: %q", acts, exps)
	}
}
