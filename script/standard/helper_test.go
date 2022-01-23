package standard

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type operatorTestInput struct {
	operator  operator
	paramters map[string]interface{}
}

type operatorTestExpected struct {
	value interface{}
	err   string
}

type operatorTest struct {
	input    operatorTestInput
	expected operatorTestExpected
}

func batchOperatorTests(t *testing.T, tests []*operatorTest) {
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := test.input.operator.Evaluate(test.input.paramters)
			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}

			assert.Equal(t, test.expected.value, actual)
		})
	}
}

func Test_findUnquotedOperators(t *testing.T) {

	type input struct {
		source    string
		subString string
	}

	tests := []struct {
		input    input
		expected int
	}{
		{
			input: input{
				source:    "operator not found",
				subString: "||",
			},
			expected: -1,
		},
		{
			input: input{
				source:    "||",
				subString: "||",
			},
			expected: 0,
		},
		{
			input: input{
				source:    "@.length",
				subString: "@",
			},
			expected: 0,
		},
		{
			input: input{
				source:    "@.email == 'admin@example.com'",
				subString: "@",
			},
			expected: 0,
		},
		{
			input: input{
				source:    "'admin@example.com'",
				subString: "@",
			},
			expected: -1,
		},
		{
			input: input{
				source:    "10<@.price",
				subString: "@",
			},
			expected: 3,
		},
		{
			input: input{
				source:    "|| and ||",
				subString: "||",
			},
			expected: 0,
		},
		{
			input: input{
				source:    "'||' and ||",
				subString: "||",
			},
			expected: 9,
		},
		{
			input: input{
				source:    `"||" and '||' and '"||"' and "'||'"`,
				subString: "||",
			},
			expected: -1,
		},
		{
			input: input{
				source:    `[1||1]`,
				subString: "||",
			},
			expected: -1,
		},
		{
			input: input{
				source:    `(1||1)||false`,
				subString: "||",
			},
			expected: 6,
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := findUnquotedOperators(test.input.source, test.input.subString)
			assert.Equal(t, test.expected, actual)
		})
	}

}
