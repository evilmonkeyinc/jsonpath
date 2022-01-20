package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test filterToken struct conforms to Token interface
var _ Token = &filterToken{}

func Test_newFilterToken(t *testing.T) {
	assert.IsType(t, &filterToken{}, newFilterToken("", nil))
}

func Test_FilterToken_String(t *testing.T) {
	tests := []*tokenStringTest{
		{
			input:    &filterToken{},
			expected: "[?()]",
		},
		{
			input:    &filterToken{expression: "true"},
			expected: "[?(true)]",
		},
		{
			input:    &filterToken{expression: "@.length<0"},
			expected: "[?(@.length<0)]",
		},
	}

	batchTokenStringTests(t, tests)
}

func Test_FilterToken_Type(t *testing.T) {
	assert.Equal(t, "filter", (&filterToken{}).Type())
}

func Test_FilterToken_Apply(t *testing.T) {

	tests := []*tokenTest{
		{
			token: &filterToken{expression: ""},
			input: input{
				current: []interface{}{"one"},
			},
			expected: expected{
				err: "invalid expression. is empty",
			},
		},
		{
			token: &filterToken{
				expression: "@.isbn",
			},
			input: input{
				current: nil,
			},
			expected: expected{
				err: "filter: invalid token target. expected [array map slice] got [nil]",
			},
		},
		{
			token: &filterToken{
				expression: "@.isbn",
			},
			input: input{
				current: "this is a string",
			},
			expected: expected{
				err: "filter: invalid token target. expected [array map slice] got [string]",
			},
		},
		{
			token: &filterToken{
				expression: "@.isbn",
			},
			input: input{
				current: []map[string]interface{}{
					{
						"isbn": 5,
						"name": "one",
					},
					{
						"isbn": "",
						"name": "two",
					},
					{
						"isbn": "string",
						"name": "three",
					},
					{
						"name": "four",
					},
				},
			},
			expected: expected{
				value: []interface{}{
					map[string]interface{}{
						"isbn": 5,
						"name": "one",
					},
					map[string]interface{}{
						"isbn": "string",
						"name": "three",
					},
				},
			},
		},
		{
			token: &filterToken{
				expression: "@.price<10",
			},
			input: input{
				current: []map[string]interface{}{
					{
						"price": 5,
						"name":  "one",
					},
					{
						"price": 9.99,
						"name":  "two",
					},
					{
						"price": 15,
						"name":  "three",
					},
				},
			},
			expected: expected{
				value: []interface{}{
					map[string]interface{}{
						"price": 5,
						"name":  "one",
					},
					map[string]interface{}{
						"price": 9.99,
						"name":  "two",
					},
				},
			},
		},
		{
			token: &filterToken{
				expression: "@.price<$.expensive",
			},
			input: input{
				root: map[string]interface{}{
					"expensive": 10,
				},
				current: []map[string]interface{}{
					{
						"price": 5,
						"name":  "one",
					},
					{
						"price": 9.99,
						"name":  "two",
					},
					{
						"price": 15,
						"name":  "three",
					},
				},
			},
			expected: expected{
				value: []interface{}{
					map[string]interface{}{
						"price": 5,
						"name":  "one",
					},
					map[string]interface{}{
						"price": 9.99,
						"name":  "two",
					},
				},
			},
		},
		{
			token: &filterToken{
				expression: "@.price<$.expensive",
			},
			input: input{
				root: map[string]interface{}{
					"expensive": 10,
				},
				current: map[string]interface{}{
					"one": map[string]interface{}{
						"price": 5,
						"name":  "one",
					},
					"two": map[string]interface{}{
						"price": 9.99,
						"name":  "two",
					},
					"three": map[string]interface{}{
						"price": 15,
						"name":  "three",
					},
					"four": map[string]interface{}{
						"name": "three",
					},
				},
			},
			expected: expected{
				value: []interface{}{
					map[string]interface{}{
						"price": 5,
						"name":  "one",
					},
					map[string]interface{}{
						"price": 9.99,
						"name":  "two",
					},
				},
			},
		},
		{
			token: &filterToken{
				expression: "@.price<$.expensive",
			},
			input: input{
				root: map[string]interface{}{
					"expensive": 10,
				},
				current: []map[string]interface{}{
					{
						"price": 5,
						"name":  "one",
					},
					{
						"price": 9.99,
						"name":  "two",
					},
					{
						"price": 15,
						"name":  "three",
					},
				},
				tokens: []Token{
					&indexToken{index: 1},
				},
			},
			expected: expected{
				value: map[string]interface{}{
					"price": 9.99,
					"name":  "two",
				},
			},
		},
		{
			token: &filterToken{
				expression: "@.price<$.expensive",
			},
			input: input{
				root: map[string]interface{}{
					"expensive": 10,
				},
				current: map[string]interface{}{
					"one": map[string]interface{}{
						"price": 5,
						"name":  "one",
					},
					"two": map[string]interface{}{
						"price": 9.99,
						"name":  "two",
					},
					"three": map[string]interface{}{
						"price": 15,
						"name":  "three",
					},
					"four": map[string]interface{}{
						"name": "three",
					},
				},
				tokens: []Token{
					&keyToken{key: "name"},
				},
			},
			expected: expected{
				value: []interface{}{
					"one",
					"two",
				},
			},
		},
	}

	batchTokenTests(t, tests)

}
