package vad

import (
	"crypto/sha256"
	"encoding/json"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/askasoft/pango/str"
)

// IsFileName is illegal file name
func IsFileName(s string) bool {
	return !str.ContainsAny(s, `\/:*?"<>|`)
}

// IsE164 is the validation function for validating if the current field's value is a valid e.164 formatted phone number.
func IsE164(s string) bool {
	return rxE164.MatchString(s)
}

// IsEmail checks if the string is an email.
func IsEmail(s string) bool {
	return rxEmail.MatchString(s)
}

// IsExistingEmail checks if the string is an email of existing domain
func IsExistingEmail(email string) bool {
	if len(email) < 6 || len(email) > 254 {
		return false
	}
	at := str.LastIndex(email, "@")
	if at <= 0 || at > len(email)-3 {
		return false
	}
	user := email[:at]
	host := email[at+1:]
	if len(user) > 64 {
		return false
	}

	if rxUserDot.MatchString(user) || !rxUser.MatchString(user) || !rxHost.MatchString(host) {
		return false
	}
	if _, err := net.LookupMX(host); err != nil {
		if _, err := net.LookupIP(host); err != nil {
			return false
		}
	}

	return true
}

// IsURI checks if the string is an URI.
func IsURI(s string) bool {
	u, _, _ := str.CutByte(s, '#')

	if len(u) == 0 {
		return false
	}

	_, err := url.ParseRequestURI(u)
	return err == nil
}

// IsURL checks if the string is an URL.
func IsURL(s string) bool {
	u, _, _ := str.CutByte(s, '#')
	if len(u) == 0 {
		return false
	}

	url, err := url.ParseRequestURI(u)
	if err != nil || url.Scheme == "" {
		return false
	}

	return true
}

// IsURLWithScheme checks if the string is an URL with specified shemes.
func IsURLWithScheme(s string, schemes ...string) bool {
	u, _, _ := str.CutByte(s, '#')
	if len(u) == 0 {
		return false
	}

	url, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}

	for _, scheme := range schemes {
		if str.EqualFold(scheme, url.Scheme) {
			return true
		}
	}

	return false
}

// IsHttpURL checks if the string is an http:// URL.
func IsHttpURL(s string) bool {
	return IsURLWithScheme(s, "http")
}

// IsHttpsURL checks if the string is an https:// URL.
func IsHttpsURL(s string) bool {
	return IsURLWithScheme(s, "https")
}

// IsHttpxURL checks if the string is an https?:// URL.
func IsHttpxURL(s string) bool {
	return IsURLWithScheme(s, "http", "https")
}

// IsRequestURL checks if the string rawurl, assuming
// it was received in an HTTP request, is a valid
// URL confirm to RFC 3986
func IsRequestURL(rawurl string) bool {
	url, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return false //Couldn't even parse the rawurl
	}
	if len(url.Scheme) == 0 {
		return false //No Scheme found
	}
	return true
}

// IsRequestURI checks if the string rawurl, assuming
// it was received in an HTTP request, is an
// absolute URI or an absolute path.
func IsRequestURI(rawurl string) bool {
	_, err := url.ParseRequestURI(rawurl)
	return err == nil
}

// IsHexColor checks if the string is a hex decimal color.
func IsHexColor(s string) bool {
	if s == "" || s[0] != '#' {
		return false
	}

	s = s[1:]
	switch len(s) {
	case 3, 6:
		return str.IsHexadecimal(s)
	default:
		return false
	}
}

// IsRGBColor checks if the string is a valid RGB color in form rgb(RRR, GGG, BBB).
func IsRGBColor(s string) bool {
	return rxRGBColor.MatchString(s)
}

// IsRGBAColor checks if the string is a valid RGBA color in form rgb(RRR, GGG, BBB).
func IsRGBAColor(s string) bool {
	return rxRGBAColor.MatchString(s)
}

// IsHSLColor  checks if the string is a valid HSLA color in form hsl(0, 100%, 50%).
func IsHSLColor(s string) bool {
	return rxHSLColor.MatchString(s)
}

// IsHSLAColor  checks if the string is a valid HSLA color in form hsla(0, 100%, 50%, 0.5).
func IsHSLAColor(s string) bool {
	return rxHSLAColor.MatchString(s)
}

// IsByteLength checks if the string's length (in bytes) falls in a range.
func IsByteLength(s string, min, max int) bool {
	return len(s) >= min && len(s) <= max
}

