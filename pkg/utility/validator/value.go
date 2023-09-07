package validator

func IsValueNil[T comparable](data T) bool {
	var t T
	if data == t {
		return true
	}
	return false
}

func IsValueNotNil(data any) bool {
	if data == nil {
		return false
	}
	return true
}
