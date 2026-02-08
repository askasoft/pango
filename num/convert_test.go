package num

import (
	"testing"
)

func TestAtoi(t *testing.T) {
	cs := []struct {
		w int
		s string
		n []int
	}{
		{1, "1", nil},
		{2, "02", nil},
		{0660, "0660", nil},
		{0xf, "0xf", nil},
		{-1, "a", []int{-1}},
	}

	for i, c := range cs {
		a := Atoi(c.s, c.n...)
		if a != c.w {
			t.Errorf("[%d] Atoi(%s, %v) = %v, want %v", i, c.s, c.n, a, c.w)
		}
	}
}

func TestAtol(t *testing.T) {
	cs := []struct {
		w int64
		s string
		n []int64
	}{
		{1, "1", nil},
		{2, "02", nil},
		{0660, "0660", nil},
		{0xf, "0xf", nil},
		{-1, "a", []int64{-1}},
	}

	for i, c := range cs {
		a := Atol(c.s, c.n...)
		if a != c.w {
			t.Errorf("[%d] Atol(%s, %v) = %v, want %v", i, c.s, c.n, a, c.w)
		}
	}
}

func TestAtof(t *testing.T) {
	cs := []struct {
		w float64
		s string
		n []float64
	}{
		{1, "1", nil},
		{-1, "a", []float64{-1}},
	}

	for i, c := range cs {
		a := Atof(c.s, c.n...)
		if a != c.w {
			t.Errorf("[%d] Atof(%s, %v) = %v, want %v", i, c.s, c.n, a, c.w)
		}
	}
}

func TestFtoa(t *testing.T) {
	cs := []struct {
		w string
		n float64
	}{
		{"200", 200},
		{"2", 2},
		{"2.2", 2.2},
		{"2.02", 2.02},
		{"200.02", 200.020},
	}

	for i, c := range cs {
		a := Ftoa(c.n)
		if a != c.w {
			t.Errorf("[%d] Ftoa(%f) = %v, want %v", i, c.n, a, c.w)
		}
	}
}

func TestFtoaWithDigits(t *testing.T) {
	cs := []struct {
		w string
		n float64
		d int
	}{
		{"1", 1.23, 0},
		{"1.2", 1.23, 1},
		{"1.3", 1.26, 1},
		{"1.23", 1.23, 2},
		{"1.23", 1.23, 3},
		{"1.234", 1.234, 3},
		{"1.235", 1.2346, 3},
	}

	for i, c := range cs {
		a := FtoaWithDigits(c.n, c.d)
		if a != c.w {
			t.Errorf("[%d] FtoaWithDigits(%f, %d) = %v, want %v", i, c.n, c.d, a, c.w)
		}
	}
}

func TestFormatRoman(t *testing.T) {
	tests := []struct {
		n       int
		want    string
		wantErr bool
	}{
		// 範囲外
		{0, "", true},
		{-1, "", true},
		{4000, "", true},

		// 単純な数字
		{1, "I", false},
		{3, "III", false},
		{4, "IV", false},
		{5, "V", false},
		{9, "IX", false},
		{10, "X", false},
		{40, "XL", false},
		{50, "L", false},
		{90, "XC", false},
		{100, "C", false},
		{400, "CD", false},
		{500, "D", false},
		{900, "CM", false},
		{1000, "M", false},

		// 複合
		{58, "LVIII", false},       // L (50) + V (5) + III (3)
		{1994, "MCMXCIV", false},   // 1000 + 900 + 90 + 4
		{2023, "MMXXIII", false},   // 2000 + 20 + 3
		{3999, "MMMCMXCIX", false}, // 最大値
	}

	for _, tt := range tests {
		got, err := FormatRoman(tt.n)
		if tt.wantErr {
			if err == nil {
				t.Errorf("FormatRoman(%d) expected error, got none", tt.n)
			}
			continue
		}
		if err != nil {
			t.Errorf("FormatRoman(%d) unexpected error: %v", tt.n, err)
			continue
		}
		if got != tt.want {
			t.Errorf("FormatRoman(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestIntToRoman(t *testing.T) {
	cs := []struct {
		w string
		n int
	}{
		{"III", 3},
		{"LVIII", 58},
		{"MCMXCIV", 1994},
	}

	for i, c := range cs {
		a := IntToRoman(c.n)
		if a != c.w {
			t.Errorf("[%d] IntToRoman(%v) = %v, want %v", i, c.n, a, c.w)
		}
	}
}

func TestParseRoman(t *testing.T) {
	tests := []struct {
		input   string
		want    int
		wantErr bool
	}{
		// 基本形
		{"I", 1, false},
		{"II", 2, false},
		{"III", 3, false},
		{"IV", 4, false},
		{"V", 5, false},
		{"IX", 9, false},
		{"X", 10, false},
		{"XL", 40, false},
		{"L", 50, false},
		{"XC", 90, false},
		{"C", 100, false},
		{"CD", 400, false},
		{"D", 500, false},
		{"CM", 900, false},
		{"M", 1000, false},

		// 複合形
		{"LVIII", 58, false},     // L + V + III = 58
		{"MCMXCIV", 1994, false}, // 1000 + (100 sub 1000) + (10 sub 100) + (1 sub 5)

		// エラーケース
		{"A", 0, true},
		{"", 0, true}, // 空文字列 → エラー扱い
	}

	for _, tt := range tests {
		got, err := ParseRoman(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseRoman(%q) expected error, got none", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseRoman(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseRoman(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestRomanToInt(t *testing.T) {
	cs := []struct {
		w int
		n string
	}{
		{3, "III"},
		{58, "LVIII"},
		{1994, "MCMXCIV"},
	}

	for i, c := range cs {
		a := RomanToInt(c.n)
		if a != c.w {
			t.Errorf("[%d] RomanToInt(%v) = %v, want %v", i, c.n, a, c.w)
		}
	}
}

func TestFormatAlpha(t *testing.T) {
	tests := []struct {
		n       int
		want    string
		wantErr bool
	}{
		// エラーケース
		{0, "", true},
		{-1, "", true},

		// 1桁列
		{1, "A", false},
		{2, "B", false},
		{26, "Z", false},

		// 2桁列
		{27, "AA", false},
		{28, "AB", false},
		{52, "AZ", false},
		{53, "BA", false},

		// 3桁列（Excel列番号の拡張）
		{701, "ZY", false},
		{702, "ZZ", false},
		{703, "AAA", false},
	}

	for _, tt := range tests {
		got, err := FormatAlpha(tt.n)
		if tt.wantErr {
			if err == nil {
				t.Errorf("FormatAlpha(%d) expected error, got none", tt.n)
			}
			continue
		}
		if err != nil {
			t.Errorf("FormatAlpha(%d) unexpected error: %v", tt.n, err)
			continue
		}
		if got != tt.want {
			t.Errorf("FormatAlpha(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestParseAlpha(t *testing.T) {
	tests := []struct {
		input   string
		want    int
		wantErr bool
	}{
		{"A", 1, false},
		{"Z", 26, false},
		{"AA", 27, false},
		{"AZ", 52, false},
		{"BA", 53, false},
		{"ZZ", 702, false},
		{"AAA", 703, false},
		{"a", 1, false},   // 小文字対応
		{"Az", 52, false}, // 大文字+小文字混在
		{"A1", 0, true},   // 数字混入 → エラー
		{"*", 0, true},    // 記号 → エラー
		{"", 0, true},     // 空文字列 → エラー（必要なら本体関数側で扱う）
	}

	for _, tt := range tests {
		got, err := ParseAlpha(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseAlpha(%q) expected error, got none", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseAlpha(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseAlpha(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
