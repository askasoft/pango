package xmw

import (
	"regexp"
	"strings"
)

// ProxiedFlag Proxied flag
type ProxiedFlag int

// Proxied option flags
const (
	ProxiedAny ProxiedFlag = 1 << iota
	ProxiedAuth
	ProxiedExpired
	ProxiedNoCache
	ProxiedNoStore
	ProxiedPrivate
	ProxiedNoLastModified
	ProxiedNoETag
	ProxiedOff = 0
)

// String return level string
func (pf ProxiedFlag) String() string {
	if pf == ProxiedOff {
		return "off"
	}

	fs := make([]string, 0, 9)
	if pf&ProxiedAny == ProxiedAny {
		fs = append(fs, "any")
	}
	if pf&ProxiedAuth == ProxiedAuth {
		fs = append(fs, "auth")
	}
	if pf&ProxiedExpired == ProxiedExpired {
		fs = append(fs, "expired")
	}
	if pf&ProxiedNoCache == ProxiedNoCache {
		fs = append(fs, "no-cache")
	}
	if pf&ProxiedNoStore == ProxiedNoStore {
		fs = append(fs, "no-store")
	}
	if pf&ProxiedPrivate == ProxiedPrivate {
		fs = append(fs, "private")
	}
	if pf&ProxiedNoLastModified == ProxiedNoLastModified {
		fs = append(fs, "no_last_modified")
	}
	if pf&ProxiedNoETag == ProxiedNoETag {
		fs = append(fs, "no_etag")
	}

	return strings.Join(fs, " ")
}

// toProxiedFlag parse proxied flag from string
func toProxiedFlag(ps ...string) (pf ProxiedFlag) {
	for _, s := range ps {
		s = strings.ToLower(s)
		switch s {
		case "off":
			return ProxiedOff
		case "any":
			pf |= ProxiedAny
		case "auth":
			pf |= ProxiedAuth
		case "expired":
			pf |= ProxiedExpired
		case "no-cache":
			pf |= ProxiedNoCache
		case "no-store":
			pf |= ProxiedNoStore
		case "private":
			pf |= ProxiedPrivate
		case "no_last_modified":
			pf |= ProxiedNoLastModified
		case "no_etag":
			pf |= ProxiedNoETag
		}
	}

	return
}

type prefixs []string

func (ps prefixs) Contains(uri string) bool {
	for _, path := range ps {
		if strings.HasPrefix(uri, path) {
			return true
		}
	}
	return false
}

type regexps []*regexp.Regexp

func (rs regexps) Contains(uri string) bool {
	for _, re := range rs {
		if re.MatchString(uri) {
			return true
		}
	}
	return false
}
