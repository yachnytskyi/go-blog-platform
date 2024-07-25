package logger

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/logger/model"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

type Zerolog struct {
	logger zerolog.Logger
}

func NewZerolog() *Zerolog {
	// Set the time format and global log level for zerolog.
	zerolog.TimeFieldFormat = constants.LoggerDateTimeFormat
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	// Create a new zerolog logger with console output.
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &Zerolog{logger: logger}
}

// logWithLevel handles the common logic for logging at different levels.
func (l *Zerolog) logWithLevel(level zerolog.Level, data error) {
	data = model.HandleError(data)
	jsonData, marshalError := json.Marshal(data)
	if validator.IsError(marshalError) {
		l.logger.Error().Err(marshalError).Msg("")
		return
	}

	// Log the data at the specified level.
	switch level {
	case zerolog.TraceLevel:
		l.logger.Trace().RawJSON("trace", jsonData).Msg("")
	case zerolog.DebugLevel:
		l.logger.Debug().RawJSON("debug", jsonData).Msg("")
	case zerolog.InfoLevel:
		l.logger.Info().RawJSON("info", jsonData).Msg("")
	case zerolog.WarnLevel:
		l.logger.Warn().RawJSON("warn", jsonData).Msg("")
	case zerolog.ErrorLevel:
		l.logger.Error().RawJSON("error", jsonData).Msg("")
	case zerolog.FatalLevel:
		l.logger.Fatal().RawJSON("fatal", jsonData).Msg("")
	case zerolog.PanicLevel:
		l.logger.Panic().RawJSON("panic", jsonData).Msg("")
	}
}

func (l *Zerolog) Trace(data error) {
	l.logWithLevel(zerolog.TraceLevel, data)
}

func (l *Zerolog) Debug(data error) {
	l.logWithLevel(zerolog.DebugLevel, data)
}

func (l *Zerolog) Info(data error) {
	l.logWithLevel(zerolog.InfoLevel, data)
}

func (l *Zerolog) Warn(data error) {
	l.logWithLevel(zerolog.WarnLevel, data)
}

func (l *Zerolog) Error(data error) {
	l.logWithLevel(zerolog.ErrorLevel, data)
}

func (l *Zerolog) Fatal(data error) {
	l.logWithLevel(zerolog.FatalLevel, data)
}

func (l *Zerolog) Panic(data error) {
	l.logWithLevel(zerolog.PanicLevel, data)
}
