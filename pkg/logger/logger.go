package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func InitLogger() *zerolog.Logger {
	Logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) //zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	return &Logger
}
