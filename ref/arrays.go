package ref

import (
	"errors"
	"fmt"
	"reflect"
)

func ArrayLen(a any) (int, error) {
	if a == nil {
		return 0, nil
	}

	rv := reflect.ValueOf(a)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		return rv.Len(), nil
	default:
		return 0, errors.New("ArrayLen(): invalid array or slice")
	}
}

// ArrayGet getting value from array or slice by index
// usage:
//
//	a := [][]string{
//	    { "0,0", "0,1" },
//	    { "1,0", "1,1" },
//	}
//
// {{ArrayGet a 0 1 }} // return "0,1"
func ArrayGet(a any, idxs ...int) (any, error) {
	if a == nil || len(idxs) == 0 {
		return nil, nil
	}

	rv := reflect.ValueOf(a)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		if idxs[0] < 0 || idxs[0] > rv.Len() {
			return nil, errors.New("ArrayGet(): invalid index")
		}

		r := rv.Index(idxs[0]).Interface()
		// if there is more keys, handle this recursively
		if len(idxs) > 1 {
			return ArrayGet(r, idxs[1:]...)
		}
		return r, nil
	default:
		return nil, errors.New("ArrayGet(): invalid array or slice")
	}
}

// ArraySet set value to the array or slice by index
func ArraySet(a any, i int, v any) (any, error) {
	rt := reflect.TypeOf(a)

	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		rv := reflect.ValueOf(a)
		if i < 0 || i > rv.Len() {
			return nil, errors.New("ArraySet(): invalid index")
		}

		iv := rv.Index(i)

		vv := reflect.ValueOf(v)
		vt := reflect.TypeOf(v)
		if vt.Kind() != rt.Elem().Kind() {
			cv, err := Convert(v, rt.Elem())
			if err != nil {
				return nil, fmt.Errorf("ArraySet(): invalid value type - %w", err)
			}

			vv = reflect.ValueOf(cv)
		}

		iv.Set(vv)
		return nil, nil

	default:
		return nil, errors.New("ArrayGet(): invalid array or slice")
	}
}
