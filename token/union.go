package token

import (
	"reflect"
	"sort"
	"strings"

	"github.com/evilmokeyinc/jsonpath/errors"
)

type unionToken struct {
	arguments []interface{}
}

func (token *unionToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	arguments := token.arguments

	for idx, arg := range arguments {
		if token, ok := arg.(Token); ok {
			result, err := token.Apply(root, current, nil)
			if err != nil {
				return nil, err
			}
			arguments[idx] = result
		}
	}

	elements, err := getUnion(current, arguments)
	if err != nil {
		return nil, err
	}

	if len(next) > 0 {
		nextToken := next[0]
		futureTokens := next[1:]

		results := make([]interface{}, 0)

		for _, item := range elements {
			result, _ := nextToken.Apply(root, item, futureTokens)
			if result != nil {
				results = append(results, result)
			}
		}

		return results, nil
	}

	return elements, nil
}

func getUnion(obj interface{}, keys []interface{}) ([]interface{}, error) {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return nil, errors.ErrGetUnionFromNilObject
	}
	switch objType.Kind() {
	case reflect.Map:
		keyMap := make(map[string]bool)
		for _, key := range keys {
			strKey, ok := key.(string)
			if !ok {
				return nil, errors.GetInvalidParameterError("expected string keys")
			}
			keyMap[strKey] = true
		}

		values := make([]interface{}, 0)
		objVal := reflect.ValueOf(obj)
		mapKeys := objVal.MapKeys()
		for _, kv := range mapKeys {
			if keyMap[kv.String()] {
				objVal := objVal.MapIndex(kv).Interface()
				values = append(values, objVal)
				delete(keyMap, kv.String())
			}
		}

		if len(keyMap) > 0 {
			keys := make([]string, 0)
			for key := range keyMap {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			return nil, errors.GetKeyNotFoundError(strings.Join(keys, ","))
		}

		return values, nil
	case reflect.Array, reflect.Slice:
		objVal := reflect.ValueOf(obj)
		length := objVal.Len()

		values := make([]interface{}, 0)

		for _, key := range keys {
			idx, ok := key.(int)
			if !ok {
				return nil, errors.GetInvalidParameterError("expected int keys")
			}

			if idx >= 0 {
				if idx >= length {
					return nil, errors.ErrIndexOutOfRange
				}
			} else {
				idx = length + idx
				if idx < 0 {
					return nil, errors.ErrIndexOutOfRange
				}
			}
			values = append(values, objVal.Index(idx).Interface())
		}
		return values, nil
	default:
		return nil, errors.ErrInvalidObjectMapOrSlice
	}
}
