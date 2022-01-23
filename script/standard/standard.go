package standard

import (
	"encoding/json"
	"strings"

	"github.com/evilmonkeyinc/jsonpath/option"
	"github.com/evilmonkeyinc/jsonpath/script"
)

// TODO : add support for bitwise operators | &^ ^ &  << >> after + and -
var defaultTokens []string = []string{
	"||", "&&",
	"==", "!=", "<=", ">=", "<", ">", "=~",
	"+", "-",
	"**", "*", "/", "%",
	"@", "$",
}

// ScriptEngine standard implementation of the script engine interface
type ScriptEngine struct {
}

// Compile returns a compiled expression that can be evaluated multiple times
func (engine *ScriptEngine) Compile(expression string, options *option.QueryOptions) (script.CompiledExpression, error) {
	tokens := defaultTokens
	operator, err := engine.buildOperators(expression, tokens, options)
	if err != nil {
		return nil, err
	}

	return &compiledExpression{
		expression:   expression,
		rootOperator: operator,
		engine:       engine,
		options:      options,
	}, nil
}

func (engine *ScriptEngine) buildOperators(expression string, tokens []string, options *option.QueryOptions) (operator, error) {
	nextToken := tokens[0]
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return nil, nil
	}
	if strings.HasPrefix(expression, "(") && strings.HasSuffix(expression, ")") {
		expression = strings.TrimSpace(expression[1 : len(expression)-1])
		// since we were in brackets, we need to try all the tokens again
		tokens = append([]string{""}, defaultTokens...)
	}

	idx := findUnquotedOperators(expression, nextToken)
	if idx < 0 {
		if len(tokens) == 1 {
			return nil, nil
		}
		// if none of these tokens, move onto next token
		return engine.buildOperators(expression, tokens[1:], options)
	}

	// check right for more tokens, or use raw string as input
	arg2Str := strings.TrimSpace(expression[idx+(len(nextToken)):])
	if arg2Str == "" && (nextToken != "$" && nextToken != "@") {
		if len(tokens) == 1 {
			return nil, nil
		}
		// if none of these tokens, move onto next token
		return engine.buildOperators(expression, tokens[1:], options)
	}

	var arg2 interface{} = arg2Str
	if op, err := engine.buildOperators(arg2Str, tokens, options); err != nil {
		return nil, err
	} else if op != nil {
		arg2 = op
	} else {
		arg2Str = strings.TrimSpace(arg2Str)
		arg2 = arg2Str
		if strings.HasPrefix(arg2Str, "[") && strings.HasSuffix(arg2Str, "]") {
			val := make([]interface{}, 0)
			if err := json.Unmarshal([]byte(arg2Str), &val); err == nil {
				arg2 = val
			}
		} else if strings.HasPrefix(arg2Str, "{") && strings.HasSuffix(arg2Str, "}") {
			val := make(map[string]interface{})
			if err := json.Unmarshal([]byte(arg2Str), &val); err == nil {
				arg2 = val
			}
		}
	}

	// check left for more tokens, or use raw string as input
	arg1Str := strings.TrimSpace(expression[0:idx])
	if arg1Str == "" && (nextToken != "$" && nextToken != "@") {
		if len(tokens) == 1 {
			return nil, nil
		}
		// if none of these tokens, move onto next token
		return engine.buildOperators(expression, tokens[1:], options)
	}

	var arg1 interface{} = arg1Str
	if op, err := engine.buildOperators(arg1Str, tokens, options); err != nil {
		return nil, err
	} else if op != nil {
		arg1 = op
	} else {
		arg1Str = strings.TrimSpace(arg1Str)
		arg1 = arg1Str
		if strings.HasPrefix(arg1Str, "[") && strings.HasSuffix(arg1Str, "]") {
			val := make([]interface{}, 0)
			if err := json.Unmarshal([]byte(arg1Str), &val); err == nil {
				arg1 = val
			}
		} else if strings.HasPrefix(arg1Str, "{") && strings.HasSuffix(arg1Str, "}") {
			val := make(map[string]interface{})
			if err := json.Unmarshal([]byte(arg1Str), &val); err == nil {
				arg1 = val
			}
		}
	}

	switch nextToken {
	case "@", "$":
		return newSelectorOperator(expression, engine, options)
	case "=~":
		return &regexOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "||":
		return &orOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "&&":
		return &andOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "==":
		return &equalsOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "!=":
		return &notEqualsOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "<=":
		return &lessThanOrEqualOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case ">=":
		return &greaterThanOrEqualOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "<":
		return &lessThanOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case ">":
		return &greaterThanOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "+":
		return &plusOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "-":
		return &subtractOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "**":
		return &powerOfOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "*":
		return &multiplyOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "/":
		return &divideOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	case "%":
		return &modulusOperator{
			arg1: arg1,
			arg2: arg2,
		}, nil
	}

	// will cover when we add a new operator token
	// but forget to update the switch/case
	return nil, errUnsupportedOperator
}

// Evaluate return the result of the expression evaluation
func (engine *ScriptEngine) Evaluate(root, current interface{}, expression string, options *option.QueryOptions) (interface{}, error) {
	compiled, err := engine.Compile(expression, options)
	if err != nil {
		return nil, err
	}
	return compiled.Evaluate(root, current)
}

type compiledExpression struct {
	expression   string
	rootOperator operator
	engine       *ScriptEngine
	options      *option.QueryOptions
}

func (compiled *compiledExpression) Evaluate(root, current interface{}) (interface{}, error) {
	expression := compiled.expression
	if expression == "" {
		return nil, getInvalidExpressionEmptyError()
	}
	parameters := map[string]interface{}{
		"$":    root,
		"@":    current,
		"nil":  nil,
		"null": nil,
	}

	if compiled.rootOperator == nil {
		if val, ok := parameters[expression]; ok {
			return val, nil
		}
		return expression, nil
	}

	value, err := compiled.rootOperator.Evaluate(parameters)
	if err != nil {
		return nil, err
	}

	return value, nil
}
