package tags

import (
	"fmt"

	"github.com/askasoft/pango/str"
)

type Attrs map[string]string

func (a Attrs) Get(name string) string {
	if v, ok := a[name]; ok {
		return v
	}
	return ""
}

func (a Attrs) Set(name string, value string) {
	a[name] = value
}

func (a Attrs) Add(name string, value string) {
	if value == "" {
		return
	}

	if v, ok := a[name]; ok {
		a[name] = fmt.Sprintf("%v %v", v, value)
	} else {
		a[name] = value
	}
}

func (a Attrs) Data(name string, value string) {
	name = "data-" + str.SnakeCase(name, '-')
	a.Set(name, value)
}

func (a Attrs) ID(s string) {
	a.Set("id", s)
}

func (a Attrs) Class(s string) {
	a.Add("class", s)
}

func (a Attrs) Style(s string) {
	a.Add("style", s)
}

func (a Attrs) Classes(cs ...string) {
	a.Class(str.Join(cs, " "))
}

func (a Attrs) Styles(ss ...string) {
	a.Style(str.Join(ss, ";"))
}
