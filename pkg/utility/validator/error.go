package validator

func IsErrorNil(err error) bool {
	if err != nil {
		return false
	}
	return true
}

func IsErrorNotNil(err error) bool {
	if err != nil {
		return true
	}
	return false
}
