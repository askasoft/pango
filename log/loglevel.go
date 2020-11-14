package log

// Log level
const (
	LevelNone = iota
	LevelFatal
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

// ParseLevel parse level from string
func ParseLevel(lvl string) int {
	if lvl != "" {
		switch lvl[0] {
		case 'f', 'F':
			return LevelFatal
		case 'e', 'E':
			return LevelError
		case 'w', 'W':
			return LevelWarn
		case 'd', 'D':
			return LevelDebug
		case 't', 'T':
			return LevelTrace
		}
	}
	return LevelNone
}
