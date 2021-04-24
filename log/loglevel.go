package log

// Level log level
type Level uint32

// Log level
const (
	LevelNone Level = iota
	LevelFatal
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

// String return level string
func (l Level) String() string {
	switch l {
	case LevelFatal:
		return "FATAL"
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelTrace:
		return "TRACE"
	default:
		return "NONE"
	}
}

// Prefix return level prefix
func (l Level) Prefix() string {
	switch l {
	case LevelFatal:
		return "F"
	case LevelError:
		return "E"
	case LevelWarn:
		return "W"
	case LevelInfo:
		return "I"
	case LevelDebug:
		return "D"
	case LevelTrace:
		return "T"
	default:
		return "N"
	}
}

// ParseLevel parse level from string
func ParseLevel(lvl string) Level {
	if lvl != "" {
		switch lvl[0] {
		case 'f', 'F':
			return LevelFatal
		case 'e', 'E':
			return LevelError
		case 'w', 'W':
			return LevelWarn
		case 'i', 'I':
			return LevelInfo
		case 'd', 'D':
			return LevelDebug
		case 't', 'T':
			return LevelTrace
		}
	}
	return LevelNone
}
