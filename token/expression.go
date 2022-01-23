package token

import (
	"fmt"

	"github.com/evilmonkeyinc/jsonpath/option"
	"github.com/evilmonkeyinc/jsonpath/script"
)

func newExpressionToken(expression string, engine script.Engine, options *option.QueryOptions) *expressionToken {
	return &expressionToken{
		expression: expression,
		engine:     engine,
		options:    options,
	}
}

type expressionToken struct {
	expression string
	engine     script.Engine
	options    *option.QueryOptions
}

func (token *expressionToken) String() string {
	return fmt.Sprintf("(%s)", token.expression)
}

func (token *expressionToken) Type() string {
	return "expression"
}

func (token *expressionToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	if token.expression == "" {
		return nil, getInvalidExpressionEmptyError()
	}

	value, err := token.engine.Evaluate(root, current, token.expression, token.options)
	if err != nil {
		return nil, getInvalidExpressionError(err)
	}

	if len(next) > 0 {
		return next[0].Apply(root, value, next[1:])
	}

	return value, nil
}
