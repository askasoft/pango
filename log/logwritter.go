package log

import (
	"bytes"

	"github.com/askasoft/pango/ref"
)

// Writer log writer interface
type Writer interface {
	Write(le *Event) error
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

type LogFilter struct {
	Filter Filter // log filter
}

// SetFilter set the log filter
func (ll *LogFilter) SetFilter(filter string) {
	ll.Filter = NewLogFilter(filter)
}

func (ll *LogFilter) Reject(le *Event) bool {
	return ll.Filter != nil && ll.Filter.Reject(le)
}

type LogFormatter struct {
	Formatter Formatter    // log formatter
	bb        bytes.Buffer // log buffer
}

// SetFormat set the log formatter
func (lf *LogFormatter) SetFormat(format string) {
	lf.Formatter = NewLogFormatter(format)
}

// Format format the log event
func (lf *LogFormatter) Format(le *Event) []byte {
	f := lf.Formatter
	if f == nil {
		f = le.Logger.GetFormatter()
		if f == nil {
			f = TextFmtDefault
		}
	}

	lf.bb.Reset()
	f.Write(&lf.bb, le)
	return lf.bb.Bytes()
}

type SubFormatter struct {
	Subjecter Formatter    // log formatter
	sbb       bytes.Buffer // log buffer
}

// SetSubject set the subject formatter
func (sf *SubFormatter) SetSubject(format string) {
	sf.Subjecter = NewLogFormatter(format)
}

// GetFormatter get Formatter
func (sf *SubFormatter) SubFormat(le *Event) []byte {
	f := sf.Subjecter
	if f == nil {
		f = TextFmtDefault
	}

	sf.sbb.Reset()
	f.Write(&sf.sbb, le)
	return sf.sbb.Bytes()
}
