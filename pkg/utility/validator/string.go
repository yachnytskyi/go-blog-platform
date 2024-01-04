package validator

import (
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
)

func IsStringNotEmpty(data string) bool {
	return data != constants.EmptyString
}

func AreStringsNotEqual(firstString string, secondString string) bool {
	return firstString != secondString
}