// IsUUIDv3 checks if the string is a UUID version 3.
func IsUUIDv3(s string) bool {
	return rxUUID3.MatchString(s)
}

// IsUUIDv4 checks if the string is a UUID version 4.
func IsUUIDv4(s string) bool {
	return rxUUID4.MatchString(s)
}

// IsUUIDv5 checks if the string is a UUID version 5.
func IsUUIDv5(s string) bool {
	return rxUUID5.MatchString(s)
}

// IsUUID checks if the string is a UUID (version 3, 4 or 5).
func IsUUID(s string) bool {
	return rxUUID.MatchString(s)
}

// IsULID is the validation function for validating if the field's value is a valid ULID.
func IsULID(s string) bool {
	return rxULID.MatchString(s)
}

// IsCreditCard checks if the string is a credit card.
func IsCreditCard(s string) bool {
	sanitized := str.RemoveAny(s, " -")
	if !rxCreditCard.MatchString(sanitized) {
		return false
	}

	var sum int64
	var digit string
	var tmpNum int64
	var shouldDouble bool
	for i := len(sanitized) - 1; i >= 0; i-- {
		digit = sanitized[i:(i + 1)]
		tmpNum, _ = strconv.ParseInt(digit, 0, 64)
		if shouldDouble {
			tmpNum *= 2
			if tmpNum >= 10 {
				sum += (tmpNum % 10) + 1
			} else {
				sum += tmpNum
			}
		} else {
			sum += tmpNum
		}
		shouldDouble = !shouldDouble
	}

	return sum%10 == 0
}

// IsISBN10 checks if the string is an ISBN version 10.
func IsISBN10(s string) bool {
	return IsISBN(s, 10)
}

// IsISBN13 checks if the string is an ISBN version 13.
func IsISBN13(s string) bool {
	return IsISBN(s, 13)
}

// IsISBN checks if the string is an ISBN (version 10 or 13).
// If version value is not equal to 10 or 13, it will be checks both variants.
func IsISBN(s string, version int) bool {
	sanitized := str.RemoveAny(s, " -")
	var checksum int32
	var i int32

	switch version {
	case 10:
		if !rxISBN10.MatchString(sanitized) {
			return false
		}
		for i = 0; i < 9; i++ {
			checksum += (i + 1) * int32(sanitized[i]-'0')
		}
		if sanitized[9] == 'X' {
			checksum += 10 * 10
		} else {
			checksum += 10 * int32(sanitized[9]-'0')
		}
		if checksum%11 == 0 {
			return true
		}
		return false
	case 13:
		if !rxISBN13.MatchString(sanitized) {
			return false
		}
		factor := []int32{1, 3}
		for i = 0; i < 12; i++ {
			checksum += factor[i%2] * int32(sanitized[i]-'0')
		}
		return (int32(sanitized[12]-'0'))-((10-(checksum%10))%10) == 0
	default:
		return IsISBN(s, 10) || IsISBN(s, 13)
	}
}

// IsJSON checks if the string is valid JSON (note: uses json.Unmarshal).
func IsJSON(s string) bool {
	return json.Valid(str.UnsafeBytes(s))
}

// IsBase64 checks if a string is base64 encoded.
func IsBase64(s string) bool {
	return rxBase64.MatchString(s)
}

// IsBase64URL checks if a string is base64 url safe encoded.
func IsBase64URL(s string) bool {
	return rxBase64URL.MatchString(s)
}

// IsDataURI checks if a string is base64 encoded data URI such as an image
func IsDataURI(s string) bool {
	if !str.StartsWith(s, "data:") {
		return false
	}

	kind, data, ok := str.Cut(s, ",")
	if !ok {
		return false
	}
	if str.EndsWith(kind, ";base64") {
		return IsBase64(data)
	}
	return true
}

// IsMagnetURI checks if a string is valid magnet URI
func IsMagnetURI(s string) bool {
	return rxMagnetURI.MatchString(s)
}

// IsDNSName will validate the given string as a DNS name
func IsDNSName(s string) bool {
	if s == "" || len(s)-str.CountByte(s, '.') > 255 {
		// constraints already violated
		return false
	}
	return !IsIP(s) && rxDNSName.MatchString(s)
}

// IsDnsRFC1035Label is the validation function
// for validating if the current field's value is
// a valid dns RFC 1035 label, defined in RFC 1035.
func IsDnsRFC1035Label(s string) bool {
	return rxDnsRFC1035Label.MatchString(s)
}

