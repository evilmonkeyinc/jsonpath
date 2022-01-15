package jsonpath

import (
	"encoding/json"
	"strings"

	"github.com/evilmonkeyinc/jsonpath/token"
)

// Compile compile the JSON path query
func Compile(queryPath string, isStrict bool) (*JSONPath, error) {
	jsonPath := &JSONPath{
		options: &token.ParseOptions{IsStrict: isStrict},
	}
	if err := jsonPath.compile(queryPath); err != nil {
		return nil, getInvalidJSONPathQueryWithReason(queryPath, err)
	}

	return jsonPath, nil
}

// Query will return the result of the JSONPath query applied against the specified JSON data.
func Query(queryPath string, jsonData interface{}) (interface{}, error) {
	jsonPath, err := Compile(queryPath, false)
	if err != nil {
		return nil, getInvalidJSONPathQueryWithReason(queryPath, err)
	}
	return jsonPath.Query(jsonData)
}

// QueryString will return the result of the JSONPath query applied against the specified JSON data.
func QueryString(queryPath string, jsonData string) (interface{}, error) {
	jsonPath, err := Compile(queryPath, false)
	if err != nil {
		return nil, getInvalidJSONPathQueryWithReason(queryPath, err)
	}
	return jsonPath.QueryString(jsonData)
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

	return nil
}

// Query will return the result of the JSONPath query applied against the specified JSON data.
func (query *JSONPath) Query(root interface{}) (interface{}, error) {
	if len(query.tokens) == 0 {
		return nil, getInvalidJSONPathQuery(query.queryString)
	}

	tokens := make([]token.Token, 0)
	if len(query.tokens) > 1 {
		tokens = query.tokens[1:]
	}

	found, err := query.tokens[0].Apply(root, root, tokens)
	if err != nil {
		return nil, err
	}
	return found, nil
}

// QueryString will return the result of the JSONPath query applied against the specified JSON data.
func (query *JSONPath) QueryString(jsonData string) (interface{}, error) {
	jsonData = strings.TrimSpace(jsonData)

	var root interface{}
	if strings.HasPrefix(jsonData, "{") && strings.HasSuffix(jsonData, "}") {
		root = make(map[string]interface{})
	} else if strings.HasPrefix(jsonData, "[") && strings.HasSuffix(jsonData, "]") {
		root = make([]interface{}, 0)
	} else {
		return nil, getInvalidJSONData(errDataIsUnexpectedTypeOrNil)
	}

	if err := json.Unmarshal([]byte(jsonData), &root); err != nil {
		return nil, getInvalidJSONData(err)
	}

	return query.Query(root)
}
