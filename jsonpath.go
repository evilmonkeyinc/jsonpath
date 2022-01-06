package jsonpath

import (
	"encoding/json"
	"fmt"

	"github.com/evilmokeyinc/jsonpath/token"
)

func Find(queryPath string, jsonData string) (interface{}, error) {
	jsonPath, err := Compile(queryPath)
	if err != nil {
		return nil, err
	}
	return jsonPath.Find(jsonData)
}

func Compile(queryPath string) (*JSONPath, error) {
	jsonPath := &JSONPath{}
	if err := jsonPath.compile(queryPath); err != nil {
		return nil, err
	}

	return jsonPath, nil
}

type JSONPath struct {
	queryPath string
	tokens    []token.Token
	// TODO : marshaller should be customizable
}

func (query *JSONPath) compile(queryPath string) error {
	query.queryPath = queryPath

	tokenStrings, _, err := token.Tokenize(queryPath)
	if err != nil {
		return err
	}

	tokens := make([]token.Token, len(tokenStrings))
	for idx, tokenString := range tokenStrings {
		token, err := token.Parse(tokenString)
		if err != nil {
			return err
		}
		tokens[idx] = token
	}
	query.tokens = tokens

	return nil
}

func (query *JSONPath) Find(jsonData string) (interface{}, error) {

	root := make(map[string]interface{})
	// TODO : should be able to pass in marshaller to JSONPath
	if err := json.Unmarshal([]byte(jsonData), &root); err != nil {
		return nil, fmt.Errorf("json marshalling error : %s", err.Error())
	}

	if len(query.tokens) == 0 {
		// TODO : var error
		return nil, fmt.Errorf("invalid query: no tokens")
	}

	found, err := query.tokens[0].Apply(root, root, query.tokens[1:])
	if err != nil {
		return nil, err
	}
	if array, ok := found.([]interface{}); ok {
		return array, nil
	}
	return found, nil
}
