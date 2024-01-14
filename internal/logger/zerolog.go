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
