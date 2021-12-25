package jsonpath

func Find(queryPath string, jsonData string) ([]interface{}, error) {
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
	tokens    []*token
}

func (query *JSONPath) compile(queryPath string) error {
	query.queryPath = queryPath

	tokenStrings, err := tokenize(queryPath)
	if err != nil {
		return err
	}

	tokens := make([]*token, len(tokenStrings))
	for idx, tokenString := range tokenStrings {
		token, err := parseToken(tokenString)
		if err != nil {
			return err
		}
		tokens[idx] = token
	}
	query.tokens = tokens

	return nil
}

func (query *JSONPath) Find(jsonData interface{}) ([]interface{}, error) {

	return nil, nil
}
