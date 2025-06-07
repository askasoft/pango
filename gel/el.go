package gel

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/askasoft/pango/ref"
)

type elCtx struct {
	Object any
	Strict bool
}

func (ec elCtx) Get(key string) (any, error) {
	if key == "$" {
		return ec.Object, nil
	}
	return ref.GetProperty(ec.Object, key)
}

type elObj struct {
	val string
}

func (eo elObj) String() string {
	return eo.val
}

func (eo elObj) Get(ec elCtx) (any, error) {
	if eo.val == "" {
		return nil, nil
	}
	return ec.Get(eo.val)
}

func (eo elObj) Invoke(ec elCtx, args []any) (any, error) {
	if eo.val == "" {
		return nil, errors.New("gel: empty function name")
	}

	fn, err := ec.Get(eo.val)
	if err != nil {
		return fn, err
	}
	if fn == nil {
		return nil, fmt.Errorf("gel: function %q is nil", eo.val)
	}

	fv := reflect.ValueOf(fn)
	if fv.Kind() != reflect.Func {
		return nil, fmt.Errorf("gel: %q is not a function", eo.val)
	}

	ret, err := invokeFunc(fv, args)
	if err != nil {
		return nil, fmt.Errorf("gel: function %s(): %w", ref.NameOfFuncValue(fv), err)
	}

	return ret, err
}

func invokeFunc(fv reflect.Value, args []any) (any, error) {
	ft := fv.Type()
	if ft.NumIn() != len(args) {
		return nil, fmt.Errorf("invalid argument count, want %d, got %d", ft.NumIn(), len(args))
	}

	if ft.NumOut() > 2 {
		return nil, fmt.Errorf("invalid return count, want 1~2, got %d", ft.NumOut())
	}

	var avs []reflect.Value
	for i, a := range args {
		t := ft.In(i)

		v, err := ref.CastTo(a, t)
		if err != nil {
			return nil, fmt.Errorf("invalid argument #%d - %w", i, err)
		}

		avs = append(avs, reflect.ValueOf(v))
	}

	rvs := fv.Call(avs)
	switch len(rvs) {
	case 0:
		return nil, nil
	case 1:
		return rvs[0].Interface(), nil
	case 2:
		ret, rer := rvs[0].Interface(), rvs[1].Interface()
		if rer == nil {
			return ret, nil
		}
		if err, ok := rer.(error); ok {
			return ret, err
		}
		return ret, fmt.Errorf("second return value '%T' is not error", rer)
	default:
		return nil, fmt.Errorf("invalid return count (%d)", len(rvs))
	}
}

type EL struct {
	expr string // expression
	rpn  reversePolishNotation
}

func (el *EL) Calculate(data any) (any, error) {
	return el.rpn.Calculate(elCtx{Object: data})
}

func (el *EL) CalculateStrict(data any) (any, error) {
	return el.rpn.Calculate(elCtx{Object: data, Strict: true})
}

func Compile(expr string) (*EL, error) {
	var sy shuntingYard
	if err := sy.ParseToRPN(expr); err != nil {
		return nil, err
	}
	return &EL{expr, newReversePolishNotation(&sy.rpn)}, nil
}

func Calculate(expr string, data any) (any, error) {
	el, err := Compile(expr)
	if err != nil {
		return nil, err
	}
	return el.Calculate(data)
}

func CalculateStrict(expr string, data any) (any, error) {
	el, err := Compile(expr)
	if err != nil {
		return nil, err
	}
	return el.CalculateStrict(data)
}
