package token

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/evilmokeyinc/jsonpath/errors"
)

type rangeToken struct {
	from, to, step interface{}
}

func (token *rangeToken) Apply(root, current interface{}, next []Token) (interface{}, error) {

	var fromInt, toInt, stepInt *int64

	if token.from != nil {
		if script, ok := token.from.(Token); ok {
			result, err := script.Apply(root, current, nil)
			if err != nil {
				return nil, err
			}

			if intVal, ok := isInteger(result); ok {
				tmp := int64(intVal)
				fromInt = &tmp
			} else {
				return nil, errors.ErrUnexpectedScriptResultInteger
			}
		} else if intVal, ok := isInteger(token.from); ok {
			fromInt = &intVal
		} else {
			return nil, errors.ErrInvalidParameterInteger
		}
	}

	if token.to != nil {
		if script, ok := token.to.(Token); ok {
			result, err := script.Apply(root, current, nil)
			if err != nil {
				return nil, err
			}

			if intVal, ok := isInteger(result); ok {
				tmp := int64(intVal)
				toInt = &tmp
			} else {
				return nil, errors.ErrUnexpectedScriptResultInteger
			}
		} else if intVal, ok := isInteger(token.to); ok {
			toInt = &intVal
		} else {
			return nil, errors.ErrInvalidParameterInteger
		}
	}

	if token.step != nil {
		if script, ok := token.step.(Token); ok {
			result, err := script.Apply(root, current, nil)
			if err != nil {
				return nil, err
			}

			if intVal, ok := isInteger(result); ok {
				tmp := int64(intVal)
				stepInt = &tmp
			} else {
				return nil, errors.ErrUnexpectedScriptResultInteger
			}
		} else if intVal, ok := isInteger(token.step); ok {
			stepInt = &intVal
		} else {
			return nil, errors.ErrInvalidParameterInteger
		}
	}

	rangeResult, err := getRange(current, fromInt, toInt, stepInt)
	if err != nil {
		return nil, err
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

// TODO : include tests for map, array, slice AND string objects

/**
expected responses
1. nil, error - if there is any errors
2. []interface{}, nil - if an array, map, or slice is processed correctly
3. string, nil - if a string is processed correctly
**/
func getRange(obj interface{}, start, end, step *int64) (interface{}, error) {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return nil, errors.ErrGetRangeFromNilArray
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
		return nil, errors.ErrInvalidObjectArrayMapOrString
	}

	var from int64 = 0
	if start != nil {
		from = *start
		if from < 0 {
			from = length + from
		}

		if from < 0 || from >= length {
			return nil, errors.ErrIndexOutOfRange
		}
	}

	to := length
	if end != nil {
		to = *end
		if to < 0 {
			to = length + to
		}

		if to < 0 || to >= length {
			return nil, errors.ErrIndexOutOfRange
		}
		// +1 to 'to' so we can handle 'from' and 'to' being the same
		to++
	}

	var stp int64 = 1
	if step != nil {
		stp = *step
		if stp < 1 {
			return nil, errors.ErrInvalidParameterRangeNegativeStep
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
