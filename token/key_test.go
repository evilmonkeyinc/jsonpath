package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_KeyToken_Apply(t *testing.T) {

	t.Run("nested", func(t *testing.T) {
		token := &keyToken{key: "one"}
		next := &currentToken{}

		actual, err := token.Apply(nil, map[string]string{"one": "two"}, []Token{next})
		assert.Nil(t, err)
		assert.Equal(t, "two", actual)
	})

	type input struct {
		obj interface{}
		key string
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
				obj: nil,
				err: "cannot get key from nil map",
			},
		},
		{
			input: input{
				obj: "",
			},
			expected: expected{
				obj: nil,
				err: "invalid object. expected map",
			},
		},
		{
			input: input{
				obj: map[string]interface{}{
					"key": true,
				},
				key: "key",
			},
			expected: expected{
				obj: true,
				err: "",
			},
		},
		{
			input: input{
				obj: map[string]interface{}{
					"key": true,
				},
				key: "missing",
			},
			expected: expected{
				obj: nil,
				err: "'missing' key not found in object",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {

			token := keyToken{
				key: test.input.key,
			}

			obj, err := token.Apply(test.input.obj, test.input.obj, nil)

			assert.Equal(t, test.expected.obj, obj)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}
		})
	}
}
