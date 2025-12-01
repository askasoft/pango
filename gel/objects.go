package gel

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/askasoft/pango/cas"
	"github.com/askasoft/pango/cog"
	"github.com/askasoft/pango/ref"
)

type accessOp struct {
	doubleOp
}

func (ao *accessOp) Category() int {
	return opObject
}

func (ao *accessOp) String() string {
	return "."
}

func (ao *accessOp) Priority() int {
	return 1
}

func (ao *accessOp) Calculate(ec elCtx) (any, error) {
	obj, err := ao.evalLeft(ec)
	if err != nil {
		return obj, err
	}
	if obj == nil {
		if ec.Strict {
			return nil, fmt.Errorf("gel: can't get nil.%v", ao.right)
		}
		return nil, nil
	}

	key, err := cas.ToString(ao.right)
	if err != nil {
		return nil, err
	}

	return ref.GetProperty(obj, key)
}

func (ao *accessOp) Invoke(ec elCtx, args []any) (any, error) {
	obj, err := ao.evalLeft(ec)
	if err != nil {
		return obj, err
	}
	if obj == nil {
		return nil, fmt.Errorf("gel: can't call nil.%v()", ao.right)
	}

	mn, err := cas.ToString(ao.right)
	if err != nil {
		return nil, err
	}

	fn, err := ref.GetProperty(obj, mn)
	if err != nil {
		return nil, err
	}

	fv := reflect.ValueOf(fn)
	if fv.Kind() != reflect.Func {
		return nil, fmt.Errorf("gel: invalid method %q of %T", mn, obj)
	}

	ret, err := invokeFunc(fv, args)
	if err != nil {
		return nil, fmt.Errorf("gel: method %T.%s(): %w", obj, mn, err)
	}

	return ret, err
}

type arrayGetOp struct {
	doubleOp
}

func (ago *arrayGetOp) Category() int {
	return opObject
}

func (ago *arrayGetOp) String() string {
	return "["
}

func (ago *arrayGetOp) Priority() int {
	return 1
}

func (ago *arrayGetOp) Calculate(ec elCtx) (any, error) {
	lval, rval, err := ago.calcLeftRight(ec)
	if err != nil {
		return nil, err
	}

	if lval == nil {
		if ec.Strict {
			return nil, fmt.Errorf("gel: can't get nil[%v]", rval)
		}
		return nil, nil
	}

	rk := reflect.Indirect(reflect.ValueOf(lval)).Kind()

	if rk == reflect.Map {
		return ref.MapGet(lval, rval)
	}

	idx, err := cas.ToInt(rval)
	if err != nil {
		return nil, err
	}
	return ref.ArrayGet(lval, idx)

}

type arrayEndOp struct {
	left any
}

func (aeo *arrayEndOp) Category() int {
	return opObject
}

func (aeo *arrayEndOp) String() string {
	return "]"
}

func (aeo *arrayEndOp) Operands() int {
	return 1
}

func (aeo *arrayEndOp) Priority() int {
	return 1
}

func (aeo *arrayEndOp) Wrap(op operator, operand cog.Queue[any]) error {
	var ok bool

	aeo.left, ok = operand.Poll()
	if !ok {
		return fmt.Errorf("gel: operator %q missing left operand", aeo)
	}
	return nil
}

func (aeo *arrayEndOp) Calculate(ec elCtx) (any, error) {
	if op, ok := aeo.left.(operator); ok {
		return op.Calculate(ec)
	}
	return nil, fmt.Errorf("gel: invalid left operator '%v'", aeo.left)
}

type arrayMakeOp struct {
}

func (amo *arrayMakeOp) Category() int {
	return opObject
}

func (amo *arrayMakeOp) String() string {
	return "{}"
}

func (amo *arrayMakeOp) Operands() int {
	return 0
}

func (amo *arrayMakeOp) Priority() int {
	return 1
}

func (amo *arrayMakeOp) Wrap(op operator, operand cog.Queue[any]) error {
	return nil
}

func (amo *arrayMakeOp) Calculate(ec elCtx) (any, error) {
	return nil, errors.New("gel: '{}' does not support Calculate()")
}

func (amo *arrayMakeOp) Invoke(ec elCtx, args []any) (any, error) {
	return args, nil
}

type commaOp struct {
	doubleOp
}

func (co *commaOp) Category() int {
	return opObject
}

func (co *commaOp) String() string {
	return ","
}

func (co *commaOp) Priority() int {
	return 90
}

func (co *commaOp) Calculate(ec elCtx) (any, error) {
	var objs []any

	if lco, ok := co.left.(*commaOp); ok {
		lvs, err := lco.Calculate(ec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, lvs.([]any)...)
	} else {
		lval, err := co.calcLeft(ec)
		if err != nil {
			return nil, err
		}
		objs = append(objs, lval)
	}

	rval, err := co.calcRight(ec)
	if err != nil {
		return nil, err
	}
	objs = append(objs, rval)

	return objs, nil
}

type funcInvokeOp struct {
	doubleOp

	params int
}

func (fio *funcInvokeOp) Category() int {
	return opObject
}

func (fio *funcInvokeOp) String() string {
	return "func("
}

func (fio *funcInvokeOp) Priority() int {
	return 1
}

func (fio *funcInvokeOp) Wrap(op operator, rpn cog.Queue[any]) error {
	var ok bool

	if fio.params <= 0 {
		fio.left, ok = rpn.Poll()
		if !ok {
			return errors.New("gel: missing function owner")
		}
	} else {
		fio.right, ok = rpn.Poll()
		if !ok {
			return errors.New("gel: missing function arguments")
		}

		fio.left, ok = rpn.Poll()
		if !ok {
			return errors.New("gel: missing function owner")
		}
	}
	return nil
}

func (fio *funcInvokeOp) Calculate(ec elCtx) (any, error) {
	args, err := fio.fetchParams(ec)
	if err != nil {
		return nil, err
	}

	if c, ok := fio.left.(invoker); ok {
		return c.Invoke(ec, args)
	}

	return nil, fmt.Errorf("gel: left is not a function '%T'", fio.left)
}

func (fio *funcInvokeOp) fetchParams(ec elCtx) (args []any, err error) {
	if fio.params <= 0 {
		args = []any{}
		return
	}

	var val any
	if fio.right != nil {
		if co, ok := fio.right.(*commaOp); ok {
			val, err = co.Calculate(ec)
			if err != nil {
				return
			}

			args = val.([]any)
		} else {
			val, err = fio.calcRight(ec)
			if err != nil {
				return
			}

			args = append(args, val)
		}
	}

	for i, a := range args {
		if op, ok := a.(operator); ok {
			val, err = op.Calculate(ec)
			if err != nil {
				return
			}
			args[i] = val
		}
	}

	return
}

type funcEndOp struct {
	left any
}

func (feo *funcEndOp) Category() int {
	return opObject
}

func (feo *funcEndOp) String() string {
	return "func)"
}

func (feo *funcEndOp) Operands() int {
	return 1
}

func (feo *funcEndOp) Priority() int {
	return 1
}

func (feo *funcEndOp) Wrap(op operator, rpn cog.Queue[any]) error {
	var ok bool

	feo.left, ok = rpn.Poll()
	if !ok {
		return errors.New("gel: missing function operator")
	}
	return nil
}

func (feo *funcEndOp) Calculate(ec elCtx) (any, error) {
	if fio, ok := feo.left.(*funcInvokeOp); ok {
		return fio.Calculate(ec)
	}

	return nil, fmt.Errorf("gel: left is not a function operator '%v'", feo.left)
}
