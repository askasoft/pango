package sqx

import (
	"math/rand"
	"testing"
)

func oldBindType(driverName string) Binder {
	switch driverName {
	case "postgres", "pgx", "pq-timeouts", "cloudsqlpostgres", "ql":
		return BindDollar
	case "mysql":
		return BindQuestion
	case "sqlite3":
		return BindQuestion
	case "oci8", "ora", "goracle", "godror":
		return BindColon
	case "sqlserver", "azuresql":
		return BindAt
	}
	return BindUnknown
}

func BenchmarkBindSpeed(b *testing.B) {
	testDrivers := []string{
		"postgres", "pgx", "mysql", "sqlite3", "ora", "sqlserver",
	}

	b.Run("old", func(b *testing.B) {
		b.StopTimer()
		var seq []int
		for i := 0; i < b.N; i++ {
			seq = append(seq, rand.Intn(len(testDrivers)))
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			s := oldBindType(testDrivers[seq[i]])
			if s == BindUnknown {
				b.Error("unknown driver")
			}
		}

	})

	b.Run("new", func(b *testing.B) {
		b.StopTimer()
		var seq []int
		for i := 0; i < b.N; i++ {
			seq = append(seq, rand.Intn(len(testDrivers)))
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			s := GetBinder(testDrivers[seq[i]])
			if s == BindUnknown {
				b.Error("unknown driver")
			}
		}
	})
}
