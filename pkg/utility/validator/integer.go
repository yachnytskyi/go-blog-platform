package validator

func IsIntegerZero(data int) bool {
	if data == 0 {
		return true
	}
	return false
}

func IsIntegerNotZero(data int) bool {
	if data == 0 {
		return false
	}
	return true
}

func IsIntegerZeroOrLess(data int) bool {
	if data == 0 || data < 0 {
		return true
	}
	return false
}

func IsIntegerNotZeroOrLess(data int) bool {
	if data == 0 || data < 0 {
		return false
	}
	return true
}
