package token

import "testing"

func Test_FirstNToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &firstNToken{},
			input: input{},
			expected: expected{
				err: "invalid parameter. expected integer",
			},
		},
		{
			token: &firstNToken{number: 2},
			input: input{},
			expected: expected{
				err: "cannot get range from nil array",
			},
		},
		{
			token: &firstNToken{number: 2},
			input: input{
				current: 123,
			},
			expected: expected{
				err: "invalid object. expected array, map, or string",
			},
		},
		{
			token: &firstNToken{
				number: &expressionToken{expression: ""},
			},
			input: input{},
			expected: expected{
				err: "invalid parameter. expression is empty",
			},
		},
		{
			token: &firstNToken{
				number: &expressionToken{expression: "\"key\""},
			},
			input: input{},
			expected: expected{
				err: "invalid parameter. expected integer",
			},
		},
		{
			token: &firstNToken{
				number: &expressionToken{expression: "2"},
			},
			input: input{
				current: []interface{}{
					"one",
					"two",
					"three",
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
				},
			},
		},
		{
			token: &firstNToken{
				number: -2,
			},
			input: input{
				current: []interface{}{
					"one",
					"two",
					"three",
					"four",
					"five",
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
					"three",
				},
			},
		},
		{
			token: &firstNToken{
				number: &expressionToken{expression: "-2"},
			},
			input: input{
				current: []interface{}{
					"one",
					"two",
					"three",
					"four",
					"five",
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
					"three",
				},
			},
		},
		{
			token: &firstNToken{
				number: 5,
			},
			input: input{
				current: []interface{}{
					"one",
					"two",
					"three",
					"four",
					"five",
				},
			},
			expected: expected{
				err: "index out of range",
			},
		},
		{
			token: &firstNToken{
				number: 2,
			},
			input: input{
				current: [3]interface{}{
					"one",
					"two",
					"three",
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
				},
			},
		},
		{
			token: &firstNToken{
				number: 3,
			},
			input: input{
				current: "substring",
			},
			expected: expected{
				value: "sub",
			},
		},
		{
			token: &firstNToken{
				number: 3,
			},
			input: input{
				current: map[string]string{
					"a": "one",
					"e": "five",
					"d": "four",
					"c": "three",
					"b": "two",
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
					"three",
				},
			},
		},
		{
			token: &firstNToken{
				number: 3,
			},
			input: input{
				current: "substring",
				tokens: []Token{
					&indexToken{index: 1},
				},
			},
			expected: expected{
				value: "u",
			},
		},
		{
			token: &firstNToken{
				number: 2,
			},
			input: input{
				current: []string{
					"one",
					"two",
					"three",
				},
				tokens: []Token{
					&indexToken{index: 1},
				},
			},
			expected: expected{
				value: "two",
			},
		},
		{
			token: &firstNToken{
				number: 2,
			},
			input: input{
				current: []map[string]interface{}{
					{"name": "one"},
					{"name": "two"},
					{"name": "three"},
				},
				tokens: []Token{
					&keyToken{key: "name"},
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
				},
			},
		},
	}

	batchTokenTests(t, tests)

}
