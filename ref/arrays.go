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
		return nil, fmt.Errorf("ref: %T is not a array or slice", a)
	}
}

// ArraySet set value to the array or slice by index
func ArraySet(a any, i int, v any) error {
	av := reflect.ValueOf(a)

	switch av.Kind() {
	case reflect.Slice, reflect.Array:
		if i < 0 || i > av.Len() {
			return fmt.Errorf("ref: index %d out of bounds [0:%d]", i, av.Len())
		}

		et := av.Type().Elem()
		vv := reflect.ValueOf(v)
		if vv.Type() != et {
			cv, err := CastTo(v, et)
			if err != nil {
				return err
			}

			vv = reflect.ValueOf(cv)
		}

		av.Index(i).Set(vv)
		return nil
	default:
		return fmt.Errorf("ref: %T is not a array or slice", a)
	}
}

// ToSlice convert array to slice.
func ToSlice(a any) (any, error) {
	av := reflect.ValueOf(a)

	switch av.Kind() {
	case reflect.Array:
		st := reflect.SliceOf(av.Type().Elem())
		sv := reflect.MakeSlice(st, av.Len(), av.Cap())
		reflect.Copy(sv, av)
		return sv.Interface(), nil
	case reflect.Slice:
		return a, nil
	default:
		return a, fmt.Errorf("ref: %T is not a array or slice", a)
	}
}

// SliceAdd add values to the slice `a`
// if `a` is a array, we convert it to slice and add `vs`.
func SliceAdd(a any, vs ...any) (any, error) {
	av := reflect.ValueOf(a)

	switch av.Kind() {
	case reflect.Array:
		s, err := ToSlice(a)
		if err != nil {
			return a, err
		}
		return SliceAdd(s, vs...)
	case reflect.Slice:
		if len(vs) == 0 {
			return a, nil
		}

		et := av.Type().Elem()
		rvs := make([]reflect.Value, len(vs))
		for i, v := range vs {
			vv := reflect.ValueOf(v)
			if vv.Type() != et {
				cv, err := CastTo(v, et)
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
		return a, fmt.Errorf("ref: %T is not a array or slice", a)
	}
}
