package jsonpath

import (
	"encoding/json"

	"github.com/evilmokeyinc/jsonpath/errors"
	"github.com/evilmokeyinc/jsonpath/token"
)

func Find(queryPath string, jsonData map[string]interface{}) (interface{}, error) {
	jsonPath, err := Compile(queryPath)
	if err != nil {
		return nil, err
	}
	return jsonPath.Find(jsonData)
}

func FindFromJSONString(queryPath string, jsonData string) (interface{}, error) {
	jsonPath, err := Compile(queryPath)
	if err != nil {
		return nil, err
	}
	return jsonPath.FindFromJSONString(jsonData)
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

	if len(tokens) == 0 {
		return errors.ErrInvalidQueryNoTokens
	}

	return nil
}

func (query *JSONPath) FindFromJSONString(jsonData string) (interface{}, error) {

	root := make(map[string]interface{})
	if err := json.Unmarshal([]byte(jsonData), &root); err != nil {
		return nil, errors.GetJSONMarshalFailedError(err)
	}

	return query.Find(root)
}

func (query *JSONPath) Find(jsonData map[string]interface{}) (interface{}, error) {

	if len(query.tokens) == 0 {
		return nil, errors.ErrInvalidQueryNoTokens
	}

	found, err := query.tokens[0].Apply(jsonData, jsonData, query.tokens[1:])
	if err != nil {
		return nil, err
	}
	if array, ok := found.([]interface{}); ok {
		return array, nil
	}
	return found, nil
}
