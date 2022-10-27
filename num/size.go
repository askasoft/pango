package num

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	KB = 1024
	MB = KB * KB
	GB = KB * MB
	TB = KB * GB // 1099511627776
	PB = KB * TB // 1125899906842624
	EB = KB * PB // 1152921504606846976
	ZB = KB * float64(EB)
	YB = KB * ZB
)

type unitMap map[byte]float64

var (
	sizeMap   = unitMap{'K': KB, 'M': MB, 'G': GB, 'T': TB, 'P': PB, 'E': EB, 'Z': ZB, 'Y': YB}
	sizeUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
)

func getSizeAndUnit(size float64, base float64, units []string) (float64, string) {
	i := 0
	l := len(units) - 1
	for size >= base && i < l {
		size /= base
		i++
	}
	return size, units[i]
}

// CustomSize returns a human-readable approximation of a size
// using custom format.
func CustomSize(format string, size float64, base float64, units []string) string {
	size, unit := getSizeAndUnit(size, base, units)
	return fmt.Sprintf(format, size, unit)
}

// HumanSizeWithPrecision allows the size to be in any precision,
// instead of 4 digit precision used in HumanSize.
func HumanSizeWithPrecision(size float64, precision int) string {
	size, unit := getSizeAndUnit(size, KB, sizeUnits)
	return fmt.Sprintf("%.*g %s", precision, size, unit)
}

// HumanSize returns a human-readable approximation of a size
// capped at 4 valid numbers (eg. "2.746 MB", "796 KB").
func HumanSize(size float64) string {
	return HumanSizeWithPrecision(size, 4)
}

// ParseSize returns an integer from a human-readable size using windows specification (KB = 1024B).
func ParseSize(size string) (int64, error) {
	return parseSize(size, sizeMap)
}

// MustParseSize returns an integer from a human-readable size using windows specification (KB = 1024B).
// panic if parse error
func MustParseSize(size string) int64 {
	sz, err := parseSize(size, sizeMap)
	if err != nil {
		panic(err)
	}
	return sz
}

// ParseSizeF returns a float64 from a human-readable size using windows specification (KB = 1024B).
func ParseSizeF(size string) (float64, error) {
	return parseSizeF(size, sizeMap)
}

// MustParseSizeF returns a float64 from a human-readable size using windows specification (KB = 1024B).
// panic if parse error
func MustParseSizeF(size string) float64 {
	sz, err := parseSizeF(size, sizeMap)
	if err != nil {
		panic(err)
	}
	return sz
}

func badSuffix(sfx string) (float64, error) {
	return -1, fmt.Errorf("invalid suffix: '%s'", sfx)
}

func badSize(sz string) (float64, error) {
	return -1, fmt.Errorf("invalid size: '%s'", sz)
}

// Parses the human-readable size string into the amount it represents.
func parseSize(size string, uMap unitMap) (int64, error) {
	sz, err := parseSizeF(size, uMap)
	if err != nil {
		return -1, err
	}
	return int64(sz), nil
}

// Parses the human-readable size string into the amount it represents.
func parseSizeF(size string, uMap unitMap) (float64, error) {
	sep := strings.LastIndexAny(size, "01234567890. ")
	if sep == -1 {
		// There should be at least a digit.
		return badSize(size)
	}

	sfx := size[sep+1:]

	// Omit the space separator.
	if size[sep] != ' ' {
		sep++
	}
	num := size[:sep]

	sz, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return -1, err
	}

	if sz < 0 {
		return badSize(size)
	}

	if len(sfx) == 0 {
		return sz, nil
	}

	// Process the suffix.

	if len(sfx) > 2 { // Too long.
		return badSuffix(sfx)
	}

	sfx = strings.ToUpper(sfx)

	// Trivial case: B suffix.
	if sfx[0] == 'B' {
		if len(sfx) > 1 { // no extra characters allowed after b.
			return badSuffix(sfx)
		}
		return sz, nil
	}

	// A suffix from the map.
	if mul, ok := uMap[sfx[0]]; ok {
		sz *= float64(mul)
	} else {
		return badSuffix(sfx)
	}

	// The suffix may have extra "B".
	if len(sfx) == 2 && sfx[1] != 'B' {
		return badSuffix(sfx)
	}

	return sz, nil
}
