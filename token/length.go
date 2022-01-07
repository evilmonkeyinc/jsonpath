package token

import (
	"reflect"

	"github.com/evilmokeyinc/jsonpath/errors"
)

type lengthToken struct {
}

func (token *lengthToken) Apply(root, current interface{}, next []Token) (interface{}, error) {

	objType := reflect.TypeOf(current)
	if objType == nil {
		return nil, errors.ErrGetElementsFromNilObject
	}

	switch objType.Kind() {
	case reflect.Map:

		objVal := reflect.ValueOf(current)
		current = objVal.Len()

		keys := objVal.MapKeys()
		for _, kv := range keys {
			if kv.String() == "length" {
				current = objVal.MapIndex(kv).Interface()
			}
		}
	case reflect.Array, reflect.Slice, reflect.String:
		objVal := reflect.ValueOf(current)
		current = objVal.Len()
	default:
		return nil, errors.ErrInvalidObjectArrayMapOrString
	}

	if len(next) > 0 {
		return next[0].Apply(root, current, next[1:])
	}
	return current, nil
}
