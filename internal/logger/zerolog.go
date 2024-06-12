/*
Package logger contains the logger initialisation with helper functions.
The logger we use is zerolog.
*/
package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// New is a function that returns a new instance of the zerolog.Logger struct.
func New() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return zerolog.New(os.Stderr).With().Timestamp().Logger()
}
