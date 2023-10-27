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

func Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
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

func DatabaseQueryToStringMapper(query any) string {
	return fmt.Sprintf("%v", query)
}
