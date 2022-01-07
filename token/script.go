package token

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"reflect"
	"strconv"
	"strings"

	"github.com/evilmokeyinc/jsonpath/errors"
)

type scriptToken struct {
	expression string
}

func (token *scriptToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	value, err := evaluateExpression(root, current, token.expression)
	if err != nil {
		return nil, errors.GetInvalidExpressionError(err)
	}

	// No current
	if current == nil {
		return value, nil
	}

	objType := reflect.TypeOf(current)
	if objType == nil {
		return nil, errors.ErrGetElementsFromNilObject
	}
	switch objType.Kind() {
	case reflect.Map:
		strValue, ok := value.(string)
		if !ok {
			return nil, errors.ErrInvalidParameterScriptExpectedToReturnString
		}

		objVal := reflect.ValueOf(current)
		keys := objVal.MapKeys()

		var found interface{} = nil

		for _, kv := range keys {
			if kv.String() == strValue {
				found = objVal.MapIndex(kv).Interface()
			}
		}

		if found == nil {
			return nil, errors.GetKeyNotFoundError(strValue)
		}
		current = found
	case reflect.Array, reflect.Slice:

		intValue, ok := value.(int64)
		if !ok {
			return nil, errors.ErrInvalidParameterScriptExpectedToReturnInteger
		}

		objVal := reflect.ValueOf(current)
		length := objVal.Len()

		if intValue < 0 {
			intValue += int64(length)
		}

		if intValue < 0 || intValue >= int64(length) {
			return nil, errors.ErrIndexOutOfRange
		}

		current = objVal.Index(int(intValue)).Interface()
	default:
		return nil, errors.ErrInvalidObjectArrayOrMap
	}

	if len(next) > 0 {
		return next[0].Apply(root, current, next[1:])
	}
	return current, nil
}

func evaluateExpression(root, current interface{}, expression string) (interface{}, error) {
	if expression == "" {
		return nil, errors.ErrInvalidParameterExpressionEmpty
	}

	rootIndex := strings.Index(expression, "$")
	currentIndex := strings.Index(expression, "@")

	for rootIndex > -1 || currentIndex > -1 {

		query := ""
		if rootIndex > -1 {
			query = expression[rootIndex:]
		} else if currentIndex > -1 {
			query = expression[currentIndex:]
		}

		tokenStrings, remainder, err := Tokenize(query)
		if err != nil {
			return nil, errors.GetInvalidExpressionError(err)
		}
		if remainder != "" {
			// shorten query to only what is being replaced
			query = query[0 : len(query)-len(remainder)]
		}
		if len(tokenStrings) > 0 {
			tokens := make([]Token, 0)
			for _, tokenString := range tokenStrings {
				token, err := Parse(tokenString)
				if err != nil {
					return nil, errors.GetInvalidExpressionError(err)
				}
				tokens = append(tokens, token)
			}

			value, err := tokens[0].Apply(root, current, tokens[1:])
			if err != nil {
				return nil, errors.GetInvalidExpressionError(err)
			}

			new := fmt.Sprintf("%v", value)
			if strValue, ok := value.(string); ok {
				new = fmt.Sprintf("\"%s\"", strValue)
			} else if intValue, ok := value.(int64); ok {
				new = fmt.Sprintf("%d", intValue)
			} else if boolValue, ok := value.(bool); ok {
				new = fmt.Sprintf("%t", boolValue)
			} else if floatValue, ok := value.(float64); ok {
				new = strconv.FormatFloat(floatValue, 'f', -1, 64)
			}

			expression = strings.ReplaceAll(expression, query, new)
		}

		rootIndex = strings.Index(expression, "$")
		currentIndex = strings.Index(expression, "@")
	}

	expression = strings.TrimSpace(expression)
	if expression == "" {
		// after replacing tokens, if empty, return false
		return false, nil
	}

	fs := token.NewFileSet()
	tv, err := types.Eval(fs, nil, token.NoPos, expression)
	if err != nil {
		return nil, errors.GetFailedToParseExpressionError(err)
	}
	if tv.Value == nil {
		return nil, nil
	}
	switch tv.Value.Kind() {
	case constant.String:
		return tv.Value.String(), nil
	case constant.Bool:
		strValue := tv.Value.String()
		boolVal, _ := strconv.ParseBool(strValue)
		return boolVal, nil
	case constant.Float:
		strValue := tv.Value.String()
		floatVal, _ := strconv.ParseFloat(strValue, 64)
		return floatVal, nil
	case constant.Int:
		strValue := tv.Value.String()
		intVal, _ := strconv.ParseInt(strValue, 10, 64)
		return intVal, nil
	case constant.Complex:
		fallthrough
	case constant.Unknown:
		fallthrough
	default:
		return tv.Value.String(), nil
	}
}
