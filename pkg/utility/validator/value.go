package validator

func IsValueNil[T comparable](data T) bool {
	var t T
	return data == t
}

func IsValueNotNil[T comparable](data T) bool {
	var t T
	return data != t
}
