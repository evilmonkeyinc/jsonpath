package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test keyToken struct conforms to Token interface
var _ Token = &keyToken{}

func Test_newKeyToken(t *testing.T) {
	assert.IsType(t, &keyToken{}, newKeyToken(""))
}

func Test_KeyToken_String(t *testing.T) {
	tests := []*tokenStringTest{
		{
			input:    &keyToken{},
			expected: "['']",
		},
		{
			input:    &keyToken{key: "with space"},
			expected: "['with space']",
		},
		{
			input:    &keyToken{key: "nospace"},
			expected: "['nospace']",
		},
		{
			input:    &keyToken{key: "key's"},
			expected: "['key\\'s']",
		},
	}

	batchTokenStringTests(t, tests)
}

func Test_KeyToken_Type(t *testing.T) {
	assert.Equal(t, "key", (&keyToken{}).Type())
}

var keyTests = []*tokenTest{
	{
		token: &keyToken{key: "key"},
		input: input{
			current: nil,
		},
		expected: expected{
			value: nil,
			err:   "key: invalid token target. expected [map] got [nil]",
		},
	},
	{
		token: &keyToken{key: "key"},
		input: input{
			current: "",
		},
		expected: expected{
			value: nil,
			err:   "key: invalid token target. expected [map] got [string]",
		},
	},
	{
		token: &keyToken{key: "key"},
		input: input{
			current: map[string]interface{}{
				"key": true,
			},
		},
		expected: expected{
			value: true,
			err:   "",
		},
	},
	{
		token: &keyToken{key: "missing"},
		input: input{
			current: map[string]interface{}{
				"key": true,
			},
		},
		expected: expected{
			value: nil,
			err:   "key: invalid token key 'missing' not found",
		},
	},
	{
		token: &keyToken{key: "key"},
		input: input{
			current: map[string]interface{}{
				"key": map[string]interface{}{
					"next": "nested target",
				},
			},
			tokens: []Token{
				&keyToken{
					key: "next",
				},
			},
		},
		expected: expected{
			value: "nested target",
			err:   "",
		},
	},
	{
		token: &keyToken{key: "key's"},
		input: input{
			current: map[string]interface{}{
				"key's": []interface{}{
					1, 2, 3,
				},
			},
		},
		expected: expected{
			value: []interface{}{
				1, 2, 3,
			},
			err: "",
		},
	},
	{
		token: &keyToken{key: "two"},
		input: input{
			current: sampleStruct{
				Two: "two's value",
			},
		},
		expected: expected{
			value: "two's value",
		},
	},
	{
		token: &keyToken{key: "Five"},
		input: input{
			current: &sampleStruct{
				Five: "value",
			},
		},
		expected: expected{
			value: "value",
		},
	},
	{
		token: &keyToken{key: "two"},
		input: input{
			current: sampleStruct{
				Two: "two's value",
			},
			tokens: []Token{
				&indexToken{index: 0, allowString: true},
			},
		},
		expected: expected{
			value: "t",
		},
	},
	{
		token: &keyToken{key: "two"},
		input: input{
			current: sampleStruct{},
		},
		expected: expected{
			value: "",
		},
	},
	{
		token: &keyToken{key: "three"},
		input: input{
			current: sampleStruct{
				Four: 100,
			},
		},
		expected: expected{
			value: int64(100),
		},
	},
	{
		token: &keyToken{key: "other"},
		input: input{
			current: sampleStruct{},
		},
		expected: expected{
			err: "key: invalid token key 'other' not found",
		},
	},
}

func Test_KeyToken_Apply(t *testing.T) {
	batchTokenTests(t, keyTests)
}

func Benchmark_KeyToken_Apply(b *testing.B) {
	batchTokenBenchmarks(b, keyTests)
}
