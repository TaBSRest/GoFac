package Helpers

import (
	"errors"
	"fmt"
)

func MakeError(location string, msg string) error {
	return errors.New(GetErrorMessage(location, msg))
}

func GetErrorMessage(location string, msg string) string {
	return fmt.Sprintf("%s: %s", location, msg)
}
