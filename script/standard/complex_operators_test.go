package standard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_regexOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: &regexOperator{
					arg1: nil,
					arg2: "",
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: operatorTestInput{
				operator: &regexOperator{
					arg1: "",
					arg2: nil,
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: operatorTestInput{
				operator: &regexOperator{
					arg1: "",
					arg2: "",
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &regexOperator{
					arg1: "1",
					arg2: `\d`,
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &regexOperator{
					arg1: "1",
					arg2: `/\d/i`,
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: &regexOperator{
					arg1: "string",
					arg2: `\d`,
				},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
		{
			input: operatorTestInput{
				operator: &regexOperator{
					arg1: "string",
					arg2: `\`,
				},
			},
			expected: operatorTestExpected{
				err: "invalid argument. expected a valid regexp",
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_selectorOperator(t *testing.T) {

	t.Run("tokenize_fail", func(t *testing.T) {
		operator, err := newSelectorOperator("@!", &ScriptEngine{}, nil)
		assert.Nil(t, operator)
		assert.EqualError(t, err, "unexpected token '!' at index 1")
	})

	t.Run("parse_fail", func(t *testing.T) {
		operator, err := newSelectorOperator("@[]", &ScriptEngine{}, nil)
		assert.Nil(t, operator)
		assert.EqualError(t, err, "invalid token. '[]' does not match any token format")
	})

	currentOperator, _ := newSelectorOperator("@", &ScriptEngine{}, nil)
	currentKeyOperator, _ := newSelectorOperator("@.key", &ScriptEngine{}, nil)

	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator: currentOperator,
				paramters: map[string]interface{}{
					"@": "value",
				},
			},
			expected: operatorTestExpected{
				value: "value",
			},
		},
		{
			input: operatorTestInput{
				operator: currentKeyOperator,
				paramters: map[string]interface{}{
					"@": map[string]interface{}{
						"key": "this",
					},
				},
			},
			expected: operatorTestExpected{
				value: "this",
			},
		},
	}
	batchOperatorTests(t, tests)
}
