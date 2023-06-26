package errs

import (
	"errors"
	"fmt"
)

type errorString string

func (e errorString) Error() string {
	return string(e)
}

func New(s string) error {
	return errorString(s)
}

var (
	Unwrap = errors.Unwrap
	Is     = errors.Is
	As     = errors.As
	Join   = errors.Join
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Catch(err *error) {
	if re := recover(); re != nil {
		if e, ok := re.(error); ok {
			*err = e
		} else {
			*err = fmt.Errorf("%v", re)
		}
	}
}
