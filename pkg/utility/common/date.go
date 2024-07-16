package common

import (
	"time"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// FormatDate formats a time.Time instance to a custom string format defined in constants.
//
// Parameters:
// - data: The time.Time instance to be formatted.
//
// Returns:
// - A string representing the formatted date.
func FormatDate(data time.Time) string {
	return data.Format(constants.DateTimeFormat)
}

// ParseDate parses a string formatted as defined in constants.DateTimeFormat into a time.Time instance.
//
// Parameters:
// - location: A string representing the location or context for error logging.
// - data: The string to be parsed into a time.Time instance.
//
// Returns:
// - commonModel.Result[time.Time]: The result containing the parsed time.Time instance if successful, or an error if the parsing fails.
func ParseDate(location, data string) commonModel.Result[time.Time] {
	parsedTime, parseError := time.Parse(constants.DateTimeFormat, data)
	if validator.IsError(parseError) {
		// Create and log an internal error with context if parsing fails.
		internalError := domainError.NewInternalError(location+".ParseDate.time.Parse", parseError.Error())
		logging.Logger(internalError)
		// Return the internal error.
		return commonModel.NewResultOnFailure[time.Time](internalError)
	}

	return commonModel.NewResultOnSuccess[time.Time](parsedTime)
}