func IsHostnameRFC952(s string) bool {
	return rxHostnameRFC952.MatchString(s)
}

func IsHostnameRFC1123(s string) bool {
	return rxHostnameRFC1123.MatchString(s)
}

func IsHostnamePort(s string) bool {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return false
	}

	// Port must be a iny <= 65535.
	if portNum, err := strconv.ParseInt(port, 10, 32); err != nil || portNum > 65535 || portNum < 1 {
		return false
	}

	// If host is specified, it should match a DNS name
	if host != "" {
		return IsHostnameRFC1123(host)
	}
	return true
}

func IsFQDN(s string) bool {
	return rxFqdnRFC1123.MatchString(s)
}

// IsHash checks if a string is a hash of type algorithm.
// Algorithm is one of ['md4', 'md5', 'sha1', 'sha256', 'sha384', 'sha512', 'ripemd128', 'ripemd160', 'tiger128', 'tiger160', 'tiger192', 'crc32', 'crc32b']
func IsHash(s string, algorithm string) bool {
	var n int

	algo := str.ToLower(algorithm)
	switch algo {
	case "crc32", "crc32b":
		n = 8
	case "md5", "md4", "ripemd128", "tiger128":
		n = 32
	case "sha1", "ripemd160", "tiger160":
		n = 40
	case "tiger192":
		n = 48
	case "sha256":
		n = 64
	case "sha384":
		n = 96
	case "sha512":
		n = 128
	default:
		return false
	}

	if len(s) != n {
		return false
	}
	return str.IsLetterNumber(s)
}

// IsSHA512 checks is a string is a SHA512 hash. Alias for `IsHash(s, "sha512")`
func IsSHA512(s string) bool {
	return IsHash(s, "sha512")
}

// IsSHA384 checks is a string is a SHA384 hash. Alias for `IsHash(s, "sha384")`
func IsSHA384(s string) bool {
	return IsHash(s, "sha384")
}

// IsSHA256 checks is a string is a SHA256 hash. Alias for `IsHash(s, "sha256")`
func IsSHA256(s string) bool {
	return IsHash(s, "sha256")
}

// IsTiger192 checks is a string is a Tiger192 hash. Alias for `IsHash(s, "tiger192")`
func IsTiger192(s string) bool {
	return IsHash(s, "tiger192")
}

// IsTiger160 checks is a string is a Tiger160 hash. Alias for `IsHash(s, "tiger160")`
func IsTiger160(s string) bool {
	return IsHash(s, "tiger160")
}

// IsRipeMD160 checks is a string is a RipeMD160 hash. Alias for `IsHash(s, "ripemd160")`
func IsRipeMD160(s string) bool {
	return IsHash(s, "ripemd160")
}

// IsSHA1 checks is a string is a SHA-1 hash. Alias for `IsHash(s, "sha1")`
func IsSHA1(s string) bool {
	return IsHash(s, "sha1")
}

// IsTiger128 checks is a string is a Tiger128 hash. Alias for `IsHash(s, "tiger128")`
func IsTiger128(s string) bool {
	return IsHash(s, "tiger128")
}

// IsRipeMD128 checks is a string is a RipeMD128 hash. Alias for `IsHash(s, "ripemd128")`
func IsRipeMD128(s string) bool {
	return IsHash(s, "ripemd128")
}

// IsCRC32 checks is a string is a CRC32 hash. Alias for `IsHash(s, "crc32")`
func IsCRC32(s string) bool {
	return IsHash(s, "crc32")
}

// IsCRC32b checks is a string is a CRC32b hash. Alias for `IsHash(s, "crc32b")`
func IsCRC32b(s string) bool {
	return IsHash(s, "crc32b")
}

// IsMD5 checks is a string is a MD5 hash. Alias for `IsHash(s, "md5")`
func IsMD5(s string) bool {
	return IsHash(s, "md5")
}

// IsMD4 checks is a string is a MD4 hash. Alias for `IsHash(s, "md4")`
func IsMD4(s string) bool {
	return IsHash(s, "md4")
}

// IsDialString validates the given string for usage with the various Dial() functions
func IsDialString(s string) bool {
	if h, p, err := net.SplitHostPort(s); err == nil && h != "" && p != "" && (IsDNSName(h) || IsIP(h)) && IsPort(p) {
		return true
	}

	return false
}

