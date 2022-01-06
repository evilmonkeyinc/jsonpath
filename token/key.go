package token

import (
	"reflect"

	"github.com/evilmokeyinc/jsonpath/errors"
)

type keyToken struct {
	key string
}

func (token *keyToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	objType := reflect.TypeOf(current)
	if objType == nil {
		return nil, errors.ErrGetKeyFromNilMap
	}
	switch objType.Kind() {
	case reflect.Map:
		objVal := reflect.ValueOf(current)
		keys := objVal.MapKeys()
		for _, kv := range keys {
			if kv.String() == token.key {
				value := objVal.MapIndex(kv).Interface()

				if len(next) > 0 {
					return next[0].Apply(root, value, next[1:])
				}

				return value, nil
			}
		}
		return nil, errors.GetKeyNotFoundError(token.key)
	default:
		return nil, errors.ErrInvalidObjectMap
	}
}
