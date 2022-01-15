package jsonpath

import (
	"encoding/json"

	"github.com/evilmonkeyinc/jsonpath/token"
)

// Query will return the result of the JSONPath query applied against the specified JSON data.
func Query(queryPath string, jsonData map[string]interface{}) (interface{}, error) {
	jsonPath, err := Compile(queryPath, false)
	if err != nil {
		return nil, err
	}
	return jsonPath.Query(jsonData)
}

// QueryString will return the result of the JSONPath query applied against the specified JSON data.
func QueryString(queryPath string, jsonData string) (interface{}, error) {
	jsonPath, err := Compile(queryPath, false)
	if err != nil {
		return nil, err
	}
	return jsonPath.QueryString(jsonData)
}

// Compile compile the JSON path query
func Compile(queryPath string, isStrict bool) (*JSONPath, error) {
	jsonPath := &JSONPath{
		options: &token.ParseOptions{IsStrict: isStrict},
	}
	if err := jsonPath.compile(queryPath); err != nil {
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

// QueryString will return the result of the JSONPath query applied against the specified JSON data.
func (query *JSONPath) QueryString(jsonData string) (interface{}, error) {
	root := make(map[string]interface{})
	if err := json.Unmarshal([]byte(jsonData), &root); err != nil {
		return nil, getInvalidJSONData(err)
	}

	return query.Query(root)
}

// Query will return the result of the JSONPath query applied against the specified JSON data.
func (query *JSONPath) Query(jsonData map[string]interface{}) (interface{}, error) {
	if len(query.tokens) == 0 {
		return nil, getInvalidJSONPathQuery(query.queryString)
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
