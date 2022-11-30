package funcs

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/pandafw/pango/ref"
)

// Array returns a []any{args[0], args[1], ...}
func Array(args ...any) []any {
	return args
}

// Map returns a map[string]any{kvs[0]: kvs[1], kvs[2]: kvs[3], ...}
func Map(kvs ...any) (map[string]any, error) {
	if len(kvs)&1 != 0 {
		return nil, errors.New("Map(): invalid arguments")
	}

	dict := make(map[string]any, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		key, ok := kvs[i].(string)
		if !ok {
			return nil, errors.New("Map(): keys must be strings")
		}
		dict[key] = kvs[i+1]
	}
	return dict, nil
}

// MapGet getting value from map by keys
// usage:
//
//	Data["m"] = map[string]any{
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
		return nil, fmt.Errorf("MapGet(): invalid map")
	}

	// check whether keys[0] type equals to dict key type
	// if they are different, make conversion
	kv := reflect.ValueOf(keys[0])
	kt := reflect.TypeOf(keys[0])
	if kt.Kind() != mt.Key().Kind() {
		cv, err := ref.Convert(keys[0], mt.Key())
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
		return nil, fmt.Errorf("MapSet(): invalid map")
	}

	kv := reflect.ValueOf(key)
	kt := reflect.TypeOf(key)
	if kt.Kind() != mt.Key().Kind() {
		cv, err := ref.Convert(key, mt.Key())
		if err != nil {
			return nil, fmt.Errorf("MapSet(): invalid key type - %w", err)
		}

		kv = reflect.ValueOf(cv)
	}

	vv := reflect.ValueOf(val)
	vt := reflect.TypeOf(val)
	if vt.Kind() != mt.Elem().Kind() {
		cv, err := ref.Convert(val, mt.Elem())
		if err != nil {
			return nil, fmt.Errorf("MapSet(): invalid value type - %w", err)
		}

		vv = reflect.ValueOf(cv)
	}

	mv := reflect.ValueOf(dict)
	mv.SetMapIndex(kv, vv)
	return nil, nil
}
