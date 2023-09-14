package validator

func IsBooleanTrue(data bool) bool {
	if data == true {
		return true
	}
	return false
}

func IsBooleanNotTrue(data bool) bool {
	if data == true {
		return false
	}
	return true
}
