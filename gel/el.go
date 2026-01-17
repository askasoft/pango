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
	if ec.Object == nil {
		if ec.Strict {
			return nil, fmt.Errorf("gel: can't get nil.%s", key)
		}
		return nil, nil
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

	rv, err := invokeFunc(fv, args)
	if err != nil {
		err = fmt.Errorf("gel: function %s(): %w", eo.val, err)
	}
	return rv, err
}

func invokeFunc(fv reflect.Value, args []any) (any, error) {
	ft := fv.Type()
	if ft.NumOut() > 2 {
		return nil, fmt.Errorf("too many return values: want 1~2, got %d", ft.NumOut())
	}

	rvs, err := ref.CallFunction(fv, args)
	if err != nil {
		return nil, err
	}

	switch len(rvs) {
	case 0:
		return nil, nil
	case 1:
		return rvs[0], nil
	case 2:
		ret, rer := rvs[0], rvs[1]
		if rer == nil {
			return ret, nil
		}
		if err, ok := rer.(error); ok {
			return ret, err
		}
		return ret, fmt.Errorf("gel: second return value '%T' is not error", rer)
	default:
		return nil, fmt.Errorf("gel: too many return values, want 1~2, got %d", len(rvs))
	}
}

type EL struct {
	expr string // expression
	root any    // reverse polish notation
}

// String returns the source text used to compile the el expression.
func (el *EL) String() string {
	return el.expr
}

func (el *EL) Calculate(data any) (any, error) {
	return el.calculate(elCtx{Object: data})
}

func (el *EL) CalculateStrict(data any) (any, error) {
	return el.calculate(elCtx{Object: data, Strict: true})
}

func (el *EL) calculate(ec elCtx) (any, error) {
	if el.root == nil {
		return nil, nil
	}

	if op, ok := el.root.(operator); ok {
		return op.Calculate(ec)
	}
	if eo, ok := el.root.(elObj); ok {
		return eo.Get(ec)
	}
	return el.root, nil
}

func Compile(expr string) (*EL, error) {
	var sy shuntingYard
	rpn, err := sy.ParseToRPN(expr)
	if err != nil {
		return nil, err
	}
	return &EL{expr, rpn}, nil
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
