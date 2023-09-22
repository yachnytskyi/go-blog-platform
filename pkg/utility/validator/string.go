package validator

func IsStringEmpty(data string) bool {
	return data == ""
}

func IsStringNotEmpty(data string) bool {
	return data != ""
}

func AreStringsEqual(firstString string, secondString string) bool {
	return firstString == secondString
}

func AreStringsNotEqual(firstString string, secondString string) bool {
	return firstString != secondString
}
