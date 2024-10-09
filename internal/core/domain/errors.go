package domain

import "errors"

var (
	// ErrInternalError is the error we return when something has gone wrong our end.
	ErrInternalError = errors.New("an internal error occurred")
	// ErrForbidden is an error that is returned when iconik returns a 403.
	ErrForbidden = errors.New("please check your app id and auth token are correct")
	// Err401Search is an error that is returned when user doesn't have correct permissions to search.
	Err401Search = errors.New("you do not have the correct permissions to search")
)
