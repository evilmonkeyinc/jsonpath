package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FilterToken_Apply(t *testing.T) {

	type input struct {
		expression    string
		root, current interface{}
	}
	type expected struct {
		value []interface{}
		err   string
	}

	tests := []struct {
		input    input
		expected expected
	}{
		{
			input: input{},
			expected: expected{
				err: "invalid parameter. expression is empty",
			},
		},
		{
			input: input{
				expression: "@.isbn",
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
			input: input{
				expression: "@.price<10",
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
				expression: "@.price<$.expensive",
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
	}

	for idx, test := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			token := &filterToken{expression: test.input.expression}
			actual, err := token.Apply(test.input.root, test.input.current, nil)

			if test.expected.err == "" {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err, test.expected.err)
			}

			if test.expected.value != nil {
				assert.ElementsMatch(t, test.expected.value, actual)
			} else {
				assert.Nil(t, actual)
			}
		})
	}

}
