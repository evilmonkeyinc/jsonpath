package token

import (
	"fmt"
	"reflect"
	"sort"
)

type sliceToken struct {
	number interface{}
}

func (token *sliceToken) String() string {
	return fmt.Sprintf("[:%v]", token.number)
}

func (token *sliceToken) Type() string {
	return "slice"
}

func (token *sliceToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	to := int64(0)

	if token.number == nil {
		return nil, getInvalidTokenArgumentNilError(token.Type(), reflect.Int)
	}

	if intValue, ok := isInteger(token.number); ok {
		to = intValue
	} else if expression, ok := token.number.(*expressionToken); ok {
		evaluate, err := expression.Apply(root, current, nil)
		if err != nil {
			return nil, getInvalidTokenError(token.Type(), err)
		}

		if intValue, ok := isInteger(evaluate); ok {
			to = intValue
		} else {
			kind := reflect.TypeOf(evaluate).Kind()
			err := getUnexpectedExpressionResultError(kind, reflect.Int)
			return nil, getInvalidTokenError(token.Type(), err)
		}
	} else {
		kind := reflect.TypeOf(token.number).Kind()
		return nil, getInvalidTokenArgumentError(token.Type(), kind, reflect.Int)
	}

	var objValue reflect.Value
	var length int64
	var mapKeys []reflect.Value
	isString := false

	objType := reflect.TypeOf(current)
	if objType == nil {
		return nil, getInvalidTokenTargetNilError(
			token.Type(),
			reflect.Array, reflect.Map, reflect.Slice, reflect.String,
		)
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
		return nil, getInvalidTokenTargetError(
			token.Type(),
			objType.Kind(),
			reflect.Array, reflect.Map, reflect.Slice, reflect.String,
		)
	}

	if to < 0 {
		to = length + to
	}

	if to < 0 || to >= length {
		return nil, getInvalidTokenOutOfRangeError(token.Type())
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
