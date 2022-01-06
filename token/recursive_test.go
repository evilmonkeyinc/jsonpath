package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RecursiveToken_Apply(t *testing.T) {

	t.Run("last", func(t *testing.T) {
		token := &recursiveToken{}

		input := map[string]interface{}{
			"key1": "one",
			"k2":   "two",
			"k3":   "three",
		}

		actual, err := token.Apply(nil, input, nil)
		assert.Nil(t, err)
		assert.ElementsMatch(t, actual, []interface{}{
			map[string]interface{}{
				"key1": "one",
				"k2":   "two",
				"k3":   "three",
			},
			"one",
			"two",
			"three",
		})
	})

	t.Run("nested", func(t *testing.T) {
		token := &recursiveToken{}
		next := &currentToken{}

		input := []interface{}{
			[]string{"one", "two", "three"},
		}

		actual, err := token.Apply(nil, input, []Token{next})
		assert.Nil(t, err)
		assert.ElementsMatch(t, actual, []interface{}{
			[]string{"one", "two", "three"},
			"one",
			"two",
			"three",
			"one",
			"two",
			"three",
		})
	})

}

func Test_flatten(t *testing.T) {

	tests := []struct {
		input    interface{}
		expected []interface{}
	}{
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
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := flatten(test.input)
			assert.ElementsMatch(t, test.expected, actual)
		})
	}
}
