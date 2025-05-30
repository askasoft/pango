package vad

import "regexp"

var (
	rxUser       = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	rxHost       = regexp.MustCompile(`^[^\s]+\.[^\s]+$`)
	rxUserDot    = regexp.MustCompile("(^[.]{1})|([.]{1}$)|([.]{2,})")
	rxEmail      = regexp.MustCompile(`^([\w\.\+\-]+)@([\w\-]+\.)+([a-zA-Z]{2,})$`)
	rxCreditCard = regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|(222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11}|6[27][0-9]{14})$`)
	rxISBN10     = regexp.MustCompile("^(?:[0-9]{9}X|[0-9]{10})$")
	rxISBN13     = regexp.MustCompile("^(?:[0-9]{13})$")
	rxUUID3      = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-3[0-9a-fA-F]{3}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")
	rxUUID4      = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$")
	rxUUID5      = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-5[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$")
	rxUUID       = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")
	rxRGBColor   = regexp.MustCompile(`^rgb\(\s*(0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])\s*,\s*(0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])\s*,\s*(0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])\s*\)$`)
	rxRGBAColor  = regexp.MustCompile(`^rgba\(\s*(?:(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])\s*,\s*(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])\s*,\s*(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])|(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])%\s*,\s*(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])%\s*,\s*(?:0|[1-9]\d?|1\d\d?|2[0-4]\d|25[0-5])%)\s*,\s*(?:(?:0.[1-9]*)|[01])\s*\)$`)
	rxHSLColor   = regexp.MustCompile(`^hsl\(\s*(?:0|[1-9]\d?|[12]\d\d|3[0-5]\d|360)\s*,\s*(?:(?:0|[1-9]\d?|100)%)\s*,\s*(?:(?:0|[1-9]\d?|100)%)\s*\)$`)
	rxHSLAColor  = regexp.MustCompile(`^hsla\(\s*(?:0|[1-9]\d?|[12]\d\d|3[0-5]\d|360)\s*,\s*(?:(?:0|[1-9]\d?|100)%)\s*,\s*(?:(?:0|[1-9]\d?|100)%)\s*,\s*(?:(?:0.[1-9]*)|[01])\s*\)$`)
	rxBase64     = regexp.MustCompile(`^(?:[A-Za-z0-9+\/]{4})*(?:[A-Za-z0-9+\/]{2}==|[A-Za-z0-9+\/]{3}=|[A-Za-z0-9+\/]{4})$`)
	rxBase64URL  = regexp.MustCompile(`^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3}=|[A-Za-z0-9-_]{4})$`)
	rxMagnetURI  = regexp.MustCompile(`^magnet:\?xt=urn:[a-zA-Z0-9]+:[a-zA-Z0-9]{32,40}&dn=.+&tr=.+$`)
	rxLatitude   = regexp.MustCompile(`^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?)$`)
	rxLongitude  = regexp.MustCompile(`^[-+]?(180(\.0+)?|((1[0-7]\d)|([1-9]?\d))(\.\d+)?)$`)
	rxDNSName    = regexp.MustCompile(`^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`)
	rxJWT        = regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$`)
	rxSSN        = regexp.MustCompile(`^\d{3}[- ]?\d{2}[- ]?\d{4}$`)
	rxIMEI       = regexp.MustCompile(`^[0-9a-f]{14}$|^\d{15}$|^\d{18}$`)
	rxIMSI       = regexp.MustCompile(`^\d{14,15}$`)
	rxSwiftCode  = regexp.MustCompile(`^[A-Za-z]{6}[A-Za-z0-9]{2}([A-Za-z0-9]{3})?$`)
)

var (
	rxE164 = regexp.MustCompile(`^\+[1-9]?[0-9]{7,14}$`)

	rxULID = regexp.MustCompile("^[A-HJKMNP-TV-Z0-9]{26}$")

	// https://tools.ietf.org/html/rfc952
	rxHostnameRFC952 = regexp.MustCompile(`^[a-zA-Z]([a-zA-Z0-9\-]+[\.]?)*[a-zA-Z0-9]$`)

	// accepts hostname starting with a digit https://tools.ietf.org/html/rfc1123
	rxHostnameRFC1123 = regexp.MustCompile(`^([a-zA-Z0-9]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*?$`)

	// same as rxHostnameRFC1123 but must contain a non numerical TLD (possibly ending with '.')
	rxFqdnRFC1123 = regexp.MustCompile(`^([a-zA-Z0-9]{1}[a-zA-Z0-9_-]{0,62})(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*?(\.[a-zA-Z]{1}[a-zA-Z0-9]{0,62})\.?$`)

	// bitcoin address
	rxBtcAddress = regexp.MustCompile(`^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$`)

	// bitcoin bech32 address https://en.bitcoin.it/wiki/Bech32
	rxBtcBech32UpperAddress = regexp.MustCompile(`^BC1[02-9AC-HJ-NP-Z]{7,76}$`)

	// bitcoin bech32 address https://en.bitcoin.it/wiki/Bech32
	rxBtcBech32LowerAddress = regexp.MustCompile(`^bc1[02-9ac-hj-np-z]{7,76}$`)

	// numbered capture groups https://semver.org/
	rxSemver = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

	rxDnsRFC1035Label = regexp.MustCompile("^[a-z]([-a-z0-9]*[a-z0-9]){0,62}$")
)
