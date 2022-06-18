package presenter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/smithy-go"
)

func WrapError(err error, code int) error {
	return fmt.Errorf("%v##%w", code, err)
}

func WrapAPIError(err smithy.APIError, code int) error {
	return fmt.Errorf("%v##%w", code, errors.New(err.ErrorMessage()))
}

func UnwrapCodeError(err error) (int, error) {
	parts := strings.SplitN(err.Error(), "##", 2)
	code, parseErr := strconv.Atoi(parts[0])
	if parseErr != nil {
		return 500, err
	}
	return code, errors.Unwrap(err)
}
