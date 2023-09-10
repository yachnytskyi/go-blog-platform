package common

import (
	"encoding/base64"

	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location string = "common.string."
)

func Encode(baseString string) string {
	return base64.StdEncoding.EncodeToString([]byte(baseString))
}

func Decode(encodedString string) (string, error) {
	data, decodeStringError := base64.StdEncoding.DecodeString(encodedString)
	if validator.IsErrorNotNil(decodeStringError) {
		decodeStringInternalError := domainError.NewInternalError(location+"Decode.DecodeString", decodeStringError.Error())
		logging.Logger(decodeStringInternalError)
		return "", decodeStringInternalError
	}
	return string(data), nil
}
