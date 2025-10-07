package urlx

import (
	"net/url"
	"testing"
)

func TestCleanURL(t *testing.T) {
	cs := []struct {
		s string
		w string
	}{
		{"http://user:pass@abc.com", "http://abc.com/"},
		{"http://user:pass@abc.com/?a=b#xyz", "http://abc.com/"},
	}

	for i, c := range cs {
		a, err := CleanURL(c.s)
		if err != nil {
			t.Errorf("#%d CleanURL(%q)\n     = %s\nWANT: %s", i, c.s, a, c.w)
		}
	}
}

func TestParentURL(t *testing.T) {
	cs := []struct {
		s string
		w string
	}{
		{"http://user:pass@abc.com/?a=b#xyz", "http://abc.com/"},
		{"http://user:pass@abc.com/abc?a=b#xyz", "http://abc.com/"},
		{"http://user:pass@abc.com/abc/?a=b#xyz", "http://abc.com/abc/"},
	}

	for i, c := range cs {
		a, err := ParentURL(c.s)
		if err != nil {
			t.Errorf("#%d ParentURL(%q)\n     = %s\nWANT: %s", i, c.s, a, c.w)
		}
	}
}

func TestEncodeQuery(t *testing.T) {
	tests := []struct {
		name string
		kvs  []string
		want string
	}{
		{
			name: "empty input",
			kvs:  nil,
			want: "",
		},
		{
			name: "single pair",
			kvs:  []string{"key", "value"},
			want: "?key=value",
		},
		{
			name: "multiple pairs",
			kvs:  []string{"a", "1", "b", "2", "c", "3"},
			want: "?a=1&b=2&c=3",
		},
		{
			name: "special characters escaped",
			kvs:  []string{"key with space", "value/with?symbols"},
			want: "?key+with+space=value%2Fwith%3Fsymbols",
		},
		{
			name: "odd number of elements (trailing key only)",
			kvs:  []string{"onlykey"},
			want: "?onlykey=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeQuery(tt.kvs...)
			if got != tt.want {
				t.Errorf("EncodeQuery(%v) = %q, want %q", tt.kvs, got, tt.want)
			}

			// validate query syntax if not empty
			if got != "" {
				_, err := url.Parse(got)
				if err != nil {
					t.Errorf("invalid query string: %v", err)
				}
			}
		})
	}
}