// IsPort checks if a string represents a valid port
func IsPort(s string) bool {
	if i, err := strconv.Atoi(s); err == nil && i > 0 && i < 65536 {
		return true
	}
	return false
}

// IsCIDRv4 checks if a string is a valid v4 CIDR address.
func IsCIDRv4(s string) bool {
	ip, _, err := net.ParseCIDR(s)
	return err == nil && ip.To4() != nil
}

// IsCIDRv6 checks if a string is a valid v6 CIDR address.
func IsCIDRv6(s string) bool {
	ip, _, err := net.ParseCIDR(s)
	return err == nil && ip.To4() == nil
}

// IsCIDR checks if a string is a valid v4 or v6 CIDR address.
func IsCIDR(s string) bool {
	_, _, err := net.ParseCIDR(s)
	return err == nil
}

// IsIP checks if a string is either IP version 4 or 6. Alias for `net.ParseIP`
func IsIP(s string) bool {
	return net.ParseIP(s) != nil
}

// IsIPv4 checks if the string is an IP version 4.
func IsIPv4(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && ip.To4() != nil
}

// IsIPv6 checks if the string is an IP version 6.
func IsIPv6(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && ip.To4() == nil
}

// IsIP4AddrResolvable checks if the string is a resolvable ip v4 address.
func IsIP4AddrResolvable(s string) bool {
	if !IsIPv4(s) {
		return false
	}

	_, err := net.ResolveIPAddr("ip4", s)
	return err == nil
}

// IsIP6AddrResolvable checks if the string is a resolvable ip v6 address.
func IsIP6AddrResolvable(s string) bool {
	if !IsIPv6(s) {
		return false
	}

	_, err := net.ResolveIPAddr("ip6", s)
	return err == nil
}

// IsIPAddrResolvable checks if the string is a resolvable ip address.
func IsIPAddrResolvable(s string) bool {
	if !IsIP(s) {
		return false
	}

	_, err := net.ResolveIPAddr("ip", s)
	return err == nil
}

// IsMAC checks if a string is valid MAC address.
// Possible MAC formats:
// 01:23:45:67:89:ab
// 01:23:45:67:89:ab:cd:ef
// 01-23-45-67-89-ab
// 01-23-45-67-89-ab-cd-ef
// 0123.4567.89ab
// 0123.4567.89ab.cdef
func IsMAC(s string) bool {
	_, err := net.ParseMAC(s)
	return err == nil
}

// IsHost checks if the string is a valid IP (both v4 and v6) or a valid DNS name
func IsHost(s string) bool {
	return IsIP(s) || IsDNSName(s)
}

// IsLatitude checks if a string is valid latitude.
func IsLatitude(s string) bool {
	return rxLatitude.MatchString(s)
}

// IsLongitude checks if a string is valid longitude.
func IsLongitude(s string) bool {
	return rxLongitude.MatchString(s)
}

// IsJWT checks if the string is a valid JWT string.
func IsJWT(s string) bool {
	return rxJWT.MatchString(s)
}

// IsSSN checks if the string is a valid SSN string.
func IsSSN(s string) bool {
	return rxSSN.MatchString(s)
}

// IsSwiftCode checks if the string is a valid Business Identifier Code (SWIFT code), defined in ISO 9362
func IsSwiftCode(s string) bool {
	return rxSwiftCode.MatchString(s)
}

// IsIMEI checks if a string is valid IMEI
func IsIMEI(s string) bool {
	return rxIMEI.MatchString(s)
}

