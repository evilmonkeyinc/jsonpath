package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WildcardToken_Apply(t *testing.T) {

	type input struct {
		obj interface{}
	}

	type expected struct {
		obj []interface{}
		err string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{
				obj: nil,
			},
			expected: expected{
				obj: nil,
				err: "cannot get elements from nil object",
			},
		},
		{
			input: input{
				obj: "not array or map",
			},
			expected: expected{
				obj: nil,
				err: "invalid object. expected array or map",
			},
		},
		{
			input: input{
				obj: []string{"one", "two", "three"},
			},
			expected: expected{
				obj: []interface{}{"one", "two", "three"},
			},
		},
		{
			input: input{
				obj: []interface{}{"one", "two", "three", 4, 5},
			},
			expected: expected{
				obj: []interface{}{"one", "two", "three", 4, 5},
			},
		},
		{
			input: input{
				obj: map[string]int{
					"one":   1,
					"two":   2,
					"three": 3,
				},
			},
			expected: expected{
				obj: []interface{}{1, 2, 3},
			},
		},
		{
			input: input{
				obj: map[string]string{
					"one":   "1",
					"two":   "2",
					"three": "3",
				},
			},
			expected: expected{
				obj: []interface{}{"1", "2", "3"},
			},
		},
		{
			input: input{
				obj: map[string]interface{}{
					"one":   "1",
					"two":   2,
					"three": "3",
				},
			},
			expected: expected{
				obj: []interface{}{"1", 2, "3"},
			},
		},
		{
			input: input{
				obj: [3]string{
					"1",
					"2",
					"3",
				},
			},
			expected: expected{
				obj: []interface{}{"1", "2", "3"},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			token := &wildcardToken{}
			obj, err := token.Apply(nil, test.input.obj, nil)

			assert.ElementsMatch(t, test.expected.obj, obj)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
		})
	}

}
