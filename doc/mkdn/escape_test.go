package mkdn

import (
	"testing"
)

func TestEscapeString(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"hello", "hello"}, // no specials
		{"#", "\\#"},       // single special
		{"a~b", "a\\~b"},   // inside string
		{"use *bold* and _italic_", "use \\*bold\\* and \\_italic\\_"},
		{"100% done!", "100% done\\!"}, // mixed with normal chars
		{"backtick `code`", "backtick \\`code\\`"},
		{"escape | pipe", "escape \\| pipe"},
		{"nested {curly} [brackets] (paren)", "nested \\{curly\\} \\[brackets\\] \\(paren\\)"},
	}

	for _, tt := range tests {
		got := EscapeString(tt.in)
		if got != tt.want {
			t.Errorf("EscapeString(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestUnescapeString(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"hello", "hello"}, // no escapes
		{"\\#", "#"},       // single escape
		{"a\\+b", "a+b"},   // unescape plus
		{"use \\*bold\\* and \\_italic\\_", "use *bold* and _italic_"},
		{"100% done\\!", "100% done!"},
		{"backtick \\`code\\`", "backtick `code`"},
		{"escape \\| pipe", "escape | pipe"},
		{"nested \\{curly\\} \\[brackets\\] \\(paren\\)", "nested {curly} [brackets] (paren)"},
		{"hello\\x", "hello\\x"}, // invalid escape (not followed by a special) -> keep backslash
		{"hello\\", "hello\\"},   // invalid escape (last character) -> keep backslash
	}

	for _, tt := range tests {
		got := UnescapeString(tt.in)
		if got != tt.want {
			t.Errorf("UnescapeString(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestEscapeUnescapeRoundTrip(t *testing.T) {
	inputs := []string{
		"plain text",
		"# heading",
		"some *emphasis* and _underline_",
		"{json} [array] (group)",
		"escape | pipe and backtick `here`",
		"100% done!",
	}
	for _, in := range inputs {
		escaped := EscapeString(in)
		unescaped := UnescapeString(escaped)
		if unescaped != in {
			t.Errorf("roundtrip failed: in=%q escaped=%q unescaped=%q", in, escaped, unescaped)
		}
	}
}
