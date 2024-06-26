package ldt

// Info represents a full outcome of language detection.
type Info struct {
	Lang       Lang
	Confidence float64
}

// IsReliable returns true if Confidence is greater than the Reliable Confidence Threshold
func (info *Info) IsReliable() bool {
	return info.Confidence > ReliableConfidenceThreshold
}
