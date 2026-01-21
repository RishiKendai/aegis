package logger

import (
	"github.com/rs/zerolog"
)

func Init(level string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(parseLevel(level))
}

func parseLevel(level string) zerolog.Level {
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		return zerolog.InfoLevel
	}
	return l
}
