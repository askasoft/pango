package num

import (
	"testing"
)

func TestHumanSize(t *testing.T) {
	cs := []struct {
		w string
		n float64
	}{
		{"1 KB", 1024},
		{"1 MB", 1024 * 1024},
		{"1 MB", 1048576},
		{"2 MB", 2 * MB},
		{"3.42 GB", 3.42 * GB},
		{"12.35 GB", 12.3456 * GB},
		{"5.372 TB", 5.372 * TB},
		{"2.22 PB", 2.22 * PB},
		{"1.049e+06 YB", KB * KB * KB * KB * KB * PB},
	}

	for i, c := range cs {
		a := HumanSize(c.n)
		if a != c.w {
			t.Errorf("[%d] HumanSize(%f) = %v, want %v", i, c.n, a, c.w)
		}
	}
}

func TestParseSize(t *testing.T) {
	cs := []struct {
		w int64
		s string
	}{
		{0, "0"},
		{0, "0b"},
		{0, "0B"},
		{0, "0 B"},
		{32, "32"},
		{32, "32b"},
		{32, "32B"},
		{32 * KB, "32k"},
		{32 * KB, "32K"},
		{32 * KB, "32kb"},
		{32 * KB, "32Kb"},
		{32 * MB, "32Mb"},
		{32 * GB, "32Gb"},
		{32 * TB, "32Tb"},
		{32 * PB, "32Pb"},
		{2 * EB, "2Eb"},

		{32.5 * KB, "32.5kB"},
		{32.5 * KB, "32.5 kB"},
		{32, "32.5 B"},
		{307, "0.3 K"},
		{307, ".3kB"},

		{0, "0."},
		{0, "0. "},
		{0, "0.b"},
		{0, "0.B"},
		{0, "-0"},
		{0, "-0b"},
		{0, "-0B"},
		{0, "-0 b"},
		{0, "-0 B"},
		{32, "32."},
		{32, "32.b"},
		{32, "32.B"},
		{32, "32. b"},
		{32, "32. B"},

		// We do not tolerate extra leading or trailing spaces
		// (except for a space after the number and a missing suffix},.
		{0, "0 "},
	}

	for i, c := range cs {
		a, _ := ParseSize(c.s)
		if a != c.w {
			t.Errorf("[%d] ParseSize(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestParseSizeF(t *testing.T) {
	cs := []struct {
		w float64
		s string
	}{
		{2 * ZB, "2Zb"},
		{2 * YB, "2Yb"},
	}

	for i, c := range cs {
		a, _ := ParseSizeF(c.s)
		if a != c.w {
			t.Errorf("[%d] ParseSizeF(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestParseSizeError(t *testing.T) {
	cs := []struct {
		s string
	}{
		{" 0"},
		{" 0b"},
		{" 0B"},
		{" 0 B"},
		{"0b "},
		{"0B "},
		{"0 B "},

		{""},
		{"hello"},
		{"."},
		{". "},
		{" "},
		{"  "},
		{" ."},
		{" . "},
		{"-32"},
		{"-32b"},
		{"-32B"},
		{"-32 b"},
		{"-32 B"},
		{"32b."},
		{"32B."},
		{"32 b."},
		{"32 B."},
		{"32 bb"},
		{"32 BB"},
		{"32 b b"},
		{"32 B B"},
		{"32  b"},
		{"32  B"},
		{" 32 "},
		{"32m b"},
		{"32bm"},
	}

	for i, c := range cs {
		a, e := ParseSize(c.s)
		if e == nil {
			t.Errorf("[%d] ParseSize(%q) = (%v,nil), want error", i, c.s, a)
		}
	}
}
