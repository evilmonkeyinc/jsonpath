package token

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/evilmokeyinc/jsonpath/errors"
)

type indexToken struct {
	index int64
}

func (token *indexToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	idx := token.index

	objType := reflect.TypeOf(current)
	if objType == nil {
		return nil, errors.ErrGetIndexFromNilArray
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
		return nil, errors.ErrInvalidObjectArrayMapOrString
	}

	if idx < 0 {
		idx = length + idx
	}

	if idx < 0 || idx >= length {
		return nil, errors.ErrIndexOutOfRange
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