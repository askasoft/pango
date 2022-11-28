package str

import "testing"

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"abacaba", false},
		{"", true},
	}
	for i, c := range cs {
		a := IsEmpty(c.s)
		if a != c.w {
			t.Errorf("[%d] IsEmpty(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsNotEmpty(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"abacaba", true},
		{"", false},
	}
	for i, c := range cs {
		a := IsNotEmpty(c.s)
		if a != c.w {
			t.Errorf("[%d] IsEmpty(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsASCII(t *testing.T) {
	t.Parallel()

	cs := []struct {
		w bool
		s string
	}{
		{false, ""},
		{false, "ｆｏｏbar"},
		{false, "ｘｙｚ０９８"},
		{false, "１２３456"},
		{false, "ｶﾀｶﾅ"},
		{true, "foobar"},
		{true, "0987654321"},
		{true, "test@example.com"},
		{true, "1234abcDEF"},
		{false, ""},
	}

	for i, c := range cs {
		a := IsASCII(c.s)
		if a != c.w {
			t.Errorf("[%d] IsASCII(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsPrintableASCII(t *testing.T) {
	t.Parallel()

	cs := []struct {
		w bool
		s string
	}{
		{false, ""},
		{false, "ｆｏｏbar"},
		{false, "ｘｙｚ０９８"},
		{false, "１２３456"},
		{false, "ｶﾀｶﾅ"},
		{true, "foobar"},
		{true, "0987654321"},
		{true, "test@example.com"},
		{true, "1234abcDEF"},
		{false, "newline\n"},
		{false, "\x19test\x7F"},
	}

	for i, c := range cs {
		a := IsPrintableASCII(c.s)
		if a != c.w {
			t.Errorf("[%d] IsPrintableASCII(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsAlpha(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", false},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"abc", true},
		{"소주", false},
		{"ABC", true},
		{"FoObAr", true},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", true},  //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
		{"\u0030", false}, //UTF-8(ASCII): 0
		{"123", false},
		{"0123", false},
		{"-00123", false},
		{"0", false},
		{"-0", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", false},
		{"ix", true},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}
	for i, c := range cs {
		a := IsAlpha(c.s)
		if a != c.w {
			t.Errorf("[%d] IsAlpha(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsUTFLetter(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", false},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"", false},
		{"abc", true},
		{"소주", true},
		{"ABC", true},
		{"FoObAr", true},
		{"소aBC", true},
		{"소", true},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", true},  //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
		{"\u0030", false}, //UTF-8(ASCII): 0
		{"123", false},
		{"0123", false},
		{"-00123", false},
		{"0", false},
		{"-0", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", true},
		{"ix", true},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}
	for i, c := range cs {
		a := IsUTFLetter(c.s)
		if a != c.w {
			t.Errorf("[%d] IsUTFLetter(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsAlphaNumber(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", false},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc123", true},
		{"ABC111", true},
		{"abc1", true},
		{"abc〩", false},
		{"abc", true},
		{"소주", false},
		{"ABC", true},
		{"FoObAr", true},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", true},  //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
		{"\u0030", true},  //UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", false},
		{"0", true},
		{"-0", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", false},
		{"ix", true},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}
	for i, c := range cs {
		a := IsAlphaNumber(c.s)
		if a != c.w {
			t.Errorf("[%d] IsAlphaNumber(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsUTFLetterNumber(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", true},
		{"", false},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", true},
		{"abc〩", true},
		{"abc", true},
		{"소주", true},
		{"ABC", true},
		{"FoObAr", true},
		{"소aBC", true},
		{"소", true},
		{"달기&Co.", false},
		{"〩Hours", true},
		{"\ufff0", false},
		{"\u0070", true},  //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
		{"\u0030", true},  //UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", false},
		{"0", true},
		{"-0", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"-1¾", false},
		{"1¾", true},
		{"〥〩", true},
		{"모자", true},
		{"ix", true},
		{"۳۵۶۰", true},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}
	for i, c := range cs {
		a := IsUTFLetterNumber(c.s)
		if a != c.w {
			t.Errorf("[%d] IsUTFLetterNumber(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsNumber(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", false},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"abc", false},
		{"소주", false},
		{"ABC", false},
		{"FoObAr", false},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", false}, //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
		{"\u0030", true},  //UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", false},
		{"+00123", false},
		{"0", true},
		{"-0", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"12𐅪3", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", false},
		{"ix", false},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}
	for i, c := range cs {
		a := IsNumber(c.s)
		if a != c.w {
			t.Errorf("[%d] IsNumber(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsNumeric(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", false},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"abc", false},
		{"소주", false},
		{"ABC", false},
		{"FoObAr", false},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", false}, //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
		{"\u0030", true},  //UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", true},
		{"+00123", true},
		{"0", true},
		{"-0", true},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"12𐅪3", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", false},
		{"ix", false},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", true},
	}
	for i, c := range cs {
		a := IsNumeric(c.s)
		if a != c.w {
			t.Errorf("[%d] IsNumeric(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsUTFNumeric(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", true},
		{"", false},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"abc", false},
		{"소주", false},
		{"ABC", false},
		{"FoObAr", false},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", false}, //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
		{"\u0030", true},  //UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", true},
		{"0", true},
		{"-0", true},
		{"--0", false},
		{"-0-", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"12𐅪3", true},
		{"-1¾", true},
		{"1¾", true},
		{"〥〩", true},
		{"모자", false},
		{"ix", false},
		{"۳۵۶۰", true},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", true},
	}
	for i, c := range cs {
		a := IsUTFNumeric(c.s)
		if a != c.w {
			t.Errorf("[%d] IsUTFNumeric(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsUTFDigit(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{

		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", false},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"abc", false},
		{"소주", false},
		{"ABC", false},
		{"FoObAr", false},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", false}, //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
		{"\u0030", true},  //UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", true},
		{"0", true},
		{"-0", true},
		{"--0", false},
		{"-0-", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"12𐅪3", false},
		{"1483920", true},
		{"", false},
		{"۳۵۶۰", true},
		{"-29", true},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", false},
		{"ix", false},
		{"۳۵۶۰", true},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", true},
	}
	for i, c := range cs {
		a := IsUTFDigit(c.s)
		if a != c.w {
			t.Errorf("[%d] IsUTFDigit(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsLowerCase(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"abc123", true},
		{"abc", true},
		{"a b c", true},
		{"abcß", true},
		{"abcẞ", false},
		{"ABCẞ", false},
		{"tr竪s 端ber", true},
		{"fooBar", false},
		{"123ABC", false},
		{"ABC123", false},
		{"ABC", false},
		{"S T R", false},
		{"fooBar", false},
		{"abacaba123", true},
	}
	for i, c := range cs {
		a := IsLowerCase(c.s)
		if a != c.w {
			t.Errorf("[%d] IsLowerCase(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsUpperCase(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"abc123", false},
		{"abc", false},
		{"a b c", false},
		{"abcß", false},
		{"abcẞ", false},
		{"ABCẞ", true},
		{"tr竪s 端ber", false},
		{"fooBar", false},
		{"123ABC", true},
		{"ABC123", true},
		{"ABC", true},
		{"S T R", true},
		{"fooBar", false},
		{"abacaba123", false},
	}
	for i, c := range cs {
		a := IsUpperCase(c.s)
		if a != c.w {
			t.Errorf("[%d] IsUpperCase(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestHasLowerCase(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"abc123", true},
		{"abc", true},
		{"a b c", true},
		{"abcß", true},
		{"abcẞ", true},
		{"ABCẞ", false},
		{"tr竪s 端ber", true},
		{"fooBar", true},
		{"123ABC", false},
		{"ABC123", false},
		{"ABC", false},
		{"S T R", false},
		{"fooBar", true},
		{"abacaba123", true},
		{"FÒÔBÀŘ", false},
		{"fòôbàř", true},
		{"fÒÔBÀŘ", true},
	}
	for i, c := range cs {
		a := HasLowerCase(c.s)
		if a != c.w {
			t.Errorf("[%d] HasLowerCase(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestHasUpperCase(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"abc123", false},
		{"abc", false},
		{"a b c", false},
		{"abcß", false},
		{"abcẞ", true},
		{"ABCẞ", true},
		{"tr竪s 端ber", false},
		{"fooBar", true},
		{"123ABC", true},
		{"ABC123", true},
		{"ABC", true},
		{"S T R", true},
		{"fooBar", true},
		{"abacaba123", false},
		{"FÒÔBÀŘ", true},
		{"fòôbàř", false},
		{"Fòôbàř", true},
	}
	for i, c := range cs {
		a := HasUpperCase(c.s)
		if a != c.w {
			t.Errorf("[%d] HasUpperCase(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsHexDecimal(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"abcdefg", false},
		{"", false},
		{"..", false},
		{"0x", false},
		{"0XdeadBEEF", true},
		{"0xff0044", true},
		{"fe0044", true},
		{"fe0x44", false},
	}
	for i, c := range cs {
		a := IsHexDecimal(c.s)
		if a != c.w {
			t.Errorf("[%d] IsHexDecimal(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}
func TestIsWhitespace(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"abacaba", false},
		{"", false},
		{"    ", true},
		{"  \r\n  ", true},
		{"\014\012\011\013\015", true},
		{"\014\012\011\013 abc  \015", false},
		{"\f\n\t\v\r\f", true},
		{"x\n\t\t\t\t", false},
		{"\f\n\t  \n\n\n   \v\r\f", true},
	}
	for i, c := range cs {
		a := IsWhitespace(c.s)
		if a != c.w {
			t.Errorf("[%d] IsWhitespace(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestHasWhitespace(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"abacaba", false},
		{"", false},
		{"    ", true},
		{"  \r\n  ", true},
		{"\014\012\011\013\015", true},
		{"\014\012\011\013 abc  \015", true},
		{"\f\n\t\v\r\f", true},
		{"x\n\t\t\t\t", true},
		{"\f\n\t  \n\n\n   \v\r\f", true},
	}
	for i, c := range cs {
		a := HasWhitespace(c.s)
		if a != c.w {
			t.Errorf("[%d] HasWhitespace(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestHasMultibyte(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"abc", false},
		{"123", false},
		{"<>@;.-=", false},
		{"ひらがな・カタカナ、．漢字", true},
		{"あいうえお foobar", true},
		{"test＠example.com", true},
		{"test＠example.com", true},
		{"1234abcDEｘｙｚ", true},
		{"ｶﾀｶﾅ", true},
		{"", false},
	}
	for i, c := range cs {
		a := HasMultibyte(c.s)
		if a != c.w {
			t.Errorf("[%d] HasMultibyte(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}
