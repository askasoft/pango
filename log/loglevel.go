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

var (
	levelStrings = []string{"NONE", "FATAL", "ERROR", "WARN", "INFO", "DEBUG", "TRACE"}
	levelPrefixs = []string{"N", "F", "E", "W", "I", "D", "T"}
)

// String return level string
func (l Level) String() string {
	if l < LevelNone || l > LevelTrace {
		return "UNKNOWN"
	}
	return levelStrings[l]
}

// Prefix return level prefix
func (l Level) Prefix() string {
	if l < LevelNone || l > LevelTrace {
		return "U"
	}
	return levelPrefixs[l]
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
