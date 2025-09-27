package wildcard

import (
	"strings"
	"testing"
)

func TestContains(t *testing.T) {
	cs := []struct {
		p string
		s string
		w bool
	}{
		// Test case - 1.
		// Test case with p "*". Expected to match any s.
		{
			p: "*",
			s: "s3:GetObject",
			w: true,
		},
		// Test case - 2.
		// Test case with empty p. This only matches empty string.
		{
			p: "",
			s: "s3:GetObject",
			w: false,
		},
		// Test case - 3.
		// Test case with empty p. This only matches empty string.
		{
			p: "",
			s: "",
			w: true,
		},
		// Test case - 4.
		// Test case with single "*" at the end.
		{
			p: "s3:*",
			s: "s3:ListMultipartUploadParts",
			w: true,
		},
		// Test case - 5.
		// Test case with a no "*". In this case the p and s should be the same.
		{
			p: "s3:ListBucket/MultipartUploads",
			s: "s3:ListBucket",
			w: true,
		},
		// Test case - 6.
		// Test case with a no "*". In this case the p and s should be the same.
		{
			p: "s3:ListBucket",
			s: "s3:ListBucket",
			w: true,
		},
		// Test case - 7.
		// Test case with a no "*". In this case the p and s should be the same.
		{
			p: "s3:ListBucketMultipartUploads",
			s: "s3:ListBucketMultipartUploads",
			w: true,
		},
		// Test case - 8.
		// Test case with p containing key name with a prefix. Should accept the same s without a "*".
		{
			p: "my-bucket/oo*",
			s: "my-bucket/oo",
			w: true,
		},
		// Test case - 9.
		// Test case with "*" at the end of the p.
		{
			p: "my-bucket/In*",
			s: "my-bucket/India/Karnataka/",
			w: true,
		},
		// Test case - 10.
		// Test case with prefixes shuffled.
		// This should fail.
		{
			p: "my-bucket/In*",
			s: "my-bucket/Karnataka/India/",
			w: false,
		},
		// Test case - 11.
		// Test case with s expanded to the wildcards in the p.
		{
			p: "my-bucket/In*/Ka*/Ban",
			s: "my-bucket/India/Karnataka/Ban",
			w: true,
		},
		// Test case - 12.
		// Test case with the  keyname part is repeated as prefix several times.
		// This is valid.
		{
			p: "my-bucket/In*/Ka*/Ban",
			s: "my-bucket/India/Karnataka/Ban/Ban/Ban/Ban/Ban",
			w: true,
		},
		// Test case - 13.
		// Test case to validate that `*` can be expanded into multiple prefixes.
		{
			p: "my-bucket/In**/Ka**/Ban",
			s: "my-bucket/India/Karnataka/Area1/Area2/Area3/Ban",
			w: true,
		},
		// Test case - 14.
		// Test case to validate that `*` can be expanded into multiple prefixes.
		{
			p: "my-bucket/In*/Ka*/Ban",
			s: "my-bucket/India/State1/State2/Karnataka/Area1/Area2/Area3/Ban",
			w: true,
		},
		// Test case - 15.
		// Test case where the keyname part of the p is expanded in the s.
		{
			p: "my-bucket/In*/Ka*/Ban",
			s: "my-bucket/India/Karnataka/Bangalore",
			w: true,
		},
		// Test case - 16.
		// Test case with prefixes and wildcard expanded for all "*".
		{
			p: "my-bucket/In*/Ka*/Ban*",
			s: "my-bucket/India/Karnataka/Bangalore",
			w: true,
		},
		// Test case - 17.
		// Test case with keyname part being a wildcard in the p.
		{
			p: "my-bucket/*",
			s: "my-bucket/India",
			w: true,
		},
		// Test case - 18.
		{
			p: "my-bucket/oo*",
			s: "my-bucket/odo",
			w: false,
		},

		// Test case with p containing wildcard '?'.
		// Test case - 19.
		// "my-bucket?/" matches "my-bucket1/", "my-bucket2/", "my-bucket3" etc...
		// doesn't match "mybucket/".
		{
			p: "my-bucket?/abc*",
			s: "mybucket/abc",
			w: false,
		},
		// Test case - 20.
		{
			p: "my-bucket?/abc*",
			s: "my-bucket1/abc",
			w: true,
		},
		// Test case - 21.
		{
			p: "my-?-bucket/abc*",
			s: "my--bucket/abc",
			w: false,
		},
		// Test case - 22.
		{
			p: "my-?-bucket/abc*",
			s: "my-1-bucket/abc",
			w: true,
		},
		// Test case - 23.
		{
			p: "my-?-bucket/abc*",
			s: "my-k-bucket/abc",
			w: true,
		},
		// Test case - 24.
		{
			p: "my??bucket/abc*",
			s: "mybucket/abc",
			w: false,
		},
		// Test case - 25.
		{
			p: "my??bucket/abc*",
			s: "my4abucket/abc",
			w: true,
		},
		// Test case - 26.
		{
			p: "my-bucket?abc*",
			s: "my-bucket/abc",
			w: true,
		},
		// Test case 27-28.
		// '?' matches '/' too. (works with s3).
		// This is because the namespace is considered flat.
		// "abc?efg" matches both "abcdefg" and "abc/efg".
		{
			p: "my-bucket/abc?efg",
			s: "my-bucket/abcdefg",
			w: true,
		},
		{
			p: "my-bucket/abc?efg",
			s: "my-bucket/abc/efg",
			w: true,
		},
		// Test case - 29.
		{
			p: "my-bucket/abc????",
			s: "my-bucket/abc",
			w: true,
		},
		// Test case - 30.
		{
			p: "my-bucket/abc????",
			s: "my-bucket/abcde",
			w: true,
		},
		// Test case - 31.
		{
			p: "my-bucket/abc????",
			s: "my-bucket/abcdefg",
			w: true,
		},
		// Test case 32-34.
		// test case with no '*'.
		{
			p: "my-bucket/abc?",
			s: "my-bucket/abc",
			w: true,
		},
		{
			p: "my-bucket/abc?",
			s: "my-bucket/abcd",
			w: true,
		},
		{
			p: "my-bucket/abc?",
			s: "my-bucket/abcde",
			w: false,
		},
		// Test case 35.
		{
			p: "my-bucket/mnop*?",
			s: "my-bucket/mnop",
			w: true,
		},
		// Test case 36.
		{
			p: "my-bucket/mnop*?",
			s: "my-bucket/mnopqrst/mnopqr",
			w: true,
		},
		// Test case 37.
		{
			p: "my-bucket/mnop*?",
			s: "my-bucket/mnopqrst/mnopqrs",
			w: true,
		},
		// Test case 38.
		{
			p: "my-bucket/mnop*?",
			s: "my-bucket/mnop",
			w: true,
		},
		// Test case 39.
		{
			p: "my-bucket/mnop*?",
			s: "my-bucket/mnopq",
			w: true,
		},
		// Test case 40.
		{
			p: "my-bucket/mnop*?",
			s: "my-bucket/mnopqr",
			w: true,
		},
		// Test case 41.
		{
			p: "my-bucket/mnop*?and",
			s: "my-bucket/mnopqand",
			w: true,
		},
		// Test case 42.
		{
			p: "my-bucket/mnop*?and",
			s: "my-bucket/mnopand",
			w: true,
		},
		// Test case 43.
		{
			p: "my-bucket/mnop*?and",
			s: "my-bucket/mnopqand",
			w: true,
		},
		// Test case 44.
		{
			p: "my-bucket/mnop*?",
			s: "my-bucket/mn",
			w: true,
		},
		// Test case 45.
		{
			p: "my-bucket/mnop*?",
			s: "my-bucket/mnopqrst/mnopqrs",
			w: true,
		},
		// Test case 46.
		{
			p: "my-bucket/mnop*??",
			s: "my-bucket/mnopqrst",
			w: true,
		},
		// Test case 47.
		{
			p: "my-bucket/mnop*qrst",
			s: "my-bucket/mnopabcdegqrst",
			w: true,
		},
		// Test case 48.
		{
			p: "my-bucket/mnop*?and",
			s: "my-bucket/mnopqand",
			w: true,
		},
		// Test case 49.
		{
			p: "my-bucket/mnop*?and",
			s: "my-bucket/mnopand",
			w: true,
		},
		// Test case 50.
		{
			p: "my-bucket/mnop*?and?",
			s: "my-bucket/mnopqanda",
			w: true,
		},
		// Test case 51.
		{
			p: "my-bucket/mnop*?and",
			s: "my-bucket/mnopqanda",
			w: true,
		},
		// Test case 52.
		{
			p: "my-?-bucket/abc*",
			s: "my-bucket/mnopqanda",
			w: false,
		},
		// Test case 53.
		{
			p: "abc*",
			s: "abc",
			w: true,
		},
	}

	for i, c := range cs {
		a := Contains(c.p, c.s)
		if c.w != a {
			t.Errorf("%d: Contains(%q, %q) = %v, want: %v", i+1, c.p, c.s, a, c.w)
		}

		if c.w {
			// fold match
			foldText := strings.ToUpper(c.s)
			if foldText == c.s {
				foldText = strings.ToLower(c.s)
			}
			a = ContainsFold(c.p, foldText)
			if c.w != a {
				t.Errorf("%d: ContainsFold(%q, %q) = %v, want: %v", i+1, c.p, foldText, a, c.w)
			}
		}
	}
}

