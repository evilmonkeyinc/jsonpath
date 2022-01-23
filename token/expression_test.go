package token

import (
	"fmt"
	"testing"

	"github.com/evilmonkeyinc/jsonpath/option"
	"github.com/evilmonkeyinc/jsonpath/script"
	"github.com/stretchr/testify/assert"
)

type testEngine struct {
	response           interface{}
	err                error
	compiledExpression *testCompiledExpression
}

func (engine *testEngine) Compile(expression string, options *option.QueryOptions) (script.CompiledExpression, error) {
	return engine.compiledExpression, engine.err
}

func (engine *testEngine) Evaluate(root, current interface{}, expression string, options *option.QueryOptions) (interface{}, error) {
	return engine.response, engine.err
}

type testCompiledExpression struct {
	response interface{}
	err      error
}

func (engine *testCompiledExpression) Evaluate(root, current interface{}) (interface{}, error) {
	return engine.response, engine.err
}

// Test expressionToken struct conforms to Token interface
var _ Token = &expressionToken{}

func Test_newExpressionToken(t *testing.T) {
	assert.IsType(t, &expressionToken{}, newExpressionToken("", nil, nil))
}

func Test_ExpressionToken_String(t *testing.T) {

	tests := []*tokenStringTest{
		{
			input:    &expressionToken{expression: ""},
			expected: "()",
		},
		{
			input:    &expressionToken{expression: "1+1"},
			expected: "(1+1)",
		},
		{
			input:    &expressionToken{expression: "true"},
			expected: "(true)",
		},
	}

	batchTokenStringTests(t, tests)
}

func Test_ExpressionToken_Type(t *testing.T) {
	assert.Equal(t, "expression", (&expressionToken{}).Type())
}

func Test_ExpressionToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &expressionToken{},
			input: input{},
			expected: expected{
				err: "invalid expression. is empty",
			},
		},
		{
			token: &expressionToken{
				expression: "any",
				engine:     &testEngine{err: fmt.Errorf("engine error")},
			},
			input: input{},
			expected: expected{
				err: "invalid expression. engine error",
			},
		},
		{
			token: &expressionToken{
				expression: "any",
				engine:     &testEngine{response: true},
			},
			input: input{},
			expected: expected{
				value: true,
			},
		},
		{
			token: &expressionToken{
				expression: "any",
				engine:     &testEngine{response: false},
			},
			input: input{
				tokens: []Token{&currentToken{}},
			},
			expected: expected{
				value: false,
			},
		},
	}

	batchTokenTests(t, tests)
}
