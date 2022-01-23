package standard

import (
	"fmt"
	"strconv"
	"strings"
)

type operator interface {
	Evaluate(parameters map[string]interface{}) (interface{}, error)
}

func getInteger(argument interface{}, parameters map[string]interface{}) (int64, error) {
	if argument == nil {
		return 0, errInvalidArgumentNil
	}
	if parameters == nil {
		parameters = make(map[string]interface{})
	}

	if sub, ok := argument.(operator); ok {
		arg, err := sub.Evaluate(parameters)
		if err != nil {
			return 0, err
		}
		argument = arg
	}

	if str, ok := argument.(string); ok {
		if arg, ok := parameters[str]; ok {
			argument = arg
		}

	}

	str := fmt.Sprintf("%v", argument)
	intVal, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, errInvalidArgumentExpectedInteger
	}
	return intVal, nil
}

func getNumber(argument interface{}, parameters map[string]interface{}) (float64, error) {
	if argument == nil {
		return 0, errInvalidArgumentNil
	}
	if parameters == nil {
		parameters = make(map[string]interface{})
	}

	if sub, ok := argument.(operator); ok {
		arg, err := sub.Evaluate(parameters)
		if err != nil {
			return 0, err
		}
		argument = arg
	}

	if str, ok := argument.(string); ok {
		if arg, ok := parameters[str]; ok {
			argument = arg
		}
	}

	str := fmt.Sprintf("%v", argument)
	floatVal, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, errInvalidArgumentExpectedNumber
	}
	return floatVal, nil
}

func getBoolean(argument interface{}, parameters map[string]interface{}) (bool, error) {
	if argument == nil {
		return false, errInvalidArgumentNil
	}
	if parameters == nil {
		parameters = make(map[string]interface{})
	}

	if sub, ok := argument.(operator); ok {
		arg, err := sub.Evaluate(parameters)
		if err != nil {
			return false, err
		}
		argument = arg
	}

	if str, ok := argument.(string); ok {
		if arg, ok := parameters[str]; ok {
			argument = arg
		}
	}

	str := fmt.Sprintf("%v", argument)
	boolValue, err := strconv.ParseBool(str)
	if err != nil {
		return false, errInvalidArgumentExpectedBoolean
	}
	return boolValue, nil
}

func getString(argument interface{}, parameters map[string]interface{}) (string, error) {
	if argument == nil {
		return "", errInvalidArgumentNil
	}
	if parameters == nil {
		parameters = make(map[string]interface{})
	}

	if sub, ok := argument.(operator); ok {
		arg, err := sub.Evaluate(parameters)
		if err != nil {
			return "", err
		}
		argument = arg
	}

	if str, ok := argument.(string); ok {
		if arg, ok := parameters[str]; ok {
			argument = arg
			if parsed, ok := arg.(string); ok {
				str = parsed

				if len(str) > 1 {
					if strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'") {
						return str, nil
					} else if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
						str = str[1 : len(str)-1]
					}
				}
				return fmt.Sprintf("'%s'", str), nil
			}
		} else {
			if len(str) > 1 {
				if strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'") {
					return str, nil
				} else if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
					str = str[1 : len(str)-1]
					return fmt.Sprintf("'%s'", str), nil
				}
			}
			return str, nil
		}
	}

	return fmt.Sprintf("%v", argument), nil
}
