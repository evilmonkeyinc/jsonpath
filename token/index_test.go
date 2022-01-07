package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IndexToken_Apply(t *testing.T) {

	t.Run("nested", func(t *testing.T) {
		token := &indexToken{index: 0}
		next := &currentToken{}

		actual, err := token.Apply(nil, []string{"one", "two"}, []Token{next})
		assert.Nil(t, err)
		assert.Equal(t, "one", actual)
	})

	type input struct {
		obj interface{}
		idx int
	}

	type expected struct {
		obj interface{}
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
				err: "cannot get index from nil array",
			},
		},
		{
			input: input{
				obj: "not a array",
			},
			expected: expected{
				err: "invalid object. expected array",
			},
		},
		{
			input: input{
				obj: []string{"one", "two", "three"},
				idx: 0,
			},
			expected: expected{
				obj: "one",
			},
		},
		{
			input: input{
				obj: []string{"one", "two", "three"},
				idx: 2,
			},
			expected: expected{
				obj: "three",
			},
		},
		{
			input: input{
				obj: []string{"one", "two", "three"},
				idx: 4,
			},
			expected: expected{
				err: "index out of range",
			},
		},
		{
			input: input{
				obj: []interface{}{"one", 2, "three"},
				idx: 1,
			},
			expected: expected{
				obj: 2,
			},
		},
		{
			input: input{
				obj: []interface{}{"one", 2, "three"},
				idx: -1,
			},
			expected: expected{
				obj: "three",
			},
		},
		{
			input: input{
				obj: []interface{}{"one", 2, "three"},
				idx: -2,
			},
			expected: expected{
				obj: 2,
			},
		},
		{
			input: input{
				obj: []interface{}{"one", 2, "three"},
				idx: -3,
			},
			expected: expected{
				obj: "one",
			},
		},
		{
			input: input{
				obj: []interface{}{"one", 2, "three"},
				idx: -4,
			},
			expected: expected{
				err: "index out of range",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			token := &indexToken{
				index: test.input.idx,
			}

			obj, err := token.Apply(nil, test.input.obj, nil)

			assert.Equal(t, test.expected.obj, obj)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
		})
	}

}
