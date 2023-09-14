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

func AreStringsEqual(firstString string, secondString string) bool {
	if firstString == secondString {
		return true
	}
	return false
}

func AreStringsNotEqual(firstString string, secondString string) bool {
	if firstString == secondString {
		return false
	}
	return true
}
