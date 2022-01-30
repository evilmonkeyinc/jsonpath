package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test recursiveToken struct conforms to Token interface
var _ Token = &recursiveToken{}

func Test_newRecursiveToken(t *testing.T) {
	assert.IsType(t, &recursiveToken{}, newRecursiveToken())
}

func Test_RecursiveToken_String(t *testing.T) {
	assert.Equal(t, "..", (&recursiveToken{}).String())
}

func Test_RecursiveToken_Type(t *testing.T) {
	assert.Equal(t, "recursive", (&recursiveToken{}).Type())
}

var recursiveTokenTests = []*tokenTest{
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
				[]interface{}{
					map[string]interface{}{
						"name": "six",
					},
					map[string]interface{}{
						"name": []interface{}{
							"seven",
							map[string]interface{}{
								"name": "eight",
							},
							map[string]interface{}{
								"name": "nine",
							},
						},
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
				"six",
				"seven",
				"eight",
				"nine",
				map[string]interface{}{
					"name": "eight",
				},
				map[string]interface{}{
					"name": "nine",
				},
			},
		},
	},
	{
		token: &recursiveToken{},
		input: input{},
		expected: expected{
			value: []interface{}{},
		},
	},
	{
		token: &recursiveToken{},
		input: input{
			current: "string",
		},
		expected: expected{
			value: []interface{}{"string"},
		},
	},
	{
		token: &recursiveToken{},
		input: input{
			current: []interface{}{"string", "array"},
		},
		expected: expected{
			value: []interface{}{
				[]interface{}{"string", "array"},
				"string",
				"array",
			},
		},
	},
	{
		token: &recursiveToken{},
		input: input{
			current: []string{"string", "array"},
		},
		expected: expected{
			value: []interface{}{
				[]string{"string", "array"},
				"string",
				"array",
			},
		},
	},
	{
		token: &recursiveToken{},
		input: input{
			current: map[string]interface{}{
				"this": "map",
				"with": []interface{}{
					"array",
				},
			},
		},
		expected: expected{
			value: []interface{}{
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
	},
	{
		token: &recursiveToken{},
		input: input{
			current: sampleStruct{},
		},
		expected: expected{
			value: []interface{}{
				sampleStruct{},
				"",
				int64(0),
				"",
			},
		},
	},
	{
		token: &recursiveToken{},
		input: input{
			current: &sampleStruct{
				One:   "one",
				Two:   "two",
				Three: 3,
				Four:  4,
				Five:  "five",
				Six:   "six",
			},
		},
		expected: expected{
			value: []interface{}{
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
	},
}

func Test_RecursiveToken_Apply(t *testing.T) {
	batchTokenTests(t, recursiveTokenTests)
}

func Benchmark_RecursiveToken_Apply(b *testing.B) {
	batchTokenBenchmarks(b, recursiveTokenTests)
}
