package log

import (
	"github.com/askasoft/pango/log/internal"
	"github.com/askasoft/pango/ref"
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
	return ref.SetProperty(w, k, v)
}

// safeWrite safe write log event
func safeWrite(lw Writer, le *Event) {
	defer func() {
		if r := recover(); r != nil {
			internal.Perror(r)
		}
	}()

	lw.Write(le)
}

// safeFlush safe flush log events
func safeFlush(lw Writer) {
	defer func() {
		if r := recover(); r != nil {
			internal.Perror(r)
		}
	}()

	lw.Flush()
}

// safeClose safe close log writer
func safeClose(lw Writer) {
	defer func() {
		if r := recover(); r != nil {
			internal.Perror(r)
		}
	}()

	lw.Close()
}
