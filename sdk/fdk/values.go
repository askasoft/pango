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

func (vs Values) SetBool(name string, value bool) {
	s := str.If(value, "true", "false")
	(url.Values)(vs).Set(name, s)
}

func (vs Values) SetString(name string, value string) {
	if value != "" {
		(url.Values)(vs).Set(name, value)
	}
}

func (vs Values) SetStrings(name string, value []string) {
	if len(value) > 0 {
		name += "[]"
		for _, s := range value {
			(url.Values)(vs).Add(name, s)
		}
	}
}

func (vs Values) SetStringsPtr(name string, value *[]string) {
	if value != nil {
		name += "[]"
		if len(*value) == 0 {
			(url.Values)(vs).Add(name, "")
		} else {
			for _, s := range *value {
				(url.Values)(vs).Add(name, s)
			}
		}
	}
}

func (vs Values) SetInts(name string, value []int) {
	if len(value) > 0 {
		name += "[]"
		for _, n := range value {
			(url.Values)(vs).Add(name, num.Itoa(n))
		}
	}
}

func (vs Values) SetInt64s(name string, value []int64) {
	if len(value) > 0 {
		name += "[]"
		for _, n := range value {
			(url.Values)(vs).Add(name, num.Ltoa(n))
		}
	}
}

func (vs Values) SetStringMap(name string, value map[string]string) {
	if len(value) > 0 {
		for k, v := range value {
			(url.Values)(vs).Add(fmt.Sprintf("%s[%s]", name, k), v)
		}
	}
}

func (vs Values) SetMap(name string, value map[string]any) {
	if len(value) > 0 {
		for k, v := range value {
			(url.Values)(vs).Add(fmt.Sprintf("%s[%s]", name, k), fmt.Sprint(v))
		}
	}
}

func (vs Values) SetDate(name string, value Date) {
	if !value.IsZero() {
		(url.Values)(vs).Set(name, value.String())
	}
}

func (vs Values) SetDatePtr(name string, value *Date) {
	if value != nil && !value.IsZero() {
		(url.Values)(vs).Set(name, value.String())
	}
}

func (vs Values) SetTime(name string, value Time) {
	if !value.IsZero() {
		(url.Values)(vs).Set(name, value.String())
	}
}

func (vs Values) SetTimePtr(name string, value *Time) {
	if value != nil && !value.IsZero() {
		(url.Values)(vs).Set(name, value.String())
	}
}

func (vs Values) SetInt(name string, value int) {
	if value != 0 {
		(url.Values)(vs).Set(name, num.Itoa(value))
	}
}

func (vs Values) SetInt64(name string, value int64) {
	if value != 0 {
		(url.Values)(vs).Set(name, num.Ltoa(value))
	}
}

func (vs Values) Encode() string {
	return (url.Values)(vs).Encode()
}
