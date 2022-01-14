package token

import (
	"fmt"
	"reflect"
	"sort"
)

type indexToken struct {
	index int64
}

func (token *indexToken) String() string {
	return fmt.Sprintf("[%d]", token.index)
}

func (token *indexToken) Type() string {
	return "index"
}

func (token *indexToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	idx := token.index

	objType := reflect.TypeOf(current)
	if objType == nil {
		return nil, getInvalidTokenTargetNilError(
			token.Type(),
			reflect.Array,
			reflect.Map,
			reflect.Slice,
			reflect.String,
		)
	}

	var objValue reflect.Value
	var length int64
	var mapKeys []reflect.Value
	isString := false

	switch objType.Kind() {
	case reflect.Map:
		objValue = reflect.ValueOf(current)
		length = int64(objValue.Len())
		mapKeys = objValue.MapKeys()

		sort.SliceStable(mapKeys, func(i, j int) bool {
			one := mapKeys[i]
			two := mapKeys[j]

			return one.String() < two.String()
		})
	case reflect.String:
		isString = true
		fallthrough
	case reflect.Array, reflect.Slice:
		objValue = reflect.ValueOf(current)
		length = int64(objValue.Len())
		mapKeys = nil
	default:
		return nil, getInvalidTokenTargetError(
			token.Type(),
			objType.Kind(),
			reflect.Array, reflect.Map, reflect.Slice, reflect.String,
		)
	}

	if idx < 0 {
		idx = length + idx
	}

	if idx < 0 || idx >= length {
		return nil, getInvalidTokenOutOfRangeError(token.Type())
	}

	var value interface{}

	if mapKeys != nil {
		key := mapKeys[idx]
		value = objValue.MapIndex(key).Interface()

	} else if isString {
		value = objValue.Index(int(idx)).Interface()

		if u, ok := value.(uint8); ok {
			value = fmt.Sprintf("%c", u)
		}
	} else {
		value = objValue.Index(int(idx)).Interface()
	}

	if len(next) > 0 {
		return next[0].Apply(root, value, next[1:])
	}
	return value, nil
}
