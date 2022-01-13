package jsonpath

import (
	"fmt"

	"github.com/evilmonkeyinc/jsonpath/errors"
)

func getInvalidJSONData(reason error) error {
	return fmt.Errorf("%w. %s", errors.ErrInvalidJSONData, reason.Error())
}

func getInvalidJSONPathQuery(query string) error {
	return fmt.Errorf("%w '%s'", errors.ErrInvalidJSONPathQuery, query)
}
