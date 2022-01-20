package token

import (
	"fmt"
	"reflect"
)

func newScriptToken(expression string, options *Options) *scriptToken {
	return &scriptToken{expression: expression, options: options}
}

type scriptToken struct {
	expression string
	options    *Options
}

func (token *scriptToken) String() string {
	return fmt.Sprintf("[(%s)]", token.expression)
}

func (token *scriptToken) Type() string {
	return "script"
}

func (token *scriptToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	if token.expression == "" {
		return nil, getInvalidExpressionEmptyError()
	}

	value, err := evaluateExpression(root, current, token.expression, token.options)
	if err != nil {
		return nil, getInvalidExpressionError(err)
	}

	if value == nil {
		return nil, getUnexpectedExpressionResultNilError(reflect.Int, reflect.String)
	}

	if strValue, ok := value.(string); ok {
		nextToken := &keyToken{key: strValue}
		return nextToken.Apply(root, current, next)
	} else if intValue, ok := isInteger(value); ok {
		nextToken := newIndexToken(intValue, token.options)
		return nextToken.Apply(root, current, next)
	}

	valueType := reflect.TypeOf(value)
	return nil, getUnexpectedExpressionResultError(valueType.Kind(), reflect.Int, reflect.String)
}
