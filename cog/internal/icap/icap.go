package icap

// MinArrayCap is smallest capacity that array may have.
const MinArrayCap = 16

// Doubleup grow up capcity by c *= 2 until c >= n
func Doubleup(c, n int) int {
	if c < MinArrayCap {
		c = MinArrayCap
	}
	for c < n {
		c <<= 1
	}
	return c
}

// Roundup round up size by the block size r
func Roundup(n, r int) int {
	r--
	if (n & r) == 0 {
		return n
	}

	return (n + r) & (^r)
}
