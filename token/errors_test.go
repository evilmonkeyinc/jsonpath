package token

import (
	"fmt"
	"testing"

	"github.com/evilmonkeyinc/jsonpath/errors"
	"github.com/stretchr/testify/assert"
)

func Test_isInvalidTokenError(t *testing.T) {
	tests := []struct {
		input    error
		expected bool
	}{
		{
			input:    fmt.Errorf("invalid token"),
			expected: false,
		},
		{
			input:    fmt.Errorf("is %w", errors.ErrInvalidToken),
			expected: true,
		},
		{
			input:    fmt.Errorf("is not %v", errors.ErrInvalidToken),
			expected: false,
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := isInvalidTokenError(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_isInvalidTokenTargetError(t *testing.T) {
	tests := []struct {
		input    error
		expected bool
	}{
		{
			input:    fmt.Errorf("invalid token"),
			expected: false,
		},
		{
			input:    fmt.Errorf("is %w", errors.ErrInvalidTokenTarget),
			expected: true,
		},
		{
			input:    fmt.Errorf("is not %v", errors.ErrInvalidTokenTarget),
			expected: false,
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := isInvalidTokenTargetError(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}
