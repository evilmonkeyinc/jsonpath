package token

import (
	"reflect"
)

func newWildcardToken() *wildcardToken {
	return &wildcardToken{}
}

type wildcardToken struct {
}

func (token *wildcardToken) String() string {
	return "[*]"
}

func (token *wildcardToken) Type() string {
	return "wildcard"
}

func (token *wildcardToken) Apply(root, current interface{}, next []Token) (interface{}, error) {

	elements := make([]interface{}, 0)

	objType, objVal := getTypeAndValue(current)
	if objType == nil {
		return nil, getInvalidTokenTargetNilError(
			token.Type(),
			reflect.Array, reflect.Map, reflect.Slice,
		)
	}

	switch objType.Kind() {
	case reflect.Map:
		keys := objVal.MapKeys()
		sortMapKeys(keys)
		for _, kv := range keys {
			value := objVal.MapIndex(kv).Interface()
			elements = append(elements, value)
		}
		break
	case reflect.Array, reflect.Slice:
		length := objVal.Len()
		for i := 0; i < length; i++ {
			value := objVal.Index(i).Interface()
			elements = append(elements, value)
		}
	case reflect.Struct:
		fields := getStructFields(objVal, true)
		for _, field := range fields {
			value := objVal.FieldByName(field.Name).Interface()
			elements = append(elements, value)
		}
		break
	default:
		return nil, getInvalidTokenTargetError(
			token.Type(),
			objType.Kind(),
			reflect.Array, reflect.Map, reflect.Slice,
		)
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
