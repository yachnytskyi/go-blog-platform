package common

import (
	"encoding/base64"
	"strings"

	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// Encode encodes the input data to a base64 string.
func Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Decode decodes a base64 encoded string and returns the original data.
func Decode(logger interfaces.Logger, location, encodedString string) common.Result[string] {
	decodedBytes, decodeStringError := base64.StdEncoding.DecodeString(encodedString)
	if validator.IsError(decodeStringError) {
		internalError := domainError.NewInternalError(location+".Decode.DecodeString", decodeStringError.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[string](internalError)
	}

	return common.NewResultOnSuccess[string](string(decodedBytes))
}

// SanitizeAndToLowerString trims leading and trailing white spaces from the input string.
func SanitizeAndToLowerString(data string) string {
	return strings.ToLower(strings.TrimSpace(data))
}
