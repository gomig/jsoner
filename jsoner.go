package jsoner

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/gomig/utils"
)

type JsonerOption struct {
	Only   []string
	Ignore []string
}

// Jsoner marshal json with dynamic ignores
//
// you must pass path of json data as key and options as value.
// e.g. option to marshal isbn and title of author book, map[string]JsonerOption{"author.book": JsonerOption{Only: []string{"isbn", "title"}}}
// "." path used for root
func Jsoner(v any, opt map[string]JsonerOption) ([]byte, error) {
	return json.Marshal(__mapper(v, ".", opt))
}
func JsonerIndent(v any, indent string, opt map[string]JsonerOption) ([]byte, error) {
	return json.MarshalIndent(__mapper(v, ".", opt), "", indent)
}

func __mapper(v any, _path string, opt map[string]JsonerOption) any {
	_option := func(path string) ([]string, []string) {
		if path != "." {
			path = strings.TrimLeft(path, ".")
		}
		if v, ok := opt[path]; ok {
			return v.Only, v.Ignore
		}
		return []string{}, []string{}
	}
	_isPointer := func(kind reflect.Kind) bool {
		return reflect.Ptr == kind
	}
	_isSimple := func(kind reflect.Kind) bool {
		return !utils.Contains([]reflect.Kind{reflect.Array, reflect.Slice, reflect.Map, reflect.Struct}, kind)
	}

	// Check nil
	if v == nil {
		return v
	}

	val := reflect.ValueOf(v)

	// Get value of pointer
	if _isPointer(val.Kind()) {
		if val.IsNil() {
			return nil
		}
		return __mapper(val.Elem().Interface(), _path, opt)
	}

	// return non-struct field
	if _isSimple(val.Kind()) {
		return v
	}

	switch val.Kind() {
	case reflect.Array:
	case reflect.Slice:
		_sliceV := make([]any, 0)
		for i := 0; i < val.Len(); i++ {
			_sliceV = append(_sliceV, __mapper(val.Index(i).Interface(), _path, opt))
		}
		return _sliceV
	case reflect.Map:
		res := make(map[string]any)
		only, ignore := _option(_path)
		for _, k := range val.MapKeys() {
			// check ignore field
			name := fmt.Sprint(k.Interface())
			if (len(only) > 0 && !utils.Contains(only, name)) || (len(ignore) > 0 && utils.Contains(ignore, name)) {
				continue
			}
			// Parse value
			value := val.MapIndex(k)
			res[name] = __mapper(value.Interface(), utils.If(_path == ".", "."+name, _path+"."+name), opt)
		}
		return res
	case reflect.Struct:
		bytes, _ := json.Marshal(val.Interface())
		var res map[string]any
		json.Unmarshal(bytes, &res)
		return __mapper(res, _path, opt)
	}
	return nil
}
