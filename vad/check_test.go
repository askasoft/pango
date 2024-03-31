package vad

import (
	"testing"
)

func TestIsFileName(t *testing.T) {
	t.Parallel()

	cs := []struct {
		w bool
		s string
	}{
		{true, "a.txT"},
		{true, "_"},
		{false, "a.tx/t"},
		{false, "a.tx\\t"},
		{false, "a.tx*t"},
		{false, "a.tx?t"},
	}

	for i, c := range cs {
		a := IsFileName(c.s)
		if a != c.w {
			t.Errorf("[%d] IsFileName(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsHash(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s    string
		algo string
		w    bool
	}{
		{"3ca25ae354e192b26879f651a51d92aa8a34d8d3", "sha1", true},
		{"3ca25ae354e192b26879f651a51d34d8d3", "sha1", false},
		{"3ca25ae354e192b26879f651a51d92aa8a34d8d3", "Tiger160", true},
		{"3ca25ae354e192b26879f651a51d34d8d3", "ripemd160", false},
		{"579282cfb65ca1f109b78536effaf621b853c9f7079664a3fbe2b519f435898c", "sha256", true},
		{"579282cfb65ca1f109b78536effaf621b853c9f7079664a3fbe2b519f435898casfdsafsadfsdf", "sha256", false},
		{"bf547c3fc5841a377eb1519c2890344dbab15c40ae4150b4b34443d2212e5b04aa9d58865bf03d8ae27840fef430b891", "sha384", true},
		{"579282cfb65ca1f109b78536effaf621b853c9f7079664a3fbe2b519f435898casfdsafsadfsdf", "sha384", false},
		{"45bc5fa8cb45ee408c04b6269e9f1e1c17090c5ce26ffeeda2af097735b29953ce547e40ff3ad0d120e5361cc5f9cee35ea91ecd4077f3f589b4d439168f91b9", "sha512", true},
		{"579282cfb65ca1f109b78536effaf621b853c9f7079664a3fbe2b519f435898casfdsafsadfsdf", "sha512", false},
		{"46fc0125a148788a3ac1d649566fc04eb84a746f1a6e4fa7", "tiger192", true},
		{"46fc0125a148788a3ac1d649566fc04eb84a746f1a6$$%@^", "TIGER192", false},
		{"46fc0125a148788a3ac1d649566fc04eb84a746f1a6$$%@^", "SOMEHASH", false},
	}
	for i, c := range cs {
		a := IsHash(c.s, c.algo)
		if a != c.w {
			t.Errorf("[%d] IsHash(%q, %q) = %v, want %v", i, c.s, c.algo, a, c.w)
		}
	}
}

func TestIsExistingEmail(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"foo@bar.com", true},
		{"foo@bar.com.au", true},
		{"foo+bar@bar.com", true},
		{"foo@bar.coffee..coffee", false},
		{"invalidemail@", false},
		{"invalid.com", false},
		{"@invalid.com", false},
		{"NathAn.daVIeS@DomaIn.cOM", true},
		{"NATHAN.DAVIES@DOMAIN.CO.UK", true},
		{"prasun.joshi@localhost", false},
		{"[prasun.joshi]@DomaIn.cOM", false},
		{"sizeofuserismorethansixtyfour0123sizeofuserismorethansixtyfour0123@DOMAIN.CO.UK", false},
		{"nosuchdomain@bar.nosuchdomainsuffix", false},
	}
	for i, c := range cs {
		a := IsExistingEmail(c.s)
		if a != c.w {
			t.Errorf("[%d] IsExistingEmail(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsEmail(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"foo@bar.com", true},
		{"x@x.x", false},
		{"foo@bar.com.au", true},
		{"foo+bar@bar.com", true},
		{"foo@bar.coffee", true},
		{"foo@bar.coffee..coffee", false},
		{"foo@bar.bar.coffee", true},
		{"invalidemail@", false},
		{"invalid.com", false},
		{"@invalid.com", false},
		{"NathAn.daVIeS@DomaIn.cOM", true},
		{"NATHAN.DAVIES@DOMAIN.CO.UK", true},
		{"la-mia-donna.tomo.@docomo.ne.jp", true},
		{"y.sawada@ht-meitetsuunyu.co.jp", true},
		{"xukio-zawamura@ja-nesco.com", true},
		// {"foo@bar.中文网", true},
		// {"test|123@m端ller.com", true},
		// {"hans@m端ller.com", true},
		// {"hans.m端ller@c.com", true},
	}
	for i, c := range cs {
		a := IsEmail(c.s)
		if a != c.w {
			t.Errorf("[%d] IsEmail(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsURL(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"http://foo.bar#com", true},
		{"http://foobar.com", true},
		{"https://foobar.com", true},
		{"foobar.com", false},
		{"http://foobar.coffee/", true},
		{"http://foobar.中文网/", true},
		{"http://foobar.org/", true},
		{"http://foobar.ORG", true},
		{"http://foobar.org:8080/", true},
		{"ftp://foobar.ru/", true},
		{"ftp.foo.bar", false},
		{"http://user:pass@www.foobar.com/", true},
		{"http://user:pass@www.foobar.com/path/file", true},
		{"http://127.0.0.1/", true},
		{"http://duckduckgo.com/?q=%2F", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/?foo=bar#baz=qux", true},
		{"http://foobar.com?foo=bar", true},
		{"http://www.xn--froschgrn-x9a.net/", true},
		{"http://foobar.com/a-", true},
		{"http://foobar.پاکستان/", true},
		{"http://foobar.c_o_m", true},
		{"http://_foobar.com", true},
		{"http://foo_bar.com", true},
		{"http://user:pass@foo_bar_bar.bar_foo.com", true},
		{"", false},
		{"xyz://foobar.com", true},
		{".com", false},
		{"rtmp://foobar.com", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com#baz=qux", true},
		{"http://foobar.com/t$-_.+!*\\'(),", true},
		{"http://www.foobar.com/~foobar", true},
		{"http://www.-foobar.com/", true},
		{"http://www.foo---bar.com/", true},
		{"http://r6---snnvoxuioq6.googlevideo.com", true},
		{"mailto:someone@example.com", true},
		{"irc://irc.server.org/channel", true},
		{"irc://#channel@network", true},
		{"/abs/test/dir", false},
		{"./rel/test/dir", false},
		{"http://foo^bar.org", false},
		{"http://foo&*bar.org", true},
		{"http://foo&bar.org", true},
		{"http://foo bar.org", false},
		{"http://foo.bar.org", true},
		{"http://www.foo.bar.org", true},
		{"http://www.foo.co.uk", true},
		{"foo", false},
		{"http://.foo.com", true},
		{"http://,foo.com", true},
		{",foo.com", false},
		{"http://myservice.:9093/", true},
		{"https://pbs.twimg.com/profile_images/560826135676588032/j8fWrmYY_normal.jpeg", true},
		{"http://prometheus-alertmanager.service.q:9093", true},
		{"aio1_alertmanager_container-63376c45:9093", false},
		{"https://www.logn-123-123.url.with.sigle.letter.d:12345/url/path/foo?bar=zzz#user", true},
		{"http://me.example.com", true},
		{"http://www.me.example.com", true},
		{"https://farm6.static.flickr.com", true},
		{"https://zh.wikipedia.org/wiki/Wikipedia:%E9%A6%96%E9%A1%B5", true},
		{"google", false},
		{"http://hyphenated-host-name.example.co.in", true},
		{"http://cant-end-with-hyphen-.example.com", true},
		{"http://-cant-start-with-hyphen.example.com", true},
		{"http://www.domain-can-have-dashes.com", true},
		{"http://m.abcd.com/c.html", true},
		{"http://m.abcd.com/a/b/c/d/c.html?args=a&b=c", true},
		{"http://[::1]:9093", true},
		{"http://[::1]:909388", true},
		{"1200::AB00:1234::2552:7777:1313", false},
		{"http://[2001:db8:a0b:12f0::1]/index.html", true},
		{"http://[1200:0000:AB00:1234:0000:2552:7777:1313]", true},
		{"http://user:pass@[::1]:9093/a/b/c/?a=v#abc", true},
		{"https://127.0.0.1/a/b/c?a=v&c=11d", true},
		{"https://foo_bar.example.com", true},
		{"http://foo_bar.example.com", true},
		{"http://foo_bar_fizz_buzz.example.com", true},
		{"http://_cant_start_with_underescore", true},
		{"http://cant_end_with_underescore_", true},
		{"foo_bar.example.com", false},
		{"foo_bar_fizz_buzz.example.com", false},
		{"http://hello_world.example.com", true},
		{"foo_bar-fizz-buzz:1313", false},
		{"foo_bar-fizz-buzz:13:13", false},
		{"foo_bar-fizz-buzz://1313", false},
	}
	for i, c := range cs {
		a := IsURL(c.s)
		if a != c.w {
			t.Errorf("[%d] IsURL(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsRequestURL(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"http://foo.bar/#com", true},
		{"http://foobar.com", true},
		{"https://foobar.com", true},
		{"foobar.com", false},
		{"http://foobar.coffee/", true},
		{"http://foobar.中文网/", true},
		{"http://foobar.org/", true},
		{"http://foobar.org:8080/", true},
		{"ftp://foobar.ru/", true},
		{"http://user:pass@www.foobar.com/", true},
		{"http://127.0.0.1/", true},
		{"http://duckduckgo.com/?q=%2F", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/?foo=bar#baz=qux", true},
		{"http://foobar.com?foo=bar", true},
		{"http://www.xn--froschgrn-x9a.net/", true},
		{"", false},
		{"xyz://foobar.com", true},
		{"invalid.", false},
		{".com", false},
		{"rtmp://foobar.com", true},
		{"http://www.foo_bar.com/", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/#baz=qux", true},
		{"http://foobar.com/t$-_.+!*\\'(),", true},
		{"http://www.foobar.com/~foobar", true},
		{"http://www.-foobar.com/", true},
		{"http://www.foo---bar.com/", true},
		{"mailto:someone@example.com", true},
		{"irc://irc.server.org/channel", true},
		{"/abs/test/dir", false},
		{"./rel/test/dir", false},
	}
	for i, c := range cs {
		a := IsRequestURL(c.s)
		if a != c.w {
			t.Errorf("[%d] IsRequestURL(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsRequestURI(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"http://foo.bar/#com", true},
		{"http://foobar.com", true},
		{"https://foobar.com", true},
		{"foobar.com", false},
		{"http://foobar.coffee/", true},
		{"http://foobar.中文网/", true},
		{"http://foobar.org/", true},
		{"http://foobar.org:8080/", true},
		{"ftp://foobar.ru/", true},
		{"http://user:pass@www.foobar.com/", true},
		{"http://127.0.0.1/", true},
		{"http://duckduckgo.com/?q=%2F", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/?foo=bar#baz=qux", true},
		{"http://foobar.com?foo=bar", true},
		{"http://www.xn--froschgrn-x9a.net/", true},
		{"xyz://foobar.com", true},
		{"invalid.", false},
		{".com", false},
		{"rtmp://foobar.com", true},
		{"http://www.foo_bar.com/", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/#baz=qux", true},
		{"http://foobar.com/t$-_.+!*\\'(),", true},
		{"http://www.foobar.com/~foobar", true},
		{"http://www.-foobar.com/", true},
		{"http://www.foo---bar.com/", true},
		{"mailto:someone@example.com", true},
		{"irc://irc.server.org/channel", true},
		{"/abs/test/dir", true},
		{"./rel/test/dir", false},
	}
	for i, c := range cs {
		a := IsRequestURI(c.s)
		if a != c.w {
			t.Errorf("[%d] IsRequestURI(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsHexColor(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"#ff", false},
		{"fff0", false},
		{"#ff12FG", false},
		{"#CCccCC", true},
		{"#f00", true},
		{"fff", false},
	}
	for i, c := range cs {
		a := IsHexColor(c.s)
		if a != c.w {
			t.Errorf("[%d] IsHexColor(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsRGBcolor(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"rgb(0,31,255)", true},
		{"rgb(1,349,275)", false},
		{"rgb(01,31,255)", false},
		{"rgb(0.6,31,255)", false},
		{"rgba(0,31,255)", false},
		{"rgb(0,  31, 255)", true},
	}
	for i, c := range cs {
		a := IsRGBColor(c.s)
		if a != c.w {
			t.Errorf("[%d] IsRGBcolor(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsIMEI(t *testing.T) {
	cs := []struct {
		s string
		w bool
	}{
		{"990000862471854", true},
		{"351756051523999", true},
		{"9900008624718541", false},
		{"1", false},
	}
	for i, c := range cs {
		a := IsIMEI(c.s)
		if a != c.w {
			t.Errorf("[%d] IsIMEI(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsByteLength(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s   string
		min int
		max int
		w   bool
	}{
		{"abacaba", 100, -1, false},
		{"abacaba", 1, 3, false},
		{"abacaba", 1, 7, true},
		{"abacaba", 0, 8, true},
		{"\ufff0", 1, 1, false},
	}
	for i, c := range cs {
		a := IsByteLength(c.s, c.min, c.max)
		if a != c.w {
			t.Errorf("[%d] IsByteLength(%q, %q, %q) = %v, want %v", i, c.s, c.min, c.max, a, c.w)
		}
	}
}

func TestIsJSON(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"145", true},
		{"asdf", false},
		{"123:f00", false},
		{"{\"Name\":\"Alice\",\"Body\":\"Hello\",\"Time\":1294706395881547000}", true},
		{"{}", true},
		{"{\"Key\":{\"Key\":{\"Key\":123}}}", true},
		{"[]", true},
		{"null", true},
	}
	for i, c := range cs {
		a := IsJSON(c.s)
		if a != c.w {
			t.Errorf("[%d] IsJSON(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsUUID(t *testing.T) {
	t.Parallel()

	// Tests without version
	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3xxx", false},
		{"a987fbc94bed3078cf079141ba07c9f3", false},
		{"934859", false},
		{"987fbc9-4bed-3078-cf07a-9141ba07c9f3", false},
		{"aaaaaaaa-1111-1111-aaag-111111111111", false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3", true},
	}
	for i, c := range cs {
		a := IsUUID(c.s)
		if a != c.w {
			t.Errorf("[%d] IsUUID(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}

	// UUID ver. 3
	cs = []struct {
		s string
		w bool
	}{
		{"", false},
		{"412452646", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"a987fbc9-4bed-4078-8f07-9141ba07c9f3", false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3", true},
	}
	for i, c := range cs {
		a := IsUUIDv3(c.s)
		if a != c.w {
			t.Errorf("[%d] IsUUIDv3(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}

	// UUID ver. 4
	cs = []struct {
		s string
		w bool
	}{
		{"", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"a987fbc9-4bed-5078-af07-9141ba07c9f3", false},
		{"934859", false},
		{"57b73598-8764-4ad0-a76a-679bb6640eb1", true},
		{"625e63f3-58f5-40b7-83a1-a72ad31acffb", true},
	}
	for i, c := range cs {
		a := IsUUIDv4(c.s)
		if a != c.w {
			t.Errorf("[%d] IsUUIDv4(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}

	// UUID ver. 5
	cs = []struct {
		s string
		w bool
	}{

		{"", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"9c858901-8a57-4791-81fe-4c455b099bc9", false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"987fbc97-4bed-5078-af07-9141ba07c9f3", true},
		{"987fbc97-4bed-5078-9f07-9141ba07c9f3", true},
	}
	for i, c := range cs {
		a := IsUUIDv5(c.s)
		if a != c.w {
			t.Errorf("[%d] IsUUIDv5(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsCreditCard(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{"empty", "", false},
		{"not numbers", "credit card", false},
		{"invalid luhn algorithm", "4220855426213389", false},

		{"visa", "4220855426222389", true},
		{"visa spaces", "4220 8554 2622 2389", true},
		{"visa dashes", "4220-8554-2622-2389", true},
		{"mastercard", "5139288802098206", true},
		{"american express", "374953669708156", true},
		{"discover", "6011464355444102", true},
		{"jcb", "3548209662790989", true},

		// below should be valid, do they respect international standards?
		// is our validator logic not correct?
		{"diners club international", "30190239451016", false},
		{"rupay", "6521674451993089", false},
		{"mir", "2204151414444676", false},
		{"china unionPay", "624356436327468104", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCreditCard(tt.number); got != tt.want {
				t.Errorf("IsCreditCard(%v) = %v, want %v", tt.number, got, tt.want)
			}
		})
	}
}

func TestIsISBN(t *testing.T) {
	t.Parallel()

	// Without version
	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"foo", false},
		{"3836221195", true},
		{"1-61729-085-8", true},
		{"3 423 21412 0", true},
		{"3 401 01319 X", true},
		{"9784873113685", true},
		{"978-4-87311-368-5", true},
		{"978 3401013190", true},
		{"978-3-8362-2119-1", true},
	}
	for i, c := range cs {
		a := IsISBN(c.s, -1)
		if a != c.w {
			t.Errorf("[%d] IsISBN(%q, -1) = %v, want %v", i, c.s, a, c.w)
		}
	}

	// ISBN 10
	cs = []struct {
		s string
		w bool
	}{
		{"", false},
		{"foo", false},
		{"3423214121", false},
		{"978-3836221191", false},
		{"3-423-21412-1", false},
		{"3 423 21412 1", false},
		{"3836221195", true},
		{"1-61729-085-8", true},
		{"3 423 21412 0", true},
		{"3 401 01319 X", true},
	}
	for i, c := range cs {
		a := IsISBN10(c.s)
		if a != c.w {
			t.Errorf("[%d] IsISBN10(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}

	// ISBN 13
	cs = []struct {
		s string
		w bool
	}{
		{"", false},
		{"foo", false},
		{"3-8362-2119-5", false},
		{"01234567890ab", false},
		{"978 3 8362 2119 0", false},
		{"9784873113685", true},
		{"978-4-87311-368-5", true},
		{"978 3401013190", true},
		{"978-3-8362-2119-1", true},
	}
	for i, c := range cs {
		a := IsISBN13(c.s)
		if a != c.w {
			t.Errorf("[%d] IsISBN13(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsDataURI(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"data:image/png;base64,TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4=", true},
		{"data:text/plain;base64,Vml2YW11cyBmZXJtZW50dW0gc2VtcGVyIHBvcnRhLg==", true},
		{"image/gif;base64,U3VzcGVuZGlzc2UgbGVjdHVzIGxlbw==", false},
		{"data:image/gif;base64,MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuMPNS1Ufof9EW/M98FNw" +
			"UAKrwflsqVxaxQjBQnHQmiI7Vac40t8x7pIb8gLGV6wL7sBTJiPovJ0V7y7oc0Ye" +
			"rhKh0Rm4skP2z/jHwwZICgGzBvA0rH8xlhUiTvcwDCJ0kc+fh35hNt8srZQM4619" +
			"FTgB66Xmp4EtVyhpQV+t02g6NzK72oZI0vnAvqhpkxLeLiMCyrI416wHm5Tkukhx" +
			"QmcL2a6hNOyu0ixX/x2kSFXApEnVrJ+/IxGyfyw8kf4N2IZpW5nEP847lpfj0SZZ" +
			"Fwrd1mnfnDbYohX2zRptLy2ZUn06Qo9pkG5ntvFEPo9bfZeULtjYzIl6K8gJ2uGZ" + "HQIDAQAB", true},
		{"data:image/png;base64,12345", false},
		{"", false},
		// {"data:text,:;base85,U3VzcGVuZGlzc2UgbGVjdHVzIGxlbw==", false},
	}
	for i, c := range cs {
		a := IsDataURI(c.s)
		if a != c.w {
			t.Errorf("[%d] IsDataURI(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsMagnetURI(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"magnet:?xt=urn:btih:06E2A9683BF4DA92C73A661AC56F0ECC9C63C5B4&dn=helloword2000&tr=udp://helloworld:1337/announce", true},
		{"magnet:?xt=urn:btih:3E30322D5BFC7444B7B1D8DD42404B75D0531DFB&dn=world&tr=udp://world.com:1337", true},
		{"magnet:?xt=urn:btih:4ODKSDJBVMSDSNJVBCBFYFBKNRU875DW8D97DWC6&dn=helloworld&tr=udp://helloworld.com:1337", true},
		{"magnet:?xt=urn:btih:1GSHJVBDVDVJFYEHKFHEFIO8573898434JBFEGHD&dn=foo&tr=udp://foo.com:1337", true},
		{"magnet:?xt=urn:btih:MCJDCYUFHEUD6E2752T7UJNEKHSUGEJFGTFHVBJS&dn=bar&tr=udp://bar.com:1337", true},
		{"magnet:?xt=urn:btih:LAKDHWDHEBFRFVUFJENBYYTEUY837562JH2GEFYH&dn=foobar&tr=udp://foobar.com:1337", true},
		{"magnet:?xt=urn:btih:MKCJBHCBJDCU725TGEB3Y6RE8EJ2U267UNJFGUID&dn=test&tr=udp://c.com:1337", true},
		{"magnet:?xt=urn:btih:UHWY2892JNEJ2GTEYOMDNU67E8ICGICYE92JDUGH&dn=baz&tr=udp://baz.com:1337", true},
		{"magnet:?xt=urn:btih:HS263FG8U3GFIDHWD7829BYFCIXB78XIHG7CWCUG&dn=foz&tr=udp://foz.com:1337", true},
		{"", false},
		{":?xt=urn:btih:06E2A9683BF4DA92C73A661AC56F0ECC9C63C5B4&dn=helloword2000&tr=udp://helloworld:1337/announce", false},
		{"magnett:?xt=urn:btih:3E30322D5BFC7444B7B1D8DD42404B75D0531DFB&dn=world&tr=udp://world.com:1337", false},
		{"xt=urn:btih:4ODKSDJBVMSDSNJVBCBFYFBKNRU875DW8D97DWC6&dn=helloworld&tr=udp://helloworld.com:1337", false},
		{"magneta:?xt=urn:btih:1GSHJVBDVDVJFYEHKFHEFIO8573898434JBFEGHD&dn=foo&tr=udp://foo.com:1337", false},
		{"magnet:?xt=uarn:btih:MCJDCYUFHEUD6E2752T7UJNEKHSUGEJFGTFHVBJS&dn=bar&tr=udp://bar.com:1337", false},
		{"magnet:?xt=urn:btihz&dn=foobar&tr=udp://foobar.com:1337", false},
		{"magnet:?xat=urn:btih:MKCJBHCBJDCU725TGEB3Y6RE8EJ2U267UNJFGUID&dn=test&tr=udp://c.com:1337", false},
		{"magnet::?xt=urn:btih:UHWY2892JNEJ2GTEYOMDNU67E8ICGICYE92JDUGH&dn=baz&tr=udp://baz.com:1337", false},
		{"magnet:?xt:btih:HS263FG8U3GFIDHWD7829BYFCIXB78XIHG7CWCUG&dn=foz&tr=udp://foz.com:1337", false},
	}
	for i, c := range cs {
		a := IsMagnetURI(c.s)
		if a != c.w {
			t.Errorf("[%d] IsMagnetURI(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsBase64(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4=", true},
		{"Vml2YW11cyBmZXJtZW50dW0gc2VtcGVyIHBvcnRhLg==", true},
		{"U3VzcGVuZGlzc2UgbGVjdHVzIGxlbw==", true},
		{"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuMPNS1Ufof9EW/M98FNw" +
			"UAKrwflsqVxaxQjBQnHQmiI7Vac40t8x7pIb8gLGV6wL7sBTJiPovJ0V7y7oc0Ye" +
			"rhKh0Rm4skP2z/jHwwZICgGzBvA0rH8xlhUiTvcwDCJ0kc+fh35hNt8srZQM4619" +
			"FTgB66Xmp4EtVyhpQV+t02g6NzK72oZI0vnAvqhpkxLeLiMCyrI416wHm5Tkukhx" +
			"QmcL2a6hNOyu0ixX/x2kSFXApEnVrJ+/IxGyfyw8kf4N2IZpW5nEP847lpfj0SZZ" +
			"Fwrd1mnfnDbYohX2zRptLy2ZUn06Qo9pkG5ntvFEPo9bfZeULtjYzIl6K8gJ2uGZ" + "HQIDAQAB", true},
		{"12345", false},
		{"", false},
		{"Vml2YW11cyBmZXJtZtesting123", false},
	}
	for i, c := range cs {
		a := IsBase64(c.s)
		if a != c.w {
			t.Errorf("[%d] IsBase64(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsIP(t *testing.T) {
	t.Parallel()

	// Without version
	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"127.0.0.1", true},
		{"0.0.0.0", true},
		{"255.255.255.255", true},
		{"1.2.3.4", true},
		{"::1", true},
		{"2001:db8:0000:1:1:1:1:1", true},
		{"300.0.0.0", false},
	}
	for i, c := range cs {
		a := IsIP(c.s)
		if a != c.w {
			t.Errorf("[%d] IsIP(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}

	// IPv4
	cs = []struct {
		s string
		w bool
	}{
		{"", false},
		{"127.0.0.1", true},
		{"0.0.0.0", true},
		{"255.255.255.255", true},
		{"1.2.3.4", true},
		{"::1", false},
		{"2001:db8:0000:1:1:1:1:1", false},
		{"300.0.0.0", false},
	}
	for i, c := range cs {
		a := IsIPv4(c.s)
		if a != c.w {
			t.Errorf("[%d] IsIPv4(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}

	// IPv6
	cs = []struct {
		s string
		w bool
	}{
		{"", false},
		{"127.0.0.1", false},
		{"0.0.0.0", false},
		{"255.255.255.255", false},
		{"1.2.3.4", false},
		{"::1", true},
		{"2001:db8:0000:1:1:1:1:1", true},
		{"300.0.0.0", false},
	}
	for i, c := range cs {
		a := IsIPv6(c.s)
		if a != c.w {
			t.Errorf("[%d] IsIPv6(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsPort(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"1", true},
		{"65535", true},
		{"0", false},
		{"65536", false},
		{"65538", false},
	}

	for i, c := range cs {
		a := IsPort(c.s)
		if a != c.w {
			t.Errorf("[%d] IsPort(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsDNSName(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"localhost", true},
		{"a.bc", true},
		{"a.b.", true},
		{"a.b..", false},
		{"localhost.local", true},
		{"localhost.localdomain.intern", true},
		{"l.local.intern", true},
		{"ru.link.n.svpncloud.com", true},
		{"-localhost", false},
		{"localhost.-localdomain", false},
		{"localhost.localdomain.-int", false},
		{"_localhost", true},
		{"localhost._localdomain", true},
		{"localhost.localdomain._int", true},
		{"lÖcalhost", false},
		{"localhost.lÖcaldomain", false},
		{"localhost.localdomain.üntern", false},
		{"__", true},
		{"localhost/", false},
		{"127.0.0.1", false},
		{"[::1]", false},
		{"50.50.50.50", false},
		{"localhost.localdomain.intern:65535", false},
		{"漢字汉字", false},
		{"www.jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6906k846pj3sulm4kiyk82ln5teqj9nsht59opr0cs5ssltx78lfyvml19lfq1wp4usbl0o36cmiykch1vywbttcus1p9yu0669h8fj4ll7a6bmop505908s1m83q2ec2qr9nbvql2589adma3xsq2o38os2z3dmfh2tth4is4ixyfasasasefqwe4t2ub2fz1rme.de", false},
	}

	for i, c := range cs {
		a := IsDNSName(c.s)
		if a != c.w {
			t.Errorf("[%d] IsDNS(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsHost(t *testing.T) {
	t.Parallel()
	cs := []struct {
		s string
		w bool
	}{
		{"localhost", true},
		{"localhost.localdomain", true},
		{"2001:db8:0000:1:1:1:1:1", true},
		{"::1", true},
		{"play.golang.org", true},
		{"localhost.localdomain.intern:65535", false},
		{"-[::1]", false},
		{"-localhost", false},
		{".localhost", false},
	}
	for i, c := range cs {
		a := IsHost(c.s)
		if a != c.w {
			t.Errorf("[%d] IsHost(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}

}

func TestIsDialString(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"localhost.local:1", true},
		{"localhost.localdomain:9090", true},
		{"localhost.localdomain.intern:65535", true},
		{"127.0.0.1:30000", true},
		{"[::1]:80", true},
		{"[1200::AB00:1234::2552:7777:1313]:22", false},
		{"-localhost:1", false},
		{"localhost.-localdomain:9090", false},
		{"localhost.localdomain.-int:65535", false},
		{"localhost.loc:100000", false},
		{"漢字汉字:2", false},
		{"www.jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6906k846pj3sulm4kiyk82ln5teqj9nsht59opr0cs5ssltx78lfyvml19lfq1wp4usbl0o36cmiykch1vywbttcus1p9yu0669h8fj4ll7a6bmop505908s1m83q2ec2qr9nbvql2589adma3xsq2o38os2z3dmfh2tth4is4ixyfasasasefqwe4t2ub2fz1rme.de:20000", false},
	}

	for i, c := range cs {
		a := IsDialString(c.s)
		if a != c.w {
			t.Errorf("[%d] IsDialString(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsMAC(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"3D:F2:C9:A6:B3:4F", true},
		{"3D-F2-C9-A6-B3:4F", false},
		{"123", false},
		{"", false},
		{"abacaba", false},
	}
	for i, c := range cs {
		a := IsMAC(c.s)
		if a != c.w {
			t.Errorf("[%d] IsMAC(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsLatitude(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"-90.000", true},
		{"+90", true},
		{"47.1231231", true},
		{"+99.9", false},
		{"108", false},
	}
	for i, c := range cs {
		a := IsLatitude(c.s)
		if a != c.w {
			t.Errorf("[%d] IsLatitude(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}

func TestIsLongitude(t *testing.T) {
	t.Parallel()

	cs := []struct {
		s string
		w bool
	}{
		{"", false},
		{"-180.000", true},
		{"180.1", false},
		{"+73.234", true},
		{"+382.3811", false},
		{"23.11111111", true},
	}
	for i, c := range cs {
		a := IsLongitude(c.s)
		if a != c.w {
			t.Errorf("[%d] IsLongitude(%q) = %v, want %v", i, c.s, a, c.w)
		}
	}
}
