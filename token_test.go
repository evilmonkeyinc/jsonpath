package jsonpath

import (
	"fmt"
	"testing"

	"github.com/evilmokeyinc/jsonpath/errors"
	"github.com/stretchr/testify/assert"
)

func Test_tokenize(t *testing.T) {

	type expected struct {
		tokens []string
		err    error
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
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			tokens, err := tokenize(test.input)
			assert.Equal(t, test.expected.err, err)

			for i, actual := range tokens {
				expected := test.expected.tokens[i]
				assert.EqualValues(t, expected, actual)
			}
		})
	}

}

func Test_parseToken(t *testing.T) {

	type expected struct {
		token *token
		err   string
	}

	tests := []struct {
		input    string
		expected expected
	}{
		{
			input: "",
			expected: expected{
				err: "invalid token: token can not be empty",
			},
		},
		{
			input: "['fail'",
			expected: expected{
				err: "invalid token: missing subscript close",
			},
		},
		{
			input: "[ ]",
			expected: expected{
				err: "invalid token: empty subscript",
			},
		},
		{
			input: "[?]",
			expected: expected{
				err: "invalid token: expected filter '?(' prefix and ')' suffix",
			},
		},
		{
			input: "[why though]",
			expected: expected{
				err: "invalid token: unexpected space",
			},
		},
		{
			input: "[1'2']",
			expected: expected{
				err: "invalid token: unexpected single quote",
			},
		},
		{
			input: "[1(@.length)]",
			expected: expected{
				err: "invalid token: invalid script format",
			},
		},
		{
			input: "$",
			expected: expected{
				token: &token{
					operation: tokenOperationRoot,
				},
			},
		},
		{
			input: "store",
			expected: expected{
				token: &token{
					operation: tokenOperationKey,
					arguments: []interface{}{"'store'"},
				},
			},
		},
		{
			input: "1",
			expected: expected{
				err: "invalid token: index specified as key",
			},
		},
		{
			input: "book",
			expected: expected{
				token: &token{
					operation: tokenOperationKey,
					arguments: []interface{}{"'book'"},
				},
			},
		},
		{
			input: "*",
			expected: expected{
				token: &token{
					operation: tokenOperationWildcard,
				},
			},
		},
		{
			input: "author",
			expected: expected{
				token: &token{
					operation: tokenOperationKey,
					arguments: []interface{}{"'author'"},
				},
			},
		},
		{
			input: "@",
			expected: expected{
				token: &token{
					operation: tokenOperationCurrent,
				},
			},
		},
		{
			input: "[*]",
			expected: expected{
				token: &token{
					operation: tokenOperationWildcard,
				},
			},
		},
		{
			input: "..",
			expected: expected{
				token: &token{
					operation: tokenOperationRecursive,
				},
			},
		},
		{
			input: "[?(@.isbn)]",
			expected: expected{
				token: &token{
					operation: tokenOperationFilter,
					arguments: []interface{}{"?(@.isbn)"},
				},
			},
		},
		{
			input: "[2]",
			expected: expected{
				token: &token{
					operation: tokenOperationIndex,
					arguments: []interface{}{2},
				},
			},
		},
		{
			input: "[(@.length-1)]",
			expected: expected{
				token: &token{
					operation: tokenOperationScript,
					arguments: []interface{}{"(@.length-1)"},
				},
			},
		},
		{
			input: "[-1:]",
			expected: expected{
				token: &token{
					operation: tokenOperationRange,
					arguments: []interface{}{-1, nil},
				},
			},
		},
		{
			input: "[1:(@.length-1)]",
			expected: expected{

				token: &token{
					operation: tokenOperationRange,
					arguments: []interface{}{1, "(@.length-1)"},
				},
			},
		},
		{
			input: "[0,1]",
			expected: expected{
				token: &token{
					operation: tokenOperationUnion,
					arguments: []interface{}{0, 1},
				},
			},
		},
		{
			input: "['first','last']",
			expected: expected{
				token: &token{
					operation: tokenOperationUnion,
					arguments: []interface{}{"'first'", "'last'"},
				},
			},
		},
		{
			input: "[0,]",
			expected: expected{
				err: "invalid token: empty argument in union",
			},
		},
		{
			input: "[,1]",
			expected: expected{
				err: "invalid token: empty argument in union",
			},
		},
		{
			input: "[(0),]",
			expected: expected{
				err: "invalid token: empty argument in union",
			},
		},
		{
			input: "[0,'1',]",
			expected: expected{
				err: "invalid token: empty argument in union",
			},
		},
		{
			input: "[0,(@.length-1)]",
			expected: expected{

				token: &token{
					operation: tokenOperationUnion,
					arguments: []interface{}{0, "(@.length-1)"},
				},
			},
		},
		{
			input: "[0,'one',2,(1+2)]",
			expected: expected{

				token: &token{
					operation: tokenOperationUnion,
					arguments: []interface{}{0, "'one'", 2, "(1+2)"},
				},
			},
		},
		{
			input: "[(@.length-2),(@.length-1),1]",
			expected: expected{

				token: &token{
					operation: tokenOperationUnion,
					arguments: []interface{}{"(@.length-2)", "(@.length-1)", 1},
				},
			},
		},
		{
			input: "[:2]",
			expected: expected{
				token: &token{
					operation: tokenOperationRange,
					arguments: []interface{}{nil, 2},
				},
			},
		},
		{
			input: "[0:(@.length-1):2]",
			expected: expected{
				token: &token{
					operation: tokenOperationRange,
					arguments: []interface{}{0, "(@.length-1)", 2},
				},
			},
		},
		{
			input: "[(@.length-1):1:2]",
			expected: expected{
				token: &token{
					operation: tokenOperationRange,
					arguments: []interface{}{"(@.length-1)", 1, 2},
				},
			},
		},
		{
			input: "['store']",
			expected: expected{
				token: &token{
					operation: tokenOperationKey,
					arguments: []interface{}{"'store'"},
				},
			},
		},
		{
			input: "[store]",
			expected: expected{
				err: "invalid token: unexpected string",
			},
		},
		{
			input: "[store,book]",
			expected: expected{
				err: "invalid token: unexpected union argument",
			},
		},
		{
			input: "[(1+2*(3+4)+5')]",
			expected: expected{
				token: &token{
					operation: tokenOperationScript,
					arguments: []interface{}{"(1+2*(3+4)+5')"},
				},
			},
		},
		{
			input: "['this key has brackets ( and colons : and commas , but is not a union, range, or script']",
			expected: expected{
				token: &token{
					operation: tokenOperationKey,
					arguments: []interface{}{"'this key has brackets ( and colons : and commas , but is not a union, range, or script'"},
				},
			},
		},
		{
			input: "[1,2:4]",
			expected: expected{
				err: "invalid token: cannot specify a range in a union",
			},
		},
		{
			input: "[1:2:3:]",
			expected: expected{
				err: "invalid token: incorrect number of arguments in range",
			},
		},
		{
			input: "[1:2:3:4]",
			expected: expected{
				err: "invalid token: incorrect number of arguments in range",
			},
		},
		{
			input: "['key':'end]",
			expected: expected{
				err: "invalid token: only integer or scripts allowed in range arguments",
			},
		},
		{
			input: "[::2]",
			expected: expected{
				token: &token{
					operation: tokenOperationRange,
					arguments: []interface{}{nil, nil, 2},
				},
			},
		},
		{
			input: "[:end:2]",
			expected: expected{
				err: "invalid token: only integer or scripts allowed in range arguments",
			},
		},
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			token, err := parseToken(test.input)

			if test.expected.err == "" {
				assert.Nil(t, err, fmt.Sprintf("input '%s' err check failed. expected 'nil' actual '%v'", test.input, err))
			} else {
				assert.EqualError(t, err, test.expected.err, fmt.Sprintf("input '%s' err check failed. expected '%s' actual '%v'", test.input, test.expected.err, err))
			}
			assert.Equal(t, test.expected.token, token, fmt.Sprintf("input '%s' token check failed. expected '%v' actual '%v'", test.input, test.expected.token, token))
		})
	}

}
