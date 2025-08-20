package taglib

import (
	"fmt"
	"reflect"

	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/cas"
	"github.com/askasoft/pango/ref"
)

func toString(a any) string {
	s, err := cas.ToString(a)
	if err != nil {
		panic(err)
	}
	return s
}

func convert[T any](a any, p *T) {
	t := reflect.ValueOf(*p).Type()

	v, err := ref.CastTo(a, t)
	if err != nil {
		panic(err)
	}

	*p = v.(T)
}

func contains(vs Values, v any) bool {
	return vs != nil && vs.Contains(v)
}

//------------------------------------------------

// List List field interface
type List interface {
	Each(func(any, string) bool)
	Get(key any) (string, bool)
}

type xstrdict[T comparable] interface {
	Each(func(T, string) bool)
	Get(key T) (string, bool)
}

type collection[T any] interface {
	Each(func(int, T) bool)
	Contains(T) bool
}

func AsList(a any) List {
	switch o := a.(type) {
	case List:
		return o
	case xstrdict[string]:
		return xstrdict2list[string]{o}
	case xstrdict[int]:
		return xstrdict2list[int]{o}
	case xstrdict[int8]:
		return xstrdict2list[int8]{o}
	case xstrdict[int16]:
		return xstrdict2list[int16]{o}
	case xstrdict[int32]:
		return xstrdict2list[int32]{o}
	case xstrdict[int64]:
		return xstrdict2list[int64]{o}
	case xstrdict[uint]:
		return xstrdict2list[uint]{o}
	case xstrdict[uint8]:
		return xstrdict2list[uint8]{o}
	case xstrdict[uint16]:
		return xstrdict2list[uint16]{o}
	case xstrdict[uint32]:
		return xstrdict2list[uint32]{o}
	case xstrdict[uint64]:
		return xstrdict2list[uint64]{o}
	case xstrdict[float32]:
		return xstrdict2list[float32]{o}
	case xstrdict[float64]:
		return xstrdict2list[float64]{o}
	case collection[string]:
		return strcol2list{o}
	case collection[int]:
		return xcol2list[int]{o}
	case collection[int8]:
		return xcol2list[int8]{o}
	case collection[int16]:
		return xcol2list[int16]{o}
	case collection[int32]:
		return xcol2list[int32]{o}
	case collection[int64]:
		return xcol2list[int64]{o}
	case collection[uint]:
		return xcol2list[uint]{o}
	case collection[uint8]:
		return xcol2list[uint8]{o}
	case collection[uint16]:
		return xcol2list[uint16]{o}
	case collection[uint32]:
		return xcol2list[uint32]{o}
	case collection[uint64]:
		return xcol2list[uint64]{o}
	case collection[float32]:
		return xcol2list[float32]{o}
	case collection[float64]:
		return xcol2list[float64]{o}
	case map[string]string:
		return strstrmap2list(o)
	case map[int]string:
		return xstrmap2list[int](o)
	case map[int8]string:
		return xstrmap2list[int8](o)
	case map[int16]string:
		return xstrmap2list[int16](o)
	case map[int32]string:
		return xstrmap2list[int32](o)
	case map[int64]string:
		return xstrmap2list[int64](o)
	case map[uint]string:
		return xstrmap2list[uint](o)
	case map[uint8]string:
		return xstrmap2list[uint8](o)
	case map[uint16]string:
		return xstrmap2list[uint16](o)
	case map[uint32]string:
		return xstrmap2list[uint32](o)
	case map[uint64]string:
		return xstrmap2list[uint64](o)
	case map[float32]string:
		return xstrmap2list[float32](o)
	case map[float64]string:
		return xstrmap2list[float64](o)
	case map[any]string:
		return anystrmap2list(o)
	case []string:
		return strslice2list(o)
	case []int:
		return xslice2list[int](o)
	case []int8:
		return xslice2list[int8](o)
	case []int16:
		return xslice2list[int16](o)
	case []int32:
		return xslice2list[int32](o)
	case []int64:
		return xslice2list[int64](o)
	case []uint:
		return xslice2list[uint](o)
	case []uint8:
		return xslice2list[uint8](o)
	case []uint16:
		return xslice2list[uint16](o)
	case []uint32:
		return xslice2list[uint32](o)
	case []uint64:
		return xslice2list[uint64](o)
	case []float32:
		return xslice2list[float32](o)
	case []float64:
		return xslice2list[float64](o)
	case []any:
		return anyslice2list(o)
	default:
		switch reflect.ValueOf(a).Kind() {
		case reflect.Map:
			if cv, err := ref.ConvertTo(a, ref.TypeStrStrMap); err == nil {
				return strstrmap2list(cv.(map[string]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeIntStrMap); err == nil {
				return xstrmap2list[int](cv.(map[int]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInt8StrMap); err == nil {
				return xstrmap2list[int8](cv.(map[int8]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInt16StrMap); err == nil {
				return xstrmap2list[int16](cv.(map[int16]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInt32StrMap); err == nil {
				return xstrmap2list[int32](cv.(map[int32]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInt64StrMap); err == nil {
				return xstrmap2list[int64](cv.(map[int64]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUintStrMap); err == nil {
				return xstrmap2list[uint](cv.(map[uint]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUint8StrMap); err == nil {
				return xstrmap2list[uint8](cv.(map[uint8]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUint16StrMap); err == nil {
				return xstrmap2list[uint16](cv.(map[uint16]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUint32StrMap); err == nil {
				return xstrmap2list[uint32](cv.(map[uint32]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUint64StrMap); err == nil {
				return xstrmap2list[uint64](cv.(map[uint64]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeFloat32StrMap); err == nil {
				return xstrmap2list[float32](cv.(map[float32]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeFloat64StrMap); err == nil {
				return xstrmap2list[float64](cv.(map[float64]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeAnyStrMap); err == nil {
				return anystrmap2list(cv.(map[any]string))
			}
		case reflect.Slice:
			if cv, err := ref.ConvertTo(a, ref.TypeStrings); err == nil {
				return strslice2list(cv.([]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInts); err == nil {
				return xslice2list[int](cv.([]int))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInt8s); err == nil {
				return xslice2list[int8](cv.([]int8))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInt16s); err == nil {
				return xslice2list[int16](cv.([]int16))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInt32s); err == nil {
				return xslice2list[int32](cv.([]int32))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInt64s); err == nil {
				return xslice2list[int64](cv.([]int64))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUints); err == nil {
				return xslice2list[uint](cv.([]uint))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUint8s); err == nil {
				return xslice2list[uint8](cv.([]uint8))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUint16s); err == nil {
				return xslice2list[uint16](cv.([]uint16))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUint32s); err == nil {
				return xslice2list[uint32](cv.([]uint32))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUint64s); err == nil {
				return xslice2list[uint64](cv.([]uint64))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeFloat32s); err == nil {
				return xslice2list[float32](cv.([]float32))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeFloat64s); err == nil {
				return xslice2list[float64](cv.([]float64))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeAnys); err == nil {
				return anyslice2list(cv.([]any))
			}
		}
		panic(fmt.Errorf("tags: invalid 'List' argument: %T", o))
	}
}

type xstrdict2list[T comparable] struct {
	xsd xstrdict[T]
}

func (xsdl xstrdict2list[T]) Each(f func(any, string) bool) {
	xsdl.xsd.Each(func(k T, v string) bool {
		return f(k, v)
	})
}

func (xsdl xstrdict2list[T]) Get(a any) (string, bool) {
	var k T

	convert(a, &k)
	if v, ok := xsdl.xsd.Get(k); ok {
		return v, true
	}
	return "", false
}

type strcol2list struct {
	col collection[string]
}

func (scl strcol2list) Each(f func(any, string) bool) {
	scl.col.Each(func(i int, v string) bool {
		return f(v, v)
	})
}

func (scl strcol2list) Get(a any) (string, bool) {
	s := toString(a)
	if scl.col.Contains(s) {
		return s, true
	}
	return "", false
}

type xcol2list[T comparable] struct {
	col collection[T]
}

func (xcl xcol2list[T]) Each(f func(any, string) bool) {
	xcl.col.Each(func(i int, v T) bool {
		s := toString(v)
		return f(v, s)
	})
}

func (xcl xcol2list[T]) Get(a any) (string, bool) {
	var k T

	convert(a, &k)
	if xcl.col.Contains(k) {
		return toString(k), true
	}
	return "", false
}

type strstrmap2list map[string]string

func (ssm strstrmap2list) Each(f func(any, string) bool) {
	for k, v := range ssm {
		if !f(k, v) {
			return
		}
	}
}

func (ssm strstrmap2list) Get(a any) (string, bool) {
	k := toString(a)
	if v, ok := ssm[k]; ok {
		return v, true
	}
	return "", false
}

type xstrmap2list[T comparable] map[T]string

func (xsm xstrmap2list[T]) Each(f func(any, string) bool) {
	for k, v := range xsm {
		if !f(k, v) {
			return
		}
	}
}

func (xsm xstrmap2list[T]) Get(a any) (string, bool) {
	var k T

	convert(a, &k)
	if v, ok := xsm[k]; ok {
		return v, true
	}
	return "", false
}

type anystrmap2list map[any]string

func (asm anystrmap2list) Each(f func(any, string) bool) {
	for k, v := range asm {
		if !f(k, v) {
			return
		}
	}
}

func (asm anystrmap2list) Get(a any) (string, bool) {
	if v, ok := asm[a]; ok {
		return v, true
	}
	return "", false
}

type strslice2list []string

func (ss strslice2list) Each(f func(any, string) bool) {
	for _, s := range ss {
		if !f(s, s) {
			return
		}
	}
}

func (ss strslice2list) Get(a any) (string, bool) {
	k := toString(a)
	if asg.Contains(ss, k) {
		return k, true
	}
	return "", false
}

type xslice2list[T comparable] []T

func (xs xslice2list[T]) Each(f func(any, string) bool) {
	for _, k := range xs {
		v := toString(k)
		if !f(k, v) {
			return
		}
	}
}

func (xs xslice2list[T]) Get(a any) (string, bool) {
	var v T

	convert(a, &v)

	if asg.Contains(xs, v) {
		s := toString(v)
		return s, true
	}
	return "", false
}

type anyslice2list []any

func (as anyslice2list) Each(f func(any, string) bool) {
	for _, k := range as {
		v := toString(k)
		if !f(k, v) {
			return
		}
	}
}

func (as anyslice2list) Get(a any) (string, bool) {
	if asg.Contains(as, a) {
		return toString(a), true
	}
	return "", false
}

//------------------------------------------------

// Values Values field interface
type Values interface {
	Each(func(any) bool)
	Contains(any) bool
}

func AsValues(a any) Values {
	switch o := a.(type) {
	case Values:
		return o
	case []string:
		return strslice2values(o)
	case []int:
		return xslice2values[int](o)
	case []int8:
		return xslice2values[int8](o)
	case []int16:
		return xslice2values[int16](o)
	case []int32:
		return xslice2values[int32](o)
	case []int64:
		return xslice2values[int64](o)
	case []uint:
		return xslice2values[uint](o)
	case []uint8:
		return xslice2values[uint8](o)
	case []uint16:
		return xslice2values[uint16](o)
	case []uint32:
		return xslice2values[uint32](o)
	case []uint64:
		return xslice2values[uint64](o)
	case []float32:
		return xslice2values[float32](o)
	case []float64:
		return xslice2values[float64](o)
	case []any:
		return anyslice2values(o)
	case collection[string]:
		return strcol2values{o}
	case collection[int]:
		return xcol2values[int]{o}
	case collection[int8]:
		return xcol2values[int8]{o}
	case collection[int16]:
		return xcol2values[int16]{o}
	case collection[int32]:
		return xcol2values[int32]{o}
	case collection[int64]:
		return xcol2values[int64]{o}
	case collection[uint]:
		return xcol2values[uint]{o}
	case collection[uint8]:
		return xcol2values[uint8]{o}
	case collection[uint16]:
		return xcol2values[uint16]{o}
	case collection[uint32]:
		return xcol2values[uint32]{o}
	case collection[uint64]:
		return xcol2values[uint64]{o}
	case collection[float32]:
		return xcol2values[float32]{o}
	case collection[float64]:
		return xcol2values[float64]{o}
	case collection[any]:
		return anycol2values{o}
	default:
		if reflect.ValueOf(a).Kind() == reflect.Slice {
			if cv, err := ref.ConvertTo(a, ref.TypeStrings); err == nil {
				return strslice2values(cv.([]string))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInts); err == nil {
				return xslice2values[int](cv.([]int))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInt8s); err == nil {
				return xslice2values[int8](cv.([]int8))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInt16s); err == nil {
				return xslice2values[int16](cv.([]int16))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInt32s); err == nil {
				return xslice2values[int32](cv.([]int32))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeInt64s); err == nil {
				return xslice2values[int64](cv.([]int64))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUints); err == nil {
				return xslice2values[uint](cv.([]uint))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUint8s); err == nil {
				return xslice2values[uint8](cv.([]uint8))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUint16s); err == nil {
				return xslice2values[uint16](cv.([]uint16))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUint32s); err == nil {
				return xslice2values[uint32](cv.([]uint32))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeUint64s); err == nil {
				return xslice2values[uint64](cv.([]uint64))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeFloat32s); err == nil {
				return xslice2values[float32](cv.([]float32))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeFloat64s); err == nil {
				return xslice2values[float64](cv.([]float64))
			}
			if cv, err := ref.ConvertTo(a, ref.TypeAnys); err == nil {
				return anyslice2values(cv.([]any))
			}
		}
		panic(fmt.Errorf("tags: invalid 'Values' argument: %T", o))
	}
}

type anycol2values struct {
	col collection[any]
}

func (acv anycol2values) Each(f func(any) bool) {
	acv.col.Each(func(i int, v any) bool {
		return f(v)
	})
}

func (acv anycol2values) Contains(a any) bool {
	return acv.col.Contains(a)
}

type strcol2values struct {
	col collection[string]
}

func (scv strcol2values) Each(f func(any) bool) {
	scv.col.Each(func(i int, v string) bool {
		return f(v)
	})
}

func (scv strcol2values) Contains(a any) bool {
	s := toString(a)
	return scv.col.Contains(s)
}

type xcol2values[T comparable] struct {
	col collection[T]
}

func (acv xcol2values[T]) Each(f func(any) bool) {
	acv.col.Each(func(i int, v T) bool {
		return f(v)
	})
}

func (acv xcol2values[T]) Contains(a any) bool {
	var k T

	convert(a, &k)

	return acv.col.Contains(k)
}

type anyslice2values []any

func (as anyslice2values) Each(f func(any) bool) {
	for _, a := range as {
		if !f(a) {
			return
		}
	}
}

func (as anyslice2values) Contains(a any) bool {
	return asg.Contains(as, a)
}

type strslice2values []string

func (ss strslice2values) Each(f func(any) bool) {
	for _, s := range ss {
		if !f(s) {
			return
		}
	}
}

func (ss strslice2values) Contains(a any) bool {
	return asg.Contains(ss, toString(a))
}

type xslice2values[T comparable] []T

func (xs xslice2values[T]) Each(f func(any) bool) {
	for _, v := range xs {
		if !f(v) {
			return
		}
	}
}

func (xs xslice2values[T]) Contains(a any) bool {
	var v T

	convert(a, &v)

	return asg.Contains(xs, v)
}
