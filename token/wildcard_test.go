package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test wildcardToken struct conforms to Token interface
var _ Token = &wildcardToken{}

func Test_newWildcardToken(t *testing.T) {
	assert.IsType(t, &wildcardToken{}, newWildcardToken())
}

func Test_WildcardToken_String(t *testing.T) {
	assert.Equal(t, "[*]", (&wildcardToken{}).String())
}

func Test_WildcardToken_Type(t *testing.T) {
	assert.Equal(t, "wildcard", (&wildcardToken{}).Type())
}

func Test_WildcardToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &wildcardToken{},
			input: input{
				current: nil,
			},
			expected: expected{
				value: nil,
				err:   "wildcard: invalid token target. expected [array map slice] got [nil]",
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: "not array or map",
			},
			expected: expected{
				value: nil,
				err:   "wildcard: invalid token target. expected [array map slice] got [string]",
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: []string{"one", "two", "three"},
			},
			expected: expected{
				value: []interface{}{"one", "two", "three"},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: []interface{}{"one", "two", "three", 4, 5},
			},
			expected: expected{
				value: []interface{}{"one", "two", "three", 4, 5},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: map[string]int64{
					"one":   1,
					"two":   2,
					"three": 3,
				},
			},
			expected: expected{
				value: []interface{}{
					int64(1),
					int64(2),
					int64(3),
				},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: map[string]string{
					"one":   "1",
					"two":   "2",
					"three": "3",
				},
			},
			expected: expected{
				value: []interface{}{"1", "2", "3"},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: map[string]interface{}{
					"one":   "1",
					"two":   2,
					"three": "3",
				},
			},
			expected: expected{
				value: []interface{}{"1", 2, "3"},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: [3]string{
					"1",
					"2",
					"3",
				},
			},
			expected: expected{
				value: []interface{}{"1", "2", "3"},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: []map[string]interface{}{
					{"name": "one"},
					{"name": "two"},
					{"name": "three"},
				},
			},
			expected: expected{
				value: []interface{}{
					map[string]interface{}{"name": "one"},
					map[string]interface{}{"name": "two"},
					map[string]interface{}{"name": "three"},
				},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: []map[string]interface{}{
					{"name": "one"},
					{"name": "two"},
					{"name": "three"},
				},
				tokens: []Token{
					&keyToken{
						key: "name",
					},
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
					"three",
				},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: sampleStruct{},
			},
			expected: expected{
				value: []interface{}{
					"",
					int64(0),
					"",
				},
			},
		},
		{
			token: &wildcardToken{},
			input: input{
				current: sampleStruct{
					One:   "one",
					Two:   "two",
					Three: 3,
					Four:  4,
					Five:  "five",
					Six:   "six",
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
					int64(4),
					"five",
					"six",
				},
			},
		},
	}

	batchTokenTests(t, tests)

}
