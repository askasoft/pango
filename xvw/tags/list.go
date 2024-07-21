package tags

import (
	"fmt"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/cog/hashmap"
	"github.com/askasoft/pango/num"
)

type List interface {
	Each(func(string, string) bool)

	Get(key string) (string, bool)
}

func AsList(o any) List {
	if list, ok := o.(List); ok {
		return list
	}
	if ism, ok := o.(cog.Map[int, string]); ok {
		return intstrcom{ism}
	}
	if ssm, ok := o.(map[string]string); ok {
		return hashmap.AsHashMap(ssm)
	}
	if ism, ok := o.(map[int]string); ok {
		return intstrmap(ism)
	}
	if ss, ok := o.([]string); ok {
		return strslice(ss)
	}
	if is, ok := o.([]int); ok {
		return intslice(is)
	}
	if sc, ok := o.(cog.Collection[string]); ok {
		return strcol{sc}
	}
	if ic, ok := o.(cog.Collection[int]); ok {
		return intcol{ic}
	}

	panic(fmt.Sprintf("Invalid List Argument: %T", o))
}

type strslice []string

func (ss strslice) Get(key string) (string, bool) {
	if asg.Contains(ss, key) {
		return key, true
	}
	return "", false
}

func (ss strslice) Each(f func(string, string) bool) {
	for _, s := range ss {
		if !f(s, s) {
			return
		}
	}
}

type intslice []int

func (is intslice) Get(key string) (string, bool) {
	if asg.Contains(is, num.Atoi(key)) {
		return key, true
	}
	return "", false
}

func (is intslice) Each(f func(string, string) bool) {
	for _, i := range is {
		s := num.Itoa(i)
		if !f(s, s) {
			return
		}
	}
}

type intstrmap map[int]string

func (ism intstrmap) Get(key string) (string, bool) {
	if v, ok := ism[num.Atoi(key)]; ok {
		return v, true
	}
	return "", false
}

func (ism intstrmap) Each(f func(string, string) bool) {
	for k, v := range ism {
		if !f(num.Itoa(k), v) {
			return
		}
	}
}

type intstrcom struct {
	com cog.Map[int, string]
}

func (ism intstrcom) Get(key string) (string, bool) {
	if v, ok := ism.com.Get(num.Atoi(key)); ok {
		return v, true
	}
	return "", false
}

func (ism intstrcom) Each(f func(string, string) bool) {
	ism.com.Each(func(k int, v string) bool {
		return f(num.Itoa(k), v)
	})
}

type strcol struct {
	col cog.Collection[string]
}

func (sc strcol) Get(key string) (string, bool) {
	if sc.col.Contain(key) {
		return key, true
	}
	return "", false
}

func (sc strcol) Each(f func(string, string) bool) {
	sc.col.Each(func(i int, v string) bool {
		return f(v, v)
	})
}

type intcol struct {
	col cog.Collection[int]
}

func (ic intcol) Get(key string) (string, bool) {
	if ic.col.Contain(num.Atoi(key)) {
		return key, true
	}
	return "", false
}

func (ic intcol) Each(f func(string, string) bool) {
	ic.col.Each(func(i int, v int) bool {
		s := num.Itoa(v)
		return f(s, s)
	})
}
