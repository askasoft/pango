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

	av := reflect.ValueOf(a)
	switch av.Kind() {
	case reflect.Slice, reflect.Array:
		return av.Len(), nil
	default:
		return 0, fmt.Errorf("ref: %T is not a array or slice", a)
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
	if len(idxs) == 0 {
		return nil, errors.New("ref: missing argument index")
	}

	av := reflect.ValueOf(a)
	switch av.Kind() {
	case reflect.Slice, reflect.Array:
		idx := idxs[0]
		if idx < 0 || idx > av.Len() {
			return nil, fmt.Errorf("ref: index %d out of bounds [0:%d]", idx, av.Len())
		}

		val := av.Index(idx).Interface()

		// if there is more keys, handle this recursively
		if len(idxs) > 1 {
			return ArrayGet(val, idxs[1:]...)
		}
		return val, nil
	default:
		return 0, fmt.Errorf("ref: %T is not a array or slice", a)
	}
}

// ArraySet set value to the array or slice by index
func ArraySet(a any, i int, v any) (any, error) {
	av := reflect.ValueOf(a)
	at := av.Type()

	switch at.Kind() {
	case reflect.Slice, reflect.Array:
		if i < 0 || i > av.Len() {
			return nil, fmt.Errorf("ref: index %d out of bounds [0:%d]", i, av.Len())
		}

		vv := reflect.ValueOf(v)
		if vv.Type() != at.Elem() {
			cv, err := CastTo(v, at.Elem())
			if err != nil {
				return nil, err
			}

			vv = reflect.ValueOf(cv)
		}

		av.Index(i).Set(vv)
		return nil, nil
	default:
		return 0, fmt.Errorf("ref: %T is not a array or slice", a)
	}
}

// SliceAdd add values to the slice
func SliceAdd(a any, vs ...any) (any, error) {
	av := reflect.ValueOf(a)
	at := av.Type()

	switch at.Kind() {
	case reflect.Slice:
		if len(vs) == 0 {
			return a, nil
		}

		rvs := make([]reflect.Value, len(vs))
		for i, v := range vs {
			vv := reflect.ValueOf(v)
			if vv.Type() != at.Elem() {
				cv, err := CastTo(v, at.Elem())
				if err != nil {
					return nil, err
				}
				vv = reflect.ValueOf(cv)
			}
			rvs[i] = vv
		}

		av = reflect.Append(av, rvs...)
		return av.Interface(), nil
	default:
		return 0, fmt.Errorf("ref: %T is not a slice", a)
	}
}
