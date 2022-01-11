package token

import (
	"github.com/evilmokeyinc/jsonpath/errors"
)

type scriptToken struct {
	expression string
}

func (token *scriptToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	if token.expression == "" {
		return nil, errors.ErrInvalidParameterExpressionEmpty
	}

	value, err := evaluateExpression(root, current, token.expression)
	if err != nil {
		return nil, errors.GetInvalidExpressionError(err)
	}

	if strValue, ok := value.(string); ok {
		nextToken := &keyToken{key: strValue}
		return nextToken.Apply(root, current, next)
	} else if intValue, ok := isInteger(value); ok {
		nextToken := &indexToken{index: int64(intValue)}
		return nextToken.Apply(root, current, next)
	}
	return nil, errors.ErrUnexpectedScriptResultIntegerOrString
}
