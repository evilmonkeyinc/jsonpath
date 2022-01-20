package token

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"regexp"
	"strconv"
	"strings"
)

func newExpressionToken(expression string, options *Options) *expressionToken {
	return &expressionToken{expression: expression, options: options}
}

type expressionToken struct {
	expression string
	options    *Options
}

func (token *expressionToken) String() string {
	return fmt.Sprintf("(%s)", token.expression)
}

func (token *expressionToken) Type() string {
	return "expression"
}

func (token *expressionToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	if token.expression == "" {
		return nil, getInvalidExpressionEmptyError()
	}

	value, err := evaluateExpression(root, current, token.expression, token.options)
	if err != nil {
		return nil, getInvalidExpressionError(err)
	}

	if len(next) > 0 {
		return next[0].Apply(root, value, next[1:])
	}

	return value, nil
}

// TODO : add extra support
/*
1. regex
*/

func evaluateExpression(root, current interface{}, expression string, options *Options) (interface{}, error) {
	if expression == "" {
		return nil, getInvalidExpressionEmptyError()
	}

	rootIndex := strings.Index(expression, "$")
	currentIndex := strings.Index(expression, "@")

	for rootIndex > -1 || currentIndex > -1 {

		query := ""
		if rootIndex > -1 {
			query = expression[rootIndex:]
		} else if currentIndex > -1 {
			query = expression[currentIndex:]
		}

		tokenStrings, remainder, err := tokenize(query, true)
		if err != nil {
			return nil, getInvalidExpressionError(err)
		}
		if remainder != "" {
			// shorten query to only what is being replaced
			query = query[0 : len(query)-len(remainder)]
		}
		if len(tokenStrings) > 0 {
			tokens := make([]Token, 0)
			for _, tokenString := range tokenStrings {
				token, err := Parse(tokenString, options)
				if err != nil {
					return nil, getInvalidExpressionError(err)
				}
				tokens = append(tokens, token)
			}

			value, err := tokens[0].Apply(root, current, tokens[1:])
			if err != nil {
				return nil, getInvalidExpressionError(err)
			}

			new := fmt.Sprintf("%v", value)
			if strValue, ok := value.(string); ok {
				new = fmt.Sprintf("'%s'", strValue)
			} else if intValue, ok := isInteger(value); ok {
				new = fmt.Sprintf("%d", intValue)
			} else if boolValue, ok := value.(bool); ok {
				new = fmt.Sprintf("%t", boolValue)
			} else if floatValue, ok := value.(float64); ok {
				new = strconv.FormatFloat(floatValue, 'f', -1, 64)
			}

			expression = strings.ReplaceAll(expression, query, new)
		}

		rootIndex = strings.Index(expression, "$")
		currentIndex = strings.Index(expression, "@")
	}

	expression = strings.TrimSpace(expression)

	// convert ' to " unless the ' is escaped \'
	regexp := regexp.MustCompile("(?:\\')+|(')")
	expression = regexp.ReplaceAllString(expression, "\"")

	fs := token.NewFileSet()
	tv, err := types.Eval(fs, nil, token.NoPos, expression)
	if err != nil {
		return nil, getInvalidExpressionError(err)
	}
	if tv.Value == nil {
		return nil, nil
	}
	switch tv.Value.Kind() {
	case constant.Bool:
		strValue := tv.Value.String()
		boolVal, _ := strconv.ParseBool(strValue)
		return boolVal, nil
	case constant.Float:
		strValue := tv.Value.String()
		floatVal, _ := strconv.ParseFloat(strValue, 64)
		return floatVal, nil
	case constant.Int:
		strValue := tv.Value.String()
		intVal, _ := strconv.ParseInt(strValue, 10, 64)
		return intVal, nil
	case constant.String, constant.Complex, constant.Unknown:
		fallthrough
	default:
		value := tv.Value.String()
		value = strings.ReplaceAll(value, "\"", "'")
		return value, nil
	}
}
