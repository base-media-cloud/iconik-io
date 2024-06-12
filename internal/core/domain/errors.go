package domain

import "errors"

var (
	// ErrTransformingHeaderValue is an error that is returned when failing to transform a header value to a string.
	ErrTransformingHeaderValue = errors.New("cannot transform header value to string")
	// ErrTransformingHeaderKey is an error that is returned when failing to transform a header key to a string.
	ErrTransformingHeaderKey = errors.New("cannot transform header key to string")
)
