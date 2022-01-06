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
		// TODO : wrap error?
		return nil, err
	}

	// No current
	if current == nil {
		return value, nil
	}

	// TODO : if current is array/slice/map value should be used to interperate that
	// if nore return/pass value
	objType := reflect.TypeOf(current)
	if objType == nil {
		return nil, errors.ErrGetElementsFromNilObject
	}
	switch objType.Kind() {
	case reflect.Map:
		strValue, ok := value.(string)
		if !ok {
			return nil, errors.GetInvalidParameterError("expected script to return key")
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
			return nil, errors.GetInvalidParameterError("expected script to return integer")
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
		return nil, errors.ErrInvalidObjectMapOrSlice
	}

	if len(next) > 0 {
		return next[0].Apply(root, current, next[1:])
	}
	return current, nil
}

func evaluateExpression(root, current interface{}, expression string) (interface{}, error) {
	if expression == "" {
		return nil, errors.GetInvalidParameterError("expression is empty")
	}
	/**
	TODO
	1. replace special tokens with evaluation
	2. run updated expression string through go/types.Eval
	**/

	// TODO : while loop around this block, update expression as we go along
	rootIndex := strings.Index(expression, "$")
	currentIndex := strings.Index(expression, "@")

	for rootIndex > -1 || currentIndex > -1 {

		query := ""
		if rootIndex > -1 {
			query = expression[rootIndex:]
		} else if currentIndex > -1 {
			query = expression[currentIndex:]
		} else {
			panic("something went wrong here")
		}

		tokenStrings, remainder, err := Tokenize(query) // TODO : tokenize remainder
		if err != nil {
			// TODO : wrap error, invalid script
			return nil, err
		}
		if remainder != "" {
			// shorten query to only what is being replaced
			query = query[0 : len(query)-len(remainder)]
		}
		if len(tokenStrings) > 0 {
			// TODO : this
			tokens := make([]Token, 0)
			for _, tokenString := range tokenStrings {
				token, err := Parse(tokenString)
				if err != nil {
					// TODO : wrap error
					return nil, err
				}
				tokens = append(tokens, token)
			}

			value, err := tokens[0].Apply(root, current, tokens[1:])
			if err != nil {
				// TODO : wrap error
				return nil, err
			}

			if strValue, ok := value.(string); ok {
				value = fmt.Sprintf("\"%s\"", strValue)
			}
			// TODO : value needs to be primitive
			// TODO : need to convert value to string safely
			expression = strings.ReplaceAll(expression, query, fmt.Sprintf("%v", value))
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
		// TODO : need to wrap error
		return nil, fmt.Errorf("failed to parse expression : %s", err.Error())
	}
	if tv.Value == nil {
		return nil, nil
	}
	switch tv.Value.Kind() {
	case constant.String:
		return tv.Value.String(), nil
	case constant.Bool:
		strValue := tv.Value.String()
		// TODO : parse error
		boolVal, _ := strconv.ParseBool(strValue)
		return boolVal, nil
	case constant.Float:
		strValue := tv.Value.String()
		// TODO : parse error
		floatVal, _ := strconv.ParseFloat(strValue, 64)
		return floatVal, nil
	case constant.Int:
		strValue := tv.Value.String()
		// TODO : parse error
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
