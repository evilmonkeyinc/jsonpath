package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LengthToken_Apply(t *testing.T) {

	type input struct {
		obj interface{}
	}

	type expected struct {
		length int
		err    string
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
				err: "cannot get elements from nil object",
			},
		},
		{
			input: input{
				obj: 1000,
			},
			expected: expected{
				err: "invalid object. expected map or slice",
			},
		},
		{
			input: input{
				obj: [3]string{"one", "two", "three"},
			},
			expected: expected{
				length: 3,
			},
		},
		{
			input: input{
				obj: []interface{}{"one", "two", "three", 4, 5},
			},
			expected: expected{
				length: 5,
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
				length: 3,
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
				length: 3,
			},
		},
		{
			input: input{
				obj: "this is 26 characters long",
			},
			expected: expected{
				length: 26,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			token := &lengthToken{}
			obj, err := token.Apply(nil, test.input.obj, nil)

			if test.expected.err == "" {
				assert.Nil(t, err)
				assert.Equal(t, test.expected.length, obj)
			} else {
				assert.EqualError(t, err, test.expected.err)
				assert.Nil(t, obj)
			}
		})
	}

}
