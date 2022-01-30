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
	t.Run("success", func(t *testing.T) {
		actual, err := newExpressionToken("", &testEngine{}, nil)
		assert.Nil(t, err)
		assert.IsType(t, &expressionToken{}, actual)
	})
	t.Run("fail", func(t *testing.T) {
		actual, err := newExpressionToken("", &testEngine{err: fmt.Errorf("fail")}, nil)
		assert.EqualError(t, err, "fail")
		assert.Nil(t, actual)
	})
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
				expression:         "any",
				compiledExpression: &testCompiledExpression{err: fmt.Errorf("engine error")},
			},
			input: input{},
			expected: expected{
				err: "invalid expression. engine error",
			},
		},
		{
			token: &expressionToken{
				expression:         "any",
				compiledExpression: &testCompiledExpression{response: true},
			},
			input: input{},
			expected: expected{
				value: true,
			},
		},
		{
			token: &expressionToken{
				expression:         "any",
				compiledExpression: &testCompiledExpression{response: false},
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
