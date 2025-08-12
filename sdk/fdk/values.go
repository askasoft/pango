package fdk

import (
	"fmt"
	"net/url"

	"github.com/askasoft/pango/num"
	"github.com/askasoft/pango/str"
)

type ListOption interface {
	IsNil() bool
	Values() Values
}

type PageOption struct {
	Page    int
	PerPage int
}

func (po *PageOption) IsNil() bool {
	return po == nil
}

func (po *PageOption) Values() Values {
	q := Values{}
	q.SetInt("page", po.Page)
	q.SetInt("per_page", po.PerPage)
	return q
}

type File interface {
	Field() string
	File() string
	Data() []byte
}

type Files []File

type WithFiles interface {
	Files() Files
	Values() Values
}

type Values url.Values

func (vs Values) Map() map[string][]string {
	return (map[string][]string)(vs)
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (vs Values) Get(key string) string {
	return url.Values(vs).Get(key)
}

// Set sets the key to value. It replaces any existing values.
func (vs Values) Set(key, value string) {
	url.Values(vs).Set(key, value)
}

// Add adds the value to key. It appends to any existing values associated with key.
func (vs Values) Add(key, value string) {
	url.Values(vs).Add(key, value)
}

// Del deletes the values associated with key.
func (vs Values) Del(key string) {
	url.Values(vs).Del(key)
}

// Has checks whether a given key is set.
func (vs Values) Has(key string) bool {
	return url.Values(vs).Has(key)
}

// Encode encodes the values into “URL encoded” form ("bar=baz&foo=quux") sorted by key.
func (vs Values) Encode() string {
	return url.Values(vs).Encode()
}

func (vs Values) SetBool(name string, value bool) {
	s := str.If(value, "true", "false")
	vs.Set(name, s)
}

func (vs Values) SetString(name string, value string) {
	if value != "" {
		vs.Set(name, value)
	}
}

func (vs Values) SetStrings(name string, value []string) {
	if len(value) > 0 {
		name += "[]"
		for _, s := range value {
			vs.Add(name, s)
		}
	}
}

func (vs Values) SetStringsPtr(name string, value *[]string) {
	if value != nil {
		name += "[]"
		if len(*value) == 0 {
			vs.Add(name, "")
		} else {
			for _, s := range *value {
				vs.Add(name, s)
			}
		}
	}
}

func (vs Values) SetInts(name string, value []int) {
	if len(value) > 0 {
		name += "[]"
		for _, n := range value {
			vs.Add(name, num.Itoa(n))
		}
	}
}

func (vs Values) SetInt64s(name string, value []int64) {
	if len(value) > 0 {
		name += "[]"
		for _, n := range value {
			vs.Add(name, num.Ltoa(n))
		}
	}
}

func (vs Values) SetStringMap(name string, value map[string]string) {
	if len(value) > 0 {
		for k, v := range value {
			vs.Add(fmt.Sprintf("%s[%s]", name, k), v)
		}
	}
}

func (vs Values) SetMap(name string, value map[string]any) {
	if len(value) > 0 {
		for k, v := range value {
			vs.Add(fmt.Sprintf("%s[%s]", name, k), fmt.Sprint(v))
		}
	}
}

func (vs Values) SetDate(name string, value Date) {
	if !value.IsZero() {
		vs.Set(name, value.String())
	}
}

func (vs Values) SetDatePtr(name string, value *Date) {
	if value != nil && !value.IsZero() {
		vs.Set(name, value.String())
	}
}

func (vs Values) SetTime(name string, value Time) {
	if !value.IsZero() {
		vs.Set(name, value.String())
	}
}

func (vs Values) SetTimePtr(name string, value *Time) {
	if value != nil && !value.IsZero() {
		vs.Set(name, value.String())
	}
}

func (vs Values) SetInt(name string, value int) {
	if value != 0 {
		vs.Set(name, num.Itoa(value))
	}
}

func (vs Values) SetInt64(name string, value int64) {
	if value != 0 {
		vs.Set(name, num.Ltoa(value))
	}
}