func TestContainsSimple(t *testing.T) {
	cs := []struct {
		p string
		s string
		w bool
	}{
		// Test case - 1.
		// Test case with p "*". Expected to match any s.
		{
			p: "*",
			s: "s3:GetObject",
			w: true,
		},
		// Test case - 2.
		// Test case with empty p. This only matches empty string.
		{
			p: "",
			s: "s3:GetObject",
			w: false,
		},
		// Test case - 3.
		// Test case with empty p. This only matches empty string.
		{
			p: "",
			s: "",
			w: true,
		},
		// Test case - 4.
		// Test case with single "*" at the end.
		{
			p: "s3:*",
			s: "s3:ListMultipartUploadParts",
			w: true,
		},
		// Test case - 5.
		// Test case with a no "*". In this case the p and s should be the same.
		{
			p: "s3:ListBucketMultipartUploads",
			s: "s3:ListBucket",
			w: true,
		},
		// Test case - 6.
		// Test case with a no "*". In this case the p and s should be the same.
		{
			p: "s3:ListBucket",
			s: "s3:ListBucket",
			w: true,
		},
		// Test case - 7.
		// Test case with a no "*". In this case the p and s should be the same.
		{
			p: "s3:ListBucketMultipartUploads",
			s: "s3:ListBucketMultipartUploads",
			w: true,
		},
		// Test case - 8.
		// Test case with p containing key name with a prefix. Should accept the same s without a "*".
		{
			p: "my-bucket/oo*",
			s: "my-bucket/oo",
			w: true,
		},
		// Test case - 9.
		// Test case with "*" at the end of the p.
		{
			p: "my-bucket/In*",
			s: "my-bucket/India/Karnataka/",
			w: true,
		},
		// Test case - 10.
		// Test case with prefixes shuffled.
		// This should fail.
		{
			p: "my-bucket/In*",
			s: "my-bucket/Karnataka/India/",
			w: false,
		},
		// Test case - 11.
		// Test case with s expanded to the wildcards in the p.
		{
			p: "my-bucket/In*/Ka*/Ban",
			s: "my-bucket/India/Karnataka/Ban",
			w: true,
		},
		// Test case - 12.
		// Test case with the  keyname part is repeated as prefix several times.
		// This is valid.
		{
			p: "my-bucket/In*/Ka*/Ban",
			s: "my-bucket/India/Karnataka/Ban/Ban/Ban/Ban/Ban",
			w: true,
		},
		// Test case - 13.
		// Test case to validate that `*` can be expanded into multiple prefixes.
		{
			p: "my-bucket/In*/Ka*/Ban",
			s: "my-bucket/India/Karnataka/Area1/Area2/Area3/Ban",
			w: true,
		},
		// Test case - 14.
		// Test case to validate that `*` can be expanded into multiple prefixes.
		{
			p: "my-bucket/In*/Ka*/Ban",
			s: "my-bucket/India/State1/State2/Karnataka/Area1/Area2/Area3/Ban",
			w: true,
		},
		// Test case - 15.
		// Test case where the keyname part of the p is expanded in the s.
		{
			p: "my-bucket/In*/Ka*/Ban",
			s: "my-bucket/India/Karnataka/Bangalore",
			w: true,
		},
		// Test case - 16.
		// Test case with prefixes and wildcard expanded for all "*".
		{
			p: "my-bucket/In*/Ka*/Ban*",
			s: "my-bucket/India/Karnataka/Bangalore",
			w: true,
		},
		// Test case - 17.
		// Test case with keyname part being a wildcard in the p.
		{
			p: "my-bucket/*",
			s: "my-bucket/India",
			w: true,
		},
		// Test case - 18.
		{
			p: "my-bucket/oo*",
			s: "my-bucket/odo",
			w: false,
		},
		// Test case - 11.
		{
			p: "my-bucket/oo?*",
			s: "my-bucket/oo???",
			w: true,
		},
		// Test case - 12:
		{
			p: "my-bucket/oo??*",
			s: "my-bucket/odo",
			w: false,
		},
		// Test case - 13:
		{
			p: "?h?*",
			s: "?h?hello",
			w: true,
		},
		// Test case 53.
		{
			p: "abc*",
			s: "abc",
			w: true,
		},
	}

	for i, c := range cs {
		a := ContainsSimple(c.p, c.s)
		if c.w != a {
			t.Errorf("%d: ContainsSimple(%q, %q) = %v, want: %v", i+1, c.p, c.s, a, c.w)
		}

		if c.w {
			// fold match
			foldText := strings.ToUpper(c.s)
			if foldText == c.s {
				foldText = strings.ToLower(c.s)
			}
			a = ContainsSimpleFold(c.p, foldText)
			if c.w != a {
				t.Errorf("%d: ContainsSimpleFold(%q, %q) = %v, want: %v", i+1, c.p, foldText, a, c.w)
			}
		}
	}
}
