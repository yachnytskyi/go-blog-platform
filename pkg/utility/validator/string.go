package validator

func IsStringEmpty(data string) bool {
	if data != "" {
		return false
	}
	return true
}

func IsStringNotEmpty(data string) bool {
	if data != "" {
		return true
	}
	return false
}
