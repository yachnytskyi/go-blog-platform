package common

import (
	"time"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// FormatDate formats a time.Time instance to a custom string format.
func FormatDate(data time.Time) string {
	return data.Format(constants.DateTimeFormat)
}

// ParseDate parses a string formatted as "02/01/06 15:04 Mon" into a time.Time instance.
func ParseDate(location, data string) (time.Time, error) {
	parsedTime, parseError := time.Parse(constants.DateTimeFormat, data)
	if validator.IsError(parseError) {
		// If an error occurs, create an internal error with context and log it.
		internalError := domainError.NewInternalError(location+".ParseDate.time.Parse", parseError.Error())
		logging.Logger(internalError)
		// Return a default date and the error.
		return time.Time{}, internalError
	}

	return parsedTime, nil
}
