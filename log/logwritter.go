package log

import (
	"bytes"

	"github.com/askasoft/pango/log/internal"
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

// safeWrite safe write log event
func safeWrite(lw Writer, le *Event) {
	defer func() {
		if r := recover(); r != nil {
			internal.Perror(r)
		}
	}()

	if err := lw.Write(le); err != nil {
		internal.Perror(err)
	}
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
	Buffer    bytes.Buffer // log buffer
}

// SetFormat set the log formatter
func (lf *LogFormatter) SetFormat(format string) {
	lf.Formatter = NewLogFormatter(format)
}

// Format format the log event
func (lf *LogFormatter) GetFormatter(le *Event, df ...Formatter) Formatter {
	f := lf.Formatter
	if f == nil {
		f = le.Logger.GetFormatter()
		if f == nil {
			if len(df) > 0 {
				f = df[0]
			} else {
				f = TextFmtDefault
			}
		}
	}
	return f
}

// Format format the log event
func (lf *LogFormatter) Format(le *Event, df ...Formatter) []byte {
	f := lf.GetFormatter(le, df...)
	lf.Buffer.Reset()
	f.Write(&lf.Buffer, le)
	return lf.Buffer.Bytes()
}

// Append format the log event and append to buffer
func (lf *LogFormatter) Append(le *Event, df ...Formatter) {
	f := lf.GetFormatter(le, df...)
	f.Write(&lf.Buffer, le)
}

type SubFormatter struct {
	Subjecter Formatter    // log formatter
	SubBuffer bytes.Buffer // log buffer
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

	sf.SubBuffer.Reset()
	f.Write(&sf.SubBuffer, le)
	return sf.SubBuffer.Bytes()
}
