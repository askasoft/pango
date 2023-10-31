package funcs

import "errors"

// Array returns a []any{args[0], args[1], ...}
func Array(args ...any) []any {
	return args
}

// Strings returns a []string{args[0], args[1], ...}
func Strings(args ...string) []string {
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
