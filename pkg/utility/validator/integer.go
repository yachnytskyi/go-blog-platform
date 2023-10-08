package validator

func IsIntegerPositive(data int) bool {
	return data > 0
}

func IsIntegerNegative(data int) bool {
	return data < 0
}

func IsIntegerZero(data int) bool {
	return data == 0
}

func IsIntegerNotZero(data int) bool {
	return data != 0

}

func IsIntegerZeroOrPositive(data int) bool {
	return data >= 0
}

func IsIntegerZeroOrNegative(data int) bool {
	return data <= 0
}
