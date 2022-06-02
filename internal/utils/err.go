package utils

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func ParseErrorsFromChannel(errChannel <-chan error) error {
	if len(errChannel) == 0 {
		return nil
	}

	var errorMessageBuilder strings.Builder
	errorNumber := 1
	for currentError := range errChannel {
		_, _ = fmt.Fprintf(&errorMessageBuilder, "%d: %s;", errorNumber, currentError.Error())
		errorNumber++
	}

	return errors.New(errorMessageBuilder.String())
}
