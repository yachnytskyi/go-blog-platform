package validator

func IsStringNotEmpty(data string) bool {
	return data != ""
}

func AreStringsNotEqual(firstString string, secondString string) bool {
	return firstString != secondString
}
