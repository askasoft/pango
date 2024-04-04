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
	case "sqlserver":
		return BindAt
	}
	return BindUnknown
}

/*
sync.Map implementation:

goos: linux
goarch: amd64
pkg: github.com/askasoft/pango/sqx
BenchmarkBindSpeed/old-4         	100000000	        11.0 ns/op
BenchmarkBindSpeed/new-4         	24575726	        50.8 ns/op


async.Value map implementation:

goos: linux
goarch: amd64
pkg: github.com/askasoft/pango/sqx
BenchmarkBindSpeed/old-4         	100000000	        11.0 ns/op
BenchmarkBindSpeed/new-4         	42535839	        27.5 ns/op
*/

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
