package token

import (
	"reflect"

	"github.com/evilmokeyinc/jsonpath/errors"
)

type wildcardToken struct {
}

func (token *wildcardToken) Apply(root, current interface{}, next []Token) (interface{}, error) {

	elements := make([]interface{}, 0)

	objType := reflect.TypeOf(current)
	if objType == nil {
		return nil, errors.ErrGetElementsFromNilObject
	}
	switch objType.Kind() {
	case reflect.Map:
		objVal := reflect.ValueOf(current)
		keys := objVal.MapKeys()
		for _, kv := range keys {
			value := objVal.MapIndex(kv).Interface()
			elements = append(elements, value)
		}
	case reflect.Array, reflect.Slice:
		objVal := reflect.ValueOf(current)
		length := objVal.Len()
		for i := 0; i < length; i++ {
			value := objVal.Index(i).Interface()
			elements = append(elements, value)
		}
	default:
		return nil, errors.ErrInvalidObjectArrayOrMap
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
