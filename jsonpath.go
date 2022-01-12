package jsonpath

import (
	"encoding/json"

	"github.com/evilmonkeyinc/jsonpath/token"
)

// Find will return the result of the JSONPath query applied against the specified JSON data.
func Find(queryPath string, jsonData map[string]interface{}) (interface{}, error) {
	jsonPath, err := Compile(queryPath)
	if err != nil {
		// TODO : wrap?
		return nil, err
	}
	return jsonPath.Find(jsonData)
}

// FindFromJSONString will return the result of the JSONPath query applied against the specified JSON data.
func FindFromJSONString(queryPath string, jsonData string) (interface{}, error) {
	jsonPath, err := Compile(queryPath)
	if err != nil {
		// TODO : wrap?
		return nil, err
	}
	return jsonPath.FindFromJSONString(jsonData)
}

// Compile compile the JSON path query
func Compile(queryPath string) (*JSONPath, error) {
	jsonPath := &JSONPath{}
	if err := jsonPath.compile(queryPath); err != nil {
		// TODO : wrap?
		return nil, err
	}

	return jsonPath, nil
}

// JSONPath i need to expand this
type JSONPath struct {
	queryString string
	tokens      []token.Token
	options     *token.ParseOptions
}

func (query *JSONPath) compile(queryString string) error {
	query.queryString = queryString

	tokenStrings, _, err := token.Tokenize(queryString)
	if err != nil {
		return err
	}

	tokens := make([]token.Token, len(tokenStrings))
	for idx, tokenString := range tokenStrings {
		token, err := token.Parse(tokenString, query.options)
		if err != nil {
			return err
		}
		tokens[idx] = token
	}
	query.tokens = tokens

	if len(tokens) == 0 {
		return getInvalidJSONPathQuery(queryString)
	}

	return nil
}

// FindFromJSONString will return the result of the JSONPath query applied against the specified JSON data.
func (query *JSONPath) FindFromJSONString(jsonData string) (interface{}, error) {
	root := make(map[string]interface{})
	if err := json.Unmarshal([]byte(jsonData), &root); err != nil {
		return nil, getInvalidJSONData(err)
	}

	return query.Find(root)
}

// Find will return the result of the JSONPath query applied against the specified JSON data.
func (query *JSONPath) Find(jsonData map[string]interface{}) (interface{}, error) {
	if len(query.tokens) == 0 {
		return nil, getInvalidJSONPathQuery(query.queryString)
	}

	found, err := query.tokens[0].Apply(jsonData, jsonData, query.tokens[1:])
	if err != nil {
		// TODO : wrap?
		return nil, err
	}
	if array, ok := found.([]interface{}); ok {
		return array, nil
	}
	return found, nil
}
