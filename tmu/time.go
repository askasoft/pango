package tmu

import "time"

// Atod convert string to time.Duration.
// if not found or convert error, returns the defs[0] or zero.
func Atod(s string, defs ...time.Duration) time.Duration {
	if d, err := time.ParseDuration(s); err == nil {
		return d
	}
	if len(defs) > 0 {
		return defs[0]
	}
	return 0
}
