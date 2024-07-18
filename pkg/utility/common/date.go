package common

import (
	"time"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// FormatDate formats a time.Time instance to a custom string format.
func FormatDate(data time.Time) string {
	return data.Format(constants.DateTimeFormat)
}

// ParseDate parses a string into a time.Time instance.
func ParseDate(location, data string) commonModel.Result[time.Time] {
	parsedTime, parseError := time.Parse(constants.DateTimeFormat, data)
	if validator.IsError(parseError) {
		internalError := domainError.NewInternalError(location+".ParseDate.time.Parse", parseError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[time.Time](internalError)
	}

	return commonModel.NewResultOnSuccess[time.Time](parsedTime)
}
