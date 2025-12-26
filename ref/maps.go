package ref

import (
	"errors"
	"fmt"
	"reflect"
)

// MapGet getting value from map by keys
// usage:
//
//	m := map[string]any{
//	    "a": 1,
//	    "1": map[string]float64{
//	        "c": 4,
//	    },
//	}
//
// {{MapGet m "a" }} // return 1
// {{MapGet m 1 "c" }} // return 4
func MapGet(dict any, keys ...any) (any, error) {
	if len(keys) == 0 {
		return nil, errors.New("ref: missing argument key")
	}

	mv := reflect.ValueOf(dict)
	if mv.Kind() != reflect.Map {
		return nil, fmt.Errorf("ref: %T is not a map", dict)
	}

	val, err := mapGet(mv, keys[0])
	if err != nil {
		return val, err
	}

	// if there is more keys, handle this recursively
	if len(keys) > 1 && val != nil {
		return MapGet(val, keys[1:]...)
	}
	return val, nil
}

func mapGet(mv reflect.Value, key any) (any, error) {
	mt := mv.Type()
	mk := mt.Key()

	// check whether keys[0] type equals to dict key type
	// if they are different, make conversion
	kv := reflect.ValueOf(key)
	if kv.Type() != mk {
		cv, err := CastTo(key, mk)
		if err != nil {
			return nil, fmt.Errorf("ref: invalid map key type - %w", err)
		}

		kv = reflect.ValueOf(cv)
	}

	vv := mv.MapIndex(kv)
	if !vv.IsValid() {
		return nil, nil
	}

	return vv.Interface(), nil
}

// MapSet setting value to the map
func MapSet(dict, key, value any) error {
	mv := reflect.ValueOf(dict)
	if mv.Kind() != reflect.Map {
		return fmt.Errorf("ref: %T is not a map", dict)
	}

	return mapSet(mv, key, value)
}

func mapSet(mv reflect.Value, key, val any) error {
	mt := mv.Type()
	mk, me := mt.Key(), mt.Elem()

	kv := reflect.ValueOf(key)
	if kv.Type() != mk {
		cv, err := CastTo(key, mk)
		if err != nil {
			return fmt.Errorf("ref: invalid map key type - %w", err)
		}
		kv = reflect.ValueOf(cv)
	}

	vv := reflect.ValueOf(val)
	if vv.Type() != me {
		cv, err := CastTo(val, me)
		if err != nil {
			return fmt.Errorf("ref: invalid map value type - %w", err)
		}
		vv = reflect.ValueOf(cv)
	}

	mv.SetMapIndex(kv, vv)
	return nil
}
