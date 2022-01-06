package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CurrentToken_Apply(t *testing.T) {

	root := "root"
	current := "current"

	token := &currentToken{}

	t.Run("only", func(t *testing.T) {
		actual, err := token.Apply(root, current, nil)
		assert.Equal(t, current, actual)
		assert.Nil(t, err)
	})

	t.Run("nested", func(t *testing.T) {
		actual, err := token.Apply(root, current, []Token{token})
		assert.Equal(t, current, actual)
		assert.Nil(t, err)
	})

}
