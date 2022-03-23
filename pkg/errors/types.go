package errors

import (
	"fmt"
)

func New(text string) error {
	return fmt.Errorf(text)
}

func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}
