package common

import (
	"encoding/base64"
	"regexp"
	"strings"

	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// regex to match two or more consecutive whitespace characters (spaces, tabs, etc.).
var whitespaceRegex = regexp.MustCompile(`\s{2,}`)

// Encode encodes the input data to a base64 string.
func Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Decode decodes a base64 encoded string and returns the original data.
func Decode(logger interfaces.Logger, location, encodedString string) common.Result[string] {
	decodedBytes, decodeStringError := base64.StdEncoding.DecodeString(encodedString)
	if validator.IsError(decodeStringError) {
		internalError := domain.NewInternalError(location+".Decode.DecodeString", decodeStringError.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[string](internalError)
	}

	return common.NewResultOnSuccess[string](string(decodedBytes))
}

// SanitizeAndCollapseWhitespace trims spaces and collapses multiple spaces into one.
func SanitizeAndCollapseWhitespace(input string) string {
	return strings.TrimSpace(whitespaceRegex.ReplaceAllString(input, " "))
}

// SanitizeAndToLowerString trims leading and trailing white spaces from the input string.
func SanitizeAndToLowerString(data string) string {
	return strings.ToLower(strings.TrimSpace(data))
}

// SanitizeAndCollapseWhitespaceAndToLower trims spaces, collapses multiple spaces, and converts to lowercase.
func SanitizeAndCollapseWhitespaceAndToLower(input string) string {
	return strings.ToLower(strings.TrimSpace(whitespaceRegex.ReplaceAllString(input, " ")))
}
