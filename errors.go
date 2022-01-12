package jsonpath

import (
	"fmt"

	"github.com/evilmonkeyinc/jsonpath/errors"
)

func getInvalidJSONPathQuery(query string) error {
	return fmt.Errorf("%w '%s'", errors.ErrInvalidJSONPathQuery, query)
}

func getInvalidJSONData(reason error) error {
	return fmt.Errorf("%s. %s", errors.ErrInvalidJSONData, reason.Error())
}
