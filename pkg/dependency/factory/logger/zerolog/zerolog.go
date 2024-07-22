package zerolog

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/logger/model"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

type Zerolog struct {
	logger zerolog.Logger
}

func NewZerolog() *Zerolog {
	zerolog.TimeFieldFormat = constants.DateTimeFormat
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: constants.DateTimeFormat})
	return &Zerolog{logger: logger}
}

func (l *Zerolog) Trace(data error) {
	data = httpModel.HandleError(data)
	jsonData, marshalError := json.Marshal(data)
	if validator.IsError(marshalError) {
		l.logger.Error().Err(marshalError)
		return
	}

	l.logger.Warn().RawJSON("", jsonData).Msg("")
}

func (l *Zerolog) Debug(data error) {
	data = httpModel.HandleError(data)
	jsonData, marshalError := json.Marshal(data)
	if validator.IsError(marshalError) {
		l.logger.Error().Err(marshalError)
		return
	}

	l.logger.Warn().RawJSON("", jsonData).Msg("")
}

func (l *Zerolog) Info(data error) {
	data = httpModel.HandleError(data)
	jsonData, marshalError := json.Marshal(data)
	if validator.IsError(marshalError) {
		l.logger.Error().Err(marshalError)
		return
	}

	l.logger.Warn().RawJSON("", jsonData).Msg("")
}

func (l *Zerolog) Warn(data error) {
	data = httpModel.HandleError(data)
	jsonData, marshalError := json.Marshal(data)
	if validator.IsError(marshalError) {
		l.logger.Error().Err(marshalError)
		return
	}

	l.logger.Warn().RawJSON("", jsonData).Msg("")
}

func (l *Zerolog) Error(data error) {
	data = httpModel.HandleError(data)
	jsonData, marshalError := json.Marshal(data)
	if validator.IsError(marshalError) {
		l.logger.Error().Err(marshalError)
		return
	}

	l.logger.Error().RawJSON("", jsonData).Msg("")
}

func (l *Zerolog) Fatal(data error) {
	data = httpModel.HandleError(data)
	jsonData, marshalError := json.Marshal(data)
	if validator.IsError(marshalError) {
		l.logger.Error().Err(marshalError)
		return
	}

	l.logger.Warn().RawJSON("", jsonData).Msg("")
}

func (l *Zerolog) Panic(data error) {
	data = httpModel.HandleError(data)
	jsonData, marshalError := json.Marshal(data)
	if validator.IsError(marshalError) {
		l.logger.Error().Err(marshalError)
		return
	}

	l.logger.Panic().RawJSON("", jsonData).Msg("")
}
