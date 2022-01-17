package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parse(t *testing.T) {

	type input struct {
		query string
	}

	type expected struct {
		token Token
		err   string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{query: ""},
			expected: expected{
				err: "invalid token. token string is empty",
			},
		},
		{
			input: input{query: "['fail'"},
			expected: expected{
				err: "invalid token. '['fail'' does not match any token format",
			},
		},
		{
			input: input{query: "[ ]"},
			expected: expected{
				err: "invalid token. '[ ]' does not match any token format",
			},
		},
		{
			input: input{query: "[?]"},
			expected: expected{
				err: "invalid token. '[?]' does not match any token format",
			},
		},
		{
			input: input{query: "[why though]"},
			expected: expected{
				err: "invalid token. '[why though]' does not match any token format",
			},
		},
		{
			input: input{query: "[1'2']"},
			expected: expected{
				err: "invalid token. '[1'2']' does not match any token format",
			},
		},
		{
			input: input{query: "[1(@.length)]"},
			expected: expected{
				err: "invalid expression. invalid format '1(@.length)'",
			},
		},
		{
			input: input{query: "$"},
			expected: expected{
				token: &rootToken{},
			},
		},
		{
			input: input{query: "store"},
			expected: expected{
				token: &keyToken{
					key: "store",
				},
			},
		},
		{
			input: input{query: "1"},
			expected: expected{
				token: &keyToken{key: "1"},
				err:   "",
			},
		},
		{
			input: input{query: "book"},
			expected: expected{
				token: &keyToken{
					key: "book",
				},
			},
		},
		{
			input: input{query: "*"},
			expected: expected{
				token: &wildcardToken{},
			},
		},
		{
			input: input{query: "author"},
			expected: expected{
				token: &keyToken{
					key: "author",
				},
			},
		},
		{
			input: input{query: "@"},
			expected: expected{
				token: &currentToken{},
			},
		},
		{
			input: input{query: "[*]"},
			expected: expected{
				token: &wildcardToken{},
			},
		},
		{
			input: input{query: ".."},
			expected: expected{
				token: &recursiveToken{},
			},
		},
		{
			input: input{query: "[?(@.isbn)]"},
			expected: expected{
				token: &filterToken{
					expression: "@.isbn",
				},
			},
		},
		{
			input: input{query: "[2]"},
			expected: expected{
				token: &indexToken{
					index: 2,
				},
			},
		},
		{
			input: input{query: "[(@.length-1)]"},
			expected: expected{
				token: &scriptToken{
					expression: "@.length-1",
				},
			},
		},
		{
			input: input{query: "[-1:]"},
			expected: expected{
				token: &rangeToken{
					from: int64(-1),
					to:   nil,
					step: nil,
				},
			},
		},
		{
			input: input{query: "[1:(@.length-1)]"},
			expected: expected{
				token: &rangeToken{
					from: int64(1),
					to: &expressionToken{
						expression: "@.length-1",
					},
					step: nil,
				},
			},
		},
		{
			input: input{query: "[0,1]"},
			expected: expected{
				token: &unionToken{
					arguments: []interface{}{int64(0), int64(1)},
				},
			},
		},
		{
			input: input{query: "['first','last']"},
			expected: expected{
				token: &unionToken{
					arguments: []interface{}{"first", "last"},
				},
			},
		},
		{
			input: input{query: "[0,]"},
			expected: expected{
				err: "invalid token. '[0,]' does not match any token format",
			},
		},
		{
			input: input{query: "[,1]"},
			expected: expected{
				err: "invalid token. '[,1]' does not match any token format",
			},
		},
		{
			input: input{query: "[(0),]"},
			expected: expected{
				err: "invalid token. '[(0),]' does not match any token format",
			},
		},
		{
			input: input{query: "[0,'1',]"},
			expected: expected{
				err: "invalid token. '[0,'1',]' does not match any token format",
			},
		},
		{
			input: input{query: "[0,(@.length-1)]"},
			expected: expected{

				token: &unionToken{
					arguments: []interface{}{
						int64(0),
						&expressionToken{
							expression: "@.length-1",
						},
					},
				},
			},
		},
		{
			input: input{query: "[0,'one',2,(1+2)]"},
			expected: expected{
				token: &unionToken{
					arguments: []interface{}{
						int64(0),
						"one",
						int64(2),
						&expressionToken{
							expression: "1+2",
						},
					},
				},
			},
		},
		{
			input: input{query: "[(@.length-2),(@.length-1),1]"},
			expected: expected{
				token: &unionToken{
					arguments: []interface{}{
						&expressionToken{
							expression: "@.length-2",
						},
						&expressionToken{
							expression: "@.length-1",
						},
						int64(1),
					},
				},
			},
		},
		{
			input: input{query: "[:2]"},
			expected: expected{
				token: &rangeToken{
					to: int64(2),
				},
			},
		},
		{
			input: input{query: "[:(@.length-1)]"},
			expected: expected{
				token: &rangeToken{
					to: &expressionToken{
						expression: "@.length-1",
					},
				},
			},
		},
		{
			input: input{query: "[:'key']"},
			expected: expected{
				err: "invalid expression. invalid format ''key''",
			},
		},
		{
			input: input{query: "[0:(@.length-1):2]"},
			expected: expected{
				token: &rangeToken{
					from: int64(0),
					to: &expressionToken{
						expression: "@.length-1",
					},
					step: int64(2),
				},
			},
		},
		{
			input: input{query: "[(@.length-1):1:2]"},
			expected: expected{
				token: &rangeToken{
					from: &expressionToken{
						expression: "@.length-1",
					},
					to:   int64(1),
					step: int64(2),
				},
			},
		},
		{
			input: input{query: "['store']"},
			expected: expected{
				token: &keyToken{
					key: "store",
				},
			},
		},
		{
			input: input{query: "[store]"},
			expected: expected{
				err: "invalid token. '[store]' does not match any token format",
			},
		},
		{
			input: input{query: "[store,book]"},
			expected: expected{
				err: "invalid token. '[store,book]' does not match any token format",
			},
		},
		{
			input: input{query: "[(1+2*(3+4)+5')]"},
			expected: expected{
				token: &scriptToken{
					expression: "1+2*(3+4)+5'",
				},
			},
		},
		{
			input: input{query: "['this key has brackets ( and colons : and commas , but is not a union, range, or script']"},
			expected: expected{
				token: &keyToken{
					key: "this key has brackets ( and colons : and commas , but is not a union, range, or script",
				},
			},
		},
		{
			input: input{query: "[1,2:4]"},
			expected: expected{
				err: "invalid token. '[1,2:4]' does not match any token format",
			},
		},
		{
			input: input{query: "[1:2:3:]"},
			expected: expected{
				err: "invalid token. '[1:2:3:]' does not match any token format",
			},
		},
		{
			input: input{query: "[1:2:3:4]"},
			expected: expected{
				err: "invalid token. '[1:2:3:4]' does not match any token format",
			},
		},
		{
			input: input{query: "['key':'end]"},
			expected: expected{
				err: "invalid expression. invalid format ''key''",
			},
		},
		{
			input: input{query: "[::2]"},
			expected: expected{
				token: &rangeToken{step: int64(2)},
			},
		},
		{
			input: input{query: "[:end:2]"},
			expected: expected{
				err: "invalid token. '[:end:2]' does not match any token format",
			},
		},
		{
			input: input{query: "length"},
			expected: expected{
				token: &lengthToken{},
			},
		},
		{
			input: input{query: "[length]"},
			expected: expected{
				err: "invalid token. '[length]' does not match any token format",
			},
		},
		{
			input: input{query: "['length']"},
			expected: expected{
				token: &keyToken{
					key: "length",
				},
			},
		},
		{
			input: input{query: "['']"},
			expected: expected{
				token: &keyToken{key: ""},
			},
		},
		{
			input: input{query: "['1':(@.length)]"},
			expected: expected{
				err: "invalid expression. invalid format ''1''",
			},
		},
		{
			input: input{query: "[0:'1']"},
			expected: expected{
				err: "invalid expression. invalid format ''1''",
			},
		},
		{
			input: input{query: "[0:1:'1']"},
			expected: expected{
				err: "invalid expression. invalid format ''1''",
			},
		},
		{
			input: input{query: "[0:100:(1+1)]"},
			expected: expected{
				token: &rangeToken{
					from: int64(0),
					to:   int64(100),
					step: &expressionToken{
						expression: "1+1",
					},
				},
			},
		},
		{
			input: input{query: "[:10:1]"},
			expected: expected{
				token: &rangeToken{to: int64(10), step: int64(1)},
			},
		},
		{
			input: input{query: "[1, 2]"},
			expected: expected{
				token: &unionToken{[]interface{}{int64(1), int64(2)}},
			},
		},
		{
			input: input{query: "[1: 2]"},
			expected: expected{
				token: &rangeToken{from: int64(1), to: int64(2)},
			},
		},
		{
			input: input{query: "['key\\'s']"},
			expected: expected{
				token: &keyToken{
					key: "key's",
				},
			},
		},
		{
			input: input{query: "[\\'key\\'s']"},
			expected: expected{
				err: "invalid token. '[\\'key\\'s']' does not match any token format",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			token, err := Parse(test.input.query)

			if test.expected.err == "" {
				assert.Nil(t, err, fmt.Sprintf("input '%s' err check failed. expected 'nil' actual '%v'", test.input.query, err))
			} else {
				assert.EqualError(t, err, test.expected.err, fmt.Sprintf("input '%s' err check failed. expected '%s' actual '%v'", test.input.query, test.expected.err, err))
			}
			assert.EqualValues(t, test.expected.token, token, fmt.Sprintf("input '%s' token check failed. expected '%v' actual '%v'", test.input.query, test.expected.token, token))
		})
	}

}

