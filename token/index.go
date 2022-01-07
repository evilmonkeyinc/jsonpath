package token

import (
	"reflect"

	"github.com/evilmokeyinc/jsonpath/errors"
)

type indexToken struct {
	index int
}

func (token *indexToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	idx := token.index

	objType := reflect.TypeOf(current)
	if objType == nil {
		return nil, errors.ErrGetIndexFromNilArray
	}
	switch objType.Kind() {
	case reflect.Array, reflect.Slice:
		objVal := reflect.ValueOf(current)
		length := objVal.Len()
		if idx < 0 {
			idx = length + idx
		}

		if idx < 0 || idx >= length {
			return nil, errors.ErrIndexOutOfRange
		}
		value := objVal.Index(idx).Interface()

		if len(next) > 0 {
			return next[0].Apply(root, value, next[1:])
		}
		return value, nil
	default:
		return nil, errors.ErrInvalidObjectArray
	}
}
