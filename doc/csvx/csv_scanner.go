package csvx

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"reflect"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
)

// CsvScanner a csv scanner to read csv record to struct.
type CsvScanner struct {
	cr   *csv.Reader
	Head []string
	Line int

	disallowUnknownFields bool
}

// NewScanner returns a new csv scanner that reads from r.
func NewScanner(cr *csv.Reader) *CsvScanner {
	return &CsvScanner{
		cr: cr,
	}
}

// ScanFile read csv file data to slice.
// Example:
//
//	var s1 []*struct{I int, B bool}
//	ScanFile("s1.csv", &s1)
//
//	var S2 []struct{I int, B bool}
//	ScanFile("s2.csv", &s2)
func ScanFile(name string, records any, strict ...bool) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return ScanReader(f, records, strict...)
}

// ScanFileFS read csv file data to slice.
// Example:
//
//	var s1 []*struct{I int, B bool}
//	ScanFile("s1.csv", &s1)
//
//	var S2 []struct{I int, B bool}
//	ScanFileFS(fsys, "s2.csv", &s2)
func ScanFileFS(fsys fs.FS, name string, records any, strict ...bool) error {
	f, err := fsys.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return ScanReader(f, records, strict...)
}

// ScanReader read csv data to slice.
// Example:
//
//	var s1 []*struct{I int, B bool}
//	ScanReader(r1, &s1)
//
//	var S2 []struct{I int, B bool}
//	ScanReader(r2, &s2)
func ScanReader(r io.Reader, records any, strict ...bool) error {
	br, err := iox.SkipBOM(r)
	if err != nil {
		return err
	}

	cr := csv.NewReader(br)
	return ScanCsv(cr, records, strict...)
}

// ScanCsv read csv data to slice.
// Example:
//
//	var s1 []*struct{I int, B bool}
//	ScanCsv(r1, &s1)
//
//	var S2 []struct{I int, B bool}
//	ScanCsv(r2, &s2)
func ScanCsv(cr *csv.Reader, records any, strict ...bool) error {
	cs := NewScanner(cr)

	if asg.First(strict) {
		cs.DisallowUnknownFields()
	}

	if err := cs.ScanHead(); err != nil {
		return err
	}

	return cs.ScanStructs(records)
}

// DisallowUnknownFields causes the CsvScanner to return an error when the destination
// is a struct and the input contains object keys which do not match any exported fields in the destination.
func (cs *CsvScanner) DisallowUnknownFields() {
	cs.disallowUnknownFields = true
}

// ScanHead reads one record from csv and treat it as header.
func (cs *CsvScanner) ScanHead() error {
	cs.Line++
	hs, err := cs.cr.Read()
	if err != nil {
		return err
	}

	for i, h := range hs {
		h = str.Strip(h)
		if h == "" {
			return fmt.Errorf("line %d column %d: empty column name", cs.Line, i+1)
		}
		hs[i] = h
	}
	cs.Head = hs

	return nil
}

// ScanStruct reads one record from csv and scan it to the parameter `rec`.
// The parameter `rec` should be a pointer to struct.
func (cs *CsvScanner) ScanStruct(rec any) error {
	if !ref.IsPtrType(rec) {
		return fmt.Errorf("%T is not a pointer", rec)
	}

	cs.Line++
	record, err := cs.cr.Read()
	if err != nil {
		return fmt.Errorf("line %d: %w", cs.Line, err)
	}

	if len(record) > len(cs.Head) {
		return fmt.Errorf("line %d: %w", cs.Line, csv.ErrFieldCount)
	}

	for i, s := range record {
		h := cs.Head[i]
		if err := ref.SetProperty(rec, h, s); err != nil {
			if !cs.disallowUnknownFields && ref.IsMissingFieldError(err) {
				continue
			}
			return fmt.Errorf("line %d column %d: %q data error: %q", cs.Line, i+1, h, s)
		}
	}
	return nil
}

func (cs *CsvScanner) ScanStructs(recs any) error {
	pv := reflect.ValueOf(recs)
	if pv.Kind() != reflect.Pointer {
		return fmt.Errorf("%T is not a pointer", recs)
	}

	sv := pv.Elem()
	if sv.Kind() != reflect.Slice {
		return fmt.Errorf("%T is not a slice", sv)
	}

	et := sv.Type().Elem()
	eptr := et.Kind() == reflect.Pointer
	if eptr {
		et = et.Elem()
	}

	for {
		pv := reflect.New(et)
		p := pv.Interface()

		if err := cs.ScanStruct(p); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		ev := pv
		if !eptr {
			ev = pv.Elem()
		}
		sv.Set(reflect.Append(sv, ev))
	}
}
