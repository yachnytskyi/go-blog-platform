package validator

func IsStringEmpty(data string) bool {
	if data == "" {
		return true
	}
	return false
}

func IsStringNotEmpty(data string) bool {
	if data == "" {
		return false
	}
	return true
}

func CheckMatchStrings(firstString string, secondString string) bool {
	if firstString == secondString {
		return true
	}
	return false
}