func Test_tokenize(t *testing.T) {

	type expected struct {
		tokens    []string
		remainder string
		err       string
	}

	tests := []struct {
		input    string
		expected expected
	}{
		{
			input: "$",
			expected: expected{
				tokens: []string{
					"$",
				},
			},
		},
		{
			input: "@",
			expected: expected{
				tokens: []string{
					"@",
				},
			},
		},
		{
			input: "%",
			expected: expected{
				err: "unexpected token '%' at index 0",
			},
		},
		{
			input: "$.store.book[*].author",
			expected: expected{
				tokens: []string{
					"$",
					"store",
					"book",
					"[*]",
					"author",
				},
			},
		},
		{
			input: "$..author",
			expected: expected{
				tokens: []string{
					"$",
					"..",
					"author",
				},
			},
		},
		{
			input: "$.store.*",
			expected: expected{
				tokens: []string{
					"$",
					"store",
					"*",
				},
			},
		},
		{
			input: "$.store..price",
			expected: expected{
				tokens: []string{
					"$",
					"store",
					"..",
					"price",
				},
			},
		},
		{
			input: "$..book[2]",
			expected: expected{
				tokens: []string{
					"$",
					"..",
					"book",
					"[2]",
				},
			},
		},
		{
			input: "$..book[(@.length-1)]",
			expected: expected{
				tokens: []string{
					"$",
					"..",
					"book",
					"[(@.length-1)]",
				},
			},
		},
		{
			input: "$..book[-1:]",
			expected: expected{
				tokens: []string{
					"$",
					"..",
					"book",
					"[-1:]",
				},
			},
		},
		{
			input: "$..book[0,1]",
			expected: expected{
				tokens: []string{
					"$",
					"..",
					"book",
					"[0,1]",
				},
			},
		},
		{
			input: "$..book[:2]",
			expected: expected{
				tokens: []string{
					"$",
					"..",
					"book",
					"[:2]",
				},
			},
		},
		{
			input: "$..book[?(@.isbn)]",
			expected: expected{
				tokens: []string{
					"$",
					"..",
					"book",
					"[?(@.isbn)]",
				},
			},
		},
		{
			input: "$..book[?(@.price<10)]",
			expected: expected{
				tokens: []string{
					"$",
					"..",
					"book",
					"[?(@.price<10)]",
				},
			},
		},
		{
			input: "$..*",
			expected: expected{
				tokens: []string{
					"$",
					"..",
					"*",
				},
			},
		},
		{
			input: "@.sub.query",
			expected: expected{
				tokens: []string{
					"@",
					"sub",
					"query",
				},
			},
		},
		{
			input: "user.email",
			expected: expected{
				tokens: nil,
				err:    "unexpected token 'u' at index 0",
			},
		},
		{
			input: "$['store']['book'][*]['author']",
			expected: expected{
				tokens: []string{
					"$",
					"['store']",
					"['book']",
					"[*]",
					"['author']",
				},
			},
		},
		{
			input: "$['store'].book[*].author",
			expected: expected{
				tokens: []string{
					"$",
					"['store']",
					"book",
					"[*]",
					"author",
				},
			},
		},
		{
			input: "$*",
			expected: expected{
				err: "unexpected token '*' at index 1",
			},
		},
		{
			input: "$['store']book.author",
			expected: expected{
				tokens: []string{
					"$",
					"['store']",
					"book",
					"author",
				},
			},
		},
		{
			input: "",
			expected: expected{
				err: "unexpected token '' at index 0",
			},
		},
		{
			input: "$.store.book[($.totals[0])].author",
			expected: expected{
				tokens: []string{
					"$",
					"store",
					"book",
					"[($.totals[0])]",
					"author",
				},
			},
		},
		{
			input: "$.['book'].[*].author",
			expected: expected{
				tokens: []string{
					"$",
					"['book']",
					"[*]",
					"author",
				},
			},
		},
		{
			input: "@.length-1",
			expected: expected{
				tokens:    []string{"@", "length"},
				remainder: "-1",
			},
		},
		{
			input: "@.length*2",
			expected: expected{
				tokens:    []string{"@", "length"},
				remainder: "*2",
			},
		},
		{
			input: "@.isbn",
			expected: expected{
				tokens: []string{"@", "isbn"},
			},
		},
		{
			input: "@.price<10",
			expected: expected{
				tokens:    []string{"@", "price"},
				remainder: "<10",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			tokens, remainder, err := tokenize(test.input, true)

			if test.expected.err != "" {
				assert.EqualError(t, err, test.expected.err, "unexpected error for %s", test.input)
			} else {
				assert.Nil(t, err)
			}

			for i, actual := range tokens {
				expected := test.expected.tokens[i]
				assert.EqualValues(t, expected, actual, "unexpected token for %s", test.input)
			}

			assert.Equal(t, test.expected.remainder, remainder, "unexpected remainder for %s", test.input)
		})
	}

}

type tokenStringTest struct {
	input    Token
	expected string
}

func batchTokenStringTests(t *testing.T, tests []*tokenStringTest) {
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual := test.input.String()
			assert.Equal(t, test.expected, actual)
		})
	}
}

type input struct {
	root, current interface{}
	tokens        []Token
}

type expected struct {
	value interface{}
	err   string
}

type tokenTest struct {
	token    Token
	input    input
	expected expected
}

func batchTokenTests(t *testing.T, tests []*tokenTest) {
	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			actual, err := test.token.Apply(test.input.root, test.input.current, test.input.tokens)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expected.err)
			}

			if test.expected.value != nil {
				if actual == nil {
					assert.Fail(t, "expected non-nil response")
					return
				}
				if array, ok := test.expected.value.([]interface{}); ok {
					assert.ElementsMatch(t, array, actual)
				} else {
					assert.Equal(t, test.expected.value, actual)
				}
			} else {
				assert.Nil(t, actual)
			}
		})
	}
}
