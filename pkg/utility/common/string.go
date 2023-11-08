package common

import (
	"encoding/base64"
	"fmt"

	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location string = "common.string."
)

// Encode encodes the input data to a base64 string.
func Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Decode decodes a base64 encoded string and returns the original data.
func Decode(encodedString string) (string, error) {
	data, decodeStringError := base64.StdEncoding.DecodeString(encodedString)
	if validator.IsErrorNotNil(decodeStringError) {
		decodeStringInternalError := domainError.NewInternalError(location+"Decode.DecodeString", decodeStringError.Error())
		logging.Logger(decodeStringInternalError)
		return "", decodeStringInternalError
	}
	return string(data), nil
}

// ConvertQueryToString converts a query to a string representation.
func ConvertQueryToString(query any) string {
	return fmt.Sprintf("%v", query)
}
