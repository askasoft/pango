package csvx

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/ref"
	"github.com/askasoft/pango/str"
)

type CsvReader struct {
	*csv.Reader

	Header []string
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *CsvReader {
	return &CsvReader{
		Reader: csv.NewReader(r),
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
func ScanFile(name string, records any) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return ScanReader(f, records)
}

// ScanFile read csv data to slice.
// Example:
//
//	var s1 []*struct{I int, B bool}
//	ScanReader(r1, &s1)
//
//	var S2 []struct{I int, B bool}
//	ScanReader(r2, &s2)
func ScanReader(r io.Reader, records any) error {
	sr, err := iox.SkipBOM(r)
	if err != nil {
		return err
	}

	cr := NewReader(sr)
	if err := cr.ReadHeader(); err != nil {
		return err
	}

	return cr.ScanStructs(records)
}

// ReadHeader reads one record from csv and treat it as header.
func (cr *CsvReader) ReadHeader() error {
	hs, err := cr.Reader.Read()
	if err != nil {
		return err
	}

	for i, h := range hs {
		hs[i] = str.PascalCase(str.Strip(h))
	}
	cr.Header = hs

	return nil
}

// ScanStruct reads one record from csv and scan it to the parameter `rec`.
// The parameter `rec` should be a pointer to struct.
func (cr *CsvReader) ScanStruct(rec any) (err error) {
	if !ref.IsPtrType(rec) {
		return fmt.Errorf("%T is not a pointer", rec)
	}

	var record []string

	record, err = cr.Reader.Read()
	if len(record) > len(cr.Header) {
		return csv.ErrFieldCount
	}

	for i, s := range record {
		h := cr.Header[i]
		if err = ref.SetProperty(rec, h, s); err != nil {
			return
		}
	}
	return
}

func (cr *CsvReader) ScanStructs(recs any) (err error) {
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

		if err = cr.ScanStruct(p); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		ev := pv
		if !eptr {
			ev = pv.Elem()
		}
		sv.Set(reflect.Append(sv, ev))
	}

	return nil
}
