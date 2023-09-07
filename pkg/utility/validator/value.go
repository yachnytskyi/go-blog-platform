package validator

func IsValueNil[T comparable](data T) bool {
	var t T
	if data == t {
		return true
	}
	return false
}

func IsValueNotNil[T comparable](data T) bool {
	var t T
	if data == t {
		return false 
	}
	return true
}
