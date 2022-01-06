package token

import (
	"reflect"
	"strings"

	"github.com/evilmokeyinc/jsonpath/errors"
)

type filterToken struct {
	expression string
}

func (token *filterToken) Apply(root, current interface{}, next []Token) (interface{}, error) {

	shouldInclude := func(evaluation interface{}) bool {
		if evaluation == nil {
			return false
		}

		if matches, ok := evaluation.(bool); ok {
			return matches
		} else if strValue, ok := evaluation.(string); ok {
			strValue = strings.Trim(strValue, "\"")
			return strValue != ""
		}

		return true
	}

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
			element := objVal.MapIndex(kv).Interface()

			evaluation, err := evaluateExpression(root, element, token.expression)
			if err != nil {
				// we ignore errors, it has failed evaluation
				evaluation = nil
			}

			if shouldInclude(evaluation) {
				elements = append(elements, element)
			}

		}

	case reflect.Array, reflect.Slice:

		objVal := reflect.ValueOf(current)
		length := objVal.Len()

		for i := 0; i < length; i++ {
			element := objVal.Index(i).Interface()
			evaluation, err := evaluateExpression(root, element, token.expression)
			if err != nil {
				// we ignore errors, it has failed evaluation
				evaluation = nil
			}

			if shouldInclude(evaluation) {
				elements = append(elements, element)
			}
		}

	default:
		return nil, errors.ErrInvalidObjectMapOrSlice
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

			// TODO : do we care about errors?
			result, _ := nextToken.Apply(root, element, futureTokens)
			if result != nil {
				results = append(results, result)
			}
		}

		return results, nil
	}

	return elements, nil
}
