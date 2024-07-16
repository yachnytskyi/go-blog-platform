package common

import (
	"encoding/base64"
	"fmt"

	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// Encode encodes the input data to a base64 string.
//
// Parameters:
// - data: The string data to be encoded.
//
// Returns:
// - A base64 encoded string.
func Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Decode decodes a base64 encoded string and returns the original data.
//
// Parameters:
// - location: A string representing the location or context for error logging.
// - encodedString: The base64 encoded string to be decoded.
//
// Returns:
// - commonModel.Result[string]: The result containing the decoded string data if successful, or an error if the decoding fails.
func Decode(location, encodedString string) commonModel.Result[string] {
	decodedBytes, decodeStringError := base64.StdEncoding.DecodeString(encodedString)
	if validator.IsError(decodeStringError) {
		// Create and log an internal error with context if decoding fails.
		internalError := domainError.NewInternalError(location+".Decode.DecodeString", decodeStringError.Error())
		logging.Logger(internalError)
		// Return the internal error.
		return commonModel.NewResultOnFailure[string](internalError)
	}

	return commonModel.NewResultOnSuccess[string](string(decodedBytes))
}

// ConvertQueryToString converts a query to a string representation.
//
// Parameters:
// - query: The query to be converted to a string.
//
// Returns:
// - A string representation of the query.
func ConvertQueryToString(query any) string {
	return fmt.Sprintf("%v", query)
}
