package token

import (
	"fmt"
	"testing"

	"github.com/evilmokeyinc/jsonpath/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Parse(t *testing.T) {

	type expected struct {
		token Token
		err   string
	}

	tests := []struct {
		input    string
		expected expected
	}{
		{
			input: "",
			expected: expected{
				err: "invalid token. token can not be empty",
			},
		},
		{
			input: "['fail'",
			expected: expected{
				err: "invalid token. missing subscript close",
			},
		},
		{
			input: "[ ]",
			expected: expected{
				err: "invalid token. empty subscript",
			},
		},
		{
			input: "[?]",
			expected: expected{
				err: "invalid token. invalid filter format",
			},
		},
		{
			input: "[why though]",
			expected: expected{
				err: "invalid token. unexpected space",
			},
		},
		{
			input: "[1'2']",
			expected: expected{
				err: "invalid token. unexpected single quote",
			},
		},
		{
			input: "[1(@.length)]",
			expected: expected{
				err: "invalid token. invalid script format",
			},
		},
		{
			input: "$",
			expected: expected{
				token: &rootToken{},
			},
		},
		{
			input: "store",
			expected: expected{
				token: &keyToken{
					key: "store",
				},
			},
		},
		{
			input: "1",
			expected: expected{
				err: "invalid token. index specified as key",
			},
		},
		{
			input: "book",
			expected: expected{
				token: &keyToken{
					key: "book",
				},
			},
		},
		{
			input: "*",
			expected: expected{
				token: &wildcardToken{},
			},
		},
		{
			input: "author",
			expected: expected{
				token: &keyToken{
					key: "author",
				},
			},
		},
		{
			input: "@",
			expected: expected{
				token: &currentToken{},
			},
		},
		{
			input: "[*]",
			expected: expected{
				token: &wildcardToken{},
			},
		},
		{
			input: "..",
			expected: expected{
				token: &recursiveToken{},
			},
		},
		{
			input: "[?(@.isbn)]",
			expected: expected{
				token: &filterToken{
					expression: "@.isbn",
				},
			},
		},
		{
			input: "[2]",
			expected: expected{
				token: &indexToken{
					index: 2,
				},
			},
		},
		{
			input: "[(@.length-1)]",
			expected: expected{
				token: &scriptToken{
					expression: "@.length-1",
				},
			},
		},
		{
			input: "[-1:]",
			expected: expected{
				token: &rangeToken{
					from: -1,
					to:   nil,
					step: 1,
				},
			},
		},
		{
			input: "[1:(@.length-1)]",
			expected: expected{
				token: &rangeToken{
					from: 1,
					to: &scriptToken{
						expression: "@.length-1",
					},
					step: 1,
				},
			},
		},
		{
			input: "[0,1]",
			expected: expected{
				token: &unionToken{
					arguments: []interface{}{0, 1},
				},
			},
		},
		{
			input: "['first','last']",
			expected: expected{
				token: &unionToken{
					arguments: []interface{}{"first", "last"},
				},
			},
		},
		{
			input: "[0,]",
			expected: expected{
				err: "invalid token. empty argument in union",
			},
		},
		{
			input: "[,1]",
			expected: expected{
				err: "invalid token. empty argument in union",
			},
		},
		{
			input: "[(0),]",
			expected: expected{
				err: "invalid token. empty argument in union",
			},
		},
		{
			input: "[0,'1',]",
			expected: expected{
				err: "invalid token. empty argument in union",
			},
		},
		{
			input: "[0,(@.length-1)]",
			expected: expected{

				token: &unionToken{
					arguments: []interface{}{
						0,
						&scriptToken{expression: "@.length-1"},
					},
				},
			},
		},
		{
			input: "[0,'one',2,(1+2)]",
			expected: expected{
				token: &unionToken{
					arguments: []interface{}{
						0,
						"one",
						2,
						&scriptToken{expression: "1+2"},
					},
				},
			},
		},
		{
			input: "[(@.length-2),(@.length-1),1]",
			expected: expected{
				token: &unionToken{
					arguments: []interface{}{
						&scriptToken{expression: "@.length-2"},
						&scriptToken{expression: "@.length-1"},
						1,
					},
				},
			},
		},
		{
			input: "[:2]",
			expected: expected{
				token: &rangeToken{
					from: nil,
					to:   2,
					step: 1,
				},
			},
		},
		{
			input: "[0:(@.length-1):2]",
			expected: expected{
				token: &rangeToken{
					from: 0,
					to:   &scriptToken{expression: "@.length-1"},
					step: 2,
				},
			},
		},
		{
			input: "[(@.length-1):1:2]",
			expected: expected{
				token: &rangeToken{
					from: &scriptToken{expression: "@.length-1"},
					to:   1,
					step: 2,
				},
			},
		},
		{
			input: "['store']",
			expected: expected{
				token: &keyToken{
					key: "store",
				},
			},
		},
		{
			input: "[store]",
			expected: expected{
				err: "invalid token. unexpected string",
			},
		},
		{
			input: "[store,book]",
			expected: expected{
				err: "invalid token. unexpected union argument",
			},
		},
		{
			input: "[(1+2*(3+4)+5')]",
			expected: expected{
				token: &scriptToken{
					expression: "1+2*(3+4)+5'",
				},
			},
		},
		{
			input: "['this key has brackets ( and colons : and commas , but is not a union, range, or script']",
			expected: expected{
				token: &keyToken{
					key: "this key has brackets ( and colons : and commas , but is not a union, range, or script",
				},
			},
		},
		{
			input: "[1,2:4]",
			expected: expected{
				err: "invalid token. cannot specify a range in a union",
			},
		},
		{
			input: "[1:2:3:]",
			expected: expected{
				err: "invalid token. incorrect number of arguments in range",
			},
		},
		{
			input: "[1:2:3:4]",
			expected: expected{
				err: "invalid token. incorrect number of arguments in range",
			},
		},
		{
			input: "['key':'end]",
			expected: expected{
				err: "invalid token. only integer or scripts allowed in range arguments",
			},
		},
		{
			input: "[::2]",
			expected: expected{
				token: &rangeToken{
					from: nil,
					to:   nil,
					step: 2,
				},
			},
		},
		{
			input: "[:end:2]",
			expected: expected{
				err: "invalid token. only integer or scripts allowed in range arguments",
			},
		},
		{
			input: "length",
			expected: expected{
				token: &lengthToken{},
			},
		},
		{
			input: "[length]",
			expected: expected{
				err: "invalid token. unexpected string",
			},
		},
		{
			input: "['length']",
			expected: expected{
				token: &keyToken{
					key: "length",
				},
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			token, err := Parse(test.input)

			if test.expected.err == "" {
				assert.Nil(t, err, fmt.Sprintf("input '%s' err check failed. expected 'nil' actual '%v'", test.input, err))
			} else {
				assert.EqualError(t, err, test.expected.err, fmt.Sprintf("input '%s' err check failed. expected '%s' actual '%v'", test.input, test.expected.err, err))
			}
			assert.Equal(t, test.expected.token, token, fmt.Sprintf("input '%s' token check failed. expected '%v' actual '%v'", test.input, test.expected.token, token))
		})
	}

}

func Test_Tokenize(t *testing.T) {

	type expected struct {
		tokens    []string
		remainder string
		err       error
	}

	tests := []struct {
		input    string
		expected expected
	}{
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
				err:    errors.ErrInvalidInitialToken,
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
				err: errors.ErrInvalidInitialToken,
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
				err: errors.ErrQueryNotSpecified,
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
			tokens, remainder, err := Tokenize(test.input)
			assert.Equal(t, test.expected.err, err, "unexpected error for %s", test.input)

			for i, actual := range tokens {
				expected := test.expected.tokens[i]
				assert.EqualValues(t, expected, actual, "unexpected token for %s", test.input)
			}

			assert.Equal(t, test.expected.remainder, remainder, "unexpected remainder for %s", test.input)
		})
	}

}
