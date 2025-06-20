package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(debug bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
