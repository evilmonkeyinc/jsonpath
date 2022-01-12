package token

import (
	"fmt"
	"reflect"
	"sort"
)

type rangeToken struct {
	from, to, step interface{}
}

func (token *rangeToken) Type() string {
	return "range"
}

func (token *rangeToken) Apply(root, current interface{}, next []Token) (interface{}, error) {

	var fromInt int64
	var toInt, stepInt *int64

	if token.from == nil {
		return nil, getInvalidTokenArgumentNilError(
			token.Type(),
			reflect.Int,
		)
	}

	if script, ok := token.from.(Token); ok {
		result, err := script.Apply(root, current, nil)
		if err != nil {
			return nil, getInvalidTokenError(token.Type(), err)
		}

		if result == nil {
			err := getUnexpectedExpressionResultNilError(reflect.Int)
			return nil, getInvalidTokenError(token.Type(), err)
		}

		if intVal, ok := isInteger(result); ok {
			fromInt = int64(intVal)
		} else {
			kind := reflect.TypeOf(result).Kind()
			err := getUnexpectedExpressionResultError(kind, reflect.Int)
			return nil, getInvalidTokenError(token.Type(), err)
		}
	} else if intVal, ok := isInteger(token.from); ok {
		fromInt = intVal
	} else {
		kind := reflect.TypeOf(token.from).Kind()
		return nil, getInvalidTokenArgumentError(token.Type(), kind, reflect.Int)
	}

	if token.to != nil {
		if script, ok := token.to.(Token); ok {
			result, err := script.Apply(root, current, nil)
			if err != nil {
				return nil, getInvalidTokenError(token.Type(), err)
			}

			if result == nil {
				err := getUnexpectedExpressionResultNilError(reflect.Int)
				return nil, getInvalidTokenError(token.Type(), err)
			}

			if intVal, ok := isInteger(result); ok {
				tmp := int64(intVal)
				toInt = &tmp
			} else {
				kind := reflect.TypeOf(result).Kind()
				err := getUnexpectedExpressionResultError(kind, reflect.Int)
				return nil, getInvalidTokenError(token.Type(), err)
			}
		} else if intVal, ok := isInteger(token.to); ok {
			toInt = &intVal
		} else {
			kind := reflect.TypeOf(token.to).Kind()
			return nil, getInvalidTokenArgumentError(token.Type(), kind, reflect.Int)
		}
	}

	if token.step != nil {
		if script, ok := token.step.(Token); ok {
			result, err := script.Apply(root, current, nil)
			if err != nil {
				return nil, getInvalidTokenError(token.Type(), err)
			}

			if result == nil {
				err := getUnexpectedExpressionResultNilError(reflect.Int)
				return nil, getInvalidTokenError(token.Type(), err)
			}

			if intVal, ok := isInteger(result); ok {
				tmp := int64(intVal)
				stepInt = &tmp
			} else {
				kind := reflect.TypeOf(result).Kind()
				err := getUnexpectedExpressionResultError(kind, reflect.Int)
				return nil, getInvalidTokenError(token.Type(), err)
			}
		} else if intVal, ok := isInteger(token.step); ok {
			stepInt = &intVal
		} else {
			kind := reflect.TypeOf(token.step).Kind()
			return nil, getInvalidTokenArgumentError(token.Type(), kind, reflect.Int)
		}
	}

	rangeResult, err := getRange(token, current, fromInt, toInt, stepInt)
	if err != nil {
		if isInvalidTokenTargetError(err) {
			return nil, err
		}
		return nil, getInvalidTokenError(token.Type(), err)
	}

	if substring, ok := rangeResult.(string); ok {
		if len(next) > 0 {
			return next[0].Apply(root, substring, next[1:])
		}
		return substring, nil
	}

	elements := rangeResult.([]interface{})

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

/**
expected responses
1. nil, error - if there is any errors
2. []interface{}, nil - if an array, map, or slice is processed correctly
3. string, nil - if a string is processed correctly
**/
func getRange(token Token, obj interface{}, start int64, end, step *int64) (interface{}, error) {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return nil, getInvalidTokenTargetNilError(
			token.Type(),
			reflect.Array, reflect.Map, reflect.Slice, reflect.String,
		)
	}

	var objValue reflect.Value
	var length int64
	var mapKeys []reflect.Value
	isString := false

	switch objType.Kind() {
	case reflect.Map:
		objValue = reflect.ValueOf(obj)
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
		objValue = reflect.ValueOf(obj)
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

	var from int64 = 0
	from = start
	if from < 0 {
		from = length + from
	}

	if from < 0 || from >= length {
		return nil, getInvalidTokenOutOfRangeError(token.Type())
	}

	to := length
	if end != nil {
		to = *end
		if to < 0 {
			to = length + to
		}

		if to < 0 || to >= length {
			return nil, getInvalidTokenOutOfRangeError(token.Type())
		}
		// +1 to 'to' so we can handle 'from' and 'to' being the same
		to++
	}

	var stp int64 = 1
	if step != nil {
		stp = *step
		if stp < 1 {
			return nil, getInvalidTokenOutOfRangeError(token.Type())
		}
	}

	array := make([]interface{}, 0)

	if mapKeys != nil {
		for i := from; i < to; i += stp {
			key := mapKeys[i]
			array = append(array, objValue.MapIndex(key).Interface())
		}
	} else if isString {
		substring := ""
		for i := from; i < to; i += stp {
			value := objValue.Index(int(i)).Uint()
			substring += fmt.Sprintf("%c", value)
		}
		return substring, nil
	} else {
		for i := from; i < to; i += stp {
			array = append(array, objValue.Index(int(i)).Interface())
		}
	}
	return array, nil
}
