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

// LevelFromString parse level from string
func LevelFromString(lvl string) int {
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
