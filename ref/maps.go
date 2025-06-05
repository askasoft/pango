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
		return nil, errors.New("MapGet(): missing argument key")
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
	// check whether keys[0] type equals to dict key type
	// if they are different, make conversion
	mt := mv.Type()
	kv := reflect.ValueOf(key)
	kt := reflect.TypeOf(key)
	if kt != mt.Key() {
		cv, err := CastTo(key, mt.Key())
		if err != nil {
			return nil, fmt.Errorf("MapGet(): invalid key type - %w", err)
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
func MapSet(dict any, key, val any) (any, error) {
	mv := reflect.ValueOf(dict)
	if mv.Kind() != reflect.Map {
		return nil, fmt.Errorf("ref: %T is not a map", dict)
	}

	return mapSet(mv, key, val)
}

func mapSet(mv reflect.Value, key, val any) (any, error) {
	mt := mv.Type()
	kv := reflect.ValueOf(key)
	kt := reflect.TypeOf(key)
	if kt != mt.Key() {
		cv, err := CastTo(key, mt.Key())
		if err != nil {
			return nil, fmt.Errorf("MapSet(): invalid key type - %w", err)
		}

		kv = reflect.ValueOf(cv)
	}

	vv := reflect.ValueOf(val)
	vt := reflect.TypeOf(val)
	if vt != mt.Elem() {
		cv, err := CastTo(val, mt.Elem())
		if err != nil {
			return nil, fmt.Errorf("MapSet(): invalid value type - %w", err)
		}

		vv = reflect.ValueOf(cv)
	}

	mv.SetMapIndex(kv, vv)
	return nil, nil
}
