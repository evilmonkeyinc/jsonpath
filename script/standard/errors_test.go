package standard

import (
	goErr "errors"
	"testing"

	"github.com/evilmonkeyinc/jsonpath/errors"
	"github.com/stretchr/testify/assert"
)

func Test_error(t *testing.T) {

	t.Run("getInvalidExpressionEmptyError", func(t *testing.T) {
		actual := getInvalidExpressionEmptyError()
		assert.EqualError(t, actual, "invalid expression. is empty")
		assert.True(t, goErr.Is(actual, errors.ErrInvalidExpression))
	})
}
