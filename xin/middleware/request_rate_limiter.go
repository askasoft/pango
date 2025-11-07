package middleware

import (
	"net"
	"net/http"
	"time"

	"github.com/askasoft/pango/imc"
	"github.com/askasoft/pango/net/netx"
	"github.com/askasoft/pango/xin"
)

// RequestRateLimiter http request limit middleware
type RequestRateLimiter struct {
	Limit           int
	TrustedClients  []*net.IPNet
	TooManyRequests func(c *xin.Context)

	counts *imc.Cache[string, int]
}

// NewRequestRateLimiter create a default RequestRateLimiter middleware
func NewRequestRateLimiter(limit int, duration, cleanupInterval time.Duration) *RequestRateLimiter {
	return &RequestRateLimiter{Limit: limit, counts: imc.New[string, int](duration, cleanupInterval)}
}

func (rrl *RequestRateLimiter) Duration() time.Duration {
	return rrl.counts.TTL()
}

func (rrl *RequestRateLimiter) SetDuration(d time.Duration) {
	rrl.counts.SetTTL(d)
}

// Handle process xin request
func (rrl *RequestRateLimiter) Handle(c *xin.Context) {
	limit := rrl.Limit
	if limit <= 0 {
		c.Next()
		return
	}

	cip := c.ClientIP()
	if rrl.isTrustedClient(cip) {
		c.Next()
		return
	}

	cnt := rrl.counts.Increment(cip, 1)
	if cnt <= limit {
		c.Next()
		return
	}

	if tmr := rrl.TooManyRequests; tmr != nil {
		tmr(c)
	} else {
		c.AbortWithStatus(http.StatusTooManyRequests)
	}
}

func (rrl *RequestRateLimiter) SetTrustedClients(cidrs []string) error {
	ipnets, err := netx.ParseCIDRs(cidrs)
	if err != nil {
		return err
	}

	rrl.TrustedClients = ipnets
	return nil
}

func (rrl *RequestRateLimiter) isTrustedClient(cip string) bool {
	ip := net.ParseIP(cip)
	if ip == nil {
		return false
	}

	cidrs := rrl.TrustedClients
	for _, cidr := range cidrs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}
