package tags

import (
	"fmt"

	"github.com/askasoft/pango/str"
)

type Attrs map[string]string

func (a Attrs) Get(k string) string {
	if v, ok := a[k]; ok {
		return v
	}
	return ""
}

func (a Attrs) Set(k string, v string) {
	a[k] = v
}

func (a Attrs) SetIfNotEmpty(k string, v string) {
	if v != "" {
		a[k] = v
	}
}

func (a Attrs) Add(k string, v string) {
	if v == "" {
		return
	}

	if ov, ok := a[k]; ok {
		a[k] = fmt.Sprintf("%v %v", ov, v)
	} else {
		a[k] = v
	}
}

func (a Attrs) Data(k string, v string) {
	if v != "" {
		k = "data-" + str.SnakeCase(k, '-')
		a.Set(k, v)
	}
}

func (a Attrs) ID(v string) {
	a.SetIfNotEmpty("id", v)
}

func (a Attrs) Name(v string) {
	a.SetIfNotEmpty("name", v)
}

func (a Attrs) Class(v string) {
	a.Add("class", v)
}

func (a Attrs) Style(v string) {
	a.Add("style", v)
}

func (a Attrs) Classes(cs ...string) {
	a.Class(str.Join(cs, " "))
}

func (a Attrs) Styles(ss ...string) {
	a.Style(str.Join(ss, ";"))
}
