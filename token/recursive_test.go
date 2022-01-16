package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test recursiveToken struct conforms to Token interface
var _ Token = &recursiveToken{}

func Test_RecursiveToken_String(t *testing.T) {
	assert.Equal(t, "..", (&recursiveToken{}).String())
}

func Test_RecursiveToken_Type(t *testing.T) {
	assert.Equal(t, "recursive", (&recursiveToken{}).Type())
}

func Test_RecursiveToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &recursiveToken{},
			input: input{
				current: []interface{}{nil, "one"},
			},
			expected: expected{
				value: []interface{}{
					[]interface{}{nil, "one"},
					"one",
				},
			},
		},
		{
			token: &recursiveToken{},
			input: input{
				current: map[string]interface{}{
					"key1": "one",
					"k2":   "two",
					"k3":   "three",
				},
			},
			expected: expected{
				value: []interface{}{
					map[string]interface{}{
						"key1": "one",
						"k2":   "two",
						"k3":   "three",
					},
					"one",
					"two",
					"three",
				},
			},
		},
		{
			token: &recursiveToken{},
			input: input{
				root: nil,
				current: []interface{}{
					map[string]interface{}{
						"name": "one",
						"nested": map[string]interface{}{
							"name": "four",
						},
					},
					map[string]interface{}{
						"name": "two",
					},
					map[string]interface{}{
						"name": "three",
					},
				},
				tokens: []Token{
					&keyToken{key: "name"},
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
					"three",
					"four",
				},
			},
		},
		{
			token: &recursiveToken{},
			input: input{
				root: nil,
				current: []interface{}{
					[]interface{}{
						map[string]interface{}{
							"name": "one",
						},
						map[string]interface{}{
							"name": "two",
						},
					},
					[]interface{}{
						map[string]interface{}{
							"name": "three",
						},
						map[string]interface{}{
							"name": []interface{}{"four", "five"},
						},
					},
				},
				tokens: []Token{
					&keyToken{key: "name"},
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
					"three",
					"four",
					"five",
				},
			},
		},
	}

	batchTokenTests(t, tests)
}

func Test_flatten(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected []interface{}
	}{
		{
			input:    nil,
			expected: []interface{}{},
		},
		{
			input:    "string",
			expected: []interface{}{"string"},
		},
		{
			input: []interface{}{"string", "array"},
			expected: []interface{}{
				[]interface{}{"string", "array"},
				"string",
				"array",
			},
		},
		{
			input: []string{"string", "array"},
			expected: []interface{}{
				[]string{"string", "array"},
				"string",
				"array",
			},
		},
		{
			input: map[string]interface{}{
				"this": "map",
				"with": []interface{}{
					"array",
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"this": "map",
					"with": []interface{}{
						"array",
					},
				},
				[]interface{}{
					"array",
				},
				"map",
				"array",
			},
		},
		{
			input: sampleStruct{},
			expected: []interface{}{
				sampleStruct{},
				"",
				int64(0),
				"",
			},
		},
		{
			input: &sampleStruct{
				One:   "one",
				Two:   "two",
				Three: 3,
				Four:  4,
				Five:  "five",
				Six:   "six",
			},
			expected: []interface{}{
				sampleStruct{
					One:   "one",
					Two:   "two",
					Three: 3,
					Four:  4,
					Five:  "five",
					Six:   "six",
				},
				"one",
				"two",
				int64(4),
				"five",
				"six",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := flatten(test.input)
			assert.ElementsMatch(t, test.expected, actual)
		})
	}
}
