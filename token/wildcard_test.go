package token

import (
	"testing"
)

func Test_WildcardToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &wildcardToken{},
			input: input{
				current: nil,
			},
			expected: expected{
				value: nil,
				err:   "cannot get elements from nil object",
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: "not array or map",
			},
			expected: expected{
				value: nil,
				err:   "invalid object. expected array or map",
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: []string{"one", "two", "three"},
			},
			expected: expected{
				value: []interface{}{"one", "two", "three"},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: []interface{}{"one", "two", "three", 4, 5},
			},
			expected: expected{
				value: []interface{}{"one", "two", "three", 4, 5},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: map[string]int64{
					"one":   1,
					"two":   2,
					"three": 3,
				},
			},
			expected: expected{
				value: []interface{}{
					int64(1),
					int64(2),
					int64(3),
				},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: map[string]string{
					"one":   "1",
					"two":   "2",
					"three": "3",
				},
			},
			expected: expected{
				value: []interface{}{"1", "2", "3"},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: map[string]interface{}{
					"one":   "1",
					"two":   2,
					"three": "3",
				},
			},
			expected: expected{
				value: []interface{}{"1", 2, "3"},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: [3]string{
					"1",
					"2",
					"3",
				},
			},
			expected: expected{
				value: []interface{}{"1", "2", "3"},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: []map[string]interface{}{
					{"name": "one"},
					{"name": "two"},
					{"name": "three"},
				},
			},
			expected: expected{
				value: []interface{}{
					map[string]interface{}{"name": "one"},
					map[string]interface{}{"name": "two"},
					map[string]interface{}{"name": "three"},
				},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: []map[string]interface{}{
					{"name": "one"},
					{"name": "two"},
					{"name": "three"},
				},
				tokens: []Token{
					&keyToken{
						key: "name",
					},
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
	}

	batchTokenTests(t, tests)

}
