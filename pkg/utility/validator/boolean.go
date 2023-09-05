package validator

func IsBooleanNotTrue(data bool) bool {
	if data == true {
		return false
	}
	return true
}
