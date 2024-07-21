package zerolog

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Zerolog struct {
	logger zerolog.Logger
}

func NewZerolog() *Zerolog {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	return &Zerolog{logger: logger}
}

func (l *Zerolog) Trace(data any) {
	l.logger.Trace().Interface("", data).Msg("")
}

func (l *Zerolog) Debug(data any) {
	l.logger.Debug().Interface("", data).Msg("")
}

func (l *Zerolog) Info(data any) {
	l.logger.Info().Interface("", data).Msg("")
}

func (l *Zerolog) Warn(data any) {
	l.logger.Warn().Interface("", data).Msg("")
}

func (l *Zerolog) Error(data any) {
	l.logger.Error().Interface("", data).Msg("")
}

func (l *Zerolog) Fatal(data any) {
	l.logger.Fatal().Interface("", data).Msg("")
}

func (l *Zerolog) Panic(data any) {
	l.logger.Panic().Interface("", data).Msg("")
}
