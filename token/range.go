package token

import (
	"reflect"

	"github.com/evilmokeyinc/jsonpath/errors"
)

type rangeToken struct {
	from, to, step interface{}
}

func (token *rangeToken) Apply(root, current interface{}, next []Token) (interface{}, error) {

	var fromInt, toInt, stepInt *int

	if token.from != nil {
		if script, ok := token.from.(Token); ok {
			result, err := script.Apply(root, current, nil)
			if err != nil {
				return nil, err
			}

			if intVal, ok := result.(int); ok {
				fromInt = &intVal
			} else {
				return nil, errors.ErrUnexpectedScriptResultInteger
			}
		} else if intVal, ok := token.from.(int); ok {
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

			if intVal, ok := result.(int); ok {
				toInt = &intVal
			} else {
				return nil, errors.ErrUnexpectedScriptResultInteger
			}
		} else if intVal, ok := token.to.(int); ok {
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

			if intVal, ok := result.(int); ok {
				stepInt = &intVal
			} else {
				return nil, errors.ErrUnexpectedScriptResultInteger
			}
		} else if intVal, ok := token.step.(int); ok {
			stepInt = &intVal
		} else {
			return nil, errors.ErrInvalidParameterInteger
		}
	}

	elements, err := getRange(current, fromInt, toInt, stepInt)
	if err != nil {
		return nil, err
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

func getRange(obj interface{}, start, end, step *int) ([]interface{}, error) {
	objType := reflect.TypeOf(obj)
	if objType == nil {
		return nil, errors.ErrGetRangeFromNilArray
	}

	switch objType.Kind() {
	case reflect.Array, reflect.Slice:
		objValue := reflect.ValueOf(obj)
		length := objValue.Len()

		from := 0
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
		}

		stp := 1
		if step != nil {
			stp = *step
			if stp < 1 {
				return nil, errors.ErrInvalidParameterRangeNegativeStep
			}
		}

		array := make([]interface{}, 0)
		for i := from; i < to; i += stp {
			obj := objValue.Index(i).Interface()
			array = append(array, obj)
		}
		return array, nil
	default:
		return nil, errors.ErrInvalidObjectArray
	}
}
