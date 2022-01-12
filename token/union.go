package token

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type unionToken struct {
	arguments []interface{}
}

func (token *unionToken) Type() string {
	return "union"
}

func (token *unionToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	arguments := token.arguments
	if len(arguments) == 0 {
		return nil, getInvalidTokenArgumentNilError(token.Type(), reflect.Array, reflect.Slice)
	}

	keys := make([]string, 0)
	indices := make([]int64, 0)

	for _, arg := range arguments {
		if argToken, ok := arg.(Token); ok {
			result, err := argToken.Apply(root, current, nil)
			if err != nil {
				return nil, getInvalidTokenError(token.Type(), err)
			}
			arg = result
		}

		if arg == nil {
			return nil, getInvalidTokenArgumentNilError(token.Type(), reflect.Int, reflect.String)
		}

		if strArg, ok := arg.(string); ok {
			keys = append(keys, strArg)
			if len(indices) > 0 {
				return nil, getInvalidTokenArgumentError(token.Type(), reflect.String, reflect.Int)
			}
		} else if intArg, ok := isInteger(arg); ok {
			indices = append(indices, intArg)
			if len(keys) > 0 {
				return nil, getInvalidTokenArgumentError(token.Type(), reflect.Int, reflect.String)
			}
		} else {
			argType := reflect.TypeOf(arg)
			return nil, getInvalidTokenArgumentError(token.Type(), argType.Kind(), reflect.Int, reflect.String)
		}
	}

	var unionValue interface{}

	if len(keys) > 0 {
		var err error
		unionValue, err = getUnionByKey(token, current, keys)
		if err != nil {
			if isInvalidTokenError(err) {
				return nil, err
			}
			return nil, getInvalidTokenError(token.Type(), err)
		}
	} else if len(indices) > 0 {
		var err error
		unionValue, err = getUnionByIndex(token, current, indices)
		if err != nil {
			if isInvalidTokenError(err) {
				return nil, err
			}
			return nil, getInvalidTokenError(token.Type(), err)
		}
	}

	if strValue, ok := unionValue.(string); ok {
		if len(next) > 0 {
			return next[0].Apply(root, strValue, next[1:])
		}
		return strValue, nil
	}

	elements := unionValue.([]interface{})

	if len(next) > 0 {
		nextToken := next[0]
		futureTokens := next[1:]

		if indexToken, ok := nextToken.(*indexToken); ok {
			// if next is asking for specific index
			return indexToken.Apply(current, elements, futureTokens)
		}
		// any other token type
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

func getUnionByKey(token Token, obj interface{}, keys []string) ([]interface{}, error) {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return nil, getInvalidTokenTargetNilError(token.Type(), reflect.Map)
	}

	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	switch objType.Kind() {
	case reflect.Map:
		objValue := reflect.ValueOf(obj)
		mapKeys := objValue.MapKeys()

		elements := make([]interface{}, 0)

		for _, key := range mapKeys {
			if keyMap[key.String()] {
				delete(keyMap, key.String())
				elements = append(elements, objValue.MapIndex(key).Interface())
			}
		}

		if len(keyMap) > 0 {
			remaining := make([]string, 0)
			for key := range keyMap {
				remaining = append(remaining, key)
			}
			sort.Strings(remaining)
			return nil, getInvalidTokenKeyNotFoundError(token.Type(), strings.Join(remaining, ","))
		}

		return elements, nil
	default:
		return nil, getInvalidTokenTargetError(
			token.Type(),
			objType.Kind(),
			reflect.Map,
		)
	}
}

func getUnionByIndex(token Token, obj interface{}, indices []int64) (interface{}, error) {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return nil, getInvalidTokenTargetNilError(
			token.Type(),
			reflect.Array, reflect.Map, reflect.Slice, reflect.String,
		)
	}

	var objValue reflect.Value
	var length int64
	var mapKeys []reflect.Value
	isString := false

	switch objType.Kind() {
	case reflect.Map:
		objValue = reflect.ValueOf(obj)
		length = int64(objValue.Len())
		mapKeys = objValue.MapKeys()

		sort.SliceStable(mapKeys, func(i, j int) bool {
			one := mapKeys[i]
			two := mapKeys[j]

			return one.String() < two.String()
		})
		break
	case reflect.String:
		isString = true
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		objValue = reflect.ValueOf(obj)
		length = int64(objValue.Len())
		mapKeys = nil
		break
	default:
		return nil, getInvalidTokenTargetError(
			token.Type(),
			objType.Kind(),
			reflect.Array, reflect.Map, reflect.Slice, reflect.String,
		)
	}

	values := make([]interface{}, 0)
	substring := ""

	for _, idx := range indices {
		if idx < 0 {
			idx = length + idx
		}
		if idx < 0 || idx >= length {
			return nil, getInvalidTokenOutOfRangeError(token.Type())
		}

		if mapKeys != nil {
			key := mapKeys[idx]
			values = append(values, objValue.MapIndex(key).Interface())
		} else if isString {
			value := objValue.Index(int(idx)).Interface()
			if u, ok := value.(uint8); ok {
				substring += fmt.Sprintf("%c", u)
			}
		} else {
			values = append(values, objValue.Index(int(idx)).Interface())
		}
	}

	if isString {
		return substring, nil
	}

	return values, nil
}
