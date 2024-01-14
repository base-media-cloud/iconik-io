package domain

import "errors"

var (
	// ErrInternalError is an error that is returned when an internal error occurs.
	ErrInternalError = errors.New("there was an internal error")
)
