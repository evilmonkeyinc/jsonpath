package token

import (
	"fmt"
	"reflect"
)

func newRangeToken(from, to, step interface{}, options *Options) *rangeToken {
	allowMap := false
	allowString := false

	if options != nil {
		allowMap = options.AllowMapReferenceByIndex || options.AllowMapReferenceByIndexInRange
		allowString = options.AllowStringReferenceByIndex || options.AllowStringReferenceByIndexInRange
	}

	return &rangeToken{
		from:        from,
		to:          to,
		step:        step,
		allowMap:    allowMap,
		allowString: allowString,
	}
}

type rangeToken struct {
	from, to, step interface{}
	allowMap       bool
	allowString    bool
}

func (token *rangeToken) String() string {
	fString := ""
	if token.from != nil {
		fString = fmt.Sprint(token.from)
	}
	tString := ""
	if token.to != nil {
		tString = fmt.Sprint(token.to)
	}
	if token.step == nil {
		return fmt.Sprintf("[%s:%s]", fString, tString)
	}

	sString := fmt.Sprint(token.step)
	return fmt.Sprintf("[%s:%s:%s]", fString, tString, sString)
}

func (token *rangeToken) Type() string {
	return "range"
}

func (token *rangeToken) Apply(root, current interface{}, next []Token) (interface{}, error) {

	var fromInt, toInt, stepInt *int64

	if token.from != nil {
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
				tmp := int64(intVal)
				fromInt = &tmp
			} else {
				kind := reflect.TypeOf(result).Kind()
				err := getUnexpectedExpressionResultError(kind, reflect.Int)
				return nil, getInvalidTokenError(token.Type(), err)
			}
		} else if intVal, ok := isInteger(token.from); ok {
			tmp := int64(intVal)
			fromInt = &tmp
		} else {
			kind := reflect.TypeOf(token.from).Kind()
			return nil, getInvalidTokenArgumentError(token.Type(), kind, reflect.Int)
		}
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

	rangeResult, err := token.getRange(current, fromInt, toInt, stepInt)
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
func (token *rangeToken) getRange(obj interface{}, start, end, step *int64) (interface{}, error) {

	allowedType := []reflect.Kind{
		reflect.Array,
		reflect.Slice,
	}
	if token.allowMap {
		allowedType = append(allowedType, reflect.Map)
	}
	if token.allowString {
		allowedType = append(allowedType, reflect.String)
	}

	objType, objVal := getTypeAndValue(obj)
	if objType == nil {
		return nil, getInvalidTokenTargetNilError(
			token.Type(),
			allowedType...,
		)
	}

	var length int64
	var mapKeys []reflect.Value
	isString := false

	switch objType.Kind() {
	case reflect.Map:
		if !token.allowMap {
			return nil, getInvalidTokenTargetError(
				token.Type(),
				objType.Kind(),
				allowedType...,
			)
		}
		length = int64(objVal.Len())
		mapKeys = objVal.MapKeys()
		sortMapKeys(mapKeys)
		break
	case reflect.String:
		if !token.allowString {
			return nil, getInvalidTokenTargetError(
				token.Type(),
				objType.Kind(),
				allowedType...,
			)
		}
		isString = true
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		length = int64(objVal.Len())
		mapKeys = nil
		break
	default:
		return nil, getInvalidTokenTargetError(
			token.Type(),
			objType.Kind(),
			allowedType...,
		)
	}

	var from int64 = 0
	if start != nil {
		from = *start
		if from < 0 {
			from = length + from
		}

		if from < 0 {
			from = 0
		}
		if from > length {
			from = length
		}
	}

	to := length
	if end != nil {
		to = *end
		if to < 0 {
			to = length + to
		}

		if to < 0 {
			to = 0
		}
		if to > length {
			to = length
		}
	}

	var stp int64 = 1
	if step != nil {
		stp = *step
		if stp == 0 {
			return nil, getInvalidTokenOutOfRangeError(token.Type())
		}
	}

	array := make([]interface{}, 0)

	if mapKeys != nil {
		if stp < 0 {
			for i := to - 1; i >= from; i += stp {
				key := mapKeys[i]
				array = append(array, objVal.MapIndex(key).Interface())
			}
		} else {
			for i := from; i < to; i += stp {
				key := mapKeys[i]
				array = append(array, objVal.MapIndex(key).Interface())
			}
		}
	} else if isString {
		substring := ""
		if stp < 0 {
			for i := to - 1; i >= from; i += stp {
				value := objVal.Index(int(i)).Uint()
				substring += fmt.Sprintf("%c", value)
			}
		} else {
			for i := from; i < to; i += stp {
				value := objVal.Index(int(i)).Uint()
				substring += fmt.Sprintf("%c", value)
			}
		}
		return substring, nil
	} else {
		if stp < 0 {
			for i := to - 1; i >= from; i += stp {
				array = append(array, objVal.Index(int(i)).Interface())
			}
		} else {
			for i := from; i < to; i += stp {
				array = append(array, objVal.Index(int(i)).Interface())
			}
		}
	}
	return array, nil
}
