package sqx

import (
	"reflect"
	"strings"
	"testing"
)

// helper to create a builder with dummy binder/quoter
func newBuilder() *Builder {
	return &Builder{
		Binder: BindQuestion,
		Quoter: QuoteDefault,
	}
}

func TestBuilder_SelectSQL(t *testing.T) {
	b := newBuilder().
		Select("id", "name").
		From("users").
		Where(`"id" = ?`, 10).
		Orders("-id,name").
		Limit(10).
		Offset(5)

	sql, args := b.Build()

	wantSQL := `SELECT "id", "name" FROM "users" WHERE "id" = ? ORDER BY "id" DESC, "name" ASC LIMIT 10 OFFSET 5`
	if sql != wantSQL {
		t.Errorf("unexpected SQL:\n got: %s\nwant: %s", sql, wantSQL)
	}
	if !reflect.DeepEqual(args, []any{10}) {
		t.Errorf("unexpected args: got %v, want %v", args, []any{10})
	}
}

func TestBuilder_InsertSQL(t *testing.T) {
	b := newBuilder().
		Insert("users").
		Columns("id", "name").
		Values("?", "?").
		Returns("id")

	sql := b.SQL()
	want := `INSERT INTO "users" ("id", "name") VALUES (?, ?) RETURNING "id"`
	if sql != want {
		t.Errorf("got %s, want %s", sql, want)
	}
}

func TestBuilder_UpdateSQL(t *testing.T) {
	b := newBuilder().
		Update("users").
		Setx("name", "?", "Alice").
		Where(`"id" = ?`, 1).
		Returns("id")

	sql := b.SQL()
	want := `UPDATE "users" SET "name" = ? WHERE "id" = ? RETURNING "id"`
	if sql != want {
		t.Errorf("got %s, want %s", sql, want)
	}

	args := b.Params()
	wantArgs := []any{"Alice", 1}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Errorf("args mismatch: got %v, want %v", args, wantArgs)
	}
}

func TestBuilder_DeleteSQL(t *testing.T) {
	b := newBuilder().
		Delete("users").
		Where(`"id" = ?`, 99)

	sql := b.SQL()
	want := `DELETE FROM "users" WHERE "id" = ?`
	if sql != want {
		t.Errorf("got %s, want %s", sql, want)
	}
}

func TestBuilder_CountDistinct(t *testing.T) {
	b := newBuilder().CountDistinct("name").From("users")
	sql := b.SQL()
	want := `SELECT COUNT(distinct name) FROM "users"`
	if sql != want {
		t.Errorf("got %s, want %s", sql, want)
	}
}

func TestBuilder_Reset(t *testing.T) {
	b := newBuilder().Select("id").From("users").Where("x = 1")
	b.Reset()

	if b.table != "" || len(b.columns) != 0 || len(b.wheres) != 0 {
		t.Errorf("builder not reset properly: %+v", b)
	}
}

func TestIn_NotSlice(t *testing.T) {
	sql, args := In("id", 10)
	wantSQL := "id IN (?)"
	wantArgs := []any{10}
	if sql != wantSQL || !reflect.DeepEqual(args, wantArgs) {
		t.Errorf("got (%s, %v), want (%s, %v)", sql, args, wantSQL, wantArgs)
	}
}

func TestIn_Slice(t *testing.T) {
	sql, args := In("id", []int{1, 2, 3})
	wantSQL := "id IN (?, ?, ?)"
	wantArgs := []any{1, 2, 3}

	// normalize spaces for comparison
	sql = strings.ReplaceAll(sql, " ", "")
	wantSQL = strings.ReplaceAll(wantSQL, " ", "")

	if sql != wantSQL {
		t.Errorf("SQL mismatch: got %q, want %q", sql, wantSQL)
	}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Errorf("args mismatch: got %v, want %v", args, wantArgs)
	}
}

func TestQuestionHelpers(t *testing.T) {
	got := Question(3)
	want := "?,?,?"
	if got != want {
		t.Errorf("Question(3) = %q, want %q", got, want)
	}

	gotSlice := Questions(3)
	wantSlice := []string{"?", "?", "?"}
	if !reflect.DeepEqual(gotSlice, wantSlice) {
		t.Errorf("Questions(3) = %v, want %v", gotSlice, wantSlice)
	}
}

func TestSQLCmd_String(t *testing.T) {
	tests := []struct {
		cmd  sqlcmd
		want string
	}{
		{cselect, "SELECT"},
		{cinsert, "INSERT"},
		{cdelete, "DELETE"},
		{cupdate, "UPDATE"},
		{99, "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.cmd.String(); got != tt.want {
			t.Errorf("%v.String() = %q, want %q", tt.cmd, got, tt.want)
		}
	}
}
