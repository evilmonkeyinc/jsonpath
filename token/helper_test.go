package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isInteger(t *testing.T) {

	type expected struct {
		value int64
		ok    bool
	}

	tests := []struct {
		input    interface{}
		expected expected
	}{
		{
			input: nil,
			expected: expected{
				value: 0,
				ok:    false,
			},
		},
		{
			input: []string{},
			expected: expected{
				value: 0,
				ok:    false,
			},
		},
		{
			input: "string",
			expected: expected{
				value: 0,
				ok:    false,
			},
		},
		{
			input: map[string]string{},
			expected: expected{
				value: 0,
				ok:    false,
			},
		},
		{
			input: 100,
			expected: expected{
				value: 100,
				ok:    true,
			},
		},
		{
			input: 0,
			expected: expected{
				value: 0,
				ok:    true,
			},
		},
		{
			input: int64(100),
			expected: expected{
				value: 100,
				ok:    true,
			},
		},
		{
			input: int(100),
			expected: expected{
				value: 100,
				ok:    true,
			},
		},
		{
			input: int32(100),
			expected: expected{
				value: 100,
				ok:    true,
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			value, ok := isInteger(test.input)

			assert.Equal(t, test.expected.ok, ok)
			assert.EqualValues(t, test.expected.value, value)
		})
	}
}
