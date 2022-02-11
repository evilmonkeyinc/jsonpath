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
		{
			input: operatorTestInput{
				operator: &regexOperator{
					arg1: "'1'",
					arg2: `"\d"`,
				},
			},
			expected: operatorTestExpected{
				value: true,
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
				value: "'value'",
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
				value: "'this'",
			},
		},
		{
			input: operatorTestInput{
				operator: currentKeyOperator,
				paramters: map[string]interface{}{
					"@": map[string]interface{}{
						"key": true,
					},
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator: currentKeyOperator,
				paramters: map[string]interface{}{
					"@": map[string]interface{}{
						"notkey": true,
					},
				},
			},
			expected: operatorTestExpected{
				err: "key: invalid token key 'key' not found",
			},
		},
		{
			input: operatorTestInput{
				operator: currentKeyOperator,
				paramters: map[string]interface{}{
					"@": map[string]interface{}{
						"key": "'value'",
					},
				},
			},
			expected: operatorTestExpected{
				value: "'value'",
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_inOperator(t *testing.T) {
	currentDSelector, _ := newSelectorOperator("@.d", &ScriptEngine{}, nil)

	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator:  &inOperator{arg1: nil, arg2: nil},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: operatorTestInput{
				operator:  &inOperator{arg1: nil, arg2: []interface{}{"one"}},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
		{
			input: operatorTestInput{
				operator:  &inOperator{arg1: "one", arg2: []interface{}{"one"}},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator:  &inOperator{arg1: "one", arg2: `["one","two"]`},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator:  &inOperator{arg1: "one", arg2: `{"1":"one","2":"two"}`},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator:  &inOperator{arg1: "1", arg2: `[1,2,3]`},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator:  &inOperator{arg1: "1", arg2: `["1","2","3"]`},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
		{
			input: operatorTestInput{
				operator: &inOperator{
					arg1: "2",
					arg2: currentDSelector,
				},
				paramters: map[string]interface{}{
					"@": map[string]interface{}{
						"d": []interface{}{
							float64(1),
							float64(2),
							float64(3),
						},
					},
				},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
	}
	batchOperatorTests(t, tests)
}

func Test_notInOperator(t *testing.T) {
	tests := []*operatorTest{
		{
			input: operatorTestInput{
				operator:  &notInOperator{arg1: nil, arg2: nil},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				err: "invalid argument. is nil",
			},
		},
		{
			input: operatorTestInput{
				operator:  &notInOperator{arg1: nil, arg2: []interface{}{"one"}},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
		{
			input: operatorTestInput{
				operator:  &notInOperator{arg1: "one", arg2: []interface{}{"one"}},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
		{
			input: operatorTestInput{
				operator:  &notInOperator{arg1: "one", arg2: `["one","two"]`},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
		{
			input: operatorTestInput{
				operator:  &notInOperator{arg1: "one", arg2: `{"1":"one","2":"two"}`},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
		{
			input: operatorTestInput{
				operator:  &notInOperator{arg1: "1", arg2: `[1,2,3]`},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				value: false,
			},
		},
		{
			input: operatorTestInput{
				operator:  &notInOperator{arg1: "1", arg2: `["1","2","3"]`},
				paramters: map[string]interface{}{},
			},
			expected: operatorTestExpected{
				value: true,
			},
		},
	}
	batchOperatorTests(t, tests)
}
