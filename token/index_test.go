package token

import (
	"testing"
)

// Test indexToken struct conforms to Token interface
var _ Token = &indexToken{}

func Test_IndexToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &indexToken{index: 0},
			input: input{
				current: nil,
			},
			expected: expected{
				err: "cannot get index from nil array",
			},
		},
		{
			token: &indexToken{index: 0},
			input: input{
				current: 123,
			},
			expected: expected{
				err: "invalid object. expected array, map, or string",
			},
		},
		{
			token: &indexToken{index: 5},
			input: input{
				current: "Find(X)",
			},
			expected: expected{
				value: "X",
			},
		},
		{
			token: &indexToken{index: 0},
			input: input{
				current: [3]string{"one", "two", "three"},
			},
			expected: expected{
				value: "one",
			},
		},
		{
			token: &indexToken{index: 0},
			input: input{
				current: []string{"one", "two", "three"},
			},
			expected: expected{
				value: "one",
			},
		},
		{
			token: &indexToken{index: 2},
			input: input{
				current: []string{"one", "two", "three"},
			},
			expected: expected{
				value: "three",
			},
		},
		{
			token: &indexToken{index: 4},
			input: input{
				current: []string{"one", "two", "three"},
			},
			expected: expected{
				err: "index out of range",
			},
		},
		{
			token: &indexToken{index: 1},
			input: input{
				current: []interface{}{"one", 2, "three"},
			},
			expected: expected{
				value: 2,
			},
		},
		{
			token: &indexToken{index: -1},
			input: input{
				current: []interface{}{"one", 2, "three"},
			},
			expected: expected{
				value: "three",
			},
		},
		{
			token: &indexToken{index: -2},
			input: input{
				current: []interface{}{"one", 2, "three"},
			},
			expected: expected{
				value: 2,
			},
		},
		{
			token: &indexToken{index: -3},
			input: input{
				current: []interface{}{"one", 2, "three"},
			},
			expected: expected{
				value: "one",
			},
		},
		{
			token: &indexToken{index: -4},
			input: input{
				current: []interface{}{"one", 2, "three"},
			},
			expected: expected{
				err: "index out of range",
			},
		},
		{
			token: &indexToken{index: 1},
			input: input{
				current: []interface{}{
					map[string]interface{}{
						"name":  "one",
						"value": 1,
					},
					map[string]interface{}{
						"name":  "two",
						"value": 2,
					},
					map[string]interface{}{
						"name":  "three",
						"value": 3,
					},
				},
				tokens: []Token{
					&keyToken{key: "name"},
				},
			},
			expected: expected{
				value: "two",
			},
		},
		{
			token: &indexToken{index: 1},
			input: input{
				current: map[string]interface{}{
					"a": map[string]interface{}{
						"name":  "one",
						"value": 1,
					},
					"c": map[string]interface{}{
						"name":  "three",
						"value": 3,
					},
					"b": map[string]interface{}{
						"name":  "two",
						"value": 2,
					},
				},
			},
			expected: expected{
				value: map[string]interface{}{
					"name":  "two",
					"value": 2,
				},
			},
		},
	}

	batchTokenTests(t, tests)
}
