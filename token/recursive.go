package token

import (
	"reflect"
)

type recursiveToken struct {
}

func (token *recursiveToken) Type() string {
	return "recursive"
}

func (token *recursiveToken) Apply(root, current interface{}, next []Token) (interface{}, error) {

	elements := flatten(current)

	if len(next) > 0 {
		nextToken := next[0]
		futureTokens := next[1:]

		results := make([]interface{}, 0)

		for _, item := range elements {
			result, _ := nextToken.Apply(root, item, futureTokens)
			objType := reflect.TypeOf(result)
			if objType == nil {
				continue
			}

			switch objType.Kind() {
			case reflect.Array, reflect.Slice:
				objVal := reflect.ValueOf(result)
				length := objVal.Len()
				for i := 0; i < length; i++ {
					results = append(results, objVal.Index(i).Interface())
				}
				break
			default:
				results = append(results, result)
				break
			}

		}

		return results, nil
	}
	return elements, nil
}

func flatten(obj interface{}) []interface{} {
	slice := make([]interface{}, 0)

	objType := reflect.TypeOf(obj)
	if objType == nil {
		return slice
	}
	slice = append(slice, obj)

	elements := make([]interface{}, 0)
	switch objType.Kind() {
	case reflect.Map:
		objVal := reflect.ValueOf(obj)
		keys := objVal.MapKeys()
		for _, kv := range keys {
			value := objVal.MapIndex(kv).Interface()
			elements = append(elements, value)
		}
	case reflect.Array, reflect.Slice:
		objVal := reflect.ValueOf(obj)
		length := objVal.Len()
		for i := 0; i < length; i++ {
			value := objVal.Index(i).Interface()
			elements = append(elements, value)
		}
	default:
		break
	}

	if len(elements) > 0 {
		for _, sObj := range elements {
			slice = append(slice, flatten(sObj)...)
		}
	}

	return slice
}
