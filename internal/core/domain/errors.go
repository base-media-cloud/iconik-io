package domain

import (
	"errors"
	"fmt"
)

var (
	// ErrInternalError is an error that is returned when an internal error occurs.
	ErrInternalError = errors.New("there was an internal error")
)

type wrappedErrs struct {
	errs interface{}
}

func (w *wrappedErrs) Error() string {
	return fmt.Sprintf("%v", w.errs)
}

func NewWrappedErrs(errs interface{}) *wrappedErrs {
	return &wrappedErrs{
		errs: errs,
	}
}
