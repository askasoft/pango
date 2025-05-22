package ref

import (
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
	if dict == nil || len(keys) == 0 {
		return nil, nil
	}

	mt := reflect.TypeOf(dict)
	if mt.Kind() != reflect.Map {
		return nil, fmt.Errorf("MapGet(): invalid map - %v", mt)
	}

	// check whether keys[0] type equals to dict key type
	// if they are different, make conversion
	kv := reflect.ValueOf(keys[0])
	kt := reflect.TypeOf(keys[0])
	if kt != mt.Key() {
		cv, err := CastTo(keys[0], mt.Key())
		if err != nil {
			return nil, fmt.Errorf("MapGet(): invalid key type - %w", err)
		}

		kv = reflect.ValueOf(cv)
	}

	mv := reflect.ValueOf(dict)
	vv := mv.MapIndex(kv)
	if !vv.IsValid() {
		return nil, nil
	}

	val := vv.Interface()

	// if there is more keys, handle this recursively
	if len(keys) > 1 {
		return MapGet(val, keys[1:]...)
	}
	return val, nil
}

// MapSet setting value to the map
func MapSet(dict any, key, val any) (any, error) {
	mt := reflect.TypeOf(dict)
	if mt.Kind() != reflect.Map {
		return nil, fmt.Errorf("MapSet(): invalid map type - %v", mt)
	}

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

	mv := reflect.ValueOf(dict)
	mv.SetMapIndex(kv, vv)
	return nil, nil
}
