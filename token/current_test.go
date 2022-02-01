package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test currentToken struct conforms to Token interface
var _ Token = &currentToken{}

func Test_newCurrentToken(t *testing.T) {
	assert.IsType(t, &currentToken{}, newCurrentToken())
}

func Test_CurrentToken_String(t *testing.T) {
	assert.Equal(t, "@", (&currentToken{}).String())
}

func Test_CurrentToken_Type(t *testing.T) {
	assert.Equal(t, "current", (&currentToken{}).Type())
}

var currentTests = []*tokenTest{
	{
		token: &currentToken{},
		input: input{
			current: map[string]interface{}{
				"name": "first",
			},
		},
		expected: expected{
			value: map[string]interface{}{
				"name": "first",
			},
		},
	},
	{
		token: &currentToken{},
		input: input{
			current: map[string]interface{}{
				"name": "first",
			},
			tokens: []Token{
				&keyToken{
					key: "name",
				},
			},
		},
		expected: expected{
			value: "first",
		},
	},
}

func Test_CurrentToken_Apply(t *testing.T) {
	batchTokenTests(t, currentTests)
}

func Benchmark_CurrentToken_Apply(b *testing.B) {
	batchTokenBenchmarks(b, currentTests)
}
