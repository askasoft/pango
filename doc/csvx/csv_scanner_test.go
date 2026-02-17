package csvx

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/askasoft/pango/bol"
	"github.com/askasoft/pango/num"
)

type record struct {
	Bool      bool
	Int       int
	String    string
	Duration  time.Duration
	CreatedAt time.Time
}

func testParseTime(s string) time.Time {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	return t
}

func testFormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func TestScanFileValues(t *testing.T) {
	exps := []record{}

	tm := testParseTime("2001-01-02 03:04:05")
	for i := 1; i < 1000; i++ {
		exps = append(exps, record{
			Bool:      i%2 == 0,
			Int:       i,
			String:    num.Itoa(rand.Int()),
			Duration:  time.Minute * time.Duration(i),
			CreatedAt: tm.Add(time.Hour * time.Duration(i)),
		})
	}

	if err := os.MkdirAll("testdata", 0777); err != nil {
		t.Fatal(err)
	}
	tf, err := os.CreateTemp("testdata", "values-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tf.Name())

	cw := csv.NewWriter(tf)
	_ = cw.Write([]string{"bool", "int", "string", "duration", "createdAt", "unknown"})
	for i, r := range exps {
		_ = cw.Write([]string{
			bol.Btoa(r.Bool),
			num.Itoa(r.Int),
			r.String,
			r.Duration.String(),
			testFormatTime(r.CreatedAt),
			fmt.Sprintf("unknown-%d", i),
		})
	}
	cw.Flush()
	tf.Close()

	var recs []record
	err = ScanFile(tf.Name(), &recs)
	if err != nil {
		t.Fatal(err)
	}

	if len(recs) != len(exps) {
		t.Fatalf("len(recs): %d != len(exps): %d", len(recs), len(exps))
	}

	for i, a := range recs {
		w := exps[i]
		if !reflect.DeepEqual(w, a) {
			t.Errorf("[%d] Error\nActual: %v\n  WANT: %v", i+1, a, w)
		}
	}

	// strict scan
	err = ScanFile(tf.Name(), &recs, true)
	if err == nil {
		t.Fatal("an error is expected but got nil.")
	}
}

func TestScanFilePointers(t *testing.T) {
	exps := []*record{}

	tm := testParseTime("2001-01-02 03:04:05")
	for i := 1; i < 1000; i++ {
		exps = append(exps, &record{
			Bool:      i%2 == 0,
			Int:       i,
			String:    num.Itoa(rand.Int()),
			Duration:  time.Minute * time.Duration(i),
			CreatedAt: tm.Add(time.Hour * time.Duration(i)),
		})
	}

	if err := os.MkdirAll("testdata", 0777); err != nil {
		t.Fatal(err)
	}
	tf, err := os.CreateTemp("testdata", "pointers-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tf.Name())

	cw := csv.NewWriter(tf)
	_ = cw.Write([]string{"bool", "int", "string", "duration", "CreatedAt", "Unknown"})
	for i, r := range exps {
		_ = cw.Write([]string{
			bol.Btoa(r.Bool),
			num.Itoa(r.Int),
			r.String,
			r.Duration.String(),
			testFormatTime(r.CreatedAt),
			fmt.Sprintf("unknown-%d", i),
		})
	}
	cw.Flush()
	tf.Close()

	var recs []*record

	err = ScanFileFS(os.DirFS("testdata"), filepath.Base(tf.Name()), &recs)
	if err != nil {
		t.Fatal(err)
	}

	if len(recs) != len(exps) {
		t.Fatalf("len(recs): %d != len(exps): %d", len(recs), len(exps))
	}

	for i, a := range recs {
		w := exps[i]
		if !reflect.DeepEqual(*w, *a) {
			t.Errorf("[%d] Error\nActual: %v\n  WANT: %v", i+1, *a, *w)
		}
	}

	// strict scan
	err = ScanFileFS(os.DirFS("testdata"), filepath.Base(tf.Name()), &recs, true)
	if err == nil {
		t.Fatal("an error is expected but got nil.")
	}
}
