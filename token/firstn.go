package token

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/evilmokeyinc/jsonpath/errors"
)

type firstNToken struct {
	number interface{}
}

// TODO : get rid of range errors

func (token *firstNToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	to := int64(0)

	if intValue, ok := isInteger(token.number); ok {
		to = intValue
	} else if token, ok := token.number.(*expressionToken); ok {
		evaluate, err := token.Apply(root, current, nil)
		if err != nil {
			// TODO : wrap error?
			return nil, err
		}

		if intValue, ok := isInteger(evaluate); ok {
			to = intValue
		} else {
			return nil, errors.ErrInvalidParameterInteger
		}
	} else {
		return nil, errors.ErrInvalidParameterInteger
	}

	var objValue reflect.Value
	var length int64
	var mapKeys []reflect.Value
	isString := false

	objType := reflect.TypeOf(current)
	if objType == nil {
		return nil, errors.ErrGetRangeFromNilArray
	}

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
		break
	case reflect.String:
		isString = true
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		objValue = reflect.ValueOf(current)
		length = int64(objValue.Len())
		mapKeys = nil
		break
	default:
		return nil, errors.ErrInvalidObjectArrayMapOrString
	}

	if to < 0 {
		to = length + to
	}

	if to < 0 || to >= length {
		return nil, errors.ErrIndexOutOfRange
	}

	elements := make([]interface{}, 0)
	substring := ""

	if mapKeys != nil {
		for i := int64(0); i < to; i++ {
			key := mapKeys[i]
			elements = append(elements, objValue.MapIndex(key).Interface())
		}
	} else if isString {
		for i := int64(0); i < to; i++ {
			value := objValue.Index(int(i)).Uint()
			substring += fmt.Sprintf("%c", value)
		}
	} else {
		for i := int64(0); i < to; i++ {
			elements = append(elements, objValue.Index(int(i)).Interface())
		}
	}

	if substring != "" {
		if len(next) > 0 {
			return next[0].Apply(root, substring, next[1:])
		}
		return substring, nil
	}

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
