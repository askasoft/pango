package elt

import "testing"

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		prefix   string
		suffix   string
		data     any
		want     string
	}{
		{
			name:     "no prefix found",
			template: "Hello world!",
			prefix:   "{{",
			suffix:   "}}",
			data:     "ignored",
			want:     "Hello world!",
		},
		{
			name:     "basic replacement",
			template: "This is a {{key}}.",
			prefix:   "{{",
			suffix:   "}}",
			data:     map[string]string{"key": "<<key>>"},
			want:     "This is a <<key>>.",
		},
		{
			name:     "same prefix suffix",
			template: "This ||is|| a ||key||.",
			prefix:   "||",
			suffix:   "||",
			data:     map[string]string{"is": "<<is>>", "key": "<<key>>"},
			want:     "This <<is>> a <<key>>.",
		},
		{
			name:     "cross prefix suffix",
			template: "This {}is}{} a {}key}{}.",
			prefix:   "{}",
			suffix:   "}{",
			data:     map[string]string{"is": "<<is>>", "key": "<<key>>"},
			want:     "This <<is>>} a <<key>>}.",
		},
		{
			name:     "multiple replacements",
			template: "{{a}} + {{b}} = {{c}}",
			prefix:   "{{",
			suffix:   "}}",
			data:     map[string]string{"a": "[a]", "b": "[b]", "c": "[c]"},
			want:     "[a] + [b] = [c]",
		},
		{
			name:     "prefix found but no suffix",
			template: "This is a {{key.",
			prefix:   "{{",
			suffix:   "}}",
			data:     map[string]string{"key": "<<key>>"},
			want:     "This is a {{key.",
		},
		{
			name:     "prefix2 found but no suffix",
			template: "This {{is}} a {{key.",
			prefix:   "{{",
			suffix:   "}}",
			data:     map[string]string{"is": "<<is>>", "key": "<<key>>"},
			want:     "This <<is>> a {{key.",
		},
		{
			name:     "empty string",
			template: "",
			prefix:   "{{",
			suffix:   "}}",
			data:     "empty",
			want:     "",
		},
		{
			name:     "empty key",
			template: "This is {{}}.",
			prefix:   "{{",
			suffix:   "}}",
			data:     "empty",
			want:     "This is .",
		},
		{
			name:     "nested-like pattern",
			template: ":{{outer start{{inner}}}}{{outer2{{inner2}}}}",
			prefix:   "{{",
			suffix:   "}}",
			data:     map[string]string{"inner": "(inner)", "inner2": "(inner2)"},
			want:     ":{{outer start(inner)}}{{outer2(inner2)}}",
		},
		{
			name:     "different prefix/suffix",
			template: "This is a <<key>>.",
			prefix:   "<<",
			suffix:   ">>",
			data:     map[string]string{"key": "[key]"},
			want:     "This is a [key].",
		},
		{
			name:     "adjacent keys",
			template: "{{a}}{{b}}{{c}}",
			prefix:   "{{",
			suffix:   "}}",
			data:     map[string]string{"a": "aa", "b": "bb", "c": "cc"},
			want:     "aabbcc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elt, err := Parse(tt.template, tt.prefix, tt.suffix)
			if err != nil {
				t.Errorf("Parse(%q) = %v", tt.template, err)
				return
			}

			got, err := elt.Evaluate(tt.data)
			if err != nil {
				t.Errorf("Evaluate(%q, %v) = %v", tt.template, tt.data, err)
				return
			}

			if got != tt.want {
				t.Errorf("Evaluate(%q, %v) = %q, want %q", tt.template, tt.data, got, tt.want)
			}
		})
	}
}
