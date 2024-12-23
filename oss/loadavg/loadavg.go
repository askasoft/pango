package loadavg

import "fmt"

// LoadAvg represents load average values
type LoadAvg struct {
	Loadavg1  float64 `json:"loadavg1"`
	Loadavg5  float64 `json:"loadavg5"`
	Loadavg15 float64 `json:"loadavg15"`
}

func (la *LoadAvg) String() string {
	return fmt.Sprintf("(%f,%f,%f)", la.Loadavg1, la.Loadavg5, la.Loadavg15)
}
