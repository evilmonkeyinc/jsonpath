package token

import (
	"fmt"
	"reflect"
	"strings"
)

func newFilterToken(expression string, options *Options) *filterToken {
	return &filterToken{expression: expression, options: options}
}

type filterToken struct {
	expression string
	options    *Options
}

func (token *filterToken) String() string {
	return fmt.Sprintf("[?(%s)]", token.expression)
}

func (token *filterToken) Type() string {
	return "filter"
}

func (token *filterToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	if token.expression == "" {
		return nil, getInvalidExpressionEmptyError()
	}

	shouldInclude := func(evaluation interface{}) bool {
		if evaluation == nil {
			return false
		}

		if matches, ok := evaluation.(bool); ok {
			return matches
		} else if strValue, ok := evaluation.(string); ok {
			strValue = strings.Trim(strValue, "'")
			return strValue != ""
		}

		return true
	}

	elements := make([]interface{}, 0)

	objType, objVal := getTypeAndValue(current)
	if objType == nil {
		return nil, getInvalidTokenTargetNilError(token.Type(), reflect.Array, reflect.Map, reflect.Slice)
	}
	switch objType.Kind() {
	case reflect.Map:
		keys := objVal.MapKeys()
		sortMapKeys(keys)

		for _, kv := range keys {
			element := objVal.MapIndex(kv).Interface()

			// TODO : we should compile expression so we don't have to tokenize each time
			evaluation, err := evaluateExpression(root, element, token.expression, token.options)
			if err != nil {
				// we ignore errors, it has failed evaluation
				evaluation = nil
			}

			if shouldInclude(evaluation) {
				elements = append(elements, element)
			}
		}
	case reflect.Array, reflect.Slice:
		length := objVal.Len()

		for i := 0; i < length; i++ {
			element := objVal.Index(i).Interface()
			// TODO : we should compile expression so we don't have to tokenize each time
			evaluation, err := evaluateExpression(root, element, token.expression, token.options)
			if err != nil {
				// we ignore errors, it has failed evaluation
				evaluation = nil
			}

			if shouldInclude(evaluation) {
				elements = append(elements, element)
			}
		}
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

		if indexToken, ok := nextToken.(*indexToken); ok {
			// if next is asking for specific index
			return indexToken.Apply(current, elements, futureTokens)
		}
		// any other token type

		results := make([]interface{}, 0)

		for _, element := range elements {

			result, _ := nextToken.Apply(root, element, futureTokens)
			if result != nil {
				results = append(results, result)
			}
		}

		return results, nil
	}

	return elements, nil
}