// IsIMSI checks if a string is valid IMSI
func IsIMSI(s string) bool {
	if !rxIMSI.MatchString(s) {
		return false
	}

	mcc, err := strconv.ParseInt(s[0:3], 10, 32)
	if err != nil {
		return false
	}

	switch mcc {
	case 202, 204, 206, 208, 212, 213, 214, 216, 218, 219:
	case 220, 221, 222, 226, 228, 230, 231, 232, 234, 235:
	case 238, 240, 242, 244, 246, 247, 248, 250, 255, 257:
	case 259, 260, 262, 266, 268, 270, 272, 274, 276, 278:
	case 280, 282, 283, 284, 286, 288, 289, 290, 292, 293:
	case 294, 295, 297, 302, 308, 310, 311, 312, 313, 314:
	case 315, 316, 330, 332, 334, 338, 340, 342, 344, 346:
	case 348, 350, 352, 354, 356, 358, 360, 362, 363, 364:
	case 365, 366, 368, 370, 372, 374, 376, 400, 401, 402:
	case 404, 405, 406, 410, 412, 413, 414, 415, 416, 417:
	case 418, 419, 420, 421, 422, 424, 425, 426, 427, 428:
	case 429, 430, 431, 432, 434, 436, 437, 438, 440, 441:
	case 450, 452, 454, 455, 456, 457, 460, 461, 466, 467:
	case 470, 472, 502, 505, 510, 514, 515, 520, 525, 528:
	case 530, 536, 537, 539, 540, 541, 542, 543, 544, 545:
	case 546, 547, 548, 549, 550, 551, 552, 553, 554, 555:
	case 602, 603, 604, 605, 606, 607, 608, 609, 610, 611:
	case 612, 613, 614, 615, 616, 617, 618, 619, 620, 621:
	case 622, 623, 624, 625, 626, 627, 628, 629, 630, 631:
	case 632, 633, 634, 635, 636, 637, 638, 639, 640, 641:
	case 642, 643, 645, 646, 647, 648, 649, 650, 651, 652:
	case 653, 654, 655, 657, 658, 659, 702, 704, 706, 708:
	case 710, 712, 714, 716, 722, 724, 730, 732, 734, 736:
	case 738, 740, 742, 744, 746, 748, 750, 995:
		return true
	default:
		return false
	}
	return true
}

func IsJSONObject(s string) bool {
	if !json.Valid(str.UnsafeBytes(s)) {
		return false
	}
	return str.StartsWithByte(str.Strip(s), '{')
}

func IsJSONArray(s string) bool {
	if !json.Valid(str.UnsafeBytes(s)) {
		return false
	}
	return str.StartsWithByte(str.Strip(s), '[')
}

// IsBitcoinAddress is the validation function for validating if the field's value is a valid btc address
func IsBitcoinAddress(address string) bool {
	if !rxBtcAddress.MatchString(address) {
		return false
	}

	alphabet := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	decode := [25]byte{}
	for i := 0; i < len(address); i++ {
		d := str.IndexByte(alphabet, address[i])
		for i := 24; i >= 0; i-- {
			d += 58 * int(decode[i])
			decode[i] = byte(d % 256)
			d /= 256
		}
	}

	h := sha256.New()
	_, _ = h.Write(decode[:21])
	d := h.Sum([]byte{})
	h = sha256.New()
	_, _ = h.Write(d)

	validchecksum := [4]byte{}
	computedchecksum := [4]byte{}

	copy(computedchecksum[:], h.Sum(d[:0]))
	copy(validchecksum[:], decode[21:])

	return validchecksum == computedchecksum
}

// IsBitcoinBech32Address is the validation function for validating if the field's value is a valid bech32 btc address
func IsBitcoinBech32Address(address string) bool {
	if !rxBtcBech32LowerAddress.MatchString(address) && !rxBtcBech32UpperAddress.MatchString(address) {
		return false
	}

	am := len(address) % 8
	if am == 0 || am == 3 || am == 5 {
		return false
	}

	address = strings.ToLower(address)

	alphabet := "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

	hr := []int{3, 3, 0, 2, 3} // the human readable part will always be bc
	addr := address[3:]
	dp := make([]int, 0, len(addr))

	for _, c := range addr {
		dp = append(dp, strings.IndexRune(alphabet, c))
	}

	ver := dp[0]

	if ver < 0 || ver > 16 {
		return false
	}

	if ver == 0 {
		if len(address) != 42 && len(address) != 62 {
			return false
		}
	}

	values := append(hr, dp...)

	GEN := []int{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

	p := 1

	for _, v := range values {
		b := p >> 25
		p = (p&0x1ffffff)<<5 ^ v

		for i := 0; i < 5; i++ {
			if (b>>uint(i))&1 == 1 {
				p ^= GEN[i]
			}
		}
	}

	if p != 1 {
		return false
	}

	b := uint(0)
	acc := 0
	mv := (1 << 5) - 1
	var sw []int

	for _, v := range dp[1 : len(dp)-6] {
		acc = (acc << 5) | v
		b += 5
		for b >= 8 {
			b -= 8
			sw = append(sw, (acc>>b)&mv)
		}
	}

	if len(sw) < 2 || len(sw) > 40 {
		return false
	}

	return true
}

// IsSemver is the validation function for validating if the current field's value is a valid semver version, defined in Semantic Versioning 2.0.0
func IsSemver(s string) bool {
	return rxSemver.MatchString(s)
}
