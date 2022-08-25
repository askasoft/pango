package log

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pandafw/pango/ref"
)

// Writer log writer interface
type Writer interface {
	Write(le *Event)
	Flush()
	Close()
}

// WriterCreator writer create function
type WriterCreator func() Writer

var writerCreators = make(map[string]WriterCreator)

// RegisterWriter register log writer type
func RegisterWriter(name string, wc WriterCreator) {
	writerCreators[name] = wc
}

// CreateWriter create a writer by name
func CreateWriter(name string) Writer {
	if f, ok := writerCreators[name]; ok {
		return f()
	}
	return nil
}

// ConfigWriter config the writer by the configuration map 'c'
func ConfigWriter(w Writer, c map[string]any) error {
	for k, v := range c {
		if k != "" && k[0] != '_' && v != nil {
			if err := setWriterProp(w, k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func setWriterProp(w Writer, k string, v any) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = fmt.Errorf("Panic for set %v: %v", k, er)
		}
	}()

	p := strings.Title(k)
	r := reflect.ValueOf(w)

	m := r.MethodByName("Set" + p)
	if m.IsValid() && m.Type().NumIn() == 1 {
		t := m.Type().In(0)

		i, err := ref.Convert(v, t)
		if err != nil {
			return err
		}

		rs := m.Call([]reflect.Value{reflect.ValueOf(i)})
		for _, r := range rs {
			if err, ok := r.Interface().(error); ok {
				return err
			}
		}
		return nil
	}

	f := r.Elem().FieldByName(p)
	if f.IsValid() && f.CanSet() {
		t := f.Type()

		i, err := ref.Convert(v, t)
		if err != nil {
			return err
		}

		f.Set(reflect.ValueOf(i))
		return nil
	}

	return fmt.Errorf("Missing property %q of %v", k, r.Type())
}
