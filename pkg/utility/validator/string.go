package validator

import (
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
)

func IsStringEmpty(data string) bool {
	return data == constants.EmptyString
}

func IsStringNotEmpty(data string) bool {
	return data != constants.EmptyString
}

func AreStringsEqual(firstString string, secondString string) bool {
	return firstString == secondString
}

func AreStringsNotEqual(firstString string, secondString string) bool {
	return firstString != secondString
}
