package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test rootToken struct conforms to Token interface
var _ Token = &rootToken{}

func Test_newRootToken(t *testing.T) {
	assert.IsType(t, &rootToken{}, newRootToken())
}

func Test_RootToken_String(t *testing.T) {
	assert.Equal(t, "$", (&rootToken{}).String())
}

func Test_RootToken_Type(t *testing.T) {
	assert.Equal(t, "root", (&rootToken{}).Type())
}

var rootTests = []*tokenTest{
	{
		token: &rootToken{},
		input: input{
			root: map[string]interface{}{
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
		token: &rootToken{},
		input: input{
			root: map[string]interface{}{
				"name": "first",
			},
			tokens: []Token{&keyToken{
				key: "name",
			}},
		},
		expected: expected{
			value: "first",
		},
	},
}

func Test_RootToken_Apply(t *testing.T) {
	batchTokenTests(t, rootTests)
}

func Benchmark_RootToken_Apply(b *testing.B) {
	batchTokenBenchmarks(b, rootTests)
}
